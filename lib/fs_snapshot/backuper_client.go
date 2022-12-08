package fs_snapshot

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/pescuma/go-fs-snapshot/lib/fs_snapshot/internal/rpc"
)

type clientBackuper struct {
	baseBackuper

	client       rpc.FsSnapshotClient
	backuperId   uint32
	timeout      time.Duration
	infoCallback InfoMessageCallback
}

func newClientBackuper(client rpc.FsSnapshotClient, backuperId uint32, caseSensitive bool, timeout time.Duration,
	listMountPoints func(volume string) ([]string, error),
	infoCallback InfoMessageCallback,
) *clientBackuper {

	result := &clientBackuper{}
	result.volumes = newVolumeInfos(caseSensitive)
	result.client = client
	result.backuperId = backuperId
	result.timeout = timeout
	result.infoCallback = infoCallback

	result.baseBackuper.listMountPoints = listMountPoints
	result.baseBackuper.createSnapshot = result.createSnapshot

	return result
}

func (b *clientBackuper) createSnapshot(m *mountPointInfo) (*Snapshot, error) {
	b.infoCallback(TraceLevel, "GRPC Sending server request: TryToCreateTemporarySnapshot(%v, \"%v\")",
		b.backuperId, m.dir)

	ctx, cancel := context.WithTimeout(context.Background(), b.timeout+time.Minute)
	defer cancel()

	stream, err := b.client.TryToCreateTemporarySnapshot(ctx, &rpc.TryToCreateTemporarySnapshotRequest{
		BackuperId: b.backuperId,
		Dir:        m.dir,
	})
	if err != nil {
		b.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return nil, err
	}

	var snapshot *Snapshot

	for {
		reply, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			b.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
			return nil, err
		}

		switch mr := reply.MessageOrResult.(type) {
		case *rpc.TryToCreateTemporarySnapshotReply_Message:
			b.infoCallback(MessageLevel(mr.Message.Level), "GRPC "+mr.Message.Message)

		case *rpc.TryToCreateTemporarySnapshotReply_Result:
			set := convertSnapshotSetToLocal(mr.Result.Snapshot.Set, false)
			snapshot = convertSnapshotToLocal(mr.Result.Snapshot, set)
		}
	}

	if snapshot == nil {
		return nil, errors.New("GRPC error: missing reply data")
	}

	return snapshot, nil
}

func (b *clientBackuper) Close() {
	b.infoCallback(TraceLevel, "GRPC Sending server request: CloseBackup(%v)", b.backuperId)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := b.client.CloseBackup(ctx, &rpc.CloseBackupRequest{
		BackuperId: b.backuperId,
	})
	if err != nil {
		b.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return
	}

	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			b.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
			return
		}

		b.infoCallback(MessageLevel(reply.Message.Level), "GRPC "+reply.Message.Message)
	}
}
