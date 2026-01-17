package windows

import "golang.org/x/sys/windows/registry"

func InstallAutostart(path string) error {
	k, _, _ := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)

	defer k.Close()

	return k.SetStringValue("KeyFlip", path)
}

func UninstallAutostart() error {
	k, _ := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)

	defer k.Close()

	return k.DeleteValue("KeyFlip")
}