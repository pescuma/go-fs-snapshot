//go:build windows

package internal_fs_snapshot_windows

import (
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/pkg/errors"
)

// asyncCallFunc is the callback type for callAndWait.
type asyncCallFunc func() (*ivssAsync, error)

// callAndWait calls an async functions and waits for it to either
// finish or timeout.
func callAndWait(function asyncCallFunc, timeout time.Duration) error {
	if timeout <= 0 {
		return errors.New("Timeout occured")
	}

	async, err := function()
	if err != nil {
		return err
	}

	err = async.WaitUntilAsyncFinished(uint32(timeout.Milliseconds()))
	async.Release()
	return err
}

// uiid_ivssAsync defines to GUID of ivssAsync.
//
//goland:noinspection GoSnakeCaseUsage
var uiid_ivssAsync = ole.NewGUID("{507C37B4-CF5B-4e95-B0AF-14EB9767467E}")

// ivssAsync VSS api interface.
type ivssAsync struct {
	ole.IUnknown
}

// ivssAsyncVTable is the vtable for ivssAsync.
type ivssAsyncVTable struct {
	ole.IUnknownVtbl
	cancel      uintptr
	wait        uintptr
	queryStatus uintptr
}

// Constants for IVSSAsync api.
//
//goland:noinspection GoSnakeCaseUsage
const (
	VSS_S_ASYNC_PENDING   = 0x00042309
	VSS_S_ASYNC_FINISHED  = 0x0004230A
	VSS_S_ASYNC_CANCELLED = 0x0004230B
)

// getVTable returns the vtable for ivssAsync.
func (vssAsync *ivssAsync) getVTable() *ivssAsyncVTable {
	return (*ivssAsyncVTable)(unsafe.Pointer(vssAsync.RawVTable))
}

// Cancel calls the equivalent VSS api.
func (vssAsync *ivssAsync) Cancel() error {
	return syscallN("Cancel()", vssAsync.getVTable().cancel,
		uintptr(unsafe.Pointer(vssAsync)))
}

// Wait calls the equivalent VSS api.
func (vssAsync *ivssAsync) Wait(millis uint32) error {
	return syscallN("Wait()", vssAsync.getVTable().wait,
		uintptr(unsafe.Pointer(vssAsync)), uintptr(millis))
}

// QueryStatus calls the equivalent VSS api.
func (vssAsync *ivssAsync) QueryStatus() (uint32, error) {
	var state uint32 = 0
	err := syscallN("QueryStatus()", vssAsync.getVTable().queryStatus,
		uintptr(unsafe.Pointer(vssAsync)), uintptr(unsafe.Pointer(&state)), 0)
	return state, err
}

// WaitUntilAsyncFinished waits until either the async call is finshed or
// the given timeout is reached.
func (vssAsync *ivssAsync) WaitUntilAsyncFinished(millis uint32) error {
	err := vssAsync.Wait(millis)
	if err != nil {
		_ = vssAsync.Cancel()
		return err
	}

	state, err := vssAsync.QueryStatus()
	if err != nil {
		_ = vssAsync.Cancel()
		return err
	}

	if state == VSS_S_ASYNC_CANCELLED {
		return errors.New("async operation cancelled")
	}

	if state == VSS_S_ASYNC_PENDING {
		_ = vssAsync.Cancel()
		return errors.New("async operation pending")
	}

	if state != VSS_S_ASYNC_FINISHED {
		return errors.Errorf("async operation failed with state 0x%x", state)
	}

	return nil
}
