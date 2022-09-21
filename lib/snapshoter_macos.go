//go:build darwin

package fs_snapshot

import (
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	prefix     = "com.apple.TimeMachine."
	suffix     = ".local"
	providerID = "tmutil-local"
)

func startServerForOS(infoCb InfoMessageCallback) error {
	return errors.New("can't start server with elevated privileges - run with sudo if needed")
}

func newSnapshoterForOS(cfg *SnapshoterConfig) (Snapshoter, error) {
	output, err := runAndReturnOutput(cfg.InfoCallback, "tmutil", "version")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`tmutil version ([0-9a-zA-Z_.]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return nil, errors.Errorf("unknown tmutil version: %v", output)
	}

	return &macosSnapshoter{
		infoCallback: cfg.InfoCallback,
		version:      matches[1],
	}, nil
}

type macosSnapshoter struct {
	infoCallback InfoMessageCallback
	version      string
}

func (s *macosSnapshoter) SimplifyID(id string) string {
	id = strings.TrimPrefix(id, prefix)
	id = strings.TrimSuffix(id, suffix)
	return id
}

func (s *macosSnapshoter) ListProviders(filterID string) ([]*Provider, error) {
	provider := s.newProvider()

	if filterID != "" && filterID != provider.ID {
		return []*Provider{}, nil
	}

	return []*Provider{provider}, nil
}

func (s *macosSnapshoter) ListSets(filterID string) ([]*SnapshotSet, error) {
	return []*SnapshotSet{}, nil
}

func (s *macosSnapshoter) ListSnapshots(filterID string) ([]*Snapshot, error) {
	mountPoints, err := s.listMountPoints()
	if err != nil {
		return nil, err
	}

	var result []*Snapshot

	provider := s.newProvider()

	for k, v := range mountPoints {
		output, err := runAndReturnOutput(s.infoCallback, "tmutil", "listlocalsnapshots", v)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(output, "\n")
		lines = lines[1:]

		for _, line := range lines {
			simple := s.SimplifyID(line)

			if filterID != "" && line != filterID && simple != filterID {
				continue
			}

			t, err := time.Parse("2006-01-02-150405", simple)
			if err != nil {
				return nil, err
			}

			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)

			result = append(result, &Snapshot{
				ID:           line,
				OriginalPath: k,
				SnapshotPath: "",
				CreationTime: t,
				Set:          nil,
				Provider:     provider,
				State:        "created",
				Attributes:   "",
			})
		}
	}

	return result, nil
}

func (s *macosSnapshoter) DeleteSet(id string, force bool) (bool, error) {
	return false, errors.New("snapshot sets not supported in MacOS")
}

func (s *macosSnapshoter) DeleteSnapshot(id string, force bool) (bool, error) {
	err := run(s.infoCallback, "tmutil", "deletelocalsnapshots", id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *macosSnapshoter) ListMountPoints(volume string) ([]string, error) {
	if volume != "" {
		return nil, errors.Errorf("unknown volume: %v", volume)
	}

	// TODO Fetch this data from OS
	return []string{
		"/",
	}, nil
}

func (s *macosSnapshoter) listMountPoints() (map[string]string, error) {
	// TODO Fetch this data from OS
	return map[string]string{
		"/": "/System/Volumes/Data",
	}, nil
}

func (s *macosSnapshoter) StartBackup(cfg *BackupConfig) (Backuper, error) {
	if cfg == nil {
		cfg = &BackupConfig{}
	}

	if cfg.ProviderID != "" && cfg.ProviderID != providerID {
		return nil, errors.Errorf("unknown provider id: %v", cfg.ProviderID)
	}

	mountPoints, err := s.listMountPoints()
	if err != nil {
		return nil, err
	}

	ic := cfg.InfoCallback
	if ic == nil {
		ic = s.infoCallback
	}

	return newMacosBackuper(mountPoints, s.ListMountPoints, ic)
}

func (s *macosSnapshoter) Close() {
}

func (s *macosSnapshoter) newProvider() *Provider {
	return &Provider{
		ID:      providerID,
		Name:    "Time Machine local snapshots",
		Version: s.version,
		Type:    "console application",
	}
}
