package fs_snapshot

func newNullBackuper() *nullBackuper {
	return &nullBackuper{}
}

type nullBackuper struct {
}

func (b *nullBackuper) TryToCreateTemporarySnapshot(dir string) (string, *Snapshot, error) {
	return dir, nil, nil
}

func (b *nullBackuper) ListSnapshotedDirectories() map[string]string {
	return nil
}

func (b *nullBackuper) Close() {
}
