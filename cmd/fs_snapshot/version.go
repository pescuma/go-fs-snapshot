package main

import (
	"runtime"
	"strings"
)

type versionCmd struct {
}

func (c *versionCmd) Run(ctx *context) error {
	ctx.console.Printf("fs_snapshot version %s %v/%v (compiled on %v with go %v)",
		version, runtime.GOOS, runtime.GOARCH,
		buildDate, strings.TrimPrefix(runtime.Version(), "go"))

	return nil
}
