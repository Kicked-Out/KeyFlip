//go:build darwin

package macos

import (
	"bytes"
	"os"
	"os/exec"
)

// C mechanism to read and write clipboard using pbcopy and pbpaste commands

func ReadClipboard() (string, error) {
	// Execute pbpaste command to read clipboard content
	cmd := exec.Command("pbpaste")

	// Set environment variables to ensure proper encoding
	cmd.Env = append(os.Environ(),
		"LANG=en_US.UTF-8",
		"LC_ALL=en_US.UTF-8",
	)
	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}


// WriteClipboard writes the given text to the system clipboard
func WriteClipboard(text string) error {
	// Execute pbcopy command to write to clipboard
	cmd := exec.Command("pbcopy")
	cmd.Stdin = bytes.NewBufferString(text)
	// Set environment variables to ensure proper encoding
	cmd.Env = append(os.Environ(),
		"LANG=en_US.UTF-8",
		"LC_ALL=en_US.UTF-8",
	)

	return cmd.Run()
}



// Works only on macOS - pbcopy/pbpaste commands