package macos

import (
	"bytes"
	"time"

	"github.com/Kicked-Out/KeyFlip/core"
)

// викликається з daemon
func ProcessWithConfig(cfg Config) {
	originalClipboard, err := ReadClipboard()
	if err != nil {
		originalClipboard = ""
	}

	time.Sleep(60 * time.Millisecond)

	CmdC()

	text, ok := waitForClipboardChange(originalClipboard, 350*time.Millisecond)
	if !ok {
		return
	}

	mapping := detectMapping(text)
	if mapping == nil {
		return
	}

	out := core.Transform(text, mapping)

	if err := WriteClipboard(out); err != nil {
		_ = WriteClipboard(originalClipboard)
		return
	}

	time.Sleep(120 * time.Millisecond)
	CmdV()

	time.Sleep(80 * time.Millisecond)
	_ = WriteClipboard(originalClipboard)
}

func detectMapping(text string) map[rune]rune {
	for _, r := range text {

		// латиниця → укр
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return core.EnToUa
		}

		// кирилиця → англ
		if (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я') || r == 'і' || r == 'І' {
			return core.UaToEn
		}
	}
	return nil
}

func waitForClipboardChange(original string, timeout time.Duration) (string, bool) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		text, err := ReadClipboard()
		if err == nil && text != "" && !bytes.Equal([]byte(text), []byte(original)) {
			return text, true
		}
		time.Sleep(30 * time.Millisecond)
	}
	return "", false
}
