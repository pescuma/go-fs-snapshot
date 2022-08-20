package cli

import (
	"fmt"

	"github.com/pkg/errors"
)

type deleteCmd struct {
	ID    string `arg:"" help:"The ID (simplified or full) of the snapshot to delete."`
	Force bool   `short:"f" help:"Do everything possible to try to delete."`
	Yes   bool   `short:"y" help:"Do not prompt for deletion confirmation."`
}

func (c *deleteCmd) Run(ctx *context) error {
	snapshot, err := findOneSnapshot(ctx, c.ID)
	if err != nil {
		return err
	}
	if snapshot == nil {
		return nil
	}

	fmt.Printf("Snapshot info:\n")
	printSnapshotInfo(snapshot, "   ")
	fmt.Printf("\n")

	if !c.Yes {
		confirm := askForConfirmation("Are you sure you want to delete this snapshot?")
		if !confirm {
			return nil
		}
		fmt.Printf("\n")
	}

	deleted, err := ctx.snapshoter.DeleteSnapshot(snapshot.ID, c.Force)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.Errorf("Snapshot not found.")
	}

	fmt.Printf("Snapshot %v deleted.\n", snapshot.ID)
	return nil
}
