package main

import (
	"github.com/pkg/errors"
)

type deleteCmd struct {
	ID    string `arg:"" help:"The ID (simplified or full) of the snapshot to delete."`
	Force bool   `short:"f" help:"Do everything possible to try to delete."`
	Yes   bool   `short:"y" help:"Do not prompt for deletion confirmation."`

	ServerArgs serverArgs `embed:""`
}

func (c *deleteCmd) Run(ctx *context) error {
	snapshot, err := findOneSnapshot(ctx, c.ID)
	if err != nil {
		return err
	}
	if snapshot == nil {
		return nil
	}

	ctx.console.Print("Snapshot info:")
	printSnapshotInfo(ctx, snapshot, "   ")
	ctx.console.Print("")

	if !c.Yes {
		confirm := ctx.console.AskForConfirmation("Are you sure you want to delete this snapshot?")
		if !confirm {
			return nil
		}
		ctx.console.Print("")
	}

	deleted, err := ctx.snapshoter.DeleteSnapshot(snapshot.ID, c.Force)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.Errorf("Snapshot not found.")
	}

	ctx.console.Printf("Snapshot %v deleted.", snapshot.ID)
	return nil
}
