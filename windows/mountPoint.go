//go:build windows

package fs_snapshot_windows

import "github.com/go-ole/go-ole"

// mountPoint wraps all information of a snapshot of a mountpoint on a volume.
type mountPoint struct {
	isSnapshotted        bool
	snapshotSetID        ole.GUID
	snapshotProperties   vssSnapshotProperties
	snapshotDeviceObject string
}

// IsSnapshotted is true if this mount point was snapshotted successfully.
func (p *mountPoint) IsSnapshotted() bool {
	return p.isSnapshotted
}

// GetSnapshotDeviceObject returns root path to access the snapshot files and folders.
func (p *mountPoint) GetSnapshotDeviceObject() string {
	return p.snapshotDeviceObject
}
