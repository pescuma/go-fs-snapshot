//go:build windows

package fs_snapshot

import "github.com/pescuma/go-fs-snapshot/windows"

// CreateSnapshoter creates a new snapshoter.
// In case of error a null snapshoter is returned, so you can use it without problem.
func CreateSnapshoter() (Snapshoter, error) {
	sn, err := fs_snapshot_windows.CreateWindowsSnapshoter()
	if err != nil {
		return createNullSnapshoter(), err
	}

	return sn, nil
}
