package cli

import (
	"fmt"
	"strings"

	"github.com/pescuma/go-fs-snapshot/lib"
)

func outputMessages(ctx *context) func(level fs_snapshot.MessageLevel, msg string) {
	return func(level fs_snapshot.MessageLevel, msg string) {
		switch level {
		case fs_snapshot.OutputLevel:
			fmt.Println(msg)
			fmt.Println()

		case fs_snapshot.InfoLevel:
			if ctx.globals.Verbose >= 1 {
				fmt.Println(msg)
				fmt.Println()
			}

		case fs_snapshot.DetailsLevel:
			if ctx.globals.Verbose >= 2 {
				fmt.Println(msg)
				fmt.Println()
			}

		case fs_snapshot.TraceLevel:
			if ctx.globals.Verbose >= 3 {
				msgs := strings.Split(msg, "\n")
				for _, m := range msgs {
					fmt.Println("[TRACE] " + m)
				}
				fmt.Println()
			}
		}
	}
}
