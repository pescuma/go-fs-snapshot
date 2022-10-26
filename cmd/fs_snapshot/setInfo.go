package main

import (
	"github.com/pescuma/go-fs-snapshot/lib/fs_snapshot"
)

type setInfoCmd struct {
	ID string `arg:"" help:"The ID (simplified or full) of the snapshot set."`

	ServerArgs serverArgs `embed:""`
}

func (c *setInfoCmd) Run(ctx *context) error {
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

	return nil
}

func findOneSet(ctx *context, id string) (*fs_snapshot.SnapshotSet, error) {
	sets, err := ctx.snapshoter.ListSets(id)
	if err != nil {
		return nil, err
	}

	switch len(sets) {
	case 0:
		ctx.console.Printf("No snapshot sets found with ID %v", id)
		return nil, nil
	case 1:
		return sets[0], nil
	default:
		ctx.console.Printf("Found %v snapshot sets with ID %v - please use full ID.", len(sets), id)
		return nil, nil
	}
}

func printSetInfo(ctx *context, set *fs_snapshot.SnapshotSet, prefix string) {
	ctx.console.Printf("%vID:                         %v", prefix, set.ID)
	ctx.console.Printf("%vCreation:                   %v", prefix, set.CreationTime.Local().Format("2006-01-02 15:04:05 -07"))
	ctx.console.Printf("%vSnapshot count:             %v", prefix, len(set.Snapshots))
	ctx.console.Printf("%vSnapshot count on creation: %v", prefix, set.SnapshotCountOnCreation)
	for i, snapshot := range set.Snapshots {
		ctx.console.Printf("")
		ctx.console.Printf("%vSnapshot %v:", prefix, i+1)
		printSnapshotInfo(ctx, snapshot, prefix+"   ")
	}
}
