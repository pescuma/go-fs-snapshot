//go:build windows

package fs_snapshot_windows

// vssContext is a custom type for the windows api vssContext type.
type vssContext uint

// vssContext constant values necessary for using VSS api.
//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	VSS_CTX_BACKUP vssContext = iota
	VSS_CTX_FILE_SHARE_BACKUP
	VSS_CTX_NAS_ROLLBACK
	VSS_CTX_APP_ROLLBACK
	VSS_CTX_CLIENT_ACCESSIBLE
	VSS_CTX_CLIENT_ACCESSIBLE_WRITERS
	VSS_CTX_ALL
)
