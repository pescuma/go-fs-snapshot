//go:build windows

package fs_snapshot_windows

// vssObjectType is a custom type for the windows api vssObjectType type.
type vssObjectType uint

// vssObjectType constant values necessary for using VSS api.
//goland:noinspection ALL
const (
	VSS_OBJECT_UNKNOWN vssObjectType = iota
	VSS_OBJECT_NONE
	VSS_OBJECT_SNAPSHOT_SET
	VSS_OBJECT_SNAPSHOT
	VSS_OBJECT_PROVIDER
	VSS_OBJECT_TYPE_COUNT
)
