//go:build darwin

package macos

import (
	"os"
	"path/filepath"
)

// Returns the path to the layouts file
func LayoutsPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), "layouts.json"), nil
}
