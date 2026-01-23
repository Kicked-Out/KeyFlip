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
	if !processing.TryLock() {
		return
	}
	defer processing.Unlock()

	layoutsPath, err := LayoutsPath()
	if err != nil {
		return
	}

	layouts, err := core.LoadLayouts(layoutsPath)
	if err != nil {
		return
	}

	fromLayout, ok1 := layouts[cfg.From]
	toLayout, ok2 := layouts[cfg.To]
	if !ok1 || !ok2 {
		return
	}

	// build forward and reverse mappings
	forward := make(map[rune]rune) // from -> to
	reverse := make(map[rune]rune) // to -> from

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

	original, _ := ReadClipboard()

	time.Sleep(50 * time.Millisecond)
	CmdC()

	text, ok := waitForClipboardChange(original, 300*time.Millisecond)
	if !ok || text == original {
		return
	}

	mapping := detectDirection(text, forward, reverse)
	out := core.Transform(text, mapping)

	if err := WriteClipboard(out); err != nil {
		_ = WriteClipboard(original)
		return
	}

	time.Sleep(80 * time.Millisecond)
	CmdV()

	time.Sleep(50 * time.Millisecond)
	_ = WriteClipboard(original)
}

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
