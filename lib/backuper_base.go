package fs_snapshot

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type baseBackuper struct {
	volumes      *volumeInfos
	infoCallback InfoMessageCallback

	listMountPoints func(volume string) ([]string, error)
	createSnapshot  func(m *mountPointInfo) (string, error)
	deleteSnapshot  func(m *mountPointInfo) error
}

func (b *baseBackuper) TryToCreateTemporarySnapshot(inputDirectory string) (string, error) {
	dir, err := filepath.Abs(inputDirectory)
	if err != nil {
		return inputDirectory, err
	}

	volume := strings.ToLower(filepath.VolumeName(dir))

	err = b.volumes.AddVolume(volume, b.listMountPoints)
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

func (b *baseBackuper) computeSnapshotPaths(dir string) (string, string, error) {
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

func (b *baseBackuper) ListSnapshotedDirectories() map[string]string {
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

func (b *baseBackuper) close() {
	// No locks used because this must be used only once after everything else ended

	for _, v := range b.volumes.volumes {
		for _, m := range v {
			if m.state == StateSuccess {
				_ = b.deleteSnapshot(m)
			}
		}
	}

	b.volumes.volumes = nil
}
