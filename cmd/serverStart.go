package cli

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	fs_snapshot "github.com/pescuma/go-fs-snapshot/lib"
)

type serverStartCmd struct {
	Bind           string        `help:"Address to bind, in the format ip:port. Both can be empty, but the : must be there."`
	InactivityTime time.Duration `help:"After how long without a request should the server shut down. Default is never."`
}

func (c *serverStartCmd) Run(ctx *context) error {
	var err error
	cfg := fs_snapshot.ServerConfig{}

	if c.Bind != "" {
		parts := strings.Split(c.Bind, ":")
		if len(parts) != 2 {
			return errors.Errorf("Invalid bind address: %v", c.Bind)
		}

		cfg.IP = parts[0]
		if parts[1] != "" {
			cfg.Port, err = strconv.Atoi(parts[1])
			if err != nil {
				return errors.Wrapf(err, "Invalid bind address: %v", c.Bind)
			}
		}
	}

	cfg.InactivityTime = c.InactivityTime
	cfg.InfoCallback = outputMessages(ctx)

	err = fs_snapshot.StartServer(ctx.snapshoter, &cfg)
	if err != nil {
		return err
	}

	return nil
}
