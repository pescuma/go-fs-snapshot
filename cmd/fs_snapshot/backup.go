package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"

	"github.com/pescuma/go-fs-snapshot/lib/fs_snapshot"
)

type backupCmd struct {
	Dirs []string `arg:"" name:"dir" help:"Directories to snapshot and prepare to backup." type:"existingdir"`

	ProviderID string        `help:"Select which provider to use."`
	Timeout    time.Duration `help:"Timeout to create snapshot."`
	Simple     bool          `help:"Try to do it as simple as possible, but not simpler. In Windows this means do not use VSS Writers."`
	Exec       string        `short:"e" help:"Command to execute after taking the snapshot. The snaphshot path(s) will be added to the end. If not set, this command waits for user input before deleting the snapshot(s)."`
	NoShell    bool          `help:"Do not pass the exec command to the shell to execute."`

	ServerArgs serverArgs `embed:""`
}

func (c *backupCmd) Run(ctx *context) error {
	cmd, err := c.createExecCommand()
	if err != nil {
		return err
	}

	backuper, err := ctx.snapshoter.StartBackup(&fs_snapshot.BackupConfig{
		ProviderID: c.ProviderID,
		Timeout:    c.Timeout,
		Simple:     c.Simple,
	})
	if err != nil {
		return err
	}

	defer backuper.Close()

	var snapshotDirs []string

	for _, dir := range c.Dirs {
		snapshotDir, _, err := backuper.TryToCreateTemporarySnapshot(dir)
		switch {
		case err != nil:
			ctx.console.Printf("%v: Error creating snapshot: %v", dir, err)
		default:
			ctx.console.Printf("%v: Snapshot path is %v", dir, snapshotDir)
		}

		snapshotDirs = append(snapshotDirs, snapshotDir)
	}

	ctx.console.Print("")

	if cmd != nil {
		cmd.Args = append(cmd.Args, snapshotDirs...)

		if ctx.globals.Verbose >= 1 {
			ctx.console.Printf("Executing: '%v'", strings.Join(cmd.Args, "' '"))
			ctx.console.Print("")
		}

		err = cmd.Run()

		ctx.console.Print("")

		if err != nil {
			return errors.Wrap(err, "Error executing command")
		}

		return nil

	} else {
		fmt.Print("Press <enter> to finish backup and delete snapshot(s)")

		var response string
		_, _ = fmt.Scanln(&response)

		return nil
	}
}

func (c *backupCmd) createExecCommand() (*exec.Cmd, error) {
	if c.Exec == "" {
		return nil, nil
	}

	envs, args, err := shellwords.ParseWithEnvs(c.Exec)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing command to execute")
	}

	if !c.NoShell {
		if runtime.GOOS == "windows" {
			args = append([]string{os.Getenv("COMSPEC"), "/c"}, args...)
		} else {
			args = append([]string{"sh", "-e"}, args...)
		}
	}

	cmd := exec.Command(args[0], args[1:]...)
	if cmd.Err != nil {
		return nil, cmd.Err
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), envs...)

	return cmd, nil
}
