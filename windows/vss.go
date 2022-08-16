//go:build windows

package fs_snapshot_windows

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"

	fs_snapshot "github.com/pescuma/go-fs-snapshot"
)

// NewVssSnapshot creates a new vss snapshot. If creating the snapshots doesn't
// finish within the timeout an error is returned.
func NewVssSnapshot(
	volume string, timeoutInSeconds uint, msgError fs_snapshot.ErrorHandler) (vssSnapshot, error) {
	is64Bit, err := isRunningOn64BitWindows()

	if err != nil {
		return vssSnapshot{}, newVssTextError(fmt.Sprintf(
			"Failed to detect windows architecture: %s", err.Error()))
	}

	if (is64Bit && runtime.GOARCH != "amd64") || (!is64Bit && runtime.GOARCH != "386") {
		return vssSnapshot{}, newVssTextError(fmt.Sprintf("executables compiled for %s can't use "+
			"VSS on other architectures. Please use an executable compiled for your platform.",
			runtime.GOARCH))
	}

	timeoutInMillis := uint32(timeoutInSeconds * 1000)

	oleIUnknown, err := initializeVssCOMInterface()
	if oleIUnknown != nil {
		defer oleIUnknown.Release()
	}
	if err != nil {
		return vssSnapshot{}, err
	}

	comInterface, err := queryInterface(oleIUnknown, uuid_ivssBackupComponents)
	if err != nil {
		return vssSnapshot{}, err
	}

	/*
		https://microsoft.public.win32.programmer.kernel.narkive.com/aObDj2dD/volume-shadow-copy-backupcomplete-and-vss-e-bad-state

		CreateVSSBackupComponents();
		InitializeForBackup();
		SetBackupState();
		GatherWriterMetadata();
		StartSnapshotSet();
		AddToSnapshotSet();
		PrepareForBackup();
		DoSnapshotSet();
		GetSnapshotProperties();
		<Backup all files>
		VssFreeSnapshotProperties();
		BackupComplete();
	*/

	backupComponents := (*ivssBackupComponents)(unsafe.Pointer(comInterface))

	if err := backupComponents.InitializeForBackup(); err != nil {
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	if err := backupComponents.SetContext(VSS_CTX_BACKUP); err != nil {
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	// see https://techcommunity.microsoft.com/t5/Storage-at-Microsoft/What-is-the-difference-between-VSS-Full-Backup-and-VSS-Copy/ba-p/423575

	if err := backupComponents.SetBackupState(false, false, VSS_BT_COPY, false); err != nil {
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	err = callAsyncFunctionAndWait("GatherWriterMetadata", backupComponents.GatherWriterMetadata, timeoutInMillis)
	if err != nil {
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	if isSupported, err := backupComponents.IsVolumeSupported(volume); err != nil {
		backupComponents.Release()
		return vssSnapshot{}, err
	} else if !isSupported {
		backupComponents.Release()
		return vssSnapshot{}, newVssTextError(fmt.Sprintf("Snapshots are not supported for volume "+
			"%s", volume))
	}

	snapshotSetID, err := backupComponents.StartSnapshotSet()
	if err != nil {
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	if err := backupComponents.AddToSnapshotSet(volume, &snapshotSetID); err != nil {
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	mountPoints, err := enumerateMountedFolders(volume)
	if err != nil {
		backupComponents.Release()
		return vssSnapshot{}, newVssTextError(fmt.Sprintf(
			"failed to enumerate mount points for volume %s: %s", volume, err))
	}

	mountPointInfo := make(map[string]mountPoint)

	for _, mp := range mountPoints {
		// ensure every mountpoint is available even without a valid
		// snapshot because we need to consider this when backing up files
		mountPointInfo[mp] = mountPoint{isSnapshotted: false}

		if isSupported, err := backupComponents.IsVolumeSupported(mp); err != nil {
			continue
		} else if !isSupported {
			continue
		}

		var mountPointSnapshotSetID ole.GUID
		err := backupComponents.AddToSnapshotSet(mp, &mountPointSnapshotSetID)
		if err != nil {
			backupComponents.Release()
			return vssSnapshot{}, err
		}

		mountPointInfo[mp] = mountPoint{
			isSnapshotted: true,
			snapshotSetID: mountPointSnapshotSetID,
		}
	}

	err = callAsyncFunctionAndWait("PrepareForBackup", backupComponents.PrepareForBackup, timeoutInMillis)
	if err != nil {
		// After calling PrepareForBackup one needs to call AbortBackup() before releasing the VSS
		// instance for proper cleanup.
		// It is not neccessary to call BackupComplete before releasing the VSS instance afterwards.
		backupComponents.AbortBackup()
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	err = callAsyncFunctionAndWait("DoSnapshotSet", backupComponents.DoSnapshotSet, timeoutInMillis)
	if err != nil {
		backupComponents.AbortBackup()
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	var snapshotProperties vssSnapshotProperties
	err = backupComponents.GetSnapshotProperties(snapshotSetID, &snapshotProperties)
	if err != nil {
		backupComponents.AbortBackup()
		backupComponents.Release()
		return vssSnapshot{}, err
	}

	for mp, info := range mountPointInfo {
		if !info.isSnapshotted {
			continue
		}

		err := backupComponents.GetSnapshotProperties(info.snapshotSetID,
			&info.snapshotProperties)
		if err != nil {
			msgError(mp, errors.Errorf(
				"VSS error: GetSnapshotProperties() for mount point %s returned error: %s",
				mp, err))
			info.isSnapshotted = false
		} else {
			info.snapshotDeviceObject = info.snapshotProperties.GetSnapshotDeviceObject()
		}

		mountPointInfo[mp] = info
	}

	return vssSnapshot{
		backupComponents,
		snapshotSetID,
		snapshotProperties,
		snapshotProperties.GetSnapshotDeviceObject(),
		mountPointInfo,
		timeoutInMillis,
	}, nil
}

// initializeCOMInterface initialize an instance of the VSS COM api
func initializeVssCOMInterface() (*ole.IUnknown, error) {
	vssInstance, err := loadIVssBackupComponentsConstructor()
	if err != nil {
		return nil, err
	}

	// ensure COM is initialized before use
	err = ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		return nil, newVssTextError("Failed to initialize COM")
	}

	var oleIUnknown *ole.IUnknown
	result, _, _ := vssInstance.Call(uintptr(unsafe.Pointer(&oleIUnknown)))
	hresult := HRESULT(result)

	switch hresult {
	case S_OK:
	case E_ACCESSDENIED:
		return oleIUnknown, newVssError(
			"The caller does not have sufficient backup privileges or is not an administrator",
			hresult)
	default:
		return oleIUnknown, newVssError("Failed to create VSS instance", hresult)
	}

	if oleIUnknown == nil {
		return nil, newVssError("Failed to initialize COM interface", hresult)
	}

	return oleIUnknown, nil
}

// hasSufficientPrivilegesForVSS returns nil if the user is allowed to use VSS.
func hasSufficientPrivilegesForVSS() error {
	oleIUnknown, err := initializeVssCOMInterface()
	if oleIUnknown != nil {
		oleIUnknown.Release()
	}

	return err
}

// loadIVssBackupComponentsConstructor finds the constructor of the VSS api
// inside the VSS dynamic library.
func loadIVssBackupComponentsConstructor() (*windows.LazyProc, error) {
	createInstanceName := "?CreateVssBackupComponents@@YAJPEAPEAVIVssBackupComponents@@@Z"

	if runtime.GOARCH == "386" {
		createInstanceName = "?CreateVssBackupComponents@@YGJPAPAVIVssBackupComponents@@@Z"
	}

	return findVssProc(createInstanceName)
}

// findVssProc find a function with the given name inside the VSS api
// dynamic library.
func findVssProc(procName string) (*windows.LazyProc, error) {
	vssDll := windows.NewLazySystemDLL("VssApi.dll")
	err := vssDll.Load()
	if err != nil {
		return &windows.LazyProc{}, err
	}

	proc := vssDll.NewProc(procName)
	err = proc.Find()
	if err != nil {
		return &windows.LazyProc{}, err
	}

	return proc, nil
}

// enumerateMountedFolders returns all mountpoints on the given volume.
func enumerateMountedFolders(volume string) ([]string, error) {
	var mountedFolders []string

	volumeNamePointer, err := syscall.UTF16PtrFromString(volume)
	if err != nil {
		return mountedFolders, err
	}

	volumeMountPointBuffer := make([]uint16, windows.MAX_LONG_PATH)
	handle, err := windows.FindFirstVolumeMountPoint(volumeNamePointer, &volumeMountPointBuffer[0],
		windows.MAX_LONG_PATH)
	if err != nil {
		// if there are no volumes an error is returned
		return mountedFolders, nil
	}

	defer windows.FindVolumeMountPointClose(handle)

	volumeMountPoint := syscall.UTF16ToString(volumeMountPointBuffer)
	mountedFolders = append(mountedFolders, cleanupVolumeMountPoint(volume, volumeMountPoint))

	for {
		err = windows.FindNextVolumeMountPoint(handle, &volumeMountPointBuffer[0],
			windows.MAX_LONG_PATH)

		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}

			return mountedFolders,
				newVssTextError("FindNextVolumeMountPoint() failed: " + err.Error())
		}

		volumeMountPoint := syscall.UTF16ToString(volumeMountPointBuffer)
		mountedFolders = append(mountedFolders, cleanupVolumeMountPoint(volume, volumeMountPoint))
	}

	return mountedFolders, nil
}

// isRunningOn64BitWindows returns true if running on 64-bit windows.
func isRunningOn64BitWindows() (bool, error) {
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

func cleanupVolumeMountPoint(volume, mp string) string {
	return strings.ToLower(filepath.Join(volume, mp) + string(filepath.Separator))
}
