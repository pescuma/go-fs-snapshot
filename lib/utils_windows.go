package fs_snapshot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-ole/go-ole"
)

func changeBaseDir(path string, oldBase string, newBase string) (string, error) {
	relative, err := filepath.Rel(oldBase, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(newBase, relative), nil
}

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

func toGuidString(id ole.GUID) string {
	result := id.String()
	result = result[1 : len(result)-1]
	result = strings.ToLower(result)
	return result
}

func toDate(filetime uint64) time.Time {
	result := time.Date(1601, 1, 1,
		0, 0, 0, 0,
		time.UTC,
	)

	// HACK Cant multiply by 100 or it overflows
	for i := 0; i < 100; i++ {
		result = result.Add(time.Duration(filetime))
	}

	return result
}

func createScheduledTaskName(username string) string {
	return fmt.Sprintf(`\fs_stapshot\server start (%v)`, strings.ReplaceAll(username, `\`, `_`))
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
