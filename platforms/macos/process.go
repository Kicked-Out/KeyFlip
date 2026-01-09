package macos
// Main logic for processing clipboard text based on keyboard layout conversion
import (
	"bytes"
	"sync"
	"time"

	"github.com/Kicked-Out/KeyFlip/core"
)
// Mutex to prevent concurrent processing - only one process at a time
var (
	processing sync.Mutex
)

//processes clipboard text based on provided configuration
func ProcessWithConfig(cfg Config) {
	// Try to acquire lock, if already processing, return immediately
	if !processing.TryLock() {
		return
	}
	// Release the lock when function exits
	defer processing.Unlock()

	// Determine mapping based on config
	var mapping map[rune]rune

	switch {
	case cfg.From == "en" && cfg.To == "ua":
		mapping = core.EnToUa
	case cfg.From == "ua" && cfg.To == "en":
		mapping = core.UaToEn
	default:
		return
	}

	// Read original clipboard content, If clipboard is empty, nothing to process
	original, _ := ReadClipboard()

	time.Sleep(50 * time.Millisecond)
	// Simulate Cmd+C(Ctrl+C)
	CmdC()

	// Wait for clipboard to change with a timeout
	text, ok := waitForClipboardChange(original, 300*time.Millisecond)
	if !ok {
		return
	}

	// Additional check: if text is same as original, nothing to process - potentially useless
	if bytes.Equal([]byte(text), []byte(original)) {
		return
	}

	//  Transform the text using the determined mapping
	out := core.Transform(text, mapping)

	// Write transformed text back to clipboard and restore original on failure
	if err := WriteClipboard(out); err != nil {
		_ = WriteClipboard(original)
		return
	}

	time.Sleep(80 * time.Millisecond)
	// Simulate Cmd+V(Ctrl+V)
	CmdV()

	time.Sleep(50 * time.Millisecond)
	// Waits and restore original clipboard content
	_ = WriteClipboard(original)
}

// detectMapping analyzes the text to determine if it's in English or Ukrainian layout (but it's not used now - maybe delete it later)
func detectMapping(text string) map[rune]rune {
	for _, r := range text {

		
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return core.EnToUa
		}

		
		if (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я') || r == 'і' || r == 'І' {
			return core.UaToEn
		}
	}
	return nil
}

//waits for clipboard content to change from the original within the specified timeout duration
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
