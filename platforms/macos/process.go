package macos

import (
	"sync"
	"time"

	"github.com/Kicked-Out/KeyFlip/core"
)

var (
	enToUa     = true
	stateMutex sync.Mutex
	processing sync.Mutex
)

func waitForClipboardChange(original string, timeout time.Duration) (string, bool) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		text, err := ReadClipboard()
		if err == nil && text != "" && text != original {
			return text, true
		}
		time.Sleep(30 * time.Millisecond)
	}
	return "", false
}

// Process — головна логіка daemon
func Process() {
	// не допускаємо паралельного виконання
	if !processing.TryLock() {
		return
	}
	defer processing.Unlock()

	// визначаємо напрямок
	stateMutex.Lock()
	currentEnToUa := enToUa
	stateMutex.Unlock()

	var mapping map[rune]rune
	if currentEnToUa {
		mapping = core.EnToUa
	} else {
		mapping = core.UaToEn
	}

	// зберігаємо clipboard
	originalClipboard, err := ReadClipboard()
	if err != nil {
		originalClipboard = ""
	}

	time.Sleep(60 * time.Millisecond)

	// копіюємо виділення
	CmdC()

	// чекаємо реальної зміни clipboard
	text, ok := waitForClipboardChange(originalClipboard, 350*time.Millisecond)
	if !ok {
		return
	}

	// трансформація
	out := core.Transform(text, mapping)

	if err := WriteClipboard(out); err != nil {
		_ = WriteClipboard(originalClipboard)
		return
	}

	time.Sleep(120 * time.Millisecond)
	CmdV()

	time.Sleep(80 * time.Millisecond)
	_ = WriteClipboard(originalClipboard)

	// toggle напрямку
	stateMutex.Lock()
	enToUa = !enToUa
	stateMutex.Unlock()
}
