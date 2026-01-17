package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kicked-Out/KeyFlip/core"
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
	
	case "lang":
		showLanguages()

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

	case "show":
		cfg, err := macos.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("from:   %s\n", cfg.From)
		fmt.Printf("to:     %s\n", cfg.To)
		fmt.Printf("hotkey: %s\n", cfg.Hotkey)

	case "set":
		cfg, err := macos.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// load layouts.json
		layoutsPath, err := macos.LayoutsPath()
		if err != nil {
			fmt.Println("Failed to find layouts.json:", err)
			return
		}

		layouts, err := core.LoadLayouts(layoutsPath)
		if err != nil {
			fmt.Println("Failed to load layouts:", err)
			return
		}

		isValidLang := func(lang string) bool {
			_, ok := layouts[lang]
			return ok
		}

		for _, a := range args[1:] {
			kv := strings.SplitN(a, "=", 2)
			if len(kv) != 2 {
				fmt.Println("Invalid config:", a)
				continue
			}

			switch kv[0] {
			case "from":
				if !isValidLang(kv[1]) {
					fmt.Println("Unknown language:", kv[1])
					return
				}
				cfg.From = kv[1]

			case "to":
				if !isValidLang(kv[1]) {
					fmt.Println("Unknown language:", kv[1])
					return
				}
				cfg.To = kv[1]

			case "hotkey":
				cfg.Hotkey = kv[1]

			default:
				fmt.Println("Unknown key:", kv[0])
			}
		}

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

func showLanguages() {
	path, err := macos.LayoutsPath()
	if err != nil {
		fmt.Println("Failed to locate layouts.json:", err)
		return
	}

	layouts, err := core.LoadLayouts(path)
	if err != nil {
		fmt.Println("Failed to load layouts:", err)
		return
	}

	fmt.Println("Available languages:")
	for lang := range layouts {
		fmt.Println(" -", lang)
	}
}
