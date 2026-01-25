//go:build darwin

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
	// blocks if another processing is ongoing
	if !processing.TryLock() {
		return
	}
	defer processing.Unlock()
	// Load layouts
	layoutsPath, err := LayoutsPath()
	if err != nil {
		return
	}
	// Load layouts from file
	layouts, err := core.LoadLayouts(layoutsPath)
	if err != nil {
		return
	}
	// Get from and to layouts
	fromLayout, ok1 := layouts[cfg.From]
	toLayout, ok2 := layouts[cfg.To]
	if !ok1 || !ok2 {
		return
	}

	// build forward and reverse mappings
	forward := make(map[rune]rune) // from -> to
	reverse := make(map[rune]rune) // to -> from

	// Create bidirectional mapping between fromLayout and toLayout
	for key, fromChar := range fromLayout {
		toChar, ok := toLayout[key]
		if !ok {
			continue
		}
		forward[fromChar] = toChar
		reverse[toChar] = fromChar
	}

	if len(forward) == 0 {
		return
	}

	// Read original clipboard content
	original, _ := ReadClipboard()

	time.Sleep(50 * time.Millisecond)
	CmdC()

	// Wait for clipboard to change
	text, ok := waitForClipboardChange(original, 300*time.Millisecond)
	if !ok || text == original {
		return
	}

	// Detect direction and transform text
	mapping := detectDirection(text, forward, reverse)
	out := core.Transform(text, mapping)

	// Write transformed text to clipboard
	if err := WriteClipboard(out); err != nil {
		_ = WriteClipboard(original)
		return
	}

	time.Sleep(80 * time.Millisecond)
	CmdV()

	time.Sleep(50 * time.Millisecond)
	_ = WriteClipboard(original)
}

//detects the most likely direction of text transformation based on character presence in layouts
func detectDirection(
	text string,
	forward map[rune]rune,
	reverse map[rune]rune,
) map[rune]rune {

	forwardScore := 0
	reverseScore := 0

	for _, r := range text {
		_, inF := forward[r]
		_, inR := reverse[r]

		switch {
		case inF && !inR:
			forwardScore++
		case inR && !inF:
			reverseScore++
		}
	}

	if reverseScore > forwardScore {
		return reverse
	}
	return forward
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
