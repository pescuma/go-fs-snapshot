//go:build windows

package internal_fs_snapshot_windows

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

// VssSnapshotProperties defines the properties of a VSS snapshot as part of the VSS api.
type VssSnapshotProperties struct {
	SnapshotID           ole.GUID
	SnapshotSetID        ole.GUID
	SnapshotsCount       uint32
	SnapshotDeviceObject *uint16
	OriginalVolumeName   *uint16
	OriginatingMachine   *uint16
	ServiceMachine       *uint16
	ExposedName          *uint16
	ExposedPath          *uint16
	ProviderID           ole.GUID
	SnapshotAttributes   VssVolumeSnapshotAttributes
	CreationTimestamp    uint64
	Status               VssSnapshotState
}

// GetSnapshotDeviceObject returns root path to access the snapshot files
// and folders.
func (p *VssSnapshotProperties) GetSnapshotDeviceObject() string {
	return ole.UTF16PtrToString(p.SnapshotDeviceObject)
}

// VssFreeSnapshotProperties calls the equivalent VSS api.
func VssFreeSnapshotProperties(properties *VssSnapshotProperties) error {
	proc, err := findVssProc("VssFreeSnapshotProperties")
	if err != nil {
		return err
	}

	_, _, _ = proc.Call(uintptr(unsafe.Pointer(properties)))

	return nil
}