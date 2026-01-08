package macos

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const launchdLabel = "com.keyflip.daemon"

func launchdPlistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(
		home,
		"Library",
		"LaunchAgents",
		launchdLabel+".plist",
	), nil
}

func InstallLaunchd(binaryPath string) error {
	path, err := launchdPlistPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
"http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>%s</string>

	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>

	<key>RunAtLoad</key>
	<true/>

	<key>KeepAlive</key>
	<true/>
</dict>
</plist>
`, launchdLabel, binaryPath)

	if err := os.WriteFile(path, []byte(plist), 0o644); err != nil {
		return err
	}

	// unload якщо був
	_ = exec.Command("launchctl", "unload", path).Run()

	// load
	cmd := exec.Command("launchctl", "load", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func UninstallLaunchd() error {
	path, err := launchdPlistPath()
	if err != nil {
		return err
	}

	_ = exec.Command("launchctl", "unload", path).Run()
	return os.Remove(path)
}

func LaunchdInstalled() bool {
	path, err := launchdPlistPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}
