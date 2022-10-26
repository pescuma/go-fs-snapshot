//go:build darwin

package fs_snapshot

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

func currentUserCanCreateSnapshotsForOS(infoCb InfoMessageCallback) (bool, error) {
	err := run(infoCb, "tmutil", "version")
	if err != nil {
		infoCb(InfoLevel, "tmutil not fount.")
		return false, nil
	}

	infoCb(InfoLevel, "tmutil fount.")

	_, err = exec.LookPath("mount_apfs")
	if err != nil {
		infoCb(InfoLevel, "mount_apfs not fount.")
		return false, nil
	}

	infoCb(InfoLevel, "mount_apfs fount.")

	has, err := currentUserHasFullDiskAccessPermission(infoCb)
	if err != nil {
		return false, err
	}

	if has {
		infoCb(InfoLevel, "Current user has Full Disk Access permission.")
	} else {
		infoCb(InfoLevel, "Current user does NOT have Full Disk Access permission.")
	}

	return has, nil
}

func currentUserHasFullDiskAccessPermission(infoCb InfoMessageCallback) (bool, error) {
	requires, err := requiresFullDiskAccessPermission(infoCb)
	if err != nil {
		return false, err
	}

	if !requires {
		return true, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	// Copied from https://github.com/MacPaw/PermissionsKit/blob/master/PermissionsKit/Private/FullDiskAccess/MPFullDiskAccessAuthorizer.m
	tests := []string{
		filepath.Join(home, "Library/Safari/CloudTabs.db"),
		filepath.Join(home, "Library/Safari/Bookmarks.plist"),
		"/Library/Application Support/com.apple.TCC/TCC.db",
		"/Library/Preferences/com.apple.TimeMachine.plist",
	}

	foundOne := false

	for _, t := range tests {
		_, err = os.Stat(t)
		if err != nil {
			continue
		}

		file, err := os.Open(t)
		if err != nil {
			continue
		}

		_ = file.Close()

		foundOne = true
		break
	}

	return foundOne, nil
}

func requiresFullDiskAccessPermission(infoCb InfoMessageCallback) (bool, error) {
	output, err := runAndReturnOutput(infoCb, "sw_vers", "-productVersion")
	if err != nil {
		return false, err
	}

	ver, err := semver.NewVersion(output)
	if err != nil {
		return false, err
	}

	return !ver.LessThan(semver.MustParse("10.14")), nil
}

// EnableSnapshotsForUser enables the current user to run snapshots.
// This generally must be run from a prompt with elevated privileges (root or administrator).
func enableSnapshotsForUserForOS(username string, infoCb InfoMessageCallback) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	if username == u.Username {
		can, err := currentUserCanCreateSnapshotsForOS(infoCb)
		if err != nil {
			return err
		}

		if can {
			infoCb(OutputLevel, "User %v is already able to create snapshots.", username)
			return nil
		}
	}

	requires, err := requiresFullDiskAccessPermission(infoCb)
	if err != nil {
		return err
	}

	if !requires {
		return errors.New("MacOS does not allow to enable snapshots for a user, but running the command as " +
			"root should be enough.")

	} else {
		return errors.New("MacOS does not allow to grant Full Disk Access permission from an application. " +
			"You need to open 'System Preferences...', go to the 'Privacy' tab, select 'Full Disk Access' in the list " +
			"on the left, click on the lock on the bottom, input your password and then add the correct application " +
			"to the list on the right. If you intend to use this app inside terminal, you must select 'Terminal.app' " +
			"in the list on the right (for some reason granting the permission to fs_snapshot does not work). In some " +
			"other cases you may need to add and grant the permission to 'fs_snapshot'.")

	}
}
