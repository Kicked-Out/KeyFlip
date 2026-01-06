package main

import (
	"bytes"
	"sync"
	"time"

	"github.com/Kicked-Out/KeyFlip/core"
	macos "github.com/Kicked-Out/KeyFlip/platforms/macos/macos"
)

var (
	enToUa     = true
	stateMutex sync.Mutex
)

func main() {
	macos.StartHotkeyListener(process)
	select {}
}

func process() {
	
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
	time.Sleep(200 * time.Millisecond)

	text, err := macos.ReadClipboard()
	if err != nil {
		_ = macos.WriteClipboard(originalClipboard)
		return
	}

	
	if text == "" || bytes.Equal([]byte(text), []byte(originalClipboard)) {
		return
	}

	println("COPIED:", text)

	out := core.Transform(text, mapping)
	println("TRANSFORMED:", out)

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
