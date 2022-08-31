package cli

import (
	"fmt"
	"os/user"

	"github.com/pescuma/go-fs-snapshot/lib"
)

type enableForCurrentUserCmd struct {
}

func (c *enableForCurrentUserCmd) Run(ctx *context) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	return fs_snapshot.EnableSnapshotsForUser(u.Username, outputMessages(ctx))
}

type enableForUserCmd struct {
	Username string `arg:"" help:"Username to enable snapshots."`
}

func (c *enableForUserCmd) Run(ctx *context) error {
	return fs_snapshot.EnableSnapshotsForUser(c.Username, outputMessages(ctx))
}

type enableTestCurrentUserCmd struct {
}

func (c *enableTestCurrentUserCmd) Run(ctx *context) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	can, err := fs_snapshot.CurrentUserCanCreateSnapshots(outputMessages(ctx))
	if err != nil {
		return err
	}

	if can {
		fmt.Printf("Current user (%v) can create snapshots.\n", u.Username)
	} else {
		fmt.Printf("Current user (%v) can NOT create snapshots.\n", u.Username)
	}

	return nil
}
