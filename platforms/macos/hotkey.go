//go:build darwin

package macos

// C mechanism to listen for global hotkey (Cmd+Shift+K)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <ApplicationServices/ApplicationServices.h>
#include <CoreFoundation/CoreFoundation.h>

extern void OnHotkey();

static CGEventRef keyflip_hotkeyCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    if (type != kCGEventKeyUp) return event;

    CGKeyCode keycode = (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
    CGEventFlags flags = CGEventGetFlags(event);

    // Cmd + Shift + K
    if (keycode == 40 &&
        (flags & kCGEventFlagMaskCommand) &&
        (flags & kCGEventFlagMaskShift)) {
        OnHotkey();
    }

    return event;
}

static void keyflip_startHotkeyListener(void) {
    CGEventMask mask = CGEventMaskBit(kCGEventKeyUp);
    CFMachPortRef tap = CGEventTapCreate(
        kCGSessionEventTap,
        kCGHeadInsertEventTap,
        0,
        mask,
        keyflip_hotkeyCallback,
        NULL
    );

    if (!tap) {
        fprintf(stderr, "KeyFlip: failed to create event tap. Check Accessibility permissions.\n");
        return;
    }


    CFRunLoopSourceRef source = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, tap, 0);
    CFRunLoopAddSource(CFRunLoopGetCurrent(), source, kCFRunLoopCommonModes);
    CGEventTapEnable(tap, true);
    CFRunLoopRun();
}
*/
import "C"

import (
	"sync/atomic"
	"time"
)

var hotkeyHandler atomic.Pointer[func()]
var lastFire int64

//export OnHotkey
func OnHotkey() {
	now := time.Now().UnixNano()
	prev := atomic.LoadInt64(&lastFire)

	// 250ms debounce
	if now-prev < int64(250*time.Millisecond) {
		return
	}
	atomic.StoreInt64(&lastFire, now)

	if h := hotkeyHandler.Load(); h != nil {
		go (*h)()
	}
}
// Starts the hotkey listener in a separate goroutine
func StartHotkeyListener(handler func()) {
	hotkeyHandler.Store(&handler)
	go C.keyflip_startHotkeyListener()
}
