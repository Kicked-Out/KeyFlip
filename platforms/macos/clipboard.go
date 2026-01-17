package macos

import (
	"bytes"
	"os"
	"os/exec"
)

// C mechanism to read and write clipboard using pbcopy and pbpaste commands

func ReadClipboard() (string, error) {
	cmd := exec.Command("pbpaste")

	cmd.Env = append(os.Environ(),
		"LANG=en_US.UTF-8",
		"LC_ALL=en_US.UTF-8",
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}


// WriteClipboard writes the given text to the system clipboard
func WriteClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = bytes.NewBufferString(text)

	cmd.Env = append(os.Environ(),
		"LANG=en_US.UTF-8",
		"LC_ALL=en_US.UTF-8",
	)

	return cmd.Run()
}



// Works only on macOS - pbcopy/pbpaste commands