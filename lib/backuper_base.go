package fs_snapshot

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type baseBackuper struct {
	volumes      *volumeInfos
	infoCallback InfoMessageCallback

	caseSensitive   bool
	listMountPoints func(volume string) ([]string, error)
	createSnapshot  func(m *mountPointInfo) (string, error)
	deleteSnapshot  func(m *mountPointInfo) error
}

func (b *baseBackuper) TryToCreateTemporarySnapshot(inputDirectory string) (string, error) {
	dir, err := absolutePath(inputDirectory)
	if err != nil {
		return inputDirectory, err
	}

	dir = addPathSeparatorAsSuffix(dir)

	s, err := os.Stat(dir)
	if err != nil {
		return inputDirectory, err
	}

	if !s.IsDir() {
		return inputDirectory, errors.New("only able to snapshot directories")
	}

	if !b.caseSensitive {
		dir = strings.ToLower(dir)
	}

	volume := filepath.VolumeName(dir)

	err = b.volumes.AddVolume(volume, func(volume string) ([]string, error) {
		mps, err := b.listMountPoints(volume)
		if err != nil {
			return nil, err
		}

		for i, m := range mps {
			if !b.caseSensitive {
				m = strings.ToLower(m)
			}
			mps[i] = addPathSeparatorAsSuffix(m)
		}

		return mps, nil
	})
	if err != nil {
		return inputDirectory, err
	}

	path, snapshotPath, err := b.computeSnapshotPath(dir)
	if err != nil {
		return inputDirectory, err
	}

	newDir, err := changeBaseDir(dir, path, snapshotPath)
	if err != nil {
		return inputDirectory, err
	}

	newDir = addPathSeparatorAsSuffix(newDir)

	return newDir, nil
}

func (b *baseBackuper) computeSnapshotPath(dir string) (string, string, error) {
	m := b.volumes.GetMountPoint(dir)

	// First use only a read lock to avoid stopping too much
	m.mutex.RLock()

	path := m.dir
	snapshotPath := m.snapshotDir
	state := m.state

	m.mutex.RUnlock()

	switch state {
	case StateSuccess:
		return path, snapshotPath, nil
	case StateFailed:
		return "", "", ErrorSnapshotFailedInPreviousAttempt
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Because we locked again, someone else may have already done it
	switch m.state {
	case StateSuccess:
		return path, m.snapshotDir, nil
	case StateFailed:
		return "", "", ErrorSnapshotFailedInPreviousAttempt
	}

	snapshotPath, err := b.createSnapshot(m)

	if err != nil {
		m.state = StateFailed
		return "", "", err
	}

	m.state = StateSuccess
	m.snapshotDir = snapshotPath

	return path, snapshotPath, nil
}

func (b *baseBackuper) ListSnapshotedDirectories() map[string]string {
	result := make(map[string]string)

	b.volumes.mutex.RLock()
	defer b.volumes.mutex.RUnlock()

	for _, v := range b.volumes.volumes {
		for _, m := range v {
			m.mutex.RLock()

			switch m.state {
			case StateSuccess:
				result[m.dir] = m.snapshotDir
			case StateFailed:
				result[m.dir] = m.dir
			case StatePending:
				result[m.dir] = ""
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
