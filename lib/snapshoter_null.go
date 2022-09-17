package fs_snapshot

func newNullSnapshoter() Snapshoter {
	return &nullSnapshoter{}
}

type nullSnapshoter struct {
}

func (s *nullSnapshoter) SimplifyID(id string) string {
	return id
}

func (s *nullSnapshoter) ListProviders(filterID string) ([]*Provider, error) {
	return nil, nil
}

func (s *nullSnapshoter) ListSets(filterID string) ([]*SnapshotSet, error) {
	return nil, nil
}

func (s *nullSnapshoter) ListSnapshots(filterID string) ([]*Snapshot, error) {
	return nil, nil
}

func (s *nullSnapshoter) DeleteSet(id string, force bool) (bool, error) {
	return false, nil
}

func (s *nullSnapshoter) DeleteSnapshot(id string, force bool) (bool, error) {
	return false, nil
}

func (s *nullSnapshoter) StartBackup(opts *BackupConfig) (Backuper, error) {
	return &nullBackuper{}, nil
}

func (s *nullSnapshoter) Close() {
}
