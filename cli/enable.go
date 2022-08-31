package cli

import "os/user"

type enableForCurrentUserCmd struct {
}

func (c *enableForCurrentUserCmd) Run(ctx *context) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	return GrantPrivileges(u.Username, ctx.globals.Verbose)
}

type enableForUserCmd struct {
	Username string `arg:"" help:"Username to enable snapshots."`
}

func (c *enableForUserCmd) Run(ctx *context) error {
	return GrantPrivileges(c.Username, ctx.globals.Verbose)
}

type enableTestCurrentUserCmd struct {
}

func (c *enableTestCurrentUserCmd) Run(ctx *context) error {
	return TestPrivilege(ctx.globals.Verbose)
}
