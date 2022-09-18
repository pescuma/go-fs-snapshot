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

	for _, m := range mountPoints {
		output, err := runAndReturnOutput(s.infoCallback, "tmutil", "listlocalsnapshots", m)
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
				OriginalPath: m,
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

func (s *macosSnapshoter) StartBackup(cfg *BackupConfig) (Backuper, error) {
	if cfg == nil {
		cfg = &BackupConfig{}
	}

	if cfg.ProviderID != "" && cfg.ProviderID != providerID {
		return nil, errors.Errorf("unknown provider id: %v", cfg.ProviderID)
	}

	ic := cfg.InfoCallback
	if ic == nil {
		ic = s.infoCallback
	}

	b := &macosBackuper{
		volumes:      newVolumeInfos(),
		infoCallback: ic,
	}

	err := b.volumes.AddVolume("", s.listMountPoints)
	if err != nil {
		return nil, err
	}

	return b, nil
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

func (s *macosSnapshoter) listMountPoints() ([]string, error) {
	output, err := runAndReturnOutput(s.infoCallback, "diskutil", "apfs", "list")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("Snapshot Mount Point: +(/[^\\r\\n]*)")
	matches := re.FindAllStringSubmatch(output, -1)

	result := make([]string, len(matches))

	for i, m := range matches {
		result[i] = m[1]
	}

	return result, nil
}
