//go:build windows

package internal_windows

import (
	"fmt"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

// IsRunningOn64BitWindows returns true if running on 64-bit windows.
func IsRunningOn64BitWindows() (bool, error) {
	if runtime.GOARCH == "amd64" {
		return true, nil
	}

	isWow64 := false
	err := windows.IsWow64Process(windows.CurrentProcess(), &isWow64)
	if err != nil {
		return false, err
	}

	return isWow64, nil
}

// InitializeCOM initializes COM in this process.
func InitializeCOM() error {
	if err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED); err != nil {
		// CoInitializeEx returns 1 if COM is already initialized
		if oleErr, ok := err.(*ole.OleError); !ok || oleErr.Code() != 1 {
			return errors.New("Failed to initialize COM")
		}
	}

	// https://docs.microsoft.com/en-us/windows/win32/vss/security-considerations-for-requestors
	modole32 := windows.NewLazySystemDLL("ole32.dll")
	procCoInitializeSecurity := modole32.NewProc("CoInitializeSecurity")
	hri, _, _ := procCoInitializeSecurity.Call(
		uintptr(0),          //  PSECURITY_DESCRIPTOR
		uintptr(0xFFFFFFFF), //  LONG
		uintptr(0),          //  SOLE_AUTHENTICATION_SERVICE
		uintptr(0),          //  void
		uintptr(6),          //  RPC_C_AUTHN_LEVEL_PKT_INTEGRITY DWORD
		uintptr(3),          //  RPC_C_IMP_LEVEL_IMPERSONATE DWORD
		uintptr(0),          //  void
		uintptr(0x20),       //  EOAC_STATIC DWORD
		uintptr(0),          //  void
	)
	hr := HRESULT(hri)
	if hr != S_OK {
		return errors.Errorf("Failed to initialize COM security (%v)", hr)
	}

	return nil
}

// SnapshotOptions needed to be duplicated here to avoid cyclical dependencies
type SnapshotOptions struct {
	ProviderID   *ole.GUID
	Timeout      time.Duration
	Writters     bool
	InfoCallback InfoMessageCallback
}

type InfoMessageCallback func(level MessageLevel, msg string)

type MessageLevel int

const (
	InfoLevel = iota
	DetailsLevel
	TraceLevel
)

// CreateSnapshots creates snapshots for the volumes especified.
// Even on fail a result is returned and must be closed.
func CreateSnapshots(volumes []string, opts *SnapshotOptions) (*SnapshotsResult, error) {
	if opts.ProviderID == nil {
		opts.ProviderID = ole.IID_NULL
	}
	if opts.Timeout == 0 {
		opts.Timeout = 2 * time.Minute
	}

	// https://docs.microsoft.com/en-us/windows/win32/vss/overview-of-processing-a-backup-under-vss
	// https://microsoft.public.win32.programmer.kernel.narkive.com/aObDj2dD/volume-shadow-copy-backupcomplete-and-vss-e-bad-state

	var err error

	r := SnapshotsResult{
		opts:    opts,
		volumes: make(map[string]*volumeSnapshotInfo),
	}
	start := time.Now()

	opts.InfoCallback(TraceLevel, "NewIVSSBackupComponents()")
	r.bc, err = NewIVSSBackupComponents()
	if err != nil {
		return &r, err
	}

	opts.InfoCallback(TraceLevel, "InitializeForBackup()")
	err = r.bc.InitializeForBackup()
	if err != nil {
		return &r, err
	}

	if opts.Writters {
		opts.InfoCallback(TraceLevel, "SetContext(VSS_CTX_BACKUP)")
		err = r.bc.SetContext(VSS_CTX_BACKUP)
	} else {
		opts.InfoCallback(TraceLevel, "SetContext(VSS_CTX_FILE_SHARE_BACKUP)")
		err = r.bc.SetContext(VSS_CTX_FILE_SHARE_BACKUP)
	}
	if err != nil {
		return &r, err
	}

	// see https://techcommunity.microsoft.com/t5/Storage-at-Microsoft/What-is-the-difference-between-VSS-Full-Backup-and-VSS-Copy/ba-p/423575
	opts.InfoCallback(TraceLevel, "SetBackupState(false, false, VSS_BT_COPY, false)")
	err = r.bc.SetBackupState(false, false, VSS_BT_COPY, false)
	if err != nil {
		return &r, err
	}

	if opts.Writters {
		opts.InfoCallback(TraceLevel, "GatherWriterMetadata()")
		err = callAndWait(r.bc.GatherWriterMetadata, opts.Timeout-time.Since(start))
		if err != nil {
			return &r, err
		}

		opts.InfoCallback(TraceLevel, "FreeWriterMetadata()")
		err = r.bc.FreeWriterMetadata()
		if err != nil {
			return &r, err
		}
	}

	var atLeastOneVolumeSupported = false
	for _, volume := range volumes {
		s, err := r.bc.IsVolumeSupported(r.opts.ProviderID, volume)
		if err != nil {
			return &r, err
		}

		if s {
			atLeastOneVolumeSupported = true
		} else {
			opts.InfoCallback(DetailsLevel, fmt.Sprintf("Snapshots not supported in volume %v", volume))
		}
	}

	if !atLeastOneVolumeSupported {
		opts.InfoCallback(TraceLevel, "Aboting snapshot because there is no supported volume")
		return &r, nil
	}

	opts.InfoCallback(TraceLevel, "StartSnapshotSet()")
	r.setID, err = r.bc.StartSnapshotSet()
	if err != nil {
		return &r, err
	}

	opts.InfoCallback(TraceLevel, fmt.Sprintf("Set ID: %v", r.setID))

	for _, volume := range volumes {
		info := &volumeSnapshotInfo{}
		r.volumes[volume] = info

		s, err := r.bc.IsVolumeSupported(r.opts.ProviderID, volume)
		if err != nil {
			return &r, err
		}

		if s {
			opts.InfoCallback(TraceLevel, fmt.Sprintf("AddToSnapshotSet(%v, %v)", r.opts.ProviderID, volume))
			info.id, err = r.bc.AddToSnapshotSet(r.opts.ProviderID, volume)
			if err != nil {
				return &r, err
			}

			opts.InfoCallback(TraceLevel, fmt.Sprintf("Volume %v snapshot ID: %v", volume, info.id))
		}
	}

	if opts.Writters {
		opts.InfoCallback(TraceLevel, "PrepareForBackup()")
		err = callAndWait(r.bc.PrepareForBackup, opts.Timeout-time.Since(start))
		r.prepareForBackupCalled = true
		if err != nil {
			return &r, err
		}

		opts.InfoCallback(TraceLevel, "GatherWriterStatus()")
		err = callAndWait(r.bc.GatherWriterStatus, opts.Timeout-time.Since(start))
		if err != nil {
			return &r, err
		}

		opts.InfoCallback(TraceLevel, "FreeWriterStatus()")
		err = r.bc.FreeWriterStatus()
		if err != nil {
			return &r, err
		}
	}

	opts.InfoCallback(TraceLevel, "DoSnapshotSet()")
	err = callAndWait(r.bc.DoSnapshotSet, opts.Timeout-time.Since(start))
	r.doSnapshotSetCalled = true
	if err != nil {
		return &r, err
	}

	if opts.Writters {
		opts.InfoCallback(TraceLevel, "GatherWriterStatus()")
		err = callAndWait(r.bc.GatherWriterStatus, opts.Timeout-time.Since(start))
		if err != nil {
			return &r, err
		}

		opts.InfoCallback(TraceLevel, "FreeWriterStatus()")
		err = r.bc.FreeWriterStatus()
		if err != nil {
			return &r, err
		}
	}

	for _, volume := range volumes {
		info := r.volumes[volume]

		if info.id == nil {
			continue
		}

		opts.InfoCallback(TraceLevel, fmt.Sprintf("GetSnapshotProperties(%v)", info.id))
		var properties VssSnapshotProperties
		err = r.bc.GetSnapshotProperties(info.id, &properties)
		if err != nil {
			return &r, err
		}

		info.properties = &properties
	}

	return &r, nil
}

type SnapshotsResult struct {
	opts                   *SnapshotOptions
	bc                     *IVSSBackupComponents
	setID                  *ole.GUID
	volumes                map[string]*volumeSnapshotInfo
	prepareForBackupCalled bool
	doSnapshotSetCalled    bool
}
type volumeSnapshotInfo struct {
	id         *ole.GUID
	properties *VssSnapshotProperties
}

func (r *SnapshotsResult) GetSnapshotPath(volume string) string {
	info := r.volumes[volume]

	if info.properties == nil {
		return volume
	}

	return info.properties.GetSnapshotDeviceObject()
}

func (r *SnapshotsResult) Close() {
	for _, volume := range r.volumes {
		if volume.properties != nil {
			r.opts.InfoCallback(TraceLevel, fmt.Sprintf("VssFreeSnapshotProperties(%v)", volume.id))
			_ = VssFreeSnapshotProperties(volume.properties)
		}
	}

	if r.doSnapshotSetCalled {
		if r.opts.Writters {
			// Use the full timeout here to at least all the methods once

			r.opts.InfoCallback(TraceLevel, "BackupComplete()")
			_ = callAndWait(r.bc.BackupComplete, r.opts.Timeout)

			r.opts.InfoCallback(TraceLevel, "GatherWriterStatus()")
			_ = callAndWait(r.bc.GatherWriterStatus, r.opts.Timeout)

			r.opts.InfoCallback(TraceLevel, "FreeWriterStatus()")
			_ = r.bc.FreeWriterStatus()
		}

	} else if r.prepareForBackupCalled {
		r.opts.InfoCallback(TraceLevel, "AbortBackup()")
		_ = r.bc.AbortBackup()
	}

	if r.setID != nil {
		r.opts.InfoCallback(TraceLevel, fmt.Sprintf("DeleteSnapshots(VSS_OBJECT_SNAPSHOT_SET, %v, true)", r.setID))
		_, _, _ = r.bc.DeleteSnapshots(VSS_OBJECT_SNAPSHOT_SET, r.setID, true)
	}

	r.bc.Close()
}

// EnumerateMountedFolders returns all mountpoints on the given volume.
func EnumerateMountedFolders(volume string) ([]string, error) {
	var result []string

	volumeNamePointer, err := syscall.UTF16PtrFromString(volume)
	if err != nil {
		return result, err
	}

	buffer := make([]uint16, windows.MAX_LONG_PATH)
	handle, err := windows.FindFirstVolumeMountPoint(volumeNamePointer, &buffer[0], windows.MAX_LONG_PATH)
	if err == windows.ERROR_NO_MORE_FILES {
		// if there are no volumes an error is returned
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	defer windows.FindVolumeMountPointClose(handle)

	for {
		volumeMountPoint := syscall.UTF16ToString(buffer)
		volumeMountPoint = filepath.Join(volume, volumeMountPoint) + `\`
		result = append(result, volumeMountPoint)

		err = windows.FindNextVolumeMountPoint(handle, &buffer[0], windows.MAX_LONG_PATH)

		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}

			return result, errors.New("FindNextVolumeMountPoint() failed: " + err.Error())
		}
	}

	return result, nil
}

// findVssProc find a function with the given name inside the VSS api
// dynamic library.
func findVssProc(procName string) (*windows.LazyProc, error) {
	vssDll := windows.NewLazySystemDLL("VssApi.dll")
	err := vssDll.Load()
	if err != nil {
		return nil, errors.New("Could not find VssApi.dll: " + err.Error())
	}

	proc := vssDll.NewProc(procName)
	err = proc.Find()
	if err != nil {
		return nil, errors.New("Could not find " + procName + " in VssApi.dll: " + err.Error())
	}

	return proc, nil
}
