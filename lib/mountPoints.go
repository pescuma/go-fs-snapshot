package fs_snapshot

import (
	"path/filepath"
	"strings"
	"sync"
)

type volumeInfos struct {
	mutex   sync.RWMutex
	volumes map[string]map[string]*mountPointInfo
}

type mountPointInfo struct {
	mutex        sync.RWMutex
	path         string
	state        mountPointState
	snapshotPath string
}

type mountPointState int

const (
	StatePending mountPointState = iota
	StateSuccess
	StateFailed
)

func newVolumeInfos() *volumeInfos {
	result := &volumeInfos{}

	result.volumes = make(map[string]map[string]*mountPointInfo)

	return result
}

func (i *volumeInfos) AddVolume(volume string, listMountPoints func(volume string) ([]string, error)) error {
	i.mutex.RLock()

	_, ok := i.volumes[volume]

	i.mutex.RUnlock()

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
			path:  p,
			state: StatePending,
		}
	}

	i.mutex.Lock()
	defer i.mutex.Unlock()

	_, ok = i.volumes[volume]
	if !ok {
		i.volumes[volume] = ms
	}

	return nil
}

func (i *volumeInfos) ComputeNeeded(dir string) []*mountPointInfo {
	var result []*mountPointInfo

	volume := strings.ToLower(filepath.VolumeName(dir))

	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, mount := range i.volumes[volume] {
		insideMount := strings.HasPrefix(dir, mount.path)
		mountInside := strings.HasPrefix(mount.path, dir)

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

	volume := strings.ToLower(filepath.VolumeName(dir))

	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for _, m := range i.volumes[volume] {
		if strings.HasPrefix(dir, m.path) && (result == nil || len(result.path) < len(m.path)) {
			result = m
		}
	}

	return result
}
