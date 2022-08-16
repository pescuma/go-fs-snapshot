//go:build windows

package fs_snapshot_windows

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

// vssSnapshotProperties defines the properties of a VSS snapshot as part of the VSS api.
type vssSnapshotProperties struct {
	snapshotID           ole.GUID
	snapshotSetID        ole.GUID
	snapshotsCount       uint32
	snapshotDeviceObject *uint16
	originalVolumeName   *uint16
	originatingMachine   *uint16
	serviceMachine       *uint16
	exposedName          *uint16
	exposedPath          *uint16
	providerID           ole.GUID
	snapshotAttributes   uint32
	creationTimestamp    uint64
	status               uint
}

// GetSnapshotDeviceObject returns root path to access the snapshot files
// and folders.
func (p *vssSnapshotProperties) GetSnapshotDeviceObject() string {
	return ole.UTF16PtrToString(p.snapshotDeviceObject)
}

// vssFreeSnapshotProperties calls the equivalent VSS api.
func vssFreeSnapshotProperties(properties *vssSnapshotProperties) error {
	proc, err := findVssProc("VssFreeSnapshotProperties")
	if err != nil {
		return err
	}

	_, _, err = proc.Call(uintptr(unsafe.Pointer(properties)))
	if err != nil {
		return err
	}

	return nil
}
