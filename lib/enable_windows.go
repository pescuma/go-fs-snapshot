//go:build windows

package fs_snapshot

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fourcorelabs/wintoken"
	"github.com/pkg/errors"
	"github.com/winlabs/gowin32"
	"golang.org/x/sys/windows"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pescuma/go-fs-snapshot/lib/internal/rpc"
)

const backupPrivilege = "SeBackupPrivilege"

// CurrentUserCanCreateSnapshots returns information if the current user can create snapshots
func CurrentUserCanCreateSnapshots(infoCb InfoMessageCallback) (bool, error) {
	if infoCb == nil {
		infoCb = func(level MessageLevel, format string, a ...interface{}) {}
	}

	infoCb(TraceLevel, "TOKEN OpenProcessToken()")
	token, err := wintoken.OpenProcessToken(0, wintoken.TokenPrimary)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get process token")
	}
	defer token.Close()

	infoCb(TraceLevel, "TOKEN UserDetails()")
	user, err := token.UserDetails()
	if err != nil {
		return false, err
	}
	username := fmt.Sprintf("%v\\%v", user.Domain, user.Username)

	infoCb(TraceLevel, "TOKEN GetPrivileges()")
	privileges, err := token.GetPrivileges()
	if err != nil {
		return false, err
	}

	infoCb(InfoLevel, "User %v has %v", username, privilegesAsSring(privileges))

	has := false
	for _, p := range privileges {
		if p.Name == backupPrivilege {
			has = true
		}
	}

	if has {
		return true, nil
	}

	addr := fmt.Sprintf("%v:%v", DefaultIP, DefaultPort)

	infoCb(InfoLevel, "")
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
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		infoCb(TraceLevel, "GRPC error: %v", err.Error())
		return false, err
	}
	defer conn.Close()

	client := rpc.NewFsSnapshotClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	infoCb(TraceLevel, "GRPC Sending server request: CanCreateSnapshots()")
	reply, err := client.CanCreateSnapshots(ctx, &rpc.CanCreateSnapshotsRequest{})
	if err != nil {
		infoCb(TraceLevel, "GRPC error: %v", err.Error())
		return false, err
	}

	return reply.Can, nil
}

// EnableSnapshotsForUser enables the current user to run snaphsots.
// This generally must be run from a prompt with elevated privileges (root or administrator).
func EnableSnapshotsForUser(username string, infoCb InfoMessageCallback) error {
	if infoCb == nil {
		infoCb = func(level MessageLevel, format string, a ...interface{}) {}
	}

	policy, err := gowin32.OpenLocalSecurityPolicy()
	if err != nil {
		return err
	}
	defer func() {
		_ = policy.Close()
	}()

	sid, domain, sidType, err := gowin32.GetLocalAccountByName(username)
	if err != nil {
		return err
	}

	infoCb(InfoLevel, "User information:\n   user: %v\n   sid: %v\n   domain: %v\n   type: %v",
		username, sid, domain, sidType)

	rights, err := policy.GetAccountRights(sid)
	if err != nil && err != windows.ERROR_FILE_NOT_FOUND { // https://stackoverflow.com/a/4615926
		return err
	}

	infoCb(DetailsLevel, "User %v has %v", username, rightsAsSring(rights))

	has := false
	for _, right := range rights {
		if right == backupPrivilege {
			has = true
		}
	}

	if !has {
		infoCb(OutputLevel, "Granting %v to user %v", backupPrivilege, username)

		err = policy.AddAccountRight(sid, backupPrivilege)
		if err != nil {
			return err
		}

		rights, err = policy.GetAccountRights(sid)
		if err != nil && err != windows.ERROR_NDIS_FILE_NOT_FOUND { // https://stackoverflow.com/a/4615926
			return err
		}

		infoCb(InfoLevel, "User %v now has %v", username, rightsAsSring(rights))

	} else {
		infoCb(OutputLevel, "User %v already has %v", username, backupPrivilege)
	}

	return nil
}

func privilegesAsSring(privileges []wintoken.Privilege) string {
	if len(privileges) == 0 {
		return "NO privileges"
	}

	var sb strings.Builder

	sort.Slice(privileges, func(i, j int) bool {
		return privileges[i].Name < privileges[j].Name
	})

	sb.WriteString("the privileges:")
	for _, p := range privileges {
		sb.WriteString(fmt.Sprintf("\n   %v", p.Name))
	}

	return sb.String()
}

func rightsAsSring(rights []gowin32.AccountRightName) string {
	if len(rights) == 0 {
		return "NO privileges"
	}

	var sb strings.Builder

	sort.Slice(rights, func(i, j int) bool {
		return rights[i] < rights[j]
	})

	sb.WriteString("the privileges:")
	for _, p := range rights {
		sb.WriteString(fmt.Sprintf("\n   %v", p))
	}

	return sb.String()
}

func initializePrivileges() error {
	token, err := wintoken.OpenProcessToken(0, wintoken.TokenPrimary)
	if err != nil {
		return errors.Wrap(err, "failed to get process token")
	}
	defer token.Close()

	privileges, err := token.GetPrivileges()
	if err != nil {
		return err
	}

	has := false
	for _, p := range privileges {
		if p.Name == backupPrivilege {
			has = true
		}
	}

	if !has {
		return errors.New("current user does not have backup privileges")
	}

	err = token.EnableTokenPrivilege(backupPrivilege)
	if err != nil {
		return err
	}

	return nil
}
