package fs_snapshot

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type baseBackuper struct {
	volumes      *volumeInfos
	infoCallback InfoMessageCallback

	listMountPoints func(volume string) ([]string, error)
	createSnapshot  func(m *mountPointInfo) (*Snapshot, error)
}

func (b *baseBackuper) TryToCreateTemporarySnapshot(inputDirectory string) (string, *Snapshot, error) {
	dir, err := absolutePath(inputDirectory)
	if err != nil {
		return inputDirectory, nil, err
	}

	dir = addPathSeparatorAsSuffix(dir)

	s, err := os.Stat(dir)
	if err != nil {
		return inputDirectory, nil, err
	}

	if !s.IsDir() {
		return inputDirectory, nil, errors.New("only able to snapshot directories")
	}

	volume := filepath.VolumeName(dir)

	err = b.volumes.AddVolume(volume, func(volume string) ([]string, error) {
		mps, err := b.listMountPoints(volume)
		if err != nil {
			return nil, err
		}

		for i, m := range mps {
			mps[i] = addPathSeparatorAsSuffix(m)
		}

		return mps, nil
	})
	if err != nil {
		return inputDirectory, nil, err
	}

	snapshot, err := b.getOrCreateSnapshot(dir)
	if err != nil {
		return inputDirectory, nil, err
	}

	newDir, err := changeBaseDir(dir, snapshot.OriginalDir, snapshot.SnapshotDir)
	if err != nil {
		return inputDirectory, nil, err
	}

	newDir = addPathSeparatorAsSuffix(newDir)

	return newDir, snapshot, nil
}

func (b *baseBackuper) getOrCreateSnapshot(dir string) (*Snapshot, error) {
	m := b.volumes.GetMountPoint(dir)

	// First use only a read lock to avoid stopping too much
	m.mutex.RLock()

	state := m.state

	m.mutex.RUnlock()

	switch state {
	case StateSuccess:
		return m.snapshot, nil
	case StateFailed:
		return nil, ErrSnapshotFailedInPreviousAttempt
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Because we locked again, someone else may have already done it
	switch m.state {
	case StateSuccess:
		return m.snapshot, nil
	case StateFailed:
		return nil, ErrSnapshotFailedInPreviousAttempt
	}

	snapshot, err := b.createSnapshot(m)
	if err != nil {
		m.state = StateFailed
		return nil, err
	}

	m.state = StateSuccess
	m.snapshot = snapshot

	return m.snapshot, nil
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
				result[m.dir] = m.snapshot.SnapshotDir
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
