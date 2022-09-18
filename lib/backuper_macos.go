//go:build darwin

package fs_snapshot

import (
	"os"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/pkg/errors"
)

type macosBackuper struct {
	baseBackuper

	snapshotDates []string
	snapshotPaths []string
}

func newMacosBackuper(infoCallback InfoMessageCallback, listMountPoints func() ([]string, error)) (*macosBackuper, error) {
	result := &macosBackuper{}
	result.volumes = newVolumeInfos()
	result.infoCallback = infoCallback

	result.baseBackuper.caseSensitive = true
	result.baseBackuper.absolutePath = filepath.Abs
	result.baseBackuper.listMountPoints = func(volume string) ([]string, error) {
		if volume != "" {
			return nil, errors.Errorf("unknown volume: %v", volume)
		}

		return listMountPoints()
	}
	result.baseBackuper.createSnapshot = result.createSnapshot
	result.baseBackuper.deleteSnapshot = result.deleteSnapshot

	return result, nil
}

func (b *macosBackuper) createSnapshot(m *mountPointInfo) (string, error) {
	output, err := runAndReturnOutput(b.infoCallback, "tmutil", "localsnapshot", m.path)
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

	snapshotPath, err := os.MkdirTemp(os.TempDir(), "fs_snapshot_")
	if err != nil {
		m.state = StateFailed
		return "", err
	}

	b.snapshotPaths = append(b.snapshotPaths, snapshotPath)

	err = run(b.infoCallback, "mount_apfs", "-o", "ro", "-s", prefix+snapshotDate+suffix, m.path, snapshotPath)
	if err != nil {
		m.state = StateFailed
		return "", errors.Errorf("error mounting local snapshot: %v", err)
	}

	return snapshotPath, nil
}

func (b *macosBackuper) deleteSnapshot(m *mountPointInfo) error {
	return run(b.infoCallback, "umount", m.snapshotPath)
}

func (b *macosBackuper) Close() {
	b.baseBackuper.close()

	for _, p := range b.snapshotPaths {
		err := syscall.Rmdir(p)
		if err != nil {
			b.infoCallback(InfoLevel, "Error removing %v : %v", p, err)
		}
	}

	for _, d := range b.snapshotDates {
		err := run(b.infoCallback, "tmutil", "deletelocalsnapshots", d)
		if err != nil {
			b.infoCallback(InfoLevel, "Error deleting local snapshot %v : %v", d, err)
		}
	}

	b.snapshotPaths = nil
	b.snapshotDates = nil
}