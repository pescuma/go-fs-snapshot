//go:build windows

package internal_fs_snapshot_windows

// vssBackup is a custom type for the Windows API VSSBackup type.
type vssBackup uint

// VSSBackup constant values necessary for using VSS api.
//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	VSS_BT_UNDEFINED vssBackup = iota
	VSS_BT_FULL
	VSS_BT_INCREMENTAL
	VSS_BT_DIFFERENTIAL
	VSS_BT_LOG
	VSS_BT_COPY
	VSS_BT_OTHER
)
