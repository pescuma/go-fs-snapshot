//go:build windows

package fs_snapshot

import (
	"strings"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"

	"github.com/pescuma/go-fs-snapshot/lib/fs_snapshot/internal/windows"
)

type windowsBackuper struct {
	baseBackuper

	parent     *windowsSnapshoter
	opts       *internal_windows.SnapshotOptions
	vssResults []*internal_windows.SnapshotsResult
}

func newWindowsBackuper(parent *windowsSnapshoter, providerID *ole.GUID, timeout time.Duration, simple bool, infoCallback InfoMessageCallback) *windowsBackuper {
	result := &windowsBackuper{}
	result.parent = parent

	result.volumes = newVolumeInfos(false)
	result.infoCallback = infoCallback

	result.baseBackuper.listMountPoints = parent.ListMountPoints
	result.baseBackuper.createSnapshot = result.createSnapshot

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

func (b *windowsBackuper) createSnapshot(m *mountPointInfo) (*Snapshot, error) {
	vsr, err := internal_windows.CreateSnapshots([]string{m.dir}, b.opts)

	b.vssResults = append(b.vssResults, vsr)

	if err != nil {
		return nil, err
	}

	props := vsr.GetProperties(m.dir)
	if props == nil {
		return nil, errors.Errorf("snapshots not supported in volume %v", m.dir)
	}

	sb, err := b.parent.newSnapshotsBuilder(nil)
	if err != nil {
		return nil, err
	}

	err = sb.AddSnapshot(props)
	if err != nil {
		return nil, err
	}

	for _, snapshot := range sb.Snapshots {
		if strings.EqualFold(snapshot.OriginalDir, m.dir) {
			return snapshot, nil
		}
	}
	return nil, errors.New("Failed after creating snapshot: original volume not found")
}

func (b *windowsBackuper) Close() {
	for _, r := range b.vssResults {
		r.Close()
	}

	b.vssResults = nil
}
