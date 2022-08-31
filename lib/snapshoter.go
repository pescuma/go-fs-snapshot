package fs_snapshot

import "time"

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

	// SimplifyId simplifies the snapshot, set and provider IDs, if possible
	SimplifyId(id string) string

	// DeleteSet deletes one snapshot set and all its snapshots.
	// Returns true if snapshot was found and deleted, false if it was not found and an
	// error if something went wrong.
	DeleteSet(id string, force bool) (bool, error)

	// DeleteSnapshot deletes one snapshot
	// Returns true if snapshot was found and deleted, false if it was not found and an
	// error if something went wrong.
	DeleteSnapshot(id string, force bool) (bool, error)

	// StartBackup creates a Backuper to allow easy backup creation.
	StartBackup(opts *CreateSnapshotOptions) (Backuper, error)

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

type CreateSnapshotOptions struct {
	ProviderID string

	Timeout time.Duration

	// Simple - try to do it as simple as possible, but not simpler.
	// In Windows this means do not use VSS Writers.
	Simple bool

	InfoCallback InfoMessageCallback
}

type InfoMessageCallback func(level MessageLevel, msg string)

type MessageLevel int

const (
	OutputLevel = iota
	InfoLevel
	DetailsLevel
	TraceLevel
)
