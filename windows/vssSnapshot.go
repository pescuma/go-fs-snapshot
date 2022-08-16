//go:build windows

package fs_snapshot_windows

import (
	"fmt"

	"github.com/go-ole/go-ole"
)

// vssSnapshot wraps windows volume shadow copy api (vss) via a simple
// interface to create and delete a vss snapshot.
type vssSnapshot struct {
	iVssBackupComponents *ivssBackupComponents
	snapshotID           ole.GUID
	snapshotProperties   vssSnapshotProperties
	snapshotDeviceObject string
	mountPointInfo       map[string]mountPoint
	timeoutInMillis      uint32
}

// GetSnapshotDeviceObject returns root path to access the snapshot files
// and folders.
func (p *vssSnapshot) GetSnapshotDeviceObject() string {
	return p.snapshotDeviceObject
}

// Delete deletes the created snapshot.
func (p *vssSnapshot) Delete() error {
	var err error
	if err = vssFreeSnapshotProperties(&p.snapshotProperties); err != nil {
		return err
	}

	for _, mountPoint := range p.mountPointInfo {
		if mountPoint.isSnapshotted {
			if err = vssFreeSnapshotProperties(&mountPoint.snapshotProperties); err != nil {
				return err
			}
		}
	}

	if p.iVssBackupComponents != nil {
		defer p.iVssBackupComponents.Release()

		err = callAsyncFunctionAndWait("BackupComplete", p.iVssBackupComponents.BackupComplete, p.timeoutInMillis)
		if err != nil {
			return err
		}

		if _, _, e := p.iVssBackupComponents.DeleteSnapshots(p.snapshotID); e != nil {
			err = newVssTextError(fmt.Sprintf("Failed to delete snapshot: %s", e.Error()))
			p.iVssBackupComponents.AbortBackup()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
