//go:build !windows

package fs_snapshot

import (
	"github.com/pkg/errors"
)

func newOSSnapshoter() (Snapshoter, error) {
	return nil, errors.New("snapshots not supported in this OS")
}
