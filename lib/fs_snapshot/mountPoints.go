package fs_snapshot

import (
	"path/filepath"
	"strings"
	"sync"
)

type volumeInfos struct {
	caseSensitive bool

	mutex   sync.RWMutex
	volumes map[string]map[string]*mountPointInfo
}

type mountPointInfo struct {
	dir string

	mutex    sync.RWMutex
	state    mountPointState
	snapshot *Snapshot
}

type mountPointState int

const (
	StatePending mountPointState = iota
	StateSuccess
	StateFailed
)

func newVolumeInfos(caseSensitive bool) *volumeInfos {
	result := &volumeInfos{
		caseSensitive: caseSensitive,
	}

	result.volumes = make(map[string]map[string]*mountPointInfo)

	return result
}

func (i *volumeInfos) AddVolume(volume string, listMountPoints func(volume string) ([]string, error)) error {
	if !i.caseSensitive {
		volume = strings.ToLower(volume)
	}

	i.mutex.RLock()

	_, ok := i.volumes[volume]

	i.mutex.RUnlock()

	if ok {
		return nil
	}

	i.mutex.Lock()
	defer i.mutex.Unlock()

	_, ok = i.volumes[volume]
	if ok {
		return nil
	}

	ps, err := listMountPoints(volume)
	if err != nil {
		return err
	}

	ms := make(map[string]*mountPointInfo, len(ps))
	for _, p := range ps {
		ms[p] = &mountPointInfo{
			dir:   p,
			state: StatePending,
		}
	}

	i.volumes[volume] = ms

	return nil
}

func (i *volumeInfos) ComputeNeeded(dir string) []*mountPointInfo {
	var result []*mountPointInfo

	volume := i.volumeName(dir)

	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, mount := range i.volumes[volume] {
		insideMount := i.hasPrefix(dir, mount.dir)
		mountInside := i.hasPrefix(mount.dir, dir)

		if !insideMount && !mountInside {
			continue
		}

		if mount.state == StatePending {
			result = append(result, mount)
		}
	}

	return result
}

func (i *volumeInfos) GetMountPoint(dir string) *mountPointInfo {
	var result *mountPointInfo

	volume := i.volumeName(dir)

	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, m := range i.volumes[volume] {
		if i.hasPrefix(dir, m.dir) && (result == nil || len(result.dir) < len(m.dir)) {
			result = m
		}
	}

	return result
}

func (i *volumeInfos) volumeName(dir string) string {
	result := filepath.VolumeName(dir)

	if !i.caseSensitive {
		result = strings.ToLower(result)
	}

	return result
}

func (i *volumeInfos) hasPrefix(s, prefix string) bool {
	if !i.caseSensitive {
		s = strings.ToLower(s)
		prefix = strings.ToLower(prefix)
	}

	return strings.HasPrefix(s, prefix)
}
