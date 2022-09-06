package fs_snapshot

import (
	"time"
)

type Snapshoter interface {
	// ListProviders list all provider available
	// filterID: filter by id if != ""
	ListProviders(filterID string) ([]*Provider, error)

	// ListSets list all snapshot sets available
	// filterID: filter by id if != ""
	ListSets(filterID string) ([]*SnapshotSet, error)

	// ListSnapshots list all snapshots available
	// filterID: filter by id if != ""
	ListSnapshots(filterID string) ([]*Snapshot, error)

	// SimplifyID simplifies the snapshot, set and provider IDs, if possible
	SimplifyID(id string) string

	// DeleteSet deletes one snapshot set and all its snapshots.
	// Returns true if snapshot was found and deleted, false if it was not found and an
	// error if something went wrong.
	DeleteSet(id string, force bool) (bool, error)

	// DeleteSnapshot deletes one snapshot
	// Returns true if snapshot was found and deleted, false if it was not found and an
	// error if something went wrong.
	DeleteSnapshot(id string, force bool) (bool, error)

	// StartBackup creates a Backuper to allow easy backup creation.
	StartBackup(opts *SnapshotOptions) (Backuper, error)

	// Close frees all resources.
	Close()
}

type Provider struct {
	ID      string
	Name    string
	Version string
	Type    string
}

type SnapshotSet struct {
	ID                      string
	CreationTime            time.Time
	SnapshotCountOnCreation int
	Snapshots               []*Snapshot
}

type Snapshot struct {
	ID           string
	OriginalPath string
	SnapshotPath string
	CreationTime time.Time
	Set          *SnapshotSet
	Provider     *Provider
	State        string
	Attributes   string
}

type ConnectionType int

const (
	LocalOrServer ConnectionType = iota
	LocalOnly
	ServerOnly
)

type InfoMessageCallback func(level MessageLevel, format string, a ...interface{})

type MessageLevel int

const (
	OutputLevel = iota
	InfoLevel
	DetailsLevel
	TraceLevel
)

type SnapshotOptions struct {
	ProviderID string

	Timeout time.Duration

	// Simple - try to do it as simple as possible, but not simpler.
	// In Windows this means do not use VSS Writers.
	Simple bool
}

// NewSnapshoter creates a new snapshoter.
// In case of error a null snapshoter is returned, so you can use it without problem.
func NewSnapshoter(cfg *SnapshoterConfig) (Snapshoter, error) {
	if cfg == nil {
		cfg = &SnapshoterConfig{}
	}
	cfg.setDefaults()

	var result Snapshoter
	var errLocal error
	var errServer error

	if cfg.ConnectionType != ServerOnly {
		result, errLocal = newSnapshoterForOS(cfg)
		if errLocal == nil {
			return result, nil
		}
	}

	if cfg.ConnectionType != LocalOnly {
		result, errServer = newClientSnapshoterStartingServer(cfg)
		if errServer == nil {
			return result, nil
		}
	}

	if cfg.ConnectionType != ServerOnly {
		return newNullSnapshoter(), errLocal
	} else {
		return newNullSnapshoter(), errServer
	}
}

func newClientSnapshoterStartingServer(cfg *SnapshoterConfig) (Snapshoter, error) {
	result, err := newClientSnapshoter(cfg)
	if err == nil {
		return result, nil
	}

	if cfg.ServerIP != DefaultIP && cfg.ServerPort != DefaultPort {
		// Not the default config, so don't try to start the server
		return nil, err
	}

	err1 := startServerForOS(cfg.InfoCallback)
	if err1 != nil {
		// Ignore error starting server and just return the connection error
		return nil, err
	}

	return newClientSnapshoter(cfg)
}

type SnapshoterConfig struct {
	ConnectionType ConnectionType
	ServerIP       string
	ServerPort     int
	InfoCallback   InfoMessageCallback
}

func (cfg *SnapshoterConfig) setDefaults() {
	if cfg.ServerIP == "" {
		cfg.ServerIP = DefaultIP
	}
	if cfg.ServerPort == 0 {
		cfg.ServerPort = DefaultPort
	}
	if cfg.InfoCallback == nil {
		cfg.InfoCallback = func(level MessageLevel, format string, a ...interface{}) {}
	}
}
