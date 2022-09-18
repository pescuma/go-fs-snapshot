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
	volumes       *volumeInfos
	snapshotDates []string
	snapshotPaths []string
	infoCallback  InfoMessageCallback
}

func (b *macosBackuper) TryToCreateTemporarySnapshot(inputDirectory string) (string, error) {
	dir, err := filepath.Abs(inputDirectory)
	if err != nil {
		return inputDirectory, err
	}

	path, snapshotPath, err := b.computeSnapshotPaths(dir)
	if err != nil {
		return inputDirectory, err
	}

	newDir, err := changeBaseDir(dir, path, snapshotPath)
	if err != nil {
		return inputDirectory, err
	}

	return newDir, nil
}

func (b *macosBackuper) computeSnapshotPaths(dir string) (string, string, error) {
	m := b.volumes.GetMountPoint(dir)

	// First use only a read lock to avoid stopping too much
	m.mutex.RLock()

	path := m.path
	snapshotPath := m.snapshotPath
	state := m.state

	m.mutex.RUnlock()

	switch state {
	case StateSuccess:
		return path, snapshotPath, nil
	case StateFailed:
		return "", "", errors.New("snapshot failed in a previous attempt")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Because we locked again, someone else may have already done it
	switch m.state {
	case StateSuccess:
		return m.path, m.snapshotPath, nil
	case StateFailed:
		return "", "", errors.New("snapshot failed in a previous attempt")
	}

	snapshotPath, err := b.createSnapshot(m)
	if err != nil {
		return "", "", err
	}

	m.state = StateSuccess
	m.snapshotPath = snapshotPath

	return path, snapshotPath, nil
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

func (b *macosBackuper) ListSnapshotedDirectories() map[string]string {
	result := make(map[string]string)

	b.volumes.mutex.RLock()
	defer b.volumes.mutex.RUnlock()

	for _, v := range b.volumes.volumes {
		for _, m := range v {
			m.mutex.RLock()

			if m.state == StateSuccess {
				result[m.path] = m.snapshotPath
			}

			m.mutex.RUnlock()
		}
	}

	return result
}

func (b *macosBackuper) Close() {
	// No locks used because this must be used only once after everything else ended

	for _, v := range b.volumes.volumes {
		for _, m := range v {
			if m.state == StateSuccess {
				_ = run(b.infoCallback, "umount", m.snapshotPath)
			}
		}
	}

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

	b.volumes.volumes = nil
	b.snapshotPaths = nil
	b.snapshotDates = nil
}
