package cli

import (
	"fmt"
	"time"

	"github.com/pescuma/go-fs-snapshot/lib"
)

type backupCmd struct {
	Dirs []string `arg:"" name:"dir" help:"Directories to snapshot and prepare to backup." type:"existingdir"`

	ProviderID string        `help:"Select which provider to use."`
	Timeout    time.Duration `help:"Timeout to create snapshot."`
	Simple     bool          `help:"Try to do it as simple as possible, but not simpler. In Windows this means do not use VSS Writers."`
}

func (c *backupCmd) Run(ctx *context) error {
	backuper, err := ctx.snapshoter.StartBackup(&fs_snapshot.CreateSnapshotOptions{
		ProviderID: c.ProviderID,
		Timeout:    c.Timeout,
		Simple:     c.Simple,
		InfoCallback: func(level fs_snapshot.MessageLevel, msg string) {
			switch level {
			case fs_snapshot.InfoLevel:
				if ctx.globals.Verbose >= 1 {
					fmt.Println("[INFO] " + msg)
				}
			case fs_snapshot.DetailsLevel:
				if ctx.globals.Verbose >= 2 {
					fmt.Println("[DETAILS] " + msg)
				}
			case fs_snapshot.TraceLevel:
				if ctx.globals.Verbose >= 3 {
					fmt.Println("[TRACE] " + msg)
				}
			}
		},
	})
	defer backuper.Close()

	if err != nil {
		return err
	}

	for _, dir := range c.Dirs {
		snapshotPath, err := backuper.TryToCreateTemporarySnapshot(dir)
		switch {
		case err != nil:
			fmt.Printf("%v: Error creating snapshot: %v\n", dir, err)
		case snapshotPath == dir:
			fmt.Printf("%v: Snapshots not supported for this folder\n", dir)
		default:
			fmt.Printf("%v: use snapshot path %v\n", dir, snapshotPath)
		}
	}

	fmt.Printf("\n")
	fmt.Printf("Press <enter> to finish backup and delete snapshot(s)\n")

	var response string
	_, _ = fmt.Scanln(&response)

	return nil
}
