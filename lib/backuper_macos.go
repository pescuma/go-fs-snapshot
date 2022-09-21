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

	snapshotDates []string
	snapshotPaths []string
	mountPoints   map[string]string
}

func newMacosBackuper(mountPoints map[string]string,
	listMountPoints func(volume string) ([]string, error),
	infoCallback InfoMessageCallback,
) *macosBackuper {

	result := &macosBackuper{}
	result.volumes = newVolumeInfos()
	result.infoCallback = infoCallback
	result.mountPoints = mountPoints

	result.baseBackuper.caseSensitive = true
	result.baseBackuper.listMountPoints = listMountPoints
	result.baseBackuper.createSnapshot = result.createSnapshot
	result.baseBackuper.deleteSnapshot = result.deleteSnapshot

	return result
}

func (b *macosBackuper) createSnapshot(m *mountPointInfo) (string, error) {
	drive := b.mountPoints[m.dir]

	b.infoCallback(DetailsLevel, "Creating local snapshot")

	output, err := runAndReturnOutput(b.infoCallback, "tmutil", "localsnapshot", drive)
	if err != nil {
		m.state = StateFailed
		return "", errors.Errorf("error creating local snapshot: %v", err)
	}

	re := regexp.MustCompile("Created local snapshot with date: ([0-9-]+)")
	matches := re.FindStringSubmatch(output)
	if len(matches) != 2 {
		m.state = StateFailed
		return "", errors.Errorf("unknown tmutil output: %v", output)
	}

	snapshotDate := matches[1]
	b.snapshotDates = append(b.snapshotDates, snapshotDate)

	b.infoCallback(DetailsLevel, "Created local snapshot with date %v", snapshotDate)

	snapshotPath, err := os.MkdirTemp(os.TempDir(), "fs_snapshot_")
	if err != nil {
		m.state = StateFailed
		return "", err
	}

	b.snapshotPaths = append(b.snapshotPaths, snapshotPath)

	b.infoCallback(DetailsLevel, "Mounting snapshot at %v", snapshotPath)

	err = run(b.infoCallback, "mount_apfs", "-o", "ro", "-s", prefix+snapshotDate+suffix, drive, snapshotPath)
	if err != nil {
		m.state = StateFailed
		return "", errors.Errorf("error mounting local snapshot: %v", err)
	}

	return snapshotPath, nil
}

func (b *macosBackuper) deleteSnapshot(m *mountPointInfo) error {
	b.infoCallback(DetailsLevel, "Unmounting snapshot at %v", m.snapshotDir)

	return run(b.infoCallback, "umount", m.snapshotDir)
}

func (b *macosBackuper) Close() {
	b.baseBackuper.close()

	for _, p := range b.snapshotPaths {
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

	b.snapshotPaths = nil
	b.snapshotDates = nil
}
