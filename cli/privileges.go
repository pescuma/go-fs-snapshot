//go:build windows

package cli

import (
	"fmt"
	"sort"

	"github.com/fourcorelabs/wintoken"
	"github.com/pkg/errors"
	"github.com/winlabs/gowin32"
	"github.com/winlabs/gowin32/wrappers"
	"golang.org/x/sys/windows"
)

const backupPrivilege = "SeBackupPrivilege"

func InitializePrivileges() error {
	toekn, err := gowin32.OpenCurrentProcessTokenWithAccess(wrappers.TOKEN_ADJUST_PRIVILEGES)
	if err != nil {
		return err
	}
	defer toekn.Close()

	privilege, err := gowin32.GetPrivilege(backupPrivilege)
	if err != nil {
		return err
	}

	err = toekn.EnablePrivilege(privilege, true)
	if err != nil {
		return err
	}

	token, err := wintoken.OpenProcessToken(0, wintoken.TokenPrimary)
	if err != nil {
		return errors.Wrap(err, "Failed to get process token")
	}
	defer token.Close()

	privileges, err := token.GetPrivileges()
	if err != nil {
		return err
	}

	fmt.Println("Has privileges:")
	has := false
	for _, p := range privileges {
		fmt.Println("   ", p.String())
		if p.Name == backupPrivilege {
			has = true
		}
	}

	if !has {
		return errors.New("Current user does not have backup privileges.")
	}

	err = token.EnableTokenPrivilege(backupPrivilege)
	if err != nil {
		return err
	}

	privileges, err = token.GetPrivileges()
	if err != nil {
		return err
	}

	sort.Slice(privileges, func(i, j int) bool {
		return privileges[i].Name < privileges[j].Name
	})
	for _, p := range privileges {
		fmt.Println(p.String())
	}

	return nil
}

func TestPrivilege(verbose int) error {
	token, err := wintoken.OpenProcessToken(0, wintoken.TokenPrimary)
	if err != nil {
		return errors.Wrap(err, "Failed to get process token")
	}
	defer token.Close()

	user, err := token.UserDetails()
	if err != nil {
		return err
	}
	username := fmt.Sprintf("%v\\%v", user.Domain, user.Username)

	privileges, err := token.GetPrivileges()
	if err != nil {
		return err
	}

	if verbose >= 1 {
		sort.Slice(privileges, func(i, j int) bool {
			return privileges[i].Name < privileges[j].Name
		})

		fmt.Printf("User %v has the privileges:\n", username)
		for _, p := range privileges {
			fmt.Println("   ", p.Name)
		}
		fmt.Println()
	}

	has := false
	for _, p := range privileges {
		if p.Name == backupPrivilege {
			has = true
		}
	}

	if has {
		fmt.Printf("Current user (%v) can create snapshots.\n", username)
	} else {
		fmt.Printf("Current user (%v) can NOT create snapshots.\n", username)
	}

	return nil
}

func GrantPrivileges(username string, verbose int) error {
	policy, err := gowin32.OpenLocalSecurityPolicy()
	if err != nil {
		return err
	}
	defer policy.Close()

	sid, domain, sidType, err := gowin32.GetLocalAccountByName(username)
	if err != nil {
		return err
	}

	if verbose >= 1 {
		fmt.Printf("User information:\n   user: %v\n   sid: %v\n   domain: %v\n   type: %v\n",
			username, sid, domain, sidType)
		fmt.Println()
	}

	rights, err := policy.GetAccountRights(sid)
	if err != nil && err != windows.ERROR_FILE_NOT_FOUND { // https://stackoverflow.com/a/4615926
		return err
	}

	has := false
	for _, right := range rights {
		if right == backupPrivilege {
			has = true
		}
	}

	if verbose >= 2 {
		sort.Slice(rights, func(i, j int) bool {
			return rights[i] < rights[j]
		})

		fmt.Printf("User %v had the privileges:\n", username)
		for _, p := range rights {
			fmt.Println("   ", p)
		}
		if len(rights) == 0 {
			fmt.Println("   ", "<none>")

		}
		fmt.Println()
	}

	if !has {
		fmt.Printf("Granting %v to user %v ...\n", backupPrivilege, username)

		err = policy.AddAccountRight(sid, backupPrivilege)
		if err != nil {
			return err
		}

	} else {
		fmt.Printf("User %v already has %v\n", username, backupPrivilege)
	}

	if verbose >= 1 {
		rights, err = policy.GetAccountRights(sid)
		if err != nil && err != windows.ERROR_NDIS_FILE_NOT_FOUND { // https://stackoverflow.com/a/4615926
			return err
		}

		sort.Slice(rights, func(i, j int) bool {
			return rights[i] < rights[j]
		})

		fmt.Println()
		fmt.Printf("User %v has the privileges:\n", username)
		for _, p := range rights {
			fmt.Println("   ", p)
		}
		if len(rights) == 0 {
			fmt.Println("   ", "<none>")

		}
	}

	return nil
}
