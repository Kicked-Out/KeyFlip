package macos
// mechanism which decides that KeyFlip should run always
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Label for the launchd agent
const launchdLabel = "com.keyflip.daemon"

// Returns the path to the launchd plist file (User home directory)
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

// Installs the launchd plist to autostart the application
func InstallLaunchd(binaryPath string) error {
	// Get the path to the plist file
	path, err := launchdPlistPath()
	if err != nil {
		return err
	}
	// Ensure the directory exists if not, create it
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Create the plist content
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

	// Write the plist file
	if err := os.WriteFile(path, []byte(plist), 0o644); err != nil {
		return err
	}

	// Load the plist with launchctl
	_ = exec.Command("launchctl", "unload", path).Run()

	cmd := exec.Command("launchctl", "load", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Uninstalls the launchd plist and unloads the agent
func UninstallLaunchd() error {
	path, err := launchdPlistPath()
	if err != nil {
		return err
	}

	_ = exec.Command("launchctl", "unload", path).Run()
	return os.Remove(path)
}

// Checks if the launchd plist is installed
func LaunchdInstalled() bool {
	path, err := launchdPlistPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}
