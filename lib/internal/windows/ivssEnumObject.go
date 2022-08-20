package internal_fs_snapshot_windows

import (
	"unsafe"

	"github.com/go-ole/go-ole"
)

// IVssEnumObject VSS api interface.
type IVssEnumObject struct {
	ole.IUnknown
}

// IVssEnumObjectVTable is the vtable for IVssEnumObject.
type IVssEnumObjectVTable struct {
	ole.IUnknownVtbl
	next  uintptr
	skip  uintptr
	reset uintptr
	clone uintptr
}

// getVTable returns the vtable for IVssEnumObject.
func (e *IVssEnumObject) getVTable() *IVssEnumObjectVTable {
	return (*IVssEnumObjectVTable)(unsafe.Pointer(e.RawVTable))
}

// Next calls the equivalent VSS api.
func (e *IVssEnumObject) Next(count uint, props unsafe.Pointer) (uint, error) {
	var fetched uint32

	err := syscallNF("Next()", e.getVTable().next,
		uintptr(unsafe.Pointer(e)), uintptr(count), uintptr(props),
		uintptr(unsafe.Pointer(&fetched)))

	return uint(fetched), err
}

func (e *IVssEnumObject) Close() {
	if e == nil {
		return
	}

	e.Release()
}
