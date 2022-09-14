package fs_snapshot

import (
	"os"
	"os/exec"
	"strings"
)

func runAndReturnOutput(infoCb InfoMessageCallback, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)

	infoCb(TraceLevel, "Running: '%v' '%v'", cmd.Path, strings.Join(cmd.Args[1:], "' '"))

	outputBytes, err := cmd.CombinedOutput()

	output := strings.TrimSpace(string(outputBytes))
	if output != "" {
		infoCb(TraceLevel, output)
	}

	return output, err
}

func run(infoCb InfoMessageCallback, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)

	infoCb(TraceLevel, "Running: '%v' '%v'", cmd.Path, strings.Join(cmd.Args[1:], "' '"))

	outputBytes, err := cmd.CombinedOutput()

	output := strings.TrimSpace(string(outputBytes))
	if output != "" {
		infoCb(TraceLevel, output)
	}

	return err
}

func runInline(infoCb InfoMessageCallback, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	infoCb(TraceLevel, "Running: '%v' '%v'", cmd.Path, strings.Join(cmd.Args[1:], "' '"))

	err := cmd.Run()

	return err
}
