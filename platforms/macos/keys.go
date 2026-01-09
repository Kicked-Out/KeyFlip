package macos

// C mechanism to simulate Cmd+C and Cmd+V key presses

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>

static void keyPress(int keyCode, int useCmd) {
    CGEventRef e1 = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)keyCode, true);
    if (useCmd) {
        CGEventSetFlags(e1, kCGEventFlagMaskCommand);
    }
    CGEventPost(kCGHIDEventTap, e1);
    CFRelease(e1);

    CGEventRef e2 = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)keyCode, false);
    if (useCmd) {
        CGEventSetFlags(e2, kCGEventFlagMaskCommand);
    }
    CGEventPost(kCGHIDEventTap, e2);
    CFRelease(e2);
}
*/
import "C"

// C = 8, V = 9
func CmdC() { C.keyPress(8, 1) }
func CmdV() { C.keyPress(9, 1) }
