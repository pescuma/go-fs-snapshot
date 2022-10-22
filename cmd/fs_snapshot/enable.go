package main

import (
	"os/user"

	"github.com/pescuma/go-fs-snapshot/fs_snapshot"
)

type enableForCurrentUserCmd struct {
}

func (c *enableForCurrentUserCmd) Run(ctx *context) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	return fs_snapshot.EnableSnapshotsForUser(u.Username, ctx.console.NewInfoMessageCallback())
}

type enableForUserCmd struct {
	Username string `arg:"" help:"Username to enable snapshots."`
}

func (c *enableForUserCmd) Run(ctx *context) error {
	return fs_snapshot.EnableSnapshotsForUser(c.Username, ctx.console.NewInfoMessageCallback())
}

type enableTestCurrentUserCmd struct {
}

func (c *enableTestCurrentUserCmd) Run(ctx *context) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	can, err := fs_snapshot.CurrentUserCanCreateSnapshots(ctx.console.NewInfoMessageCallback())
	if err != nil {
		return err
	}

	if can {
		ctx.console.Printf("Current user (%v) can create snapshots.", u.Username)
	} else {
		ctx.console.Printf("Current user (%v) can NOT create snapshots.", u.Username)
	}

	return nil
}
