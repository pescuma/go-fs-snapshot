package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/pescuma/go-fs-snapshot/lib"
)

func Execute() {
	var cmds commands
	ctx := kong.Parse(&cmds, kong.ShortUsageOnError())

	var s fs_snapshot.Snapshoter
	var err error

	if ctx.Path[1].Command.Name != "enable" {
		s, err = fs_snapshot.NewSnapshoter()
		defer s.Close()

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
			return
		}
	}

	err = ctx.Run(&context{
		globals:    &cmds.Globals,
		snapshoter: s,
	})

	ctx.FatalIfErrorf(err)
}

type context struct {
	globals    *globals
	snapshoter fs_snapshot.Snapshoter
}

type commands struct {
	Globals globals `embed:""`

	List   listCmd   `cmd:"" help:"List snapshots."`
	Info   infoCmd   `cmd:"" help:"Show information of a snapshot."`
	Delete deleteCmd `cmd:"" help:"Delete a snapshot."`
	Backup backupCmd `cmd:"" help:"Create snapshots to do a backup."`

	Provider struct {
		List providerListCmd `cmd:"" help:"List available snapshot providers."`
	} `cmd:"" help:"Commands related to snapshot providers."`

	Set struct {
		List   setListCmd   `cmd:"" help:"List available snapshot sets."`
		Info   setInfoCmd   `cmd:"" help:"Show information of a snapshot set."`
		Delete setDeleteCmd `cmd:"" help:"Delete a snapshot set."`
	} `cmd:"" help:"Commands related to snapshot sets."`

	Enable struct {
		For struct {
			CurrentUser enableForCurrentUserCmd `cmd:""`
			User        enableForUserCmd        `cmd:""`
		} `cmd:""`
		Test enableTestCurrentUserCmd `cmd:"" help:"Test if the current user can create snapshots."`
	} `cmd:"" help:"Enable users to create snapshots."`
}

type globals struct {
	Verbose int `short:"v" type:"counter" help:"Show more detailed information."`
}

func askForConfirmation(message string) bool {
	fmt.Printf("%s [y/N] ", message)

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		return true
	} else if response == "" || response == "n" || response == "no" {
		return false
	} else {
		fmt.Printf("Unknown answer, aborting")
		return false
	}
}
