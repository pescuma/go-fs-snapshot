package cli

import (
	"github.com/pkg/errors"
)

type setDeleteCmd struct {
	ID    string `arg:"" help:"The ID (simplified or full) of the snapshot set to delete."`
	Force bool   `short:"f" help:"Do everything possible to try to delete."`
	Yes   bool   `short:"y" help:"Do not prompt for deletion confirmation."`

	ServerArgs serverArgs `embed:""`
}

func (c *setDeleteCmd) Run(ctx *context) error {
	set, err := findOneSet(ctx, c.ID)
	if err != nil {
		return err
	}
	if set == nil {
		return nil
	}

	ctx.console.Print("Snapshot Set info:")
	printSetInfo(ctx, set, "   ")
	ctx.console.Print("")

	if !c.Yes {
		confirm := ctx.console.AskForConfirmation("Are you sure you want to delete this snapshot set and all its snapshots?")

		if !confirm {
			return nil
		}

		ctx.console.Print("")
	}

	deleted, err := ctx.snapshoter.DeleteSet(set.ID, c.Force)
	if err != nil {
		return err
	}

	if !deleted {
		return errors.Errorf("Snapshot set not found.")
	}

	ctx.console.Printf("Snapshot set %v deleted.", set.ID)
	return nil
}
