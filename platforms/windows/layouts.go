//go:build windows

package windows

import (
	"os"
	"path/filepath"
)

func LayoutsPath() (string, error) {
	exe, err := os.Executable()

	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(exe), "layouts.json"), nil
}