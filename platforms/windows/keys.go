package windows

import (
	"unsafe"
)

var (
	sendInput = user32.NewProc("SendInput")
)

const (
	INPUT_KEYBOARD = 1
	KEYEVENTF_KEYUP = 0x2
	VK_CONTROL = 0x11
	VK_C = 0x43
	VK_V = 0x56
)

type input struct {
	typ uint32
	ki keybdinput
}

type keybdinput struct {
	wVk uint16
	wScan uint16
	dwFlags uint32
	time uint32
	dwExtraInfo uintptr
}

func key(vk uint16, up bool) {
	var flags uint32
	
	if up {
		flags = KEYEVENTF_KEYUP
	}

	in := input{INPUT_KEYBOARD, keybdinput{wVk: vk, dwFlags: flags}}

	sendInput.Call(1, uintptr(unsafe.Pointer(&in)), unsafe.Sizeof(in))
}

func CtrlC() {
	key(VK_CONTROL, false)
	key(VK_C, false)
	key(VK_C, true)
	key(VK_CONTROL, true)
}

func CtrlV() {
	key(VK_CONTROL, false)
	key(VK_V, false)
	key(VK_V, true)
	key(VK_CONTROL, true)
}