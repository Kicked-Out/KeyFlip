package main

import (
	"time"

	"github.com/Kicked-Out/KeyFlip/core"
	macos "github.com/Kicked-Out/KeyFlip/platforms/macos/macos"
)

func main() {
	macos.StartHotkeyListener(process)

	
	select {}
}

func process() {
	originalClipboard, _ := macos.ReadClipboard()

	time.Sleep(60 * time.Millisecond)

	macos.CmdC()
	time.Sleep(200 * time.Millisecond)

	text, err := macos.ReadClipboard()
	if err != nil {
		println("READ ERR:", err.Error())
		_ = macos.WriteClipboard(originalClipboard)
		return
	}

	println("COPIED:", text)

	if text == "" {
		println("EMPTY SELECTION (nothing copied)")
		_ = macos.WriteClipboard(originalClipboard)
		return
	}

	out := core.Transform(text, core.EnToUa)
	println("TRANSFORMED:", out)

	_ = macos.WriteClipboard(out)
	time.Sleep(120 * time.Millisecond)
	macos.CmdV()

	time.Sleep(80 * time.Millisecond)
	_ = macos.WriteClipboard(originalClipboard)
}
