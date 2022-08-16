//go:build windows

package fs_snapshot_windows

import "fmt"

// vssError encapsulates errors returned from calling VSS api.
type vssError struct {
	text    string
	hresult HRESULT
}

// newVssError creates a new VSS api error.
func newVssError(text string, hresult HRESULT) error {
	return &vssError{text: text, hresult: hresult}
}

// newVssError creates a new VSS api error.
func newVssErrorIfResultNotOK(text string, hresult HRESULT) error {
	if hresult != S_OK {
		return newVssError(text, hresult)
	}
	return nil
}

// Error implements the error interface.
func (e *vssError) Error() string {
	return fmt.Sprintf("VSS error: %s: %s (%#x)", e.text, e.hresult.Str(), e.hresult)
}

// vssError encapsulates errors retruned from calling VSS api.
type vssTextError struct {
	text string
}

// newVssTextError creates a new VSS api error.
func newVssTextError(text string) error {
	return &vssTextError{text: text}
}

// Error implements the error interface.
func (e *vssTextError) Error() string {
	return fmt.Sprintf("VSS error: %s", e.text)
}
