//go:build windows

package internal_fs_snapshot_windows

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
)

func CreateIVSSBackupComponents() (*IVSSBackupComponents, error) {
	var createInstanceName string
	if runtime.GOARCH == "386" {
		createInstanceName = "?CreateVssBackupComponents@@YGJPAPAVIVssBackupComponents@@@Z"
	} else {
		createInstanceName = "?CreateVssBackupComponents@@YAJPEAPEAVIVssBackupComponents@@@Z"
	}

	vssInstance, err := findVssProc(createInstanceName)
	if err != nil {
		return nil, err
	}

	var oleIUnknown *ole.IUnknown
	result, _, _ := vssInstance.Call(uintptr(unsafe.Pointer(&oleIUnknown)))
	hresult := HRESULT(result)

	switch hresult {
	case S_OK:
	case E_ACCESSDENIED:
		if oleIUnknown != nil {
			oleIUnknown.Release()
		}
		return nil, newVssError(hresult, "The caller does not have sufficient backup privileges or is not an administrator")
	default:
		if oleIUnknown != nil {
			oleIUnknown.Release()
		}
		return nil, newVssError(hresult, "Failed to create VSS instance")
	}

	if oleIUnknown == nil {
		return nil, errors.New("Failed to create VSS instance: received nil")
	}

	comInterface, err := queryInterface(oleIUnknown, uuid_ivssBackupComponents)
	if err != nil {
		oleIUnknown.Release()
		return nil, errors.Errorf("Failed to create VSS instance: %v", err)
	}

	bc := &IVSSBackupComponents{}
	bc.iunknown = oleIUnknown
	bc.com = (*ivssBackupComponentsOle)(unsafe.Pointer(comInterface))

	return bc, nil
}

// uuid_ivssBackupComponents defines the GUID of IVSSBackupComponents.
//
//goland:noinspection GoSnakeCaseUsage
var uuid_ivssBackupComponents = ole.NewGUID("{665c1d5f-c218-414d-a05d-7fef5f9d5c86}")

type IVSSBackupComponents struct {
	com      *ivssBackupComponentsOle
	iunknown *ole.IUnknown
}

// IVSSBackupComponents VSS api interface.
type ivssBackupComponentsOle struct {
	ole.IUnknown
}

// ivssBackupComponentsVTable is the vtable for ivssBackupComponentsOle.
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

// getVTable returns the vtable for IVSSBackupComponents.
func (bc *ivssBackupComponentsOle) getVTable() *ivssBackupComponentsVTable {
	return (*ivssBackupComponentsVTable)(unsafe.Pointer(bc.RawVTable))
}

// Query calls the equivalent VSS api.
func (bc *IVSSBackupComponents) Query(objectType vssObjectType) (*IVssEnumObject, error) {
	var enum *IVssEnumObject
	var err error

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(ole.IID_NULL))

		err = syscallNF("Query()", bc.com.getVTable().query,
			uintptr(unsafe.Pointer(bc.com)), id[0], id[1], id[2], id[3],
			uintptr(VSS_OBJECT_NONE), uintptr(objectType), uintptr(unsafe.Pointer(&enum)))

	} else {
		err = syscallNF("Query()", bc.com.getVTable().query,
			uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(ole.IID_NULL)),
			uintptr(VSS_OBJECT_NONE), uintptr(objectType), uintptr(unsafe.Pointer(&enum)))
	}

	return enum, err
}

// AbortBackup calls the equivalent VSS api.
func (bc *IVSSBackupComponents) AbortBackup() error {
	return syscallN("AbortBackup()", bc.com.getVTable().abortBackup,
		uintptr(unsafe.Pointer(bc.com)))
}

// InitializeForBackup calls the equivalent VSS api.
func (bc *IVSSBackupComponents) InitializeForBackup() error {
	return syscallN("InitializeForBackup()", bc.com.getVTable().initializeForBackup,
		uintptr(unsafe.Pointer(bc.com)), 0)
}

// SetContext calls the equivalent VSS api.
func (bc *IVSSBackupComponents) SetContext(context VssContext) error {
	return syscallN("SetContext()", bc.com.getVTable().setContext,
		uintptr(unsafe.Pointer(bc.com)), uintptr(context))
}

// GatherWriterMetadata calls the equivalent VSS api.
func (bc *IVSSBackupComponents) GatherWriterMetadata() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown

	err := syscallN("GatherWriterMetadata()", bc.com.getVTable().gatherWriterMetadata,
		uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(&oleIUnknown)))

	return bc.convertToVSSAsync("GatherWriterMetadata()", oleIUnknown, err)
}

// FreeWriterMetadata calls the equivalent VSS api.
func (bc *IVSSBackupComponents) FreeWriterMetadata() error {
	return syscallN("IVSSBackupComponents()", bc.com.getVTable().freeWriterMetadata,
		uintptr(unsafe.Pointer(bc.com)))
}

// GatherWriterStatus calls the equivalent VSS api.
func (bc *IVSSBackupComponents) GatherWriterStatus() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown

	err := syscallN("GatherWriterStatus()", bc.com.getVTable().gatherWriterStatus,
		uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(&oleIUnknown)))

	return bc.convertToVSSAsync("GatherWriterStatus()", oleIUnknown, err)
}

// FreeWriterStatus calls the equivalent VSS api.
func (bc *IVSSBackupComponents) FreeWriterStatus() error {
	return syscallN("IVSSBackupComponents()", bc.com.getVTable().freeWriterStatus,
		uintptr(unsafe.Pointer(bc.com)))
}

// IsVolumeSupported calls the equivalent VSS api.
func (bc *IVSSBackupComponents) IsVolumeSupported(providerID *ole.GUID, volumeName string) (bool, error) {
	volumeNamePointer, err := syscall.UTF16PtrFromString(volumeName)
	if err != nil {
		return false, err
	}

	var isSupportedRaw uint32

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(providerID))

		err = syscallN("IsVolumeSupported()", bc.com.getVTable().isVolumeSupported,
			uintptr(unsafe.Pointer(bc.com)), id[0], id[1], id[2], id[3],
			uintptr(unsafe.Pointer(volumeNamePointer)), uintptr(unsafe.Pointer(&isSupportedRaw)))
	} else {
		err = syscallN("IsVolumeSupported()", bc.com.getVTable().isVolumeSupported,
			uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(providerID)),
			uintptr(unsafe.Pointer(volumeNamePointer)), uintptr(unsafe.Pointer(&isSupportedRaw)))
	}

	if err != nil {
		return false, err
	}

	return apiIntToBool(isSupportedRaw), nil
}

// StartSnapshotSet calls the equivalent VSS api.
func (bc *IVSSBackupComponents) StartSnapshotSet() (*ole.GUID, error) {
	var snapshotSetID ole.GUID

	err := syscallN("StartSnapshotSet()", bc.com.getVTable().startSnapshotSet,
		uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(&snapshotSetID)))

	return &snapshotSetID, err
}

// AddToSnapshotSet calls the equivalent VSS api.
func (bc *IVSSBackupComponents) AddToSnapshotSet(providerID *ole.GUID, volumeName string) (*ole.GUID, error) {
	volumeNamePointer, err := syscall.UTF16PtrFromString(volumeName)
	if err != nil {
		return nil, err
	}

	var snapshotID ole.GUID

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(providerID))

		err = syscallN("AddToSnapshotSet()", bc.com.getVTable().addToSnapshotSet,
			uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(volumeNamePointer)),
			id[0], id[1], id[2], id[3], uintptr(unsafe.Pointer(&snapshotID)))
	} else {
		err = syscallN("AddToSnapshotSet()", bc.com.getVTable().addToSnapshotSet,
			uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(volumeNamePointer)),
			uintptr(unsafe.Pointer(providerID)), uintptr(unsafe.Pointer(&snapshotID)))
	}

	return &snapshotID, err
}

// PrepareForBackup calls the equivalent VSS api.
func (bc *IVSSBackupComponents) PrepareForBackup() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown

	err := syscallN("PrepareForBackup()", bc.com.getVTable().prepareForBackup,
		uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(&oleIUnknown)))

	return bc.convertToVSSAsync("PrepareForBackup()", oleIUnknown, err)
}

// SetBackupState calls the equivalent VSS api.
func (bc *IVSSBackupComponents) SetBackupState(selectComponents bool, backupBootableSystemState bool,
	backupType vssBackup, partialFileSupport bool,
) error {
	selectComponentsVal := apiBoolToInt(selectComponents)
	backupBootableSystemStateVal := apiBoolToInt(backupBootableSystemState)
	partialFileSupportVal := apiBoolToInt(partialFileSupport)

	return syscallN("SetBackupState()", bc.com.getVTable().setBackupState,
		uintptr(unsafe.Pointer(bc.com)), uintptr(selectComponentsVal),
		uintptr(backupBootableSystemStateVal), uintptr(backupType), uintptr(partialFileSupportVal))
}

// DoSnapshotSet calls the equivalent VSS api.
func (bc *IVSSBackupComponents) DoSnapshotSet() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown

	err := syscallN("DoSnapshotSet()", bc.com.getVTable().doSnapshotSet,
		uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(&oleIUnknown)))

	return bc.convertToVSSAsync("DoSnapshotSet()", oleIUnknown, err)
}

// DeleteSnapshots calls the equivalent VSS api.
func (bc *IVSSBackupComponents) DeleteSnapshots(objectType vssObjectType, snapshotID *ole.GUID, force bool) (int32, *ole.GUID, error) {
	var deletedSnapshots int32 = 0
	var nonDeletedSnapshotID ole.GUID
	var err error

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(snapshotID))

		err = syscallN("DeleteSnapshots()", bc.com.getVTable().deleteSnapshots,
			uintptr(unsafe.Pointer(bc.com)), id[0], id[1], id[2], id[3],
			uintptr(objectType), uintptr(apiBoolToInt(force)),
			uintptr(unsafe.Pointer(&deletedSnapshots)),
			uintptr(unsafe.Pointer(&nonDeletedSnapshotID)),
		)
	} else {
		err = syscallN("DeleteSnapshots()", bc.com.getVTable().deleteSnapshots,
			uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(snapshotID)),
			uintptr(objectType), uintptr(apiBoolToInt(force)),
			uintptr(unsafe.Pointer(&deletedSnapshots)),
			uintptr(unsafe.Pointer(&nonDeletedSnapshotID)))
	}

	if err != nil && err.(*VssError).HResult == VSS_E_OBJECT_NOT_FOUND {
		return 0, snapshotID, nil
	}

	return deletedSnapshots, &nonDeletedSnapshotID, err
}

// GetSnapshotProperties calls the equivalent VSS api.
func (bc *IVSSBackupComponents) GetSnapshotProperties(snapshotID *ole.GUID, properties *VssSnapshotProperties) error {
	var err error

	if runtime.GOARCH == "386" {
		id := (*[4]uintptr)(unsafe.Pointer(snapshotID))

		err = syscallN("GetSnapshotProperties()", bc.com.getVTable().getSnapshotProperties,
			uintptr(unsafe.Pointer(bc.com)), id[0], id[1], id[2], id[3],
			uintptr(unsafe.Pointer(properties)))
	} else {
		err = syscallN("GetSnapshotProperties()", bc.com.getVTable().getSnapshotProperties,
			uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(snapshotID)),
			uintptr(unsafe.Pointer(properties)))
	}

	return err
}

// BackupComplete calls the equivalent VSS api.
func (bc *IVSSBackupComponents) BackupComplete() (*ivssAsync, error) {
	var oleIUnknown *ole.IUnknown

	err := syscallN("BackupComplete()", bc.com.getVTable().backupComplete,
		uintptr(unsafe.Pointer(bc.com)), uintptr(unsafe.Pointer(&oleIUnknown)))

	return bc.convertToVSSAsync("BackupComplete()", oleIUnknown, err)
}

// convertToVSSAsync looks up ivssAsync interface if given result is a success.
func (bc *IVSSBackupComponents) convertToVSSAsync(name string, oleIUnknown *ole.IUnknown, err error) (*ivssAsync, error) {
	if err != nil {
		return nil, err
	}

	comInterface, err := queryInterface(oleIUnknown, uiid_ivssAsync)
	if err != nil {
		return nil, errors.New(name + ": " + err.Error())
	}

	result := (*ivssAsync)(unsafe.Pointer(comInterface))
	if result == nil {
		return nil, errors.New(name + ": conversion to IVSSAsync returned nil")
	}

	return result, nil
}

func (bc *IVSSBackupComponents) Close() {
	if bc == nil {
		return
	}

	bc.com.Release()
	bc.iunknown.Release()
}
