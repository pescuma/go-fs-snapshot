package fs_snapshot

type nullBackuper struct {
}

func (b *nullBackuper) TryToCreateTemporarySnapshot(path string) (string, error) {
	return path, nil
}

func (b *nullBackuper) ListSnapshotedDirectories() map[string]string {
	return nil
}

func (b *nullBackuper) Close() {
}
