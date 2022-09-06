package fs_snapshot

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pescuma/go-fs-snapshot/lib/internal/rpc"
)

func newClientSnapshoter(cfg *SnapshoterConfig) (Snapshoter, error) {
	var err error

	result := &clientSnapshoter{
		infoCallback: cfg.InfoCallback,
	}

	addr := fmt.Sprintf("%v:%v", cfg.ServerIP, cfg.ServerPort)

	connCtx, _ := context.WithTimeout(context.Background(), time.Second)

	result.conn, err = grpc.DialContext(connCtx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to server: %v", addr)
	}

	cfg.InfoCallback(DetailsLevel, "Connected to server at: %v", addr)

	result.client = rpc.NewFsSnapshotClient(result.conn)
	result.ctx, result.cancel = context.WithCancel(context.Background())

	return result, nil
}

type clientSnapshoter struct {
	conn         *grpc.ClientConn
	client       rpc.FsSnapshotClient
	ctx          context.Context
	cancel       context.CancelFunc
	infoCallback InfoMessageCallback
}

func (s *clientSnapshoter) ListProviders(filterID string) ([]*Provider, error) {
	s.infoCallback(TraceLevel, "GRPC Sending server request: ListProviders(\"%v\")", filterID)

	reply, err := s.client.ListProviders(s.ctx, &rpc.ListProvidersRequest{
		FilterId: filterID,
	})
	if err != nil {
		s.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return nil, err
	}

	result := make([]*Provider, len(reply.Providers))
	for i, p := range reply.Providers {
		result[i] = convertProviderToLocal(p)
	}

	return result, nil
}

func (s *clientSnapshoter) ListSets(filterID string) ([]*SnapshotSet, error) {
	s.infoCallback(TraceLevel, "GRPC Sending server request: ListSets(\"%v\")", filterID)

	reply, err := s.client.ListSets(s.ctx, &rpc.ListSetsRequest{
		FilterId: filterID,
	})
	if err != nil {
		s.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return nil, err
	}

	result := make([]*SnapshotSet, len(reply.Sets))
	for i, set := range reply.Sets {
		result[i] = convertSnapshotSetToLocal(set, true)
	}

	return result, nil
}

func (s *clientSnapshoter) ListSnapshots(filterID string) ([]*Snapshot, error) {
	s.infoCallback(TraceLevel, "GRPC Sending server request: ListSnapshots(\"%v\")", filterID)

	reply, err := s.client.ListSnapshots(s.ctx, &rpc.ListSnapshotsRequest{
		FilterId: filterID,
	})
	if err != nil {
		s.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return nil, err
	}

	result := make([]*Snapshot, len(reply.Snapshots))
	setsById := make(map[string]*SnapshotSet)

	for i, snap := range reply.Snapshots {
		set, exists := setsById[snap.Set.Id]
		if !exists {
			set = convertSnapshotSetToLocal(snap.Set, false)
			setsById[snap.Set.Id] = set
		}

		result[i] = convertSnapshotToLocal(snap, set)

		set.Snapshots = append(set.Snapshots, result[i])
	}

	return result, nil
}

func (s *clientSnapshoter) SimplifyID(id string) string {
	s.infoCallback(TraceLevel, "GRPC Sending server request: SimplifyId(\"%v\")", id)

	reply, err := s.client.SimplifyId(s.ctx, &rpc.SimplifyIdRequest{
		Id: id,
	})
	if err != nil {
		s.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return id // Nothing else that can be done here
	}

	return reply.SimpleId
}

func (s *clientSnapshoter) DeleteSet(id string, force bool) (bool, error) {
	s.infoCallback(TraceLevel, "GRPC Sending server request: DeleteSet(\"%v\", %v)", id, force)

	reply, err := s.client.DeleteSet(s.ctx, &rpc.DeleteRequest{
		Id:    id,
		Force: force,
	})
	if err != nil {
		s.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return false, err
	}

	return reply.Deleted, nil
}

func (s *clientSnapshoter) DeleteSnapshot(id string, force bool) (bool, error) {
	s.infoCallback(TraceLevel, "GRPC Sending server request: DeleteSnapshot(\"%v\", %v)", id, force)

	reply, err := s.client.DeleteSnapshot(s.ctx, &rpc.DeleteRequest{
		Id:    id,
		Force: force,
	})
	if err != nil {
		s.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return false, err
	}

	return reply.Deleted, nil
}

func (s *clientSnapshoter) StartBackup(opts *SnapshotOptions) (Backuper, error) {
	//TODO implement me
	panic("implement me")
}

func (s *clientSnapshoter) Close() {
	s.cancel()
	_ = s.conn.Close()
}
