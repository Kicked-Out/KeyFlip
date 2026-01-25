//go:build darwin
// +build darwin

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
	// Restart launchd agent
	case "restart":
		if err := macos.UninstallLaunchd(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}

		if err := macos.InstallLaunchd(exe); err != nil {
			fmt.Fprintln(os.Stderr, "Restart failed:", err)
			os.Exit(1)
		}

		fmt.Println("KeyFlip restarted.")
	// Install as launchd agent
	case "install":
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}

		already := macos.LaunchdInstalled()

		if err := macos.InstallLaunchd(exe); err != nil {
			fmt.Fprintln(os.Stderr, "Install failed:", err)
			os.Exit(1)
		}

		if already {
			fmt.Println("KeyFlip restarted.")
		} else {
			fmt.Println("KeyFlip installed.")
			fmt.Println("Grant Accessibility permissions:")
			fmt.Println("System Settings → Privacy & Security → Accessibility → KeyFlip")
			fmt.Println("After granting permissions, run:")
			fmt.Println("  keyflip restart")
		}

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
	// Show available languages
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
		fmt.Println("  keyflip config set from=xx to=yy")
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

		layoutsPath, err := macos.LayoutsPath()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to find layouts.json:", err)
			os.Exit(1)
		}

		layouts, err := core.LoadLayouts(layoutsPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to load layouts:", err)
			os.Exit(1)
		}

		isValidLang := func(lang string) bool {
			_, ok := layouts[lang]
			return ok
		}

		hasError := false

		for _, a := range args[1:] {
			kv := strings.SplitN(a, "=", 2)
			if len(kv) != 2 {
				fmt.Fprintln(os.Stderr, "Invalid config:", a)
				hasError = true
				continue
			}

			switch kv[0] {
			case "from":
				if !isValidLang(kv[1]) {
					fmt.Fprintln(os.Stderr, "Unknown language:", kv[1])
					hasError = true
					continue
				}
				cfg.From = kv[1]

			case "to":
				if !isValidLang(kv[1]) {
					fmt.Fprintln(os.Stderr, "Unknown language:", kv[1])
					hasError = true
					continue
				}
				cfg.To = kv[1]

			
			default:
				fmt.Fprintln(os.Stderr, "Unknown key:", kv[0])
				hasError = true
			}
		}

		if hasError {
			os.Exit(1)
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
  ./keyflip              Run daemon
  ./keyflip install      Install autostart
  ./keyflip uninstall    Remove autostart
  ./keyflip config       Configure behavior
  ./keyflip lang         Show available languages
  ./keyflip help         Show this message

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
