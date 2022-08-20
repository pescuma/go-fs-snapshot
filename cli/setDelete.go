package cli

import (
	"fmt"

	"github.com/pkg/errors"
)

type setDeleteCmd struct {
	ID    string `arg:"" help:"The ID (simplified or full) of the snapshot set to delete."`
	Force bool   `short:"f" help:"Do everything possible to try to delete."`
	Yes   bool   `short:"y" help:"Do not prompt for deletion confirmation."`
}

func (c *setDeleteCmd) Run(ctx *context) error {
	set, err := findOneSet(ctx, c.ID)
	if err != nil {
		return err
	}
	if set == nil {
		return nil
	}

	fmt.Printf("Snapshot Set info:\n")
	printSetInfo(set, "   ")
	fmt.Printf("\n")

	if !c.Yes {
		confirm := askForConfirmation("Are you sure you want to delete this snapshot set and all its snapshots?")
		if !confirm {
			return nil
		}
		fmt.Printf("\n")
	}

	deleted, err := ctx.snapshoter.DeleteSet(set.ID, c.Force)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.Errorf("Snapshot set not found.")
	}

	fmt.Printf("Snapshot set %v deleted.\n", set.ID)
	return nil
}
