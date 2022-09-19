//go:build ignore

package main

import (
	"github.com/pescuma/go-build"
)

func main() {
	cfg := build.NewBuilderConfig()
	cfg.Archs = []string{
		"windows/386", "windows/amd64",
		"darwin",
		//"linux/amd64",
	}

	b, err := build.NewBuilder(cfg)
	if err != nil {
		panic(err)
	}

	err = b.RunTarget("all")
	if err != nil {
		panic(err)
	}
}
