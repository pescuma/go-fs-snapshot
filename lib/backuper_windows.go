//go:build windows

package fs_snapshot

import (
	"time"

	"github.com/go-ole/go-ole"

	"github.com/pescuma/go-fs-snapshot/lib/internal/windows"
)

type windowsBackuper struct {
	baseBackuper

	opts       *internal_windows.SnapshotOptions
	vssResults []*internal_windows.SnapshotsResult
}

func newWindowsBackuper(providerID *ole.GUID, timeout time.Duration, simple bool,
	listMountPoints func(volume string) ([]string, error),
	infoCallback InfoMessageCallback,
) *windowsBackuper {

	result := &windowsBackuper{}

	result.volumes = newVolumeInfos()
	result.infoCallback = infoCallback

	result.baseBackuper.caseSensitive = false
	result.baseBackuper.listMountPoints = listMountPoints
	result.baseBackuper.createSnapshot = result.createSnapshot
	result.baseBackuper.deleteSnapshot = result.deleteSnapshot

	result.opts = &internal_windows.SnapshotOptions{
		ProviderID: providerID,
		Timeout:    timeout,
		Writters:   !simple,
		InfoCallback: func(level internal_windows.MessageLevel, format string, a ...interface{}) {
			infoCallback(MessageLevel(level), format, a...)
		},
	}

	return result
}

func (b *windowsBackuper) createSnapshot(m *mountPointInfo) (string, error) {
	vsr, err := internal_windows.CreateSnapshots([]string{m.dir}, b.opts)

	b.vssResults = append(b.vssResults, vsr)

	if err != nil {
		return "", err
	}

	return vsr.GetSnapshotPath(m.dir), nil
}

func (b *windowsBackuper) deleteSnapshot(m *mountPointInfo) error {
	return nil
}

func (b *windowsBackuper) Close() {
	b.baseBackuper.close()

	for _, r := range b.vssResults {
		r.Close()
	}

	b.vssResults = nil
}
