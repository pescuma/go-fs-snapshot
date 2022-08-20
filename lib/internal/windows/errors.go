//go:build windows

package internal_fs_snapshot_windows

import (
	"fmt"
)

// VssError encapsulates errors returned from calling VSS api.
type VssError struct {
	Text    string
	HResult HRESULT
}

// newVssError creates a new VSS api error.
func newVssError(hresult HRESULT, text string) *VssError {
	return &VssError{Text: text, HResult: hresult}
}

// newVssErrorF creates a new VSS api error.
func newVssErrorF(hresult HRESULT, format string, a ...interface{}) *VssError {
	return &VssError{Text: fmt.Sprintf(format, a), HResult: hresult}
}

// Error implements the error interface.
func (e *VssError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Text, e.HResult.Str())
}
