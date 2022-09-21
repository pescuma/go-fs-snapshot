//go:build windows

package fs_snapshot

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-ole/go-ole"
)

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
