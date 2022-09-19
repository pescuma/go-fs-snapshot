package fs_snapshot

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pescuma/go-fs-snapshot/lib/internal/rpc"
)

// CurrentUserCanCreateSnapshots returns information if the current user can create snapshots
func CurrentUserCanCreateSnapshots(infoCb InfoMessageCallback) (bool, error) {
	if infoCb == nil {
		infoCb = func(level MessageLevel, format string, a ...interface{}) {}
	}

	can, err := currentUserCanCreateSnapshotsForOS(infoCb)
	if err == nil {
		return can, nil
	}

	can, err2 := serverCanCreateSnapshots(infoCb)
	if err2 == nil {
		return can, nil
	}

	return false, err
}

func serverCanCreateSnapshots(infoCb InfoMessageCallback) (bool, error) {
	addr := fmt.Sprintf("%v:%v", DefaultIP, DefaultPort)

	infoCb(InfoLevel, "Trying to open connection to server at: %v", addr)

	can, err := testServerCanCreateSnapshots(addr, infoCb)
	if err == nil {
		return can, nil
	}

	err = startServerForOS(infoCb)
	if err != nil {
		return false, nil
	}

	can, err = testServerCanCreateSnapshots(addr, infoCb)
	if err != nil {
		return false, nil
	}

	return can, nil
}

func testServerCanCreateSnapshots(addr string, infoCb InfoMessageCallback) (bool, error) {
	infoCb(TraceLevel, "GRPC Connecting to server at: %v", addr)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		infoCb(TraceLevel, "GRPC error: %v", err.Error())
		return false, err
	}
	defer conn.Close()

	client := rpc.NewFsSnapshotClient(conn)

	infoCb(TraceLevel, "GRPC Sending server request: CanCreateSnapshots()")

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := client.CanCreateSnapshots(ctx, &rpc.CanCreateSnapshotsRequest{})
	if err != nil {
		infoCb(TraceLevel, "GRPC error: %v", err.Error())
		return false, err
	}

	return reply.Can, nil
}

// EnableSnapshotsForUser enables the current user to run snapshots.
// This generally must be run from a prompt with elevated privileges (root or administrator).
func EnableSnapshotsForUser(username string, infoCb InfoMessageCallback) error {
	if infoCb == nil {
		infoCb = func(level MessageLevel, format string, a ...interface{}) {}
	}

	return enableSnapshotsForUserForOS(username, infoCb)
}
