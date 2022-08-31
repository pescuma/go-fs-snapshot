//go:build !windows

package fs_snapshot

// CurrentUserCanCreateSnapshots returns information if the current user can create snapshots
func CurrentUserCanCreateSnapshots(infoCb InfoMessageCallback) (can bool, username string, err error) {
	return false, "", errors.New("snapshots not supported in this OS")
}

// EnableSnapshotsForUser enables the current user to run snaphsots.
// This generally must be run from a prompt with elevated privileges (root or administrator).
func EnableSnapshotsForUser(username string, infoCb InfoMessageCallback) error {
	return errors.New("snapshots not supported in this OS")
}
