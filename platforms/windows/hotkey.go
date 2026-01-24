//go:build windows

package windows

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

var (
	// user32 = syscall.NewLazyDLL("user32.dll")
	registerHotKey = user32.NewProc("RegisterHotKey")
	unregisterHotKey = user32.NewProc("UnregisterHotKey")
	getMessageW = user32.NewProc("GetMessageW")
)


const (
	MOD_CONTROL = 0x0002
	MOD_SHIFT = 0x0004
	WM_HOTKEY = 0x0312
	VK_K = 0x4B
)

type MSG struct {
	Hwnd uintptr
	Message uint32
	WParam uintptr
	LParam uintptr
	Time uint32
	Pt struct {
		X, Y int32
	}
}

func StartHotkeyListener(handler func()) error {
	errChan := make(chan error, 1)

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		r, _, err := registerHotKey.Call(
			0,
			1,
			MOD_CONTROL|MOD_SHIFT,
			VK_K,
		)

		if r == 0 {
			errChan <- fmt.Errorf("RegisterHotKey failed: %v", err)
			return
		}

		errChan <- nil

		defer unregisterHotKey.Call(0, 1)

		var msg MSG
		var lastFire time.Time

		for {
			ret, _, _ := getMessageW.Call(
				uintptr(unsafe.Pointer(&msg)),
				0,
				0,
				0,
			)

			if ret == 0 {
				return
			}

			if msg.Message == WM_HOTKEY {
				if time.Since(lastFire) < 500*time.Millisecond {
					continue
				}

				lastFire = time.Now()

				go handler()
			}
		}
	}()

	return <- errChan
}