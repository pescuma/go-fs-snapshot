//go:build !windows

package fs_snapshot

import (
	"github.com/pkg/errors"
)

func startServerForOS(infoCb InfoMessageCallback) error {
	return errors.New("snapshots not supported in this OS")
}

func newSnapshoterForOS(cfg *SnapshoterConfig) (Snapshoter, error) {
	return nil, errors.New("snapshots not supported in this OS")
}
