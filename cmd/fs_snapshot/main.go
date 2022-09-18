package main

import (
	"reflect"

	"github.com/alecthomas/kong"

	"github.com/pescuma/go-fs-snapshot/lib"
)

var (
	version   = "?"
	buildDate = "?"
)

func main() {
	var cmds commands
	ctx := kong.Parse(&cmds, kong.ShortUsageOnError())

	err := execute(ctx, &cmds.Globals)

	ctx.FatalIfErrorf(err)
}

func execute(ctx *kong.Context, gs *globals) error {
	c := newConsole(gs.Verbose)

	var s fs_snapshot.Snapshoter

	sa := getServerArgs(ctx)
	if sa != nil {
		ip, port, err := parseAddr(sa.Server)
		if err != nil {
			return err
		}

		ct := fs_snapshot.LocalOrServer
		if sa.Server != "" && !sa.ServerOnlyAsFallback {
			ct = fs_snapshot.ServerOnly
		}

		s, err = fs_snapshot.NewSnapshoter(&fs_snapshot.SnapshoterConfig{
			InfoCallback:   c.NewInfoMessageCallback(),
			ConnectionType: ct,
			ServerIP:       ip,
			ServerPort:     port,
		})
		defer s.Close()

		if err != nil {
			return err
		}
	}

	return ctx.Run(&context{
		globals:    gs,
		snapshoter: s,
		console:    c,
	})
}

func getServerArgs(ctx *kong.Context) *serverArgs {
	var cmd reflect.Value

	for _, p := range ctx.Path {
		if p.Command != nil {
			cmd = p.Command.Target
		}
	}

	if !cmd.IsValid() {
		return nil
	}

	r := cmd.FieldByName("ServerArgs")
	if !r.IsValid() {
		return nil
	}

	sa := r.Interface().(serverArgs)
	return &sa
}

type context struct {
	globals    *globals
	snapshoter fs_snapshot.Snapshoter
	console    *console
}

type commands struct {
	Globals globals `embed:""`

	Version versionCmd `cmd:"" help:"Print version information."`
	List    listCmd    `cmd:"" help:"List snapshots."`
	Info    infoCmd    `cmd:"" help:"Show information of a snapshot."`
	Delete  deleteCmd  `cmd:"" help:"Delete a snapshot."`
	Backup  backupCmd  `cmd:"" help:"Create snapshots to do a backup."`

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

	Server struct {
		Start serverStartCmd `cmd:"" help:"Start a server and wait for commands."`
	} `cmd:"" help:"Commands related to a fs_snapshot server."`
}

type globals struct {
	Verbose int `short:"v" type:"counter" help:"Show more detailed information."`
}
