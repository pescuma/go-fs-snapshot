package cli

import (
	"fmt"

	"github.com/pescuma/go-fs-snapshot/lib"
)

type infoCmd struct {
	ID string `arg:"" help:"The ID (simplified or full) of the snapshot."`
}

func (c *infoCmd) Run(ctx *context) error {
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

	return nil
}

func findOneSnapshot(ctx *context, id string) (*fs_snapshot.Snapshot, error) {
	snapshots, err := ctx.snapshoter.ListSnapshots(id)
	if err != nil {
		return nil, err
	}

	switch len(snapshots) {
	case 0:
		fmt.Printf("No snapshot found with ID %v\n", id)
		return nil, nil
	case 1:
		return snapshots[0], nil
	default:
		fmt.Printf("Found %v snapshots with ID %v - please use full ID.\n", len(snapshots), id)
		return nil, nil
	}
}

func printSnapshotInfo(snapshot *fs_snapshot.Snapshot, prefix string) {
	fmt.Printf("%vID:            %v\n", prefix, snapshot.ID)
	fmt.Printf("%vSet ID:        %v\n", prefix, snapshot.Set.ID)
	fmt.Printf("%vOriginal path: %v\n", prefix, snapshot.OriginalPath)
	fmt.Printf("%vSnapshot path: %v\n", prefix, snapshot.SnapshotPath)
	fmt.Printf("%vCreation:      %v\n", prefix, snapshot.CreationTime.Local().Format("2006-01-02 15:04:05 -07"))
	fmt.Printf("%vProvider:      %v\n", prefix, snapshot.Provider.Name)
	fmt.Printf("%vState:         %v\n", prefix, snapshot.State)
	fmt.Printf("%vAttributes:    %v\n", prefix, snapshot.Attributes)
}
