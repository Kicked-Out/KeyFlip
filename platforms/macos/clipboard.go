package macos

import (
	"bytes"
	"os/exec"
)

// C mechanism to read and write clipboard using pbcopy and pbpaste commands

func ReadClipboard() (string, error) {
	// Use pbpaste to read clipboard content
	cmd := exec.Command("pbpaste")
	// Capture the output in a buffer
	var out bytes.Buffer
	// write output to buffer
	cmd.Stdout = &out
	// Run the command and return output as string
	err := cmd.Run()
	return out.String(), err
}

// WriteClipboard writes the given text to the system clipboard
func WriteClipboard(text string) error {
	// Use pbcopy to write text to clipboard
	cmd := exec.Command("pbcopy")
	// makes a buffer-rider from the string text connects it as stdin to pbcopy
	cmd.Stdin = bytes.NewBufferString(text)
	// Run the command
	return cmd.Run()
}


// Works only on macOS - pbcopy/pbpaste commands