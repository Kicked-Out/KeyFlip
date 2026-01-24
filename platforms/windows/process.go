//go:build windows

package windows

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/Kicked-Out/KeyFlip/core"
)

var (
	processing sync.Mutex
)

func ProcessWithConfig(cfg Config) {
	if !processing.TryLock() {
		fmt.Println("Already processing, skipping...")
		return
	}

	defer func() {
		processing.Unlock()
		fmt.Println("Processing finished, mutex unlocked.")
	}()

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

	forward := make(map[rune]rune)
	reverse := make(map[rune]rune)

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

	original, err := ReadClipboard()

	if err != nil {
		return
	}

	time.Sleep(500 * time.Millisecond)
	CtrlC()

	fmt.Println("Waiting for clipboard...")

	text, ok := waitForClipboardChange(original, 500*time.Millisecond)

	if !ok || text == original {
		fmt.Println("Clipboard change timeout or failed.")
		return
	}

	mapping := detectDirection(text, forward, reverse)
	out := core.Transform(text, mapping)

	if err := WriteClipboard(out); err != nil {
		_ = WriteClipboard(original)

		return
	}

	time.Sleep(80 * time.Millisecond)
	CtrlV()

	time.Sleep(50 * time.Millisecond)
	_ = WriteClipboard(original)
}

func detectDirection(
	text string,
	forward map[rune]rune,
	reverse map[rune]rune,
)map[rune]rune {
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
