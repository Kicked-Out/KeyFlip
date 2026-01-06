package main

import (
	"sync"
	"time"

	"github.com/Kicked-Out/KeyFlip/core"
	macos "github.com/Kicked-Out/KeyFlip/platforms/macos/macos"
)

var (
	enToUa     = true
	stateMutex sync.Mutex
	processing sync.Mutex
)


func waitForClipboardChange(original string, timeout time.Duration) (string, bool) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		text, err := macos.ReadClipboard()
		if err == nil && text != "" && text != original {
			return text, true
		}
		time.Sleep(30 * time.Millisecond)
	}

	return "", false
}

func main() {
	macos.StartHotkeyListener(process)
	select {}
}

func process() {
	
	if !processing.TryLock() {
		return
	}
	defer processing.Unlock()

	
	stateMutex.Lock()
	currentEnToUa := enToUa
	stateMutex.Unlock()

	var mapping map[rune]rune
	if currentEnToUa {
		mapping = core.EnToUa
	} else {
		mapping = core.UaToEn
	}

	
	originalClipboard, err := macos.ReadClipboard()
	if err != nil {
		originalClipboard = ""
	}

	
	time.Sleep(60 * time.Millisecond)

	
	macos.CmdC()


	text, ok := waitForClipboardChange(originalClipboard, 350*time.Millisecond)
	if !ok {
		return
	}

	
	out := core.Transform(text, mapping)


	if err := macos.WriteClipboard(out); err != nil {
		_ = macos.WriteClipboard(originalClipboard)
		return
	}

	
	time.Sleep(120 * time.Millisecond)
	macos.CmdV()

	
	time.Sleep(80 * time.Millisecond)
	_ = macos.WriteClipboard(originalClipboard)

	
	stateMutex.Lock()
	enToUa = !enToUa
	stateMutex.Unlock()
}
