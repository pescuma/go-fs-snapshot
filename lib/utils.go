package fs_snapshot

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// absolutePath is only needed on windows, but because of the server we need it to always be there.
func absolutePath(path string) (string, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return path, err
	}

	// If starts with \\?\ and is not a file share, remove it
	if strings.HasPrefix(abspath, `\\?\`) && !strings.HasPrefix(abspath, `\\?\UNC\`) {
		return abspath[4:], nil
	}

	return abspath, nil
}

func changeBaseDir(path string, oldBase string, newBase string) (string, error) {
	relative, err := filepath.Rel(oldBase, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(newBase, relative), nil
}

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
