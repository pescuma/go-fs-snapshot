package cli

import (
	"fmt"
	"strings"

	"github.com/pescuma/go-fs-snapshot/lib"
)

type console struct {
	verbosity int

	lastLevel fs_snapshot.MessageLevel
}

func newConsole(verbosity int) *console {
	return &console{
		verbosity: verbosity,
		lastLevel: -1,
	}
}

func (c *console) NewInfoMessageCallback() func(level fs_snapshot.MessageLevel, format string, a ...interface{}) {
	return func(level fs_snapshot.MessageLevel, format string, a ...interface{}) {
		if int(level) > c.verbosity {
			return
		}

		c.Printlf(level, format, a)
	}
}

func (c *console) Print(msg string) {
	c.Printl(fs_snapshot.OutputLevel, msg)
}

func (c *console) Printf(format string, a ...interface{}) {
	c.Printl(fs_snapshot.OutputLevel, fmt.Sprintf(format, a...))
}

func (c *console) Printlf(level fs_snapshot.MessageLevel, format string, a ...interface{}) {
	c.Printl(level, fmt.Sprintf(format, a...))
}

func (c *console) Printl(level fs_snapshot.MessageLevel, msg string) {
	c.printLevelSeparation(level)

	switch level {
	case fs_snapshot.OutputLevel:
		fmt.Println(msg)

	case fs_snapshot.InfoLevel:
		fmt.Println(msg)

	case fs_snapshot.DetailsLevel:
		msgs := strings.Split(msg, "\n")
		for _, m := range msgs {
			fmt.Println("[DETAIL] " + m)
		}

	case fs_snapshot.TraceLevel:
		msgs := strings.Split(msg, "\n")
		for _, m := range msgs {
			fmt.Println("[TRACE] " + m)
		}
	}
}

func (c *console) AskForConfirmation(message string) bool {
	c.printLevelSeparation(fs_snapshot.OutputLevel)

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
		fmt.Printf("Unknown answer, considering it as NO")
		return false
	}
}

func (c *console) printLevelSeparation(level fs_snapshot.MessageLevel) {
	if c.lastLevel != -1 && c.lastLevel != level {
		fmt.Println()
	}
	c.lastLevel = level
}
