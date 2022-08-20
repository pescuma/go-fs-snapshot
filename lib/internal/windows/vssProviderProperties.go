package internal_fs_snapshot_windows

import (
	"fmt"
	"unsafe"

	"github.com/go-ole/go-ole"
)

// VssProviderProperties defines the properties of a VSS provider as part of the VSS api.
type VssProviderProperties struct {
	ProviderID        ole.GUID
	ProviderName      *uint16
	ProviderType      VssProviderType
	ProviderVersion   *uint16
	ProviderVersionID ole.GUID
	ClassID           ole.GUID
	padding           [52]byte // Has to have the same size as VssSnapshotProperties
}

func (p *VssProviderProperties) Close() {
	ole.CoTaskMemFree(uintptr(unsafe.Pointer(p.ProviderName)))
	p.ProviderName = nil
	ole.CoTaskMemFree(uintptr(unsafe.Pointer(p.ProviderVersion)))
	p.ProviderName = nil
}

type VssProviderType uint32

//goland:noinspection GoSnakeCaseUsage
const (
	VSS_PROV_UNKNOWN   VssProviderType = 0
	VSS_PROV_SYSTEM    VssProviderType = 1
	VSS_PROV_SOFTWARE  VssProviderType = 2
	VSS_PROV_HARDWARE  VssProviderType = 3
	VSS_PROV_FILESHARE VssProviderType = 4
)

func (t VssProviderType) Str() string {
	switch t {
	case VSS_PROV_UNKNOWN:
		return "Unknown"
	case VSS_PROV_SYSTEM:
		return "System"
	case VSS_PROV_SOFTWARE:
		return "Software"
	case VSS_PROV_HARDWARE:
		return "Hardware"
	case VSS_PROV_FILESHARE:
		return "File share"
	default:
		panic(fmt.Sprintf("Unknown type: %v", int(t)))
	}
}
