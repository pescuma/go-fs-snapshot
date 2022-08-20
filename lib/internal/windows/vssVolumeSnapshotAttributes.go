package internal_fs_snapshot_windows

import "strings"

type VssVolumeSnapshotAttributes uint32

//goland:noinspection GoSnakeCaseUsage
const (
	VSS_VOLSNAP_ATTR_PERSISTENT           VssVolumeSnapshotAttributes = 0x1
	VSS_VOLSNAP_ATTR_NO_AUTORECOVERY      VssVolumeSnapshotAttributes = 0x2
	VSS_VOLSNAP_ATTR_CLIENT_ACCESSIBLE    VssVolumeSnapshotAttributes = 0x4
	VSS_VOLSNAP_ATTR_NO_AUTO_RELEASE      VssVolumeSnapshotAttributes = 0x8
	VSS_VOLSNAP_ATTR_NO_WRITERS           VssVolumeSnapshotAttributes = 0x10
	VSS_VOLSNAP_ATTR_TRANSPORTABLE        VssVolumeSnapshotAttributes = 0x20
	VSS_VOLSNAP_ATTR_NOT_SURFACED         VssVolumeSnapshotAttributes = 0x40
	VSS_VOLSNAP_ATTR_NOT_TRANSACTED       VssVolumeSnapshotAttributes = 0x80
	VSS_VOLSNAP_ATTR_HARDWARE_ASSISTED    VssVolumeSnapshotAttributes = 0x10000
	VSS_VOLSNAP_ATTR_DIFFERENTIAL         VssVolumeSnapshotAttributes = 0x20000
	VSS_VOLSNAP_ATTR_PLEX                 VssVolumeSnapshotAttributes = 0x40000
	VSS_VOLSNAP_ATTR_IMPORTED             VssVolumeSnapshotAttributes = 0x80000
	VSS_VOLSNAP_ATTR_EXPOSED_LOCALLY      VssVolumeSnapshotAttributes = 0x100000
	VSS_VOLSNAP_ATTR_EXPOSED_REMOTELY     VssVolumeSnapshotAttributes = 0x200000
	VSS_VOLSNAP_ATTR_AUTORECOVER          VssVolumeSnapshotAttributes = 0x400000
	VSS_VOLSNAP_ATTR_ROLLBACK_RECOVERY    VssVolumeSnapshotAttributes = 0x800000
	VSS_VOLSNAP_ATTR_DELAYED_POSTSNAPSHOT VssVolumeSnapshotAttributes = 0x1000000
	VSS_VOLSNAP_ATTR_TXF_RECOVERY         VssVolumeSnapshotAttributes = 0x2000000
	VSS_VOLSNAP_ATTR_FILE_SHARE           VssVolumeSnapshotAttributes = 0x400000
)

func (a VssVolumeSnapshotAttributes) Str() string {
	var result []string

	if a&VSS_VOLSNAP_ATTR_PERSISTENT > 0 {
		result = append(result, "persistent")
	}
	if a&VSS_VOLSNAP_ATTR_NO_AUTORECOVERY > 0 {
		result = append(result, "no auto recovery")
	}
	if a&VSS_VOLSNAP_ATTR_CLIENT_ACCESSIBLE > 0 {
		result = append(result, "client-accessible")
	}
	if a&VSS_VOLSNAP_ATTR_NO_AUTO_RELEASE > 0 {
		result = append(result, "no auto release")
	}
	if a&VSS_VOLSNAP_ATTR_NO_WRITERS > 0 {
		result = append(result, "no writers")
	}
	if a&VSS_VOLSNAP_ATTR_TRANSPORTABLE > 0 {
		result = append(result, "transportable")
	}
	if a&VSS_VOLSNAP_ATTR_NOT_SURFACED > 0 {
		result = append(result, "not surfaced")
	}
	if a&VSS_VOLSNAP_ATTR_NOT_TRANSACTED > 0 {
		result = append(result, "not transacted")
	}
	if a&VSS_VOLSNAP_ATTR_HARDWARE_ASSISTED > 0 {
		result = append(result, "hardware assisted")
	}
	if a&VSS_VOLSNAP_ATTR_DIFFERENTIAL > 0 {
		result = append(result, "differential")
	}
	if a&VSS_VOLSNAP_ATTR_PLEX > 0 {
		result = append(result, "plex")
	}
	if a&VSS_VOLSNAP_ATTR_IMPORTED > 0 {
		result = append(result, "imported")
	}
	if a&VSS_VOLSNAP_ATTR_EXPOSED_LOCALLY > 0 {
		result = append(result, "exposed locally")
	}
	if a&VSS_VOLSNAP_ATTR_EXPOSED_REMOTELY > 0 {
		result = append(result, "exposed remotely")
	}
	if a&VSS_VOLSNAP_ATTR_AUTORECOVER > 0 {
		result = append(result, "autorecover")
	}
	if a&VSS_VOLSNAP_ATTR_ROLLBACK_RECOVERY > 0 {
		result = append(result, "rollback recovery")
	}
	if a&VSS_VOLSNAP_ATTR_DELAYED_POSTSNAPSHOT > 0 {
		result = append(result, "delayed postsnapshot")
	}
	if a&VSS_VOLSNAP_ATTR_TXF_RECOVERY > 0 {
		result = append(result, "txf recovery")
	}
	if a&VSS_VOLSNAP_ATTR_FILE_SHARE > 0 {
		result = append(result, "file share")
	}

	return strings.Join(result, ", ")
}
