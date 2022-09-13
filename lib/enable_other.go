//go:build !windows

package fs_snapshot

import "errors"

// CurrentUserCanCreateSnapshots returns information if the current user can create snapshots
func CurrentUserCanCreateSnapshots(infoCb InfoMessageCallback) (bool, error) {
	return false, errors.New("snapshots not supported in this OS")
}

// EnableSnapshotsForUser enables the current user to run snapshots.
// This generally must be run from a prompt with elevated privileges (root or administrator).
func EnableSnapshotsForUser(username string, infoCb InfoMessageCallback) error {
	return errors.New("snapshots not supported in this OS")
}
