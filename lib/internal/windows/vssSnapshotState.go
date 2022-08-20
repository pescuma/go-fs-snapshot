package internal_fs_snapshot_windows

import "fmt"

type VssSnapshotState uint

//goland:noinspection GoSnakeCaseUsage
const (
	VSS_SS_UNKNOWN VssSnapshotState = iota
	VSS_SS_PREPARING
	VSS_SS_PROCESSING_PREPARE
	VSS_SS_PREPARED
	VSS_SS_PROCESSING_PRECOMMIT
	VSS_SS_PRECOMMITTED
	VSS_SS_PROCESSING_COMMIT
	VSS_SS_COMMITTED
	VSS_SS_PROCESSING_POSTCOMMIT
	VSS_SS_PROCESSING_PREFINALCOMMIT
	VSS_SS_PREFINALCOMMITTED
	VSS_SS_PROCESSING_POSTFINALCOMMIT
	VSS_SS_CREATED
	VSS_SS_ABORTED
	VSS_SS_DELETED
	VSS_SS_POSTCOMMITTED
	VSS_SS_COUNT
)

func (s VssSnapshotState) Str() string {
	switch s {
	case VSS_SS_UNKNOWN:
		return "Unknown"
	case VSS_SS_PREPARING:
		return "Preparing"
	case VSS_SS_PROCESSING_PREPARE:
		return "Processing prepare"
	case VSS_SS_PREPARED:
		return "Prepared"
	case VSS_SS_PROCESSING_PRECOMMIT:
		return "Processing precommit"
	case VSS_SS_PRECOMMITTED:
		return "Precommitted"
	case VSS_SS_PROCESSING_COMMIT:
		return "Processing commit"
	case VSS_SS_COMMITTED:
		return "Committed"
	case VSS_SS_PROCESSING_POSTCOMMIT:
		return "Processing postcommit"
	case VSS_SS_PROCESSING_PREFINALCOMMIT:
		return "Processing prefinalcommit"
	case VSS_SS_PREFINALCOMMITTED:
		return "Prefinalcommitted"
	case VSS_SS_PROCESSING_POSTFINALCOMMIT:
		return "Processing postfinalcommit"
	case VSS_SS_CREATED:
		return "Created"
	case VSS_SS_ABORTED:
		return "Aborted"
	case VSS_SS_DELETED:
		return "Deleted"
	case VSS_SS_POSTCOMMITTED:
		return "Postcommitted"
	case VSS_SS_COUNT:
		return "Count"
	default:
		panic(fmt.Sprintf("unknown state: %v", uint(s)))
	}
}
