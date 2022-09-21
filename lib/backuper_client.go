package fs_snapshot

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/pescuma/go-fs-snapshot/lib/internal/rpc"
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
	result.volumes = newVolumeInfos()
	result.client = client
	result.backuperId = backuperId
	result.timeout = timeout
	result.infoCallback = infoCallback

	result.baseBackuper.caseSensitive = caseSensitive
	result.baseBackuper.listMountPoints = listMountPoints
	result.baseBackuper.createSnapshot = result.createSnapshot
	result.baseBackuper.deleteSnapshot = func(m *mountPointInfo) error { return nil }

	return result
}

func (b *clientBackuper) createSnapshot(m *mountPointInfo) (string, error) {
	b.infoCallback(TraceLevel, "GRPC Sending server request: TryToCreateTemporarySnapshot(%v, \"%v\")",
		b.backuperId, m.path)

	ctx, cancel := context.WithTimeout(context.Background(), b.timeout+time.Minute)
	defer cancel()

	stream, err := b.client.TryToCreateTemporarySnapshot(ctx, &rpc.TryToCreateTemporarySnapshotRequest{
		BackuperId: b.backuperId,
		Dir:        m.path,
	})
	if err != nil {
		b.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
		return "", err
	}

	snapshotDir := ""

	for {
		reply, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			b.infoCallback(TraceLevel, "GRPC error: %v", err.Error())
			return "", err
		}

		switch mr := reply.MessageOrResult.(type) {
		case *rpc.TryToCreateTemporarySnapshotReply_Message:
			b.infoCallback(MessageLevel(mr.Message.Level), mr.Message.Message)

		case *rpc.TryToCreateTemporarySnapshotReply_Result:
			snapshotDir = mr.Result.SnapshotDir
		}
	}

	if snapshotDir == "" {
		return "", errors.New("GRPC error: missing reply data")
	}

	return snapshotDir, nil
}

func (b *clientBackuper) Close() {
	b.baseBackuper.close()

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

		b.infoCallback(MessageLevel(reply.Message.Level), reply.Message.Message)
	}
}
