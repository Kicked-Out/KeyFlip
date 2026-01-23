//go:build windows

package windows

import (
	"fmt"
	"time"
	"unsafe"
)

var (
	sendInput = user32.NewProc("SendInput")
)

const (
	INPUT_KEYBOARD = 1
	KEYEVENTF_KEYUP = 0x0002

	VK_CONTROL = 0x11
	VK_SHIFT = 0x10
	VK_C = 0x43
	VK_V = 0x56
)

type input struct {
	Type uint32
	_ uint32
	Ki KEYBDINPUT
	_ uint64
}

type KEYBDINPUT struct {
	WVk uint16
	WScan uint16
	DwFlags uint32
	Time uint32
	DwExtraInfo uintptr
}

func keyVK(vk uint16, up bool) {
	var flags uint32

	if up {
		flags |= KEYEVENTF_KEYUP
	}

	in := input{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			WVk: vk,
			WScan: 0,
			DwFlags: flags,
			Time: 0,
			DwExtraInfo: 0,
		},
	}

	ret, _, err := sendInput.Call(
		1,
		uintptr(unsafe.Pointer(&in)),
		unsafe.Sizeof(in),
	)

	if ret == 0 {
		fmt.Println("SendInput failed:", err)
	}
}

func CtrlC() {
	keyVK(VK_SHIFT, true)
	keyVK(VK_CONTROL, true)
	time.Sleep(10 * time.Millisecond)

	keyVK(VK_CONTROL, false)
	keyVK(VK_C, false)
	time.Sleep(20 * time.Millisecond)
	keyVK(VK_C, true)
	keyVK(VK_CONTROL, true)
}

func CtrlV() {
	keyVK(VK_SHIFT, true)
	keyVK(VK_CONTROL, true)
	time.Sleep(10 * time.Millisecond)

	keyVK(VK_CONTROL, false)
	keyVK(VK_V, false)
	time.Sleep(20 * time.Millisecond)
	keyVK(VK_V, true)
	keyVK(VK_CONTROL, true)
}