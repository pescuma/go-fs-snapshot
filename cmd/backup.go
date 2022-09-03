package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"

	"github.com/pescuma/go-fs-snapshot/lib"
)

type backupCmd struct {
	Dirs []string `arg:"" name:"dir" help:"Directories to snapshot and prepare to backup." type:"existingdir"`

	ProviderID string        `help:"Select which provider to use."`
	Timeout    time.Duration `help:"Timeout to create snapshot."`
	Simple     bool          `help:"Try to do it as simple as possible, but not simpler. In Windows this means do not use VSS Writers."`
	Exec       string        `short:"e" help:"Command to execute after taking the snapshot. The snaphshot path(s) will be added to the end. If not set, this command waits for user input before deleting the snapshot(s)."`
	NoShell    bool          `help:"Do not pass the exec command to the shell to execute."`
}

func (c *backupCmd) Run(ctx *context) error {
	cmd, err := c.createExecCommand()
	if err != nil {
		return err
	}

	backuper, err := ctx.snapshoter.StartBackup(&fs_snapshot.SnapshotOptions{
		ProviderID:   c.ProviderID,
		Timeout:      c.Timeout,
		Simple:       c.Simple,
		InfoCallback: outputMessages(ctx),
	})
	defer backuper.Close()
	if err != nil {
		return err
	}

	var snapshotPaths []string

	for _, dir := range c.Dirs {
		snapshotPath, err := backuper.TryToCreateTemporarySnapshot(dir)
		switch {
		case err != nil:
			fmt.Printf("%v: Error creating snapshot: %v\n", dir, err)
		case snapshotPath == dir:
			fmt.Printf("%v: Snapshots not supported for this folder\n", dir)
		default:
			fmt.Printf("%v: Snapshot path is %v\n", dir, snapshotPath)
		}

		snapshotPaths = append(snapshotPaths, snapshotPath)
	}

	fmt.Printf("\n")

	if cmd != nil {
		cmd.Args = append(cmd.Args, snapshotPaths...)

		if ctx.globals.Verbose >= 1 {
			fmt.Printf("Executing: '%v'\n", strings.Join(cmd.Args, "' '"))
			fmt.Printf("\n")
		}

		err = cmd.Run()
		fmt.Printf("\n")

		return errors.Wrapf(err, "Error executing command")

	} else {
		fmt.Printf("Press <enter> to finish backup and delete snapshot(s)\n")

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
