//go:build ignore

package main

import (
	"github.com/pescuma/go-build"
)

func main() {
	cfg := build.NewBuilderConfig()
	cfg.Archs = []string{"windows", "darwin"}

	b, err := build.NewBuilder(cfg)
	if err != nil {
		panic(err)
	}

	err = b.RunTarget("build")
	if err != nil {
		panic(err)
	}
}
