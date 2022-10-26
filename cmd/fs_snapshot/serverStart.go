package main

import (
	"time"

	"github.com/pescuma/go-fs-snapshot/lib/fs_snapshot"
)

type serverStartCmd struct {
	Bind           string        `help:"Address to bind, in the format ip:port. Both can be empty, but the : must be there."`
	InactivityTime time.Duration `help:"After how long without a request should the server shut down. Default is never."`
	Force          bool          `short:"f" help:"Start the server even if it's not supported in this OS. For tests.'"`
}

func (c *serverStartCmd) Run(ctx *context) error {
	var err error

	var ip string
	var port int
	if c.Bind != "" {
		ip, port, err = parseAddr(c.Bind)
		if err != nil {
			return err
		}
	}

	s, err := fs_snapshot.NewSnapshoter(&fs_snapshot.SnapshoterConfig{
		ConnectionType: fs_snapshot.LocalOnly,
		InfoCallback:   ctx.console.NewInfoMessageCallback(),
	})
	defer s.Close()
	if err != nil {
		if !c.Force {
			return err
		} else {
			ctx.console.Printf("Error: %v. Starting no-op server anyway.", err)
		}
	}

	err = fs_snapshot.StartServer(s, &fs_snapshot.ServerConfig{
		InactivityTime: c.InactivityTime,
		InfoCallback:   ctx.console.NewInfoMessageCallback(),
		IP:             ip,
		Port:           port,
	})
	if err != nil {
		return err
	}

	return nil
}
