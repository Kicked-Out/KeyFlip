//go:build windows

package windows

import (
	"syscall"
	"unsafe"
)

var (
	user32 =  syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	
	openClipboard = user32.NewProc("OpenClipboard")
	closeClipboard = user32.NewProc("CloseClipboard")
	getClipboardData = user32.NewProc("GetClipboardData")
	setClipboardData = user32.NewProc("SetClipboardData")
	emptyClipboard = user32.NewProc("EmptyClipboard")

	globalLock = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	globalSize = kernel32.NewProc("GlobalSize")
	globalAlloc = kernel32.NewProc("GlobalAlloc")
)

const CF_UNICODETEXT = 13
const GMEM_MOVEABLE = 0x0002

func ReadClipboard() (string, error) {
	var r, _, err = openClipboard.Call(0)

	if r == 0 {
		return "", err
	}

	defer closeClipboard.Call()

	var h, _, _ = getClipboardData.Call(CF_UNICODETEXT);

	if h == 0 {
		return "", nil
	}

	var ptr, _, _ = globalLock.Call(h)

	if ptr == 0 {
		return "", nil
	}

	defer globalUnlock.Call(h)

	var size, _, _ = globalSize.Call(h)

	if (size == 0) {
		return "", nil
	}

	var buf = unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), size/2)

	return syscall.UTF16ToString(buf), nil
}

func WriteClipboard(text string) error {
	r, _, err := openClipboard.Call(0)

	if r == 0 {
		return err
	}

	defer closeClipboard.Call()

	emptyClipboard.Call()

	data, err := syscall.UTF16FromString(text)

	if err != nil {
		return err
	}

	var size = uintptr(len(data) * 2)

	h, _, err := globalAlloc.Call(GMEM_MOVEABLE, size)

	if h == 0 {
		return err
	}

	var ptr, _, _ = globalLock.Call(h)

	if ptr == 0 {
		return err
	}

	var buf = unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), len(data))

	copy(buf, data)
	globalUnlock.Call(h)

	setClipboardData.Call(CF_UNICODETEXT, h)

	return nil
}