//go:build !windows && !darwin

package fs_snapshot

func startServerForOS(infoCb InfoMessageCallback) error {
	return ErrNotSupportedInThisOS
}

func newSnapshoterForOS(cfg *SnapshoterConfig) (Snapshoter, error) {
	return nil, ErrNotSupportedInThisOS
}
