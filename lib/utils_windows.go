package fs_snapshot

import (
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
