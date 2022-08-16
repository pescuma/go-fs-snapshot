package fs_snapshot

// ErrorHandler is used to report errors via callback
type ErrorHandler func(item string, err error) error

type Snapshoter interface {
	// CreateSnapshot returns the path inside a snapshot.
	// If the path is not yet available as a snapshot, a snapshot is created.
	// If creation of a snapshot fails the file's original path is returned as
	// a fallback.
	CreateSnapshot(path string) string

	// Destroy deletes all snapshots that had to be created and free other resources.
	Destroy() error
}
