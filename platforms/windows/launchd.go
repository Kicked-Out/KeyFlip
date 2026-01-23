//go:build windows

package windows

import "golang.org/x/sys/windows/registry"

func InstallAutostart(path string) error {
	k, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)

	if err != nil {
		return err
	}

	defer k.Close()

	return k.SetStringValue("KeyFlip", path)
}

func AutostartInstalled() bool {
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE,
	)

	if err != nil {
		return false
	}

	defer k.Close()

	_, _, err = k.GetStringValue("KeyFlip")

	return err == nil
}
func UninstallAutostart() error {
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)

	if err != nil {
		return err
	}

	defer k.Close()

	return k.DeleteValue("KeyFlip")
}