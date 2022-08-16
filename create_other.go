//go:build !windows

package fs_snapshot

// CreateSnapshoter creates a new snapshoter.
// In case of error a null snapshoter is returned, so you can use it without problem
func CreateSnapshoter() (Snapshoter, error) {
	return createNullSnapshoter(), errors.New("snapshots not supported in this OS")
}
