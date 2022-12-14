//go:build windows

package fs_snapshot

import (
	"fmt"
	"os"
	"os/user"
	"sort"
	"strings"

	"github.com/fourcorelabs/wintoken"
	"github.com/pkg/errors"
	"github.com/winlabs/gowin32"
	"golang.org/x/sys/windows"
)

const backupPrivilege = "SeBackupPrivilege"
const batchLogonPrivilege = "SeBatchLogonRight"

func currentUserCanCreateSnapshotsForOS(infoCb InfoMessageCallback) (bool, error) {
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

	return has, nil
}

func enableSnapshotsForUserForOS(username string, infoCb InfoMessageCallback) error {
	u, err := user.Lookup(username)
	if err != nil {
		return err
	}

	username = u.Username

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

	err = grantPrivilege(policy, username, sid, backupPrivilege, rights, infoCb)
	if err != nil {
		return err
	}

	infoCb(OutputLevel, "")
	err = grantPrivilege(policy, username, sid, batchLogonPrivilege, rights, infoCb)
	if err != nil {
		return err
	}

	infoCb(OutputLevel, "")
	infoCb(OutputLevel, "Creating scheduled task \"%v\"", createScheduledTaskName(username))
	err = createScheduledTask(username, infoCb)
	if err != nil {
		return err
	}

	return nil
}

func grantPrivilege(policy *gowin32.SecurityPolicy, username string, sid gowin32.SecurityID,
	toGrant gowin32.AccountRightName, rights []gowin32.AccountRightName, infoCb InfoMessageCallback) error {

	has := false
	for _, right := range rights {
		if right == toGrant {
			has = true
		}
	}

	if !has {
		infoCb(OutputLevel, "Granting %v to user %v", toGrant, username)

		err := policy.AddAccountRight(sid, toGrant)
		if err != nil {
			return err
		}

		rights, err := policy.GetAccountRights(sid)
		if err != nil && err != windows.ERROR_NDIS_FILE_NOT_FOUND { // https://stackoverflow.com/a/4615926
			return err
		}

		infoCb(InfoLevel, "User %v now has %v", username, rightsAsSring(rights))

	} else {
		infoCb(OutputLevel, "User %v already has %v", username, toGrant)
	}

	return nil
}

func createScheduledTask(username string, infoCb InfoMessageCallback) error {
	f, err := createScheduledTaskXML()
	if err != nil {
		return err
	}
	defer os.Remove(f)

	u, err := user.Current()
	if err != nil {
		return err
	}

	runMode := runInline
	if u.Username == username {
		runMode = run
	}

	err = runMode(infoCb, "schtasks", "/Create",
		"/TN", createScheduledTaskName(username),
		"/RU", username,
		"/NP",
		"/XML", f,
		"/F",
		"/HRESULT")
	if err != nil {
		return errors.Wrap(err, "error creating scheduled task")
	}

	return nil
}

func createScheduledTaskXML() (string, error) {
	f, err := os.CreateTemp("", "fs_snapshot-task-*.xml")
	if err != nil {
		return "", err
	}
	defer f.Close()

	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	data := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <Triggers />
  <Principals>
    <Principal id="Author">
      <LogonType>S4U</LogonType>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>false</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>false</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>%v</Command>
      <Arguments>server start --inactivity-time=5m</Arguments>
    </Exec>
  </Actions>
</Task>`, exe))
	if _, err = f.Write(data); err != nil {
		return "", err
	}

	return f.Name(), nil
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
		return errors.New("the current user does not have sufficient backup privileges or is not an administrator")
	}

	err = token.EnableTokenPrivilege(backupPrivilege)
	if err != nil {
		return err
	}

	return nil
}
