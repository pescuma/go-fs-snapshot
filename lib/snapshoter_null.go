package fs_snapshot

func createNullSnapshoter() Snapshoter {
	return &nullSnapshoter{}
}

type nullSnapshoter struct {
}

func (s *nullSnapshoter) SimplifyId(id string) string {
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

func (s *nullSnapshoter) StartBackup(opts *CreateSnapshotOptions) (Backuper, error) {
	return &nullBackuper{}, nil
}

func (s *nullSnapshoter) Close() {
}
