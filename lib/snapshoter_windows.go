//go:build windows

package fs_snapshot

import (
	"os"
	"os/user"
	"runtime"
	"sort"
	"strings"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"

	internal_windows "github.com/pescuma/go-fs-snapshot/lib/internal/windows"
)

const simpleIdLength = 7

func startServerForOS(infoCb InfoMessageCallback) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	err = run(infoCb, "schtasks", "/Run",
		"/TN", createScheduledTaskName(u.Username),
		"/HRESULT")
	if err != nil {
		infoCb(TraceLevel, "error running scheduled task: %v", err.Error())
		return err
	}

	return nil
}

func newSnapshoterForOS(cfg *SnapshoterConfig) (Snapshoter, error) {
	err := initializePrivileges()
	if err != nil {
		return nil, err
	}

	is64Bit, err := internal_windows.IsRunningOn64BitWindows()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to detect windows architecture: %s", err.Error())
	}

	if (is64Bit && runtime.GOARCH != "amd64") || (!is64Bit && runtime.GOARCH != "386") {
		return nil, errors.Errorf("executables compiled for %v can't use VSS on other architectures. "+
			"Please use an executable compiled for your platform.", runtime.GOARCH)
	}

	err = internal_windows.InitializeCOM()
	if err != nil {
		return nil, err
	}

	bc, err := internal_windows.NewIVSSBackupComponents()
	bc.Close()
	if err != nil {
		return nil, err
	}

	infoCb := cfg.InfoCallback
	if infoCb == nil {
		infoCb = func(level MessageLevel, format string, a ...interface{}) {}
	}

	return &windowsSnapshoter{
		infoCallback: infoCb,
	}, nil
}

type windowsSnapshoter struct {
	infoCallback InfoMessageCallback
}

func (s *windowsSnapshoter) SimplifyID(id string) string {
	id = strings.ReplaceAll(id, "-", "")
	return id[:simpleIdLength]
}

func (s *windowsSnapshoter) ListProviders(filterID string) ([]*Provider, error) {
	var result []*Provider

	bc, err := s.NewBackupComponentsForManagement()
	defer bc.Close()
	if err != nil {
		return nil, err
	}

	enum, err := bc.Query(internal_windows.VSS_OBJECT_PROVIDER)
	defer enum.Close()
	if err != nil {
		return nil, err
	}

	for {
		var props struct {
			objectType uint32
			provider   internal_windows.VssProviderProperties
		}

		count, err := enum.Next(1, unsafe.Pointer(&props))
		if err != nil {
			return nil, err
		}

		if count < 1 {
			break
		}

		id := toGuidString(props.provider.ProviderID)

		add := true
		if filterID != "" && filterID != id && filterID != s.SimplifyID(id) {
			add = false

		}

		if add {
			result = append(result, &Provider{
				ID:      id,
				Name:    ole.UTF16PtrToString(props.provider.ProviderName),
				Version: ole.UTF16PtrToString(props.provider.ProviderVersion),
				Type:    props.provider.ProviderType.Str(),
			})
		}

		props.provider.Close()
	}

	return result, nil
}

func (s *windowsSnapshoter) ListSets(filterID string) ([]*SnapshotSet, error) {
	_, result, err := s.listSnapshotsAndSets("", filterID)
	return result, err
}

func (s *windowsSnapshoter) ListSnapshots(filterID string) ([]*Snapshot, error) {
	result, _, err := s.listSnapshotsAndSets(filterID, "")
	return result, err
}

func (s *windowsSnapshoter) listSnapshotsAndSets(filterSnapshotID string, filterSetID string) ([]*Snapshot, []*SnapshotSet, error) {
	bc, err := s.NewBackupComponentsForManagement()
	defer bc.Close()
	if err != nil {
		return nil, nil, err
	}

	enum, err := bc.Query(internal_windows.VSS_OBJECT_SNAPSHOT)
	defer enum.Close()
	if err != nil {
		return nil, nil, err
	}

	providers, err := s.ListProviders("")
	if err != nil {
		return nil, nil, err
	}

	providersById := make(map[string]*Provider)
	for _, p := range providers {
		providersById[p.ID] = p
	}

	setsById := make(map[string]*SnapshotSet)

	var sets []*SnapshotSet
	var snapshots []*Snapshot

	for {
		var props struct {
			objectType uint32
			snapshot   internal_windows.VssSnapshotProperties
		}

		count, err := enum.Next(1, unsafe.Pointer(&props))
		if err != nil {
			return nil, nil, err
		}

		if count < 1 {
			break
		}

		volumes, err := getVolumeNames(props.snapshot.OriginalVolumeName)
		if err != nil {
			return nil, nil, err
		}

		provider := providersById[toGuidString(props.snapshot.ProviderID)]

		setID := toGuidString(props.snapshot.SnapshotSetID)
		snapshotID := toGuidString(props.snapshot.SnapshotID)

		add := true
		if filterSetID != "" && filterSetID != setID && filterSetID != s.SimplifyID(setID) {
			add = false
		} else if filterSnapshotID != "" && filterSnapshotID != snapshotID && filterSnapshotID != s.SimplifyID(snapshotID) {
			add = false
		}

		if add {
			set, exists := setsById[setID]
			if !exists {
				set = &SnapshotSet{
					ID:                      setID,
					CreationTime:            toDate(props.snapshot.CreationTimestamp),
					SnapshotCountOnCreation: int(props.snapshot.SnapshotsCount),
				}
				setsById[setID] = set
			}

			snapshot := &Snapshot{
				ID:           snapshotID,
				OriginalPath: strings.Join(volumes, ", "),
				SnapshotPath: ole.UTF16PtrToString(props.snapshot.SnapshotDeviceObject),
				CreationTime: toDate(props.snapshot.CreationTimestamp),
				Provider:     provider,
				Set:          set,
				State:        props.snapshot.Status.Str(),
				Attributes:   props.snapshot.SnapshotAttributes.Str(),
			}

			set.Snapshots = append(set.Snapshots, snapshot)
			if set.CreationTime.After(snapshot.CreationTime) {
				set.CreationTime = snapshot.CreationTime
			}

			snapshots = append(snapshots, snapshot)
			sets = append(sets, set)
		}

		err = internal_windows.VssFreeSnapshotProperties(&props.snapshot)
		if err != nil {
			return nil, nil, err
		}
	}

	return snapshots, sets, nil
}

func getVolumeNames(volume *uint16) ([]string, error) {
	buffer := make([]uint16, windows.MAX_LONG_PATH)
	var required uint32

	err := windows.GetVolumePathNamesForVolumeName(volume, &buffer[0], uint32(len(buffer)),
		&required)

	if err != nil && required > windows.MAX_LONG_PATH {
		buffer = make([]uint16, required)
		err = windows.GetVolumePathNamesForVolumeName(volume, &buffer[0], uint32(len(buffer)),
			&required)
	}

	var result []string
	i := 0

	for i < len(buffer) {
		name := ole.UTF16PtrToString(&buffer[i])
		i += len(name)*2 + 2

		if name == "" {
			break
		}

		if name[len(name)-1] == os.PathSeparator {
			name = name[:len(name)-1]
		}

		result = append(result, name)
	}

	sort.Slice(result, func(a, b int) bool {
		return result[a] < result[b]
	})

	return result, nil
}

func (s *windowsSnapshoter) DeleteSet(id string, force bool) (bool, error) {
	guid, err := s.getSetID(id)
	if err != nil {
		return false, err
	}

	bc, err := s.NewBackupComponentsForManagement()
	defer bc.Close()
	if err != nil {
		return false, err
	}

	deleted, _, err := bc.DeleteSnapshots(internal_windows.VSS_OBJECT_SNAPSHOT_SET, guid, force)

	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, nil
	}

	return true, nil
}

func (s *windowsSnapshoter) DeleteSnapshot(id string, force bool) (bool, error) {
	guid, err := s.getSnapshotID(id)
	if err != nil {
		return false, err
	}

	bc, err := s.NewBackupComponentsForManagement()
	defer bc.Close()
	if err != nil {
		return false, err
	}

	deleted, _, err := bc.DeleteSnapshots(internal_windows.VSS_OBJECT_SNAPSHOT, guid, force)

	if err != nil {
		return false, err
	}
	if deleted == 0 {
		return false, nil
	}

	return true, nil
}

func (s *windowsSnapshoter) NewBackupComponentsForManagement() (*internal_windows.IVSSBackupComponents, error) {
	bc, err := internal_windows.NewIVSSBackupComponents()
	if err != nil {
		return bc, err
	}

	err = bc.InitializeForBackup()
	if err != nil {
		return bc, err
	}

	err = bc.SetContext(internal_windows.VSS_CTX_ALL)
	if err != nil {
		return bc, err
	}

	return bc, nil
}

func (s *windowsSnapshoter) StartBackup(cfg *BackupConfig) (Backuper, error) {
	if cfg == nil {
		cfg = &BackupConfig{}
	}

	providerID, err := s.getProviderID(cfg.ProviderID)
	if err != nil {
		return nil, err
	}

	ic := cfg.InfoCallback
	if ic == nil {
		ic = s.infoCallback
	}

	return &windowsBackuper{
		cfg: &internal_windows.SnapshotOptions{
			ProviderID:   providerID,
			Timeout:      cfg.Timeout,
			Writters:     !cfg.Simple,
			InfoCallback: ic,
		},
		infoCallback: s.infoCallback,
		volumes:      make(map[string]*volumeInfo),
	}, nil
}

func (s *windowsSnapshoter) getProviderID(id string) (*ole.GUID, error) {
	if id == "" {
		return nil, nil
	}

	result := ole.NewGUID(id)

	if result == nil {
		// It's a simplified one, so we must get the full ID

		providers, err := s.ListProviders(id)
		if err != nil {
			return nil, err
		}

		switch len(providers) {
		case 0:
			return nil, errors.Errorf("Unknown provider ID: %v", id)
		case 1:
			// continue
		default:
			return nil, errors.Errorf("found %v providers with ID %v - please use full ID", len(providers), id)
		}

		result = ole.NewGUID(providers[0].ID)
	}

	return result, nil
}

func (s *windowsSnapshoter) getSnapshotID(id string) (*ole.GUID, error) {
	if id == "" {
		return nil, nil
	}

	result := ole.NewGUID(id)

	if result == nil {
		// It's a simplified one, so we must get the full ID

		snapshots, err := s.ListSnapshots(id)
		if err != nil {
			return nil, err
		}

		switch len(snapshots) {
		case 0:
			return nil, errors.Errorf("Unknown provider ID: %v", id)
		case 1:
			// continue
		default:
			return nil, errors.Errorf("found %v snapshots with ID %v - please use full ID", len(snapshots), id)
		}

		result = ole.NewGUID(snapshots[0].ID)
	}

	return result, nil
}

func (s *windowsSnapshoter) getSetID(id string) (*ole.GUID, error) {
	if id == "" {
		return nil, nil
	}

	result := ole.NewGUID(id)

	if result == nil {
		// It's a simplified one, so we must get the full ID

		sets, err := s.ListSets(id)
		if err != nil {
			return nil, err
		}

		switch len(sets) {
		case 0:
			return nil, errors.Errorf("Unknown provider ID: %v", id)
		case 1:
			// continue
		default:
			return nil, errors.Errorf("found %v snapshot sets with ID %v - please use full ID", len(sets), id)
		}

		result = ole.NewGUID(sets[0].ID)
	}

	return result, nil
}

func (s *windowsSnapshoter) Close() {
}
