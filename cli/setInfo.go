package cli

import (
	"fmt"

	"github.com/pescuma/go-fs-snapshot/lib"
)

type setInfoCmd struct {
	ID string `arg:"" help:"The ID (simplified or full) of the snapshot set."`
}

func (c *setInfoCmd) Run(ctx *context) error {
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

	return nil
}

func findOneSet(ctx *context, id string) (*fs_snapshot.SnapshotSet, error) {
	sets, err := ctx.snapshoter.ListSets(id)
	if err != nil {
		return nil, err
	}

	switch len(sets) {
	case 0:
		fmt.Printf("No snapshot sets found with ID %v\n", id)
		return nil, nil
	case 1:
		return sets[0], nil
	default:
		fmt.Printf("Found %v snapshot sets with ID %v - please use full ID.\n", len(sets), id)
		return nil, nil
	}
}

func printSetInfo(set *fs_snapshot.SnapshotSet, prefix string) {
	fmt.Printf("%vID:                         %v\n", prefix, set.ID)
	fmt.Printf("%vCreation:                   %v\n", prefix, set.CreationTime.Local().Format("2006-01-02 15:04:05 -07"))
	fmt.Printf("%vSnapshot count:             %v\n", prefix, len(set.Snapshots))
	fmt.Printf("%vSnapshot count on creation: %v\n", prefix, set.SnapshotCountOnCreation)
	for i, snapshot := range set.Snapshots {
		fmt.Printf("\n")
		fmt.Printf("%vSnapshot %v:\n", prefix, i+1)
		printSnapshotInfo(snapshot, prefix+"   ")
	}
}
