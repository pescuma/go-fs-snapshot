package main

import (
	"github.com/pescuma/go-fs-snapshot/lib"
)

type infoCmd struct {
	ID string `arg:"" help:"The ID (simplified or full) of the snapshot."`

	ServerArgs serverArgs `embed:""`
}

func (c *infoCmd) Run(ctx *context) error {
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

	return nil
}

func findOneSnapshot(ctx *context, id string) (*fs_snapshot.Snapshot, error) {
	snapshots, err := ctx.snapshoter.ListSnapshots(id)
	if err != nil {
		return nil, err
	}

	switch len(snapshots) {
	case 0:
		ctx.console.Printf("No snapshot found with ID %v", id)
		return nil, nil
	case 1:
		return snapshots[0], nil
	default:
		ctx.console.Printf("Found %v snapshots with ID %v - please use full ID.", len(snapshots), id)
		return nil, nil
	}
}

func printSnapshotInfo(ctx *context, snapshot *fs_snapshot.Snapshot, prefix string) {
	ctx.console.Printf("%vID:            %v", prefix, snapshot.ID)
	ctx.console.Printf("%vSet ID:        %v", prefix, snapshot.Set.ID)
	ctx.console.Printf("%vOriginal path: %v", prefix, snapshot.OriginalPath)
	ctx.console.Printf("%vSnapshot path: %v", prefix, snapshot.SnapshotPath)
	ctx.console.Printf("%vCreation:      %v", prefix, snapshot.CreationTime.Local().Format("2006-01-02 15:04:05 -07"))
	ctx.console.Printf("%vProvider:      %v", prefix, snapshot.Provider.Name)
	ctx.console.Printf("%vState:         %v", prefix, snapshot.State)
	ctx.console.Printf("%vAttributes:    %v", prefix, snapshot.Attributes)
}
