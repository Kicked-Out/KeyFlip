package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	macos "github.com/Kicked-Out/KeyFlip/platforms/macos"
)
// This file is entry point for KeyFlip application. Decides whether to run as daemon or CLI based on arguments.
func main() {
	// If there are any command-line arguments, run as CLI, if not, run as daemon
	if len(os.Args) > 1 {
		runCLI(os.Args[1:])
		return
	}
	runDaemon()
}

func runDaemon() {
	fmt.Println("KeyFlip daemon starting...")
	// Load config and start hotkey listener
	macos.StartHotkeyListener(func() {
		cfg, err := macos.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			return
		}
		macos.ProcessWithConfig(cfg)
	})
	// Keep the daemon running
	select {}
}

// =======================
// CLI
// =======================

func runCLI(args []string) {
	switch args[0] {

	// Show help message
	case "help", "--help", "-h":
		printHelp()
	// Install as launchd agent
	case "install":
		// Get the path to the current executable
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		// Resolve any symlinks to get the actual path
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}
		// Install launchd plist - autostart
		if err := macos.InstallLaunchd(exe); err != nil {
			fmt.Fprintln(os.Stderr, "Install failed:", err)
			os.Exit(1)
		}

		fmt.Println("KeyFlip installed as launchd agent.")
		fmt.Println("Grant Accessibility permissions:")
		fmt.Println("System Settings → Privacy & Security → Accessibility → KeyFlip")
	
	// Uninstall launchd agent
	case "uninstall":
		if err := macos.UninstallLaunchd(); err != nil {
			fmt.Fprintln(os.Stderr, "Uninstall failed:", err)
			os.Exit(1)
		}
		fmt.Println("KeyFlip uninstalled.")

	// Configure KeyFlip settings
	case "config":
		handleConfig(args[1:])

	default:
		fmt.Println("Unknown command:", args[0])
		printHelp()
	}
}

// Handle 'config' CLI commands
func handleConfig(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage:")
		fmt.Println("  keyflip config show")
		fmt.Println("  keyflip config set from=en to=ua")
		return
	}

	switch args[0] {
	// Show current configuration
	case "show":
		cfg, err := macos.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("from:   %s\n", cfg.From)
		fmt.Printf("to:     %s\n", cfg.To)
		fmt.Printf("hotkey: %s\n", cfg.Hotkey)

	// set configuration values
	case "set":
		// Load existing config
		cfg, err := macos.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Parse key=value pairs from arguments for example:from=en to=ua hotkey=cmd+shift+k
		for _, a := range args[1:] {
			// Split into key and value "from=en" → ["from","en"]
			kv := strings.SplitN(a, "=", 2)
			if len(kv) != 2 {
				fmt.Println("Invalid config:", a)
				continue
			}
			// Update config fields based on key
			switch kv[0] {
			case "from":
				cfg.From = kv[1]
			case "to":
				cfg.To = kv[1]
			case "hotkey":
				cfg.Hotkey = kv[1]
			default:
				fmt.Println("Unknown key:", kv[0])
			}
		}
		// Save updated config
		if err := macos.SaveConfig(cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Config saved.")
	}
}
// Print help message for CLI usage
func printHelp() {
	fmt.Println(`
KeyFlip — keyboard layout fixer

USAGE:
  keyflip              Run daemon
  keyflip install      Install autostart
  keyflip uninstall    Remove autostart
  keyflip config       Configure behavior
`)
}
