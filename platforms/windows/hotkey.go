package windows

import (
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	registerHotKey = user32.NewProc("RegisterHotKey")
	unregisterHotKey = user32.NewProc("UnregisterHotKey")
	getMessageW = user32.NewProc("GetMessageW")
)


const (
	MOD_ALT = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT = 0x0004
	MOD_WIN = 0x0008

	WM_HOTKEY = 0x0312
	VK_K = 0x48
)

type msg struct {
	hwnd uintptr
	message uint32
	wParam uintptr
	lParam uintptr
	time uint32
	pt struct {
		x, y int32
	}
}

var hotkeyHandler atomic.Pointer[func()]
var lastFire int64

func StartHotkeyListener(handler func()) {
	hotkeyHandler.Store(&handler)
	
	go messageLoop()
}

func messageLoop() {
	ret, _, err := registerHotKey.Call(
		0,
		1,
		MOD_CONTROL|MOD_SHIFT,
		VK_K,
	)

	if ret == 0 {
		panic("RegisterHotKey failed: " + err.Error())
	}

	defer unregisterHotKey.Call(0, 1)

	var m msg
	
	for {
		r, _, _ := getMessageW.Call(
			uintptr(unsafe.Pointer(&m)),
			0,
			0,
			0,
		)

		if r == 0 {
			return
		}

		if m.message == WM_HOTKEY {
			onHotkey()
		}
	}
}

func onHotkey() {
	now := time.Now().UnixNano()
	prev := atomic.LoadInt64(&lastFire)

	if (now - prev < int64(250 * time.Millisecond)) {
		return
	}

	atomic.StoreInt64(&lastFire, now)

	if h := hotkeyHandler.Load(); h != nil {
		go (*h)()
	}
}