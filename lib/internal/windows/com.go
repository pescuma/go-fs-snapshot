//go:build windows

package internal_fs_snapshot_windows

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

// HRESULT is a custom type for the Windows api HRESULT type.
type HRESULT uint

// HRESULT constant values necessary for using VSS api.
//
//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
const (
	S_OK                                            HRESULT = 0x00000000
	S_FALSE                                         HRESULT = 0x00000001
	E_ACCESSDENIED                                  HRESULT = 0x80070005
	E_OUTOFMEMORY                                   HRESULT = 0x8007000E
	E_INVALIDARG                                    HRESULT = 0x80070057
	VSS_E_BAD_STATE                                 HRESULT = 0x80042301
	VSS_E_UNEXPECTED                                HRESULT = 0x80042302
	VSS_E_PROVIDER_ALREADY_REGISTERED               HRESULT = 0x80042303
	VSS_E_PROVIDER_NOT_REGISTERED                   HRESULT = 0x80042304
	VSS_E_PROVIDER_VETO                             HRESULT = 0x80042306
	VSS_E_PROVIDER_IN_USE                           HRESULT = 0x80042307
	VSS_E_OBJECT_NOT_FOUND                          HRESULT = 0x80042308
	VSS_E_VOLUME_NOT_SUPPORTED                      HRESULT = 0x8004230C
	VSS_E_VOLUME_NOT_SUPPORTED_BY_PROVIDER          HRESULT = 0x8004230E
	VSS_E_OBJECT_ALREADY_EXISTS                     HRESULT = 0x8004230D
	VSS_E_UNEXPECTED_PROVIDER_ERROR                 HRESULT = 0x8004230F
	VSS_E_CORRUPT_XML_DOCUMENT                      HRESULT = 0x80042310
	VSS_E_INVALID_XML_DOCUMENT                      HRESULT = 0x80042311
	VSS_E_MAXIMUM_NUMBER_OF_VOLUMES_REACHED         HRESULT = 0x80042312
	VSS_E_FLUSH_WRITES_TIMEOUT                      HRESULT = 0x80042313
	VSS_E_HOLD_WRITES_TIMEOUT                       HRESULT = 0x80042314
	VSS_E_UNEXPECTED_WRITER_ERROR                   HRESULT = 0x80042315
	VSS_E_SNAPSHOT_SET_IN_PROGRESS                  HRESULT = 0x80042316
	VSS_E_MAXIMUM_NUMBER_OF_SNAPSHOTS_REACHED       HRESULT = 0x80042317
	VSS_E_WRITER_INFRASTRUCTURE                     HRESULT = 0x80042318
	VSS_E_WRITER_NOT_RESPONDING                     HRESULT = 0x80042319
	VSS_E_WRITER_ALREADY_SUBSCRIBED                 HRESULT = 0x8004231A
	VSS_E_UNSUPPORTED_CONTEXT                       HRESULT = 0x8004231B
	VSS_E_VOLUME_IN_USE                             HRESULT = 0x8004231D
	VSS_E_MAXIMUM_DIFFAREA_ASSOCIATIONS_REACHED     HRESULT = 0x8004231E
	VSS_E_INSUFFICIENT_STORAGE                      HRESULT = 0x8004231F
	VSS_E_NO_SNAPSHOTS_IMPORTED                     HRESULT = 0x80042320
	VSS_E_SOME_SNAPSHOTS_NOT_IMPORTED               HRESULT = 0x80042321
	VSS_E_MAXIMUM_NUMBER_OF_REMOTE_MACHINES_REACHED HRESULT = 0x80042322
	VSS_E_REMOTE_SERVER_UNAVAILABLE                 HRESULT = 0x80042323
	VSS_E_REMOTE_SERVER_UNSUPPORTED                 HRESULT = 0x80042324
	VSS_E_REVERT_IN_PROGRESS                        HRESULT = 0x80042325
	VSS_E_REVERT_VOLUME_LOST                        HRESULT = 0x80042326
	VSS_E_REBOOT_REQUIRED                           HRESULT = 0x80042327
	VSS_E_TRANSACTION_FREEZE_TIMEOUT                HRESULT = 0x80042328
	VSS_E_TRANSACTION_THAW_TIMEOUT                  HRESULT = 0x80042329
	VSS_E_VOLUME_NOT_LOCAL                          HRESULT = 0x8004232D
	VSS_E_CLUSTER_TIMEOUT                           HRESULT = 0x8004232E
	VSS_E_WRITERERROR_INCONSISTENTSNAPSHOT          HRESULT = 0x800423F0
	VSS_E_WRITERERROR_OUTOFRESOURCES                HRESULT = 0x800423F1
	VSS_E_WRITERERROR_TIMEOUT                       HRESULT = 0x800423F2
	VSS_E_WRITERERROR_RETRYABLE                     HRESULT = 0x800423F3
	VSS_E_WRITERERROR_NONRETRYABLE                  HRESULT = 0x800423F4
	VSS_E_WRITERERROR_RECOVERY_FAILED               HRESULT = 0x800423F5
	VSS_E_BREAK_REVERT_ID_FAILED                    HRESULT = 0x800423F6
	VSS_E_LEGACY_PROVIDER                           HRESULT = 0x800423F7
	VSS_E_MISSING_DISK                              HRESULT = 0x800423F8
	VSS_E_MISSING_HIDDEN_VOLUME                     HRESULT = 0x800423F9
	VSS_E_MISSING_VOLUME                            HRESULT = 0x800423FA
	VSS_E_AUTORECOVERY_FAILED                       HRESULT = 0x800423FB
	VSS_E_DYNAMIC_DISK_ERROR                        HRESULT = 0x800423FC
	VSS_E_NONTRANSPORTABLE_BCD                      HRESULT = 0x800423FD
	VSS_E_CANNOT_REVERT_DISKID                      HRESULT = 0x800423FE
	VSS_E_RESYNC_IN_PROGRESS                        HRESULT = 0x800423FF
	VSS_E_CLUSTER_ERROR                             HRESULT = 0x80042400
	VSS_E_UNSELECTED_VOLUME                         HRESULT = 0x8004232A
	VSS_E_SNAPSHOT_NOT_IN_SET                       HRESULT = 0x8004232B
	VSS_E_NESTED_VOLUME_LIMIT                       HRESULT = 0x8004232C
	VSS_E_NOT_SUPPORTED                             HRESULT = 0x8004232F
	VSS_E_WRITERERROR_PARTIAL_FAILURE               HRESULT = 0x80042336
	VSS_E_WRITER_STATUS_NOT_AVAILABLE               HRESULT = 0x80042409
)

// hresultToString maps a HRESULT value to a human-readable string.
//
//goland:noinspection SpellCheckingInspection
var hresultToString = map[HRESULT]string{
	S_OK:                                            "S_OK",
	E_ACCESSDENIED:                                  "E_ACCESSDENIED",
	E_OUTOFMEMORY:                                   "E_OUTOFMEMORY",
	E_INVALIDARG:                                    "E_INVALIDARG",
	VSS_E_BAD_STATE:                                 "VSS_E_BAD_STATE",
	VSS_E_UNEXPECTED:                                "VSS_E_UNEXPECTED",
	VSS_E_PROVIDER_ALREADY_REGISTERED:               "VSS_E_PROVIDER_ALREADY_REGISTERED",
	VSS_E_PROVIDER_NOT_REGISTERED:                   "VSS_E_PROVIDER_NOT_REGISTERED",
	VSS_E_PROVIDER_VETO:                             "VSS_E_PROVIDER_VETO",
	VSS_E_PROVIDER_IN_USE:                           "VSS_E_PROVIDER_IN_USE",
	VSS_E_OBJECT_NOT_FOUND:                          "VSS_E_OBJECT_NOT_FOUND",
	VSS_E_VOLUME_NOT_SUPPORTED:                      "VSS_E_VOLUME_NOT_SUPPORTED",
	VSS_E_VOLUME_NOT_SUPPORTED_BY_PROVIDER:          "VSS_E_VOLUME_NOT_SUPPORTED_BY_PROVIDER",
	VSS_E_OBJECT_ALREADY_EXISTS:                     "VSS_E_OBJECT_ALREADY_EXISTS",
	VSS_E_UNEXPECTED_PROVIDER_ERROR:                 "VSS_E_UNEXPECTED_PROVIDER_ERROR",
	VSS_E_CORRUPT_XML_DOCUMENT:                      "VSS_E_CORRUPT_XML_DOCUMENT",
	VSS_E_INVALID_XML_DOCUMENT:                      "VSS_E_INVALID_XML_DOCUMENT",
	VSS_E_MAXIMUM_NUMBER_OF_VOLUMES_REACHED:         "VSS_E_MAXIMUM_NUMBER_OF_VOLUMES_REACHED",
	VSS_E_FLUSH_WRITES_TIMEOUT:                      "VSS_E_FLUSH_WRITES_TIMEOUT",
	VSS_E_HOLD_WRITES_TIMEOUT:                       "VSS_E_HOLD_WRITES_TIMEOUT",
	VSS_E_UNEXPECTED_WRITER_ERROR:                   "VSS_E_UNEXPECTED_WRITER_ERROR",
	VSS_E_SNAPSHOT_SET_IN_PROGRESS:                  "VSS_E_SNAPSHOT_SET_IN_PROGRESS",
	VSS_E_MAXIMUM_NUMBER_OF_SNAPSHOTS_REACHED:       "VSS_E_MAXIMUM_NUMBER_OF_SNAPSHOTS_REACHED",
	VSS_E_WRITER_INFRASTRUCTURE:                     "VSS_E_WRITER_INFRASTRUCTURE",
	VSS_E_WRITER_NOT_RESPONDING:                     "VSS_E_WRITER_NOT_RESPONDING",
	VSS_E_WRITER_ALREADY_SUBSCRIBED:                 "VSS_E_WRITER_ALREADY_SUBSCRIBED",
	VSS_E_UNSUPPORTED_CONTEXT:                       "VSS_E_UNSUPPORTED_CONTEXT",
	VSS_E_VOLUME_IN_USE:                             "VSS_E_VOLUME_IN_USE",
	VSS_E_MAXIMUM_DIFFAREA_ASSOCIATIONS_REACHED:     "VSS_E_MAXIMUM_DIFFAREA_ASSOCIATIONS_REACHED",
	VSS_E_INSUFFICIENT_STORAGE:                      "VSS_E_INSUFFICIENT_STORAGE",
	VSS_E_NO_SNAPSHOTS_IMPORTED:                     "VSS_E_NO_SNAPSHOTS_IMPORTED",
	VSS_E_SOME_SNAPSHOTS_NOT_IMPORTED:               "VSS_E_SOME_SNAPSHOTS_NOT_IMPORTED",
	VSS_E_MAXIMUM_NUMBER_OF_REMOTE_MACHINES_REACHED: "VSS_E_MAXIMUM_NUMBER_OF_REMOTE_MACHINES_REACHED",
	VSS_E_REMOTE_SERVER_UNAVAILABLE:                 "VSS_E_REMOTE_SERVER_UNAVAILABLE",
	VSS_E_REMOTE_SERVER_UNSUPPORTED:                 "VSS_E_REMOTE_SERVER_UNSUPPORTED",
	VSS_E_REVERT_IN_PROGRESS:                        "VSS_E_REVERT_IN_PROGRESS",
	VSS_E_REVERT_VOLUME_LOST:                        "VSS_E_REVERT_VOLUME_LOST",
	VSS_E_REBOOT_REQUIRED:                           "VSS_E_REBOOT_REQUIRED",
	VSS_E_TRANSACTION_FREEZE_TIMEOUT:                "VSS_E_TRANSACTION_FREEZE_TIMEOUT",
	VSS_E_TRANSACTION_THAW_TIMEOUT:                  "VSS_E_TRANSACTION_THAW_TIMEOUT",
	VSS_E_VOLUME_NOT_LOCAL:                          "VSS_E_VOLUME_NOT_LOCAL",
	VSS_E_CLUSTER_TIMEOUT:                           "VSS_E_CLUSTER_TIMEOUT",
	VSS_E_WRITERERROR_INCONSISTENTSNAPSHOT:          "VSS_E_WRITERERROR_INCONSISTENTSNAPSHOT",
	VSS_E_WRITERERROR_OUTOFRESOURCES:                "VSS_E_WRITERERROR_OUTOFRESOURCES",
	VSS_E_WRITERERROR_TIMEOUT:                       "VSS_E_WRITERERROR_TIMEOUT",
	VSS_E_WRITERERROR_RETRYABLE:                     "VSS_E_WRITERERROR_RETRYABLE",
	VSS_E_WRITERERROR_NONRETRYABLE:                  "VSS_E_WRITERERROR_NONRETRYABLE",
	VSS_E_WRITERERROR_RECOVERY_FAILED:               "VSS_E_WRITERERROR_RECOVERY_FAILED",
	VSS_E_BREAK_REVERT_ID_FAILED:                    "VSS_E_BREAK_REVERT_ID_FAILED",
	VSS_E_LEGACY_PROVIDER:                           "VSS_E_LEGACY_PROVIDER",
	VSS_E_MISSING_DISK:                              "VSS_E_MISSING_DISK",
	VSS_E_MISSING_HIDDEN_VOLUME:                     "VSS_E_MISSING_HIDDEN_VOLUME",
	VSS_E_MISSING_VOLUME:                            "VSS_E_MISSING_VOLUME",
	VSS_E_AUTORECOVERY_FAILED:                       "VSS_E_AUTORECOVERY_FAILED",
	VSS_E_DYNAMIC_DISK_ERROR:                        "VSS_E_DYNAMIC_DISK_ERROR",
	VSS_E_NONTRANSPORTABLE_BCD:                      "VSS_E_NONTRANSPORTABLE_BCD",
	VSS_E_CANNOT_REVERT_DISKID:                      "VSS_E_CANNOT_REVERT_DISKID",
	VSS_E_RESYNC_IN_PROGRESS:                        "VSS_E_RESYNC_IN_PROGRESS",
	VSS_E_CLUSTER_ERROR:                             "VSS_E_CLUSTER_ERROR",
	VSS_E_UNSELECTED_VOLUME:                         "VSS_E_UNSELECTED_VOLUME",
	VSS_E_SNAPSHOT_NOT_IN_SET:                       "VSS_E_SNAPSHOT_NOT_IN_SET",
	VSS_E_NESTED_VOLUME_LIMIT:                       "VSS_E_NESTED_VOLUME_LIMIT",
	VSS_E_NOT_SUPPORTED:                             "VSS_E_NOT_SUPPORTED",
	VSS_E_WRITERERROR_PARTIAL_FAILURE:               "VSS_E_WRITERERROR_PARTIAL_FAILURE",
	VSS_E_WRITER_STATUS_NOT_AVAILABLE:               "VSS_E_WRITER_STATUS_NOT_AVAILABLE",
}

// Str converts a HRESULT to a human-readable string.
func (h HRESULT) Str() string {
	if i, ok := hresultToString[h]; ok {
		return i
	}

	return "UNKNOWN"
}

// apiBoolToInt converts a bool for use calling the VSS api
func apiBoolToInt(input bool) uint {
	switch input {
	case true:
		return 1
	default:
		return 0
	}
}

// apiBoolToInt converts an int from calling the VSS api to a bool
func apiIntToBool(input uint32) bool {
	switch input {
	case 0:
		return false
	default:
		return true
	}
}

// queryInterface is a wrapper around the windows QueryInterface api.
func queryInterface(oleIUnknown *ole.IUnknown, guid *ole.GUID) (*interface{}, error) {
	var ivss *interface{}

	err := syscallN("QueryInterface", oleIUnknown.VTable().QueryInterface,
		uintptr(unsafe.Pointer(oleIUnknown)), uintptr(unsafe.Pointer(guid)),
		uintptr(unsafe.Pointer(&ivss)))

	return ivss, err
}

// syscallN is a wapper for syscall.SyscallN that creates a nice error as return
func syscallN(name string, trap uintptr, args ...uintptr) error {
	result, _, _ := syscall.SyscallN(trap, args...)

	hresult := HRESULT(result)
	if hresult != S_OK {
		return newVssErrorF(HRESULT(result), "%s failed", name)
	}

	return nil
}

// syscallNF is a wapper for syscall.SyscallN that creates a nice error as return
func syscallNF(name string, trap uintptr, args ...uintptr) error {
	result, _, _ := syscall.SyscallN(trap, args...)

	hresult := HRESULT(result)
	if hresult != S_OK && hresult != S_FALSE {
		return newVssErrorF(HRESULT(result), "%s failed", name)
	}

	return nil
}
