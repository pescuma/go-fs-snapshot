//go:build darwin

package fs_snapshot

import (
	"os"
	"regexp"
	"syscall"

	"github.com/pkg/errors"
)

type macosBackuper struct {
	baseBackuper

	parent         *macosSnapshoter
	snapshotDates  []string
	snapshotDirs   []string
	snapshotMounts []string
	mountPoints    map[string]string
}

func newMacosBackuper(parent *macosSnapshoter,
	mountPoints map[string]string,
	infoCallback InfoMessageCallback,
) *macosBackuper {

	result := &macosBackuper{}
	result.parent = parent
	result.volumes = newVolumeInfos()
	result.infoCallback = infoCallback
	result.mountPoints = mountPoints

	result.baseBackuper.caseSensitive = true
	result.baseBackuper.listMountPoints = parent.ListMountPoints
	result.baseBackuper.createSnapshot = result.createSnapshot

	return result
}

func (b *macosBackuper) createSnapshot(m *mountPointInfo) (*Snapshot, error) {
	drive := b.mountPoints[m.dir]

	b.infoCallback(DetailsLevel, "Creating local snapshot")

	output, err := runAndReturnOutput(b.infoCallback, "tmutil", "localsnapshot", drive)
	if err != nil {
		return nil, errors.Errorf("error creating local snapshot: %v", err)
	}

	re := regexp.MustCompile("Created local snapshot with date: ([0-9-]+)")
	matches := re.FindStringSubmatch(output)
	if len(matches) != 2 {
		return nil, errors.Errorf("unknown tmutil output: %v", output)
	}

	snapshotDate := matches[1]
	id := prefix + snapshotDate + suffix

	b.snapshotDates = append(b.snapshotDates, snapshotDate)

	b.infoCallback(DetailsLevel, "Created local snapshot with date %v", snapshotDate)

	snapshotDir, err := os.MkdirTemp(os.TempDir(), "fs_snapshot_")
	if err != nil {
		return nil, err
	}

	b.snapshotDirs = append(b.snapshotDirs, snapshotDir)

	b.infoCallback(DetailsLevel, "Mounting snapshot at %v", snapshotDir)

	err = run(b.infoCallback, "mount_apfs", "-o", "rdonly,nobrowse", "-s", id, drive, snapshotDir)
	if err != nil {
		return nil, errors.Errorf("error mounting local snapshot: %v", err)
	}

	b.snapshotMounts = append(b.snapshotMounts, snapshotDir)

	snapshot, err := b.parent.newSnapshot(id, snapshotDate, m.dir, snapshotDir, nil)
	if err != nil {
		return nil, errors.Errorf("error creating snapshot object: %v", err)
	}

	return snapshot, nil
}

func (b *macosBackuper) Close() {
	for _, m := range b.snapshotMounts {
		b.infoCallback(DetailsLevel, "Unmounting snapshot at %v", m)
		err := run(b.infoCallback, "umount", m)
		if err != nil {
			b.infoCallback(InfoLevel, "Error unmounting %v : %v", m, err)
		}
	}

	for _, p := range b.snapshotDirs {
		b.infoCallback(DetailsLevel, "Deleting snapshot mount folder %v", p)
		err := syscall.Rmdir(p)
		if err != nil {
			b.infoCallback(InfoLevel, "Error removing %v : %v", p, err)
		}
	}

	for _, d := range b.snapshotDates {
		b.infoCallback(DetailsLevel, "Deleting local snapshot with date %v", d)
		err := run(b.infoCallback, "tmutil", "deletelocalsnapshots", d)
		if err != nil {
			b.infoCallback(InfoLevel, "Error deleting local snapshot %v : %v", d, err)
		}
	}

	b.snapshotMounts = nil
	b.snapshotDirs = nil
	b.snapshotDates = nil
}
