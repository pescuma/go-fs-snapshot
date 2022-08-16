package fs_snapshot

type nullSnapshoter struct {
}

func createNullSnapshoter() Snapshoter {
	return &nullSnapshoter{}
}

func (*nullSnapshoter) CreateSnapshot(path string) string {
	return path
}

func (*nullSnapshoter) Destroy() error {
	return nil
}
