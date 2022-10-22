package fs_snapshot

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pescuma/go-fs-snapshot/fs_snapshot/internal/rpc"
)

func newClientSnapshoter(cfg *SnapshoterConfig) (Snapshoter, error) {
	var err error

	result := &clientSnapshoter{
		infoCallback: cfg.InfoCallback,
	}

	addr := fmt.Sprintf("%v:%v", cfg.ServerIP, cfg.ServerPort)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result.conn, err = grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to server: %v", addr)
	}

	cfg.InfoCallback(DetailsLevel, "Connected to server at: %v", addr)

	result.client = rpc.NewFsSnapshotClient(result.conn)

	return result, nil
}

type clientSnapshoter struct {
	conn         *grpc.ClientConn
	client       rpc.FsSnapshotClient
	infoCallback InfoMessageCallback
}

func (s *clientSnapshoter) ListProviders(filterID string) ([]*Provider, error) {
	s.infoCallback(TraceLevel, "GRPC Sending server request: ListProviders(\"%v\")", filterID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reply, err := s.client.ListProviders(ctx, &rpc.ListProvidersRequest{
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reply, err := s.client.ListSets(ctx, &rpc.ListSetsRequest{
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reply, err := s.client.ListSnapshots(ctx, &rpc.ListSnapshotsRequest{
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reply, err := s.client.SimplifyId(ctx, &rpc.SimplifyIdRequest{
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reply, err := s.client.DeleteSet(ctx, &rpc.DeleteRequest{
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reply, err := s.client.DeleteSnapshot(ctx, &rpc.DeleteRequest{
		Id:    id,
		Force: force,
	})
	if err != nil {
		s.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return false, err
	}

	return reply.Deleted, nil
}

func (s *clientSnapshoter) ListMountPoints(volume string) ([]string, error) {
	s.infoCallback(TraceLevel, "GRPC Sending server request: ListMountPoints(\"%v\")", volume)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reply, err := s.client.ListMountPoints(ctx, &rpc.ListMountPointsRequest{
		Volume: volume,
	})
	if err != nil {
		s.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return nil, err
	}

	return reply.MountPoints, nil
}

func (s *clientSnapshoter) StartBackup(cfg *BackupConfig) (Backuper, error) {
	ic := cfg.InfoCallback
	if ic == nil {
		ic = s.infoCallback
	}

	ic(TraceLevel, "GRPC Sending server request: StartBackup(\"%v\", %v, %v)",
		cfg.ProviderID, int32(cfg.Timeout.Seconds()), cfg.Simple)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := s.client.StartBackup(ctx, &rpc.StartBackupRequest{
		ProviderId:   cfg.ProviderID,
		TimeoutInSec: int32(cfg.Timeout.Seconds()),
		Simple:       cfg.Simple,
	})
	if err != nil {
		ic(TraceLevel, "GRPC error: %v", err.Error())
		return nil, err
	}

	received := false
	backuperId := uint32(0)
	caseSensitive := true

	for {
		reply, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			ic(TraceLevel, "GRPC error: %v", err.Error())
			return nil, err
		}

		switch mr := reply.MessageOrResult.(type) {
		case *rpc.StartBackupReply_Message:
			ic(MessageLevel(mr.Message.Level), "GRPC "+mr.Message.Message)

		case *rpc.StartBackupReply_Result:
			received = true
			backuperId = mr.Result.BackuperId
			caseSensitive = mr.Result.CaseSensitive
		}
	}

	if !received {
		return nil, errors.New("GRPC error: missing reply data")
	}

	return newClientBackuper(s.client, backuperId, caseSensitive, cfg.Timeout, s.ListMountPoints, ic), nil
}

func (s *clientSnapshoter) Close() {
	_ = s.conn.Close()
}
