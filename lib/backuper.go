package fs_snapshot

// Backuper is a class that allows easy temporary snapshot creation.
// All snapshots created will be deleted when calling Close.
type Backuper interface {
	// TryToCreateTemporarySnapshot returns the snapshoted diretory volume if a snapshot could be made,
	// or the original directory otherwise.
	// If the directory is not yet available as a snapshot, a snapshot is created.
	// If the directory is inside an existing snapshot, the snapshot is re-used.
	// An error is returned iff some problem occurred while creating the snapshot, but not
	// if the directory does not support snapshots.
	TryToCreateTemporarySnapshot(directory string) (string, error)

	// Close frees all resources.
	Close()
}
