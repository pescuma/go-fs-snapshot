package fs_snapshot

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pescuma/go-fs-snapshot/lib/internal/windows"
)

type windowsBackuper struct {
	opts         *internal_fs_snapshot_windows.SnapshotOptions
	mutex        sync.RWMutex
	volumes      map[string]*volumeInfo // The keys are all lower case
	vssResults   []*internal_fs_snapshot_windows.SnapshotsResult
	infoCallback InfoMessageCallback
}

type volumeInfo struct {
	volume       string
	state        volumeState
	snapshotPath string
}

type volumeState int

const (
	StatePending volumeState = iota
	StateSuccess
	StateFailed
)

func (b *windowsBackuper) TryToCreateTemporarySnapshot(inputDirectory string) (string, error) {
	dir, err := absolutePath(inputDirectory)
	if err != nil {
		return inputDirectory, err
	}

	dir = strings.ToLower(dir) + `\`
	dirInfo := b.getPathInfo(dir)

	switch {
	case dirInfo != nil && dirInfo.state == StateFailed:
		return inputDirectory, nil

	case dirInfo != nil && dirInfo.state == StateSuccess:
		newDir, err := changeBaseDir(dir, dirInfo.volume, dirInfo.snapshotPath)
		if err != nil {
			return inputDirectory, err
		}

		return newDir, nil
	}

	volume := filepath.VolumeName(dir) + `\`

	mounts, err := internal_fs_snapshot_windows.EnumerateMountedFolders(volume)
	if err != nil {
		return inputDirectory, err
	}

	for i := range mounts {
		mounts[i] = strings.ToLower(mounts[i])
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	needed := b.computeNeededInsideLock(dir, volume, mounts)

	if len(needed) == 0 {
		panic("Should not happen")
	}

	if b.infoCallback != nil {
		b.infoCallback(InfoLevel, "Creating VSS snapshot of "+strings.Join(needed, " ; "))
	}

	vsr, err := internal_fs_snapshot_windows.CreateSnapshots(needed, b.opts)
	b.vssResults = append(b.vssResults, vsr)
	if err != nil {
		b.markFailed(needed...)
		b.markPending(mounts)
		return inputDirectory, err
	}

	for i := range needed {
		orig := needed[i]
		snap := vsr.GetSnapshotPath(orig)

		if snap == orig {
			b.markFailed(orig)
		} else {
			b.markSuccess(orig, snap)
		}
	}
	b.markPending(mounts)

	dirInfo = b.getPathInfoInsideLock(dir)
	if dirInfo.state != StateSuccess {
		return inputDirectory, nil
	}

	newDir, err := changeBaseDir(dir, dirInfo.volume, dirInfo.snapshotPath)
	if err != nil {
		return inputDirectory, err
	}

	return newDir, nil
}

func (b *windowsBackuper) computeNeededInsideLock(dir string, volume string, mounts []string) []string {
	var needed []string

	var needVolume = true
	for _, mount := range mounts {
		insideMount := strings.HasPrefix(dir, mount)
		mountInside := strings.HasPrefix(mount, dir)

		if !insideMount && !mountInside {
			continue
		}

		if insideMount {
			needVolume = false
		}

		if b.getVolumeStateInsideLock(mount) == StatePending {
			needed = append(needed, mount)

			if b.infoCallback != nil {
				b.infoCallback(DetailsLevel, fmt.Sprintf("Detected mount point %v inside snapshot dir %v", mount, dir))
			}
		}
	}

	if needVolume && b.getVolumeStateInsideLock(volume) == StatePending {
		needed = append([]string{volume}, needed...)
	}
	return needed
}

func (b *windowsBackuper) getVolumeStateInsideLock(volume string) volumeState {
	mountInfo, ok := b.volumes[volume]
	if !ok {
		return StatePending
	}

	return mountInfo.state
}

func (b *windowsBackuper) getPathInfo(dir string) *volumeInfo {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.getPathInfoInsideLock(dir)
}

func (b *windowsBackuper) getPathInfoInsideLock(dir string) *volumeInfo {
	var dirInfo *volumeInfo
	for _, p := range b.volumes {
		if strings.HasPrefix(dir, p.volume) && (dirInfo == nil || len(dirInfo.volume) < len(p.volume)) {
			dirInfo = p
		}
	}

	return dirInfo
}

func (b *windowsBackuper) markPending(volumes []string) {
	for _, v := range volumes {
		_, ok := b.volumes[v]
		if !ok {
			b.volumes[v] = &volumeInfo{
				volume: v,
				state:  StatePending,
			}
		}
	}
}

func (b *windowsBackuper) markFailed(volumes ...string) {
	for _, v := range volumes {
		b.volumes[v] = &volumeInfo{
			volume: v,
			state:  StateFailed,
		}
	}
}

func (b *windowsBackuper) markSuccess(volume, snapshotPath string) {
	b.volumes[volume] = &volumeInfo{
		volume:       volume,
		state:        StateSuccess,
		snapshotPath: snapshotPath,
	}
}

func (b *windowsBackuper) ListSnapshotedDirectories() map[string]string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	result := make(map[string]string)

	for _, v := range b.volumes {
		if v.state != StateSuccess {
			continue
		}

		result[v.volume] = v.snapshotPath
	}

	return result
}

func (b *windowsBackuper) Close() {
	for _, r := range b.vssResults {
		r.Close()
	}
}
