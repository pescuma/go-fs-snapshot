//go:build windows

package fs_snapshot_windows

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

// uuid_ivssBackupComponents defines the GUID of ivssBackupComponents.
//goland:noinspection GoSnakeCaseUsage
var uuid_ivssBackupComponents = ole.NewGUID("{665c1d5f-c218-414d-a05d-7fef5f9d5c86}")

// ivssBackupComponents VSS api interface.
type ivssBackupComponents struct {
	ole.IUnknown
}

// ivssBackupComponentsVTable is the vtable for ivssBackupComponents.
type ivssBackupComponentsVTable struct {
	ole.IUnknownVtbl
	getWriterComponentsCount      uintptr
	getWriterComponents           uintptr
	initializeForBackup           uintptr
	setBackupState                uintptr
	initializeForRestore          uintptr
	setRestoreState               uintptr
	gatherWriterMetadata          uintptr
	getWriterMetadataCount        uintptr
	getWriterMetadata             uintptr
	freeWriterMetadata            uintptr
	addComponent                  uintptr
	prepareForBackup              uintptr
	abortBackup                   uintptr
	gatherWriterStatus            uintptr
	getWriterStatusCount          uintptr
	freeWriterStatus              uintptr
	getWriterStatus               uintptr
	setBackupSucceeded            uintptr
	setBackupOptions              uintptr
	setSelectedForRestore         uintptr
	setRestoreOptions             uintptr
	setAdditionalRestores         uintptr
	setPreviousBackupStamp        uintptr
	saveAsXML                     uintptr
	backupComplete                uintptr
	addAlternativeLocationMapping uintptr
	addRestoreSubcomponent        uintptr
	setFileRestoreStatus          uintptr
	addNewTarget                  uintptr
	setRangesFilePath             uintptr
	preRestore                    uintptr
	postRestore                   uintptr
	setContext                    uintptr
	startSnapshotSet              uintptr
	addToSnapshotSet              uintptr
	doSnapshotSet                 uintptr
	deleteSnapshots               uintptr
	importSnapshots               uintptr
	breakSnapshotSet              uintptr
	getSnapshotProperties         uintptr
	query                         uintptr
	isVolumeSupported             uintptr
	disableWriterClasses          uintptr
	enableWriterClasses           uintptr
	disableWriterInstances        uintptr
	exposeSnapshot                uintptr
	revertToSnapshot              uintptr
	queryRevertStatus             uintptr
}

// getVTable returns the vtable for ivssBackupComponents.
func (vss *ivssBackupComponents) getVTable() *ivssBackupComponentsVTable {
	return (*ivssBackupComponentsVTable)(unsafe.Pointer(vss.RawVTable))
}

// AbortBackup calls the equivalent VSS api.
func (vss *ivssBackupComponents) AbortBackup() error {
	return syscallN("AbortBackup()", vss.getVTable().abortBackup,
		uintptr(unsafe.Pointer(vss)))
}

// InitializeForBackup calls the equivalent VSS api.
func (vss *ivssBackupComponents) InitializeForBackup() error {
	return syscallN("InitializeForBackup()", vss.getVTable().initializeForBackup,
		uintptr(unsafe.Pointer(vss)), 0)
}

// SetContext calls the equivalent VSS api.
func (vss *ivssBackupComponents) SetContext(context vssContext) error {
	return syscallN("SetContext()", vss.getVTable().setContext,
		uintptr(unsafe.Pointer(vss)), uintptr(context))
}

// GatherWriterMetadata calls the equivalent VSS api.
func (vss *ivssBackupComponents) GatherWriterMetadata() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown
	err := syscallN("GatherWriterMetadata()", vss.getVTable().gatherWriterMetadata,
		uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(&oleIUnknown)))
	return vss.convertToVSSAsync(oleIUnknown, err)
}

// convertToVSSAsync looks up ivssAsync interface if given result
// is a success.
func (vss *ivssBackupComponents) convertToVSSAsync(oleIUnknown *ole.IUnknown, err error) (*ivssAsync, error) {
	if err != nil {
		return nil, err
	}

	comInterface, err := queryInterface(oleIUnknown, uiid_ivssAsync)
	if err != nil {
		return nil, err
	}

	result := (*ivssAsync)(unsafe.Pointer(comInterface))
	return result, nil
}

// IsVolumeSupported calls the equivalent VSS api.
func (vss *ivssBackupComponents) IsVolumeSupported(volumeName string) (bool, error) {
	volumeNamePointer, err := syscall.UTF16PtrFromString(volumeName)
	if err != nil {
		panic(err)
	}

	var isSupportedRaw uint32
	var result uintptr

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(ole.IID_NULL))

		result, _, _ = syscall.SyscallN(vss.getVTable().isVolumeSupported,
			uintptr(unsafe.Pointer(vss)), id[0], id[1], id[2], id[3],
			uintptr(unsafe.Pointer(volumeNamePointer)), uintptr(unsafe.Pointer(&isSupportedRaw)))
	} else {
		result, _, _ = syscall.SyscallN(vss.getVTable().isVolumeSupported,
			uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(ole.IID_NULL)),
			uintptr(unsafe.Pointer(volumeNamePointer)), uintptr(unsafe.Pointer(&isSupportedRaw)))
	}

	var isSupported bool
	if isSupportedRaw == 0 {
		isSupported = false
	} else {
		isSupported = true
	}

	return isSupported, newVssErrorIfResultNotOK("IsVolumeSupported() failed", HRESULT(result))
}

// StartSnapshotSet calls the equivalent VSS api.
func (vss *ivssBackupComponents) StartSnapshotSet() (ole.GUID, error) {
	var snapshotSetID ole.GUID
	result, _, _ := syscall.SyscallN(vss.getVTable().startSnapshotSet,
		uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(&snapshotSetID)))

	return snapshotSetID, newVssErrorIfResultNotOK("StartSnapshotSet() failed", HRESULT(result))
}

// AddToSnapshotSet calls the equivalent VSS api.
func (vss *ivssBackupComponents) AddToSnapshotSet(volumeName string, idSnapshot *ole.GUID) error {
	volumeNamePointer, err := syscall.UTF16PtrFromString(volumeName)
	if err != nil {
		panic(err)
	}

	var result uintptr = 0

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(ole.IID_NULL))

		result, _, _ = syscall.SyscallN(vss.getVTable().addToSnapshotSet,
			uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(volumeNamePointer)), id[0], id[1],
			id[2], id[3], uintptr(unsafe.Pointer(idSnapshot)))
	} else {
		result, _, _ = syscall.SyscallN(vss.getVTable().addToSnapshotSet,
			uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(volumeNamePointer)),
			uintptr(unsafe.Pointer(ole.IID_NULL)), uintptr(unsafe.Pointer(idSnapshot)))
	}

	return newVssErrorIfResultNotOK("AddToSnapshotSet() failed", HRESULT(result))
}

// PrepareForBackup calls the equivalent VSS api.
func (vss *ivssBackupComponents) PrepareForBackup() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown
	result, _, _ := syscall.SyscallN(vss.getVTable().prepareForBackup,
		uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(&oleIUnknown)))

	err := newVssErrorIfResultNotOK("PrepareForBackup() failed", HRESULT(result))
	return vss.convertToVSSAsync(oleIUnknown, err)
}

// SetBackupState calls the equivalent VSS api.
func (vss *ivssBackupComponents) SetBackupState(selectComponents bool,
	backupBootableSystemState bool, backupType vssBackup, partialFileSupport bool,
) error {
	selectComponentsVal := apiBoolToInt(selectComponents)
	backupBootableSystemStateVal := apiBoolToInt(backupBootableSystemState)
	partialFileSupportVal := apiBoolToInt(partialFileSupport)

	result, _, _ := syscall.SyscallN(vss.getVTable().setBackupState,
		uintptr(unsafe.Pointer(vss)), uintptr(selectComponentsVal),
		uintptr(backupBootableSystemStateVal), uintptr(backupType), uintptr(partialFileSupportVal))

	return newVssErrorIfResultNotOK("SetBackupState() failed", HRESULT(result))
}

// DoSnapshotSet calls the equivalent VSS api.
func (vss *ivssBackupComponents) DoSnapshotSet() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown
	result, _, _ := syscall.SyscallN(vss.getVTable().doSnapshotSet,
		uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(&oleIUnknown)))

	err := newVssErrorIfResultNotOK("DoSnapshotSet() failed", HRESULT(result))
	return vss.convertToVSSAsync(oleIUnknown, err)
}

// DeleteSnapshots calls the equivalent VSS api.
func (vss *ivssBackupComponents) DeleteSnapshots(snapshotID ole.GUID) (int32, ole.GUID, error) {
	var deletedSnapshots int32 = 0
	var nondeletedSnapshotID ole.GUID
	var result uintptr = 0

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(&snapshotID))

		result, _, _ = syscall.SyscallN(vss.getVTable().deleteSnapshots,
			uintptr(unsafe.Pointer(vss)), id[0], id[1], id[2], id[3],
			uintptr(VSS_OBJECT_SNAPSHOT), uintptr(1), uintptr(unsafe.Pointer(&deletedSnapshots)),
			uintptr(unsafe.Pointer(&nondeletedSnapshotID)),
		)
	} else {
		result, _, _ = syscall.SyscallN(vss.getVTable().deleteSnapshots,
			uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(&snapshotID)),
			uintptr(VSS_OBJECT_SNAPSHOT), uintptr(1), uintptr(unsafe.Pointer(&deletedSnapshots)),
			uintptr(unsafe.Pointer(&nondeletedSnapshotID)))
	}

	err := newVssErrorIfResultNotOK("DeleteSnapshots() failed", HRESULT(result))
	return deletedSnapshots, nondeletedSnapshotID, err
}

// GetSnapshotProperties calls the equivalent VSS api.
func (vss *ivssBackupComponents) GetSnapshotProperties(snapshotID ole.GUID,
	properties *vssSnapshotProperties) error {
	var result uintptr = 0

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(&snapshotID))

		result, _, _ = syscall.SyscallN(vss.getVTable().getSnapshotProperties,
			uintptr(unsafe.Pointer(vss)), id[0], id[1], id[2], id[3],
			uintptr(unsafe.Pointer(properties)))
	} else {
		result, _, _ = syscall.SyscallN(vss.getVTable().getSnapshotProperties,
			uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(&snapshotID)),
			uintptr(unsafe.Pointer(properties)))
	}

	return newVssErrorIfResultNotOK("GetSnapshotProperties() failed", HRESULT(result))
}

// BackupComplete calls the equivalent VSS api.
func (vss *ivssBackupComponents) BackupComplete() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown
	result, _, _ := syscall.SyscallN(vss.getVTable().backupComplete,
		uintptr(unsafe.Pointer(vss)), uintptr(unsafe.Pointer(&oleIUnknown)))

	err := newVssErrorIfResultNotOK("BackupComplete() failed", HRESULT(result))
	return vss.convertToVSSAsync(oleIUnknown, err)
}
