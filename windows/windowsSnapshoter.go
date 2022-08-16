//go:build windows

package fs_snapshot_windows

import "github.com/pescuma/go-fs-snapshot"

type windowsSnapshoter struct {
}

func CreateWindowsSnapshoter() (fs_snapshot.Snapshoter, error) {
	if err := hasSufficientPrivilegesForVSS(); err != nil {
		return nil, err
	}

	return &windowsSnapshoter{}, nil
}

func (w windowsSnapshoter) CreateSnapshot(path string) string {
	//TODO implement me
	panic("implement me")
}

func (w windowsSnapshoter) Destroy() error {
	return nil
}
