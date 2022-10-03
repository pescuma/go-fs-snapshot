//go:build !windows && !darwin

package fs_snapshot

func startServerForOS(infoCb InfoMessageCallback) error {
	return ErrorNotSupportedInThisOS
}

func newSnapshoterForOS(cfg *SnapshoterConfig) (Snapshoter, error) {
	return nil, ErrorNotSupportedInThisOS
}
