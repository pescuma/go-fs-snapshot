package fs_snapshot

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/pescuma/go-fs-snapshot/lib/internal/rpc"
)

const DefaultIP = "localhost"
const DefaultPort = 33721

func StartServer(snapshoter Snapshoter, cfg *ServerConfig) error {
	if cfg == nil {
		cfg = &ServerConfig{}
	}
	cfg.setDefaults()

	lis, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		return errors.Wrapf(err, "failed to listen to %v", cfg.Address())
	}

	s := grpc.NewServer()
	rpc.RegisterFsSnapshotServer(s, &server{
		snapshoter:   snapshoter,
		infoCallback: cfg.InfoCallback,
	})

	cfg.InfoCallback(OutputLevel, "fs_snapshot server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		return err
	}

	return nil
}

type ServerConfig struct {
	IP             string
	Port           int
	InactivityTime time.Duration
	InfoCallback   InfoMessageCallback
}

func (cfg *ServerConfig) setDefaults() {
	if cfg.IP == "" {
		cfg.IP = DefaultIP
	}
	if cfg.Port == 0 {
		cfg.Port = DefaultPort
	}
	if cfg.InfoCallback == nil {
		cfg.InfoCallback = func(level MessageLevel, format string, a ...interface{}) {}
	}
}

func (cfg *ServerConfig) Address() string {
	return fmt.Sprintf("%v:%d", cfg.IP, cfg.Port)
}

type server struct {
	rpc.UnimplementedFsSnapshotServer

	snapshoter   Snapshoter
	infoCallback InfoMessageCallback
}

func (s *server) CanCreateSnapshots(ctx context.Context, request *rpc.CanCreateSnapshotsRequest) (*rpc.CanCreateSnapshotsReply, error) {
	s.infoCallback(TraceLevel, "GRPC Received request: Snapshoter.CanCreateSnapshots()")

	// If the server has started it can create snapshots, at least for now
	reply := rpc.CanCreateSnapshotsReply{
		Can: true,
	}

	return &reply, nil
}

func (s *server) ListProviders(ctx context.Context, request *rpc.ListProvidersRequest) (*rpc.ListProvidersReply, error) {
	s.infoCallback(TraceLevel, "GRPC Received request: Snapshoter.ListProviders(\"%v\")", request.FilterId)

	providers, err := s.snapshoter.ListProviders(request.FilterId)
	if err != nil {
		return nil, err
	}

	reply := rpc.ListProvidersReply{
		Providers: make([]*rpc.Provider, len(providers)),
	}
	for i, p := range providers {
		reply.Providers[i] = convertProvidertToRPC(p)
	}

	return &reply, nil
}

func (s *server) ListSets(ctx context.Context, request *rpc.ListSetsRequest) (*rpc.ListSetsReply, error) {
	s.infoCallback(TraceLevel, "GRPC Received request: Snapshoter.ListSets(\"%v\")", request.FilterId)

	sets, err := s.snapshoter.ListSets(request.FilterId)
	if err != nil {
		return nil, err
	}

	reply := rpc.ListSetsReply{
		Sets: make([]*rpc.SnapshotSet, len(sets)),
	}
	for i, set := range sets {
		reply.Sets[i] = convertSnapshotSetToRPC(set, true)
	}

	return &reply, nil
}

func (s *server) ListSnapshots(ctx context.Context, request *rpc.ListSnapshotsRequest) (*rpc.ListSnapshotsReply, error) {
	s.infoCallback(TraceLevel, "GRPC Received request: Snapshoter.ListSnapshots(\"%v\")", request.FilterId)

	snaps, err := s.snapshoter.ListSnapshots(request.FilterId)
	if err != nil {
		return nil, err
	}

	reply := rpc.ListSnapshotsReply{
		Snapshots: make([]*rpc.Snapshot, len(snaps)),
	}
	for i, snap := range snaps {
		reply.Snapshots[i] = convertSnapshotToRPC(snap, true)
	}

	return &reply, nil
}

func (s *server) SimplifyId(ctx context.Context, request *rpc.SimplifyIdRequest) (*rpc.SimplifyIdReply, error) {
	s.infoCallback(TraceLevel, "GRPC Received request: Snapshoter.SimplifyID(\"%v\")", request.Id)

	simpleId := s.snapshoter.SimplifyID(request.Id)

	return &rpc.SimplifyIdReply{
		SimpleId: simpleId,
	}, nil
}

func (s *server) DeleteSet(ctx context.Context, request *rpc.DeleteRequest) (*rpc.DeleteReply, error) {
	s.infoCallback(TraceLevel, "GRPC Received request: Snapshoter.DeleteSet(\"%v\", %v)", request.Id, request.Force)

	deleted, err := s.snapshoter.DeleteSet(request.Id, request.Force)
	if err != nil {
		return nil, err
	}

	return &rpc.DeleteReply{
		Deleted: deleted,
	}, nil
}

func (s *server) DeleteSnapshot(ctx context.Context, request *rpc.DeleteRequest) (*rpc.DeleteReply, error) {
	s.infoCallback(TraceLevel, "GRPC Received request: Snapshoter.DeleteSnapshot(\"%v\", %v)", request.Id, request.Force)

	deleted, err := s.snapshoter.DeleteSnapshot(request.Id, request.Force)
	if err != nil {
		return nil, err
	}

	return &rpc.DeleteReply{
		Deleted: deleted,
	}, nil
}

func (s *server) StartBackup(request *rpc.StartBackupRequest, backupServer rpc.FsSnapshot_StartBackupServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *server) TryToCreateTemporarySnapshot(request *rpc.TryToCreateTemporarySnapshotRequest, snapshotServer rpc.FsSnapshot_TryToCreateTemporarySnapshotServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *server) CloseBackup(request *rpc.CloseBackupRequest, backupServer rpc.FsSnapshot_CloseBackupServer) error {
	//TODO implement me
	panic("implement me")
}

func convertProvidertToRPC(p *Provider) *rpc.Provider {
	return &rpc.Provider{
		Id:      p.ID,
		Name:    p.Name,
		Version: p.Version,
		Type:    p.Type,
	}
}

func convertProviderToLocal(p *rpc.Provider) *Provider {
	return &Provider{
		ID:      p.Id,
		Name:    p.Name,
		Version: p.Version,
		Type:    p.Type,
	}
}

func convertSnapshotSetToRPC(set *SnapshotSet, includeSnapshots bool) *rpc.SnapshotSet {
	result := &rpc.SnapshotSet{
		Id:                      set.ID,
		CreationTime:            set.CreationTime.In(time.UTC).Unix(),
		SnapshotCountOnCreation: int32(set.SnapshotCountOnCreation),
	}

	if includeSnapshots {
		result.Snapshots = make([]*rpc.Snapshot, len(set.Snapshots))
		for i, snap := range set.Snapshots {
			result.Snapshots[i] = convertSnapshotToRPC(snap, false)
		}
	}

	return result
}

func convertSnapshotSetToLocal(set *rpc.SnapshotSet, includeSnapshots bool) *SnapshotSet {
	result := &SnapshotSet{
		ID:                      set.Id,
		CreationTime:            time.Unix(set.CreationTime, 0).UTC().In(time.Local),
		SnapshotCountOnCreation: int(set.SnapshotCountOnCreation),
	}

	if includeSnapshots {
		result.Snapshots = make([]*Snapshot, len(set.Snapshots))
		for i, snap := range set.Snapshots {
			result.Snapshots[i] = convertSnapshotToLocal(snap, result)
		}
	}

	return result
}

func convertSnapshotToRPC(snap *Snapshot, includeSet bool) *rpc.Snapshot {
	result := &rpc.Snapshot{
		Id:           snap.ID,
		OriginalPath: snap.OriginalPath,
		SnapshotPath: snap.SnapshotPath,
		CreationTime: snap.CreationTime.In(time.UTC).Unix(),
		Provider:     convertProvidertToRPC(snap.Provider),
		State:        snap.State,
		Attributes:   snap.Attributes,
	}

	if includeSet {
		result.Set = convertSnapshotSetToRPC(snap.Set, false)
	}

	return result
}

func convertSnapshotToLocal(snap *rpc.Snapshot, set *SnapshotSet) *Snapshot {
	return &Snapshot{
		ID:           snap.Id,
		OriginalPath: snap.OriginalPath,
		SnapshotPath: snap.SnapshotPath,
		CreationTime: time.Unix(snap.CreationTime, 0).UTC().In(time.Local),
		Set:          set,
		Provider:     convertProviderToLocal(snap.Provider),
		State:        snap.State,
		Attributes:   snap.Attributes,
	}
}
