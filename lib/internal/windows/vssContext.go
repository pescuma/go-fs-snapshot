//go:build windows

package internal_windows

// VssContext is a custom type for the Windows API VSSContext type.
type VssContext uint

// VSSContext constant values necessary for using VSS api.
//
//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	VSS_CTX_BACKUP                    VssContext = 0
	VSS_CTX_FILE_SHARE_BACKUP         VssContext = 0x00000010
	VSS_CTX_NAS_ROLLBACK              VssContext = 0x00000019
	VSS_CTX_APP_ROLLBACK              VssContext = 0x00000009
	VSS_CTX_CLIENT_ACCESSIBLE         VssContext = 0x0000001d
	VSS_CTX_CLIENT_ACCESSIBLE_WRITERS VssContext = 0x0000000d
	VSS_CTX_ALL                       VssContext = 0xffffffff
)
