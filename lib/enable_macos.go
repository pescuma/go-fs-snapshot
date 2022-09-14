//go:build darwin

package fs_snapshot

import "github.com/pkg/errors"

func currentUserCanCreateSnapshotsForOS(infoCb InfoMessageCallback) (bool, error) {
	err := run(infoCb, "tmutil", "version")
	if err != nil {
		return false, nil
	}

	return true, nil
}

// EnableSnapshotsForUser enables the current user to run snapshots.
// This generally must be run from a prompt with elevated privileges (root or administrator).
func EnableSnapshotsForUser(username string, infoCb InfoMessageCallback) error {
	return errors.New("can't enable for MacOS - run with sudo if needed")
}
