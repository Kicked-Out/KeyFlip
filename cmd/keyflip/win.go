//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kicked-Out/KeyFlip/core"
	windows "github.com/Kicked-Out/KeyFlip/platforms/windows"
)

func main() {
	if len(os.Args) > 1 {
		runCLI(os.Args[1:])
		return
	}

	runDaemon()
}

func runDaemon() {
	fmt.Println("KeyFlip daemon starting (Windows)...")

	err := windows.StartHotkeyListener(func() {
		cfg, err := windows.LoadConfig()

		if err != nil {
			fmt.Println("Failed to load config:", err)
			return
		}

		windows.ProcessWithConfig(cfg)
	})

	if err != nil {
		fmt.Println("Hotkey error:", err)
		return
	}

	select {}
}

func runCLI(args []string) {
	switch args[0] {
	case "help", "--help", "-h":
		printHelp()

	case "restart":
		if err := windows.UninstallAutostart(); err != nil {
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

		if err := windows.InstallAutostart(exe); err != nil {
			fmt.Fprintln(os.Stderr, "Restart failed:", err)
			os.Exit(1)
		}

		fmt.Println("KeyFlip restarted.")

	case "install":
		exe, err := os.Executable()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}

		already := windows.AutostartInstalled()

		if err := windows.InstallAutostart(exe); err != nil {
			fmt.Fprintln(os.Stderr, "Install failed:", err)
			os.Exit(1)
		}

		if already {
			fmt.Println("KeyFlip restarted.")
		} else {
			fmt.Println("KeyFlip installed.")
			fmt.Println("If hotkeys do not work, run as Administrator.")
		}
	
	case "uninstall":
		if err := windows.UninstallAutostart(); err != nil {
			fmt.Fprintln(os.Stderr, "Unistall failed:", err)
			os.Exit(1)
		}

		fmt.Println("KeyFlip uninstalled.")

	case "config":
		handleConfig(args[1:])

	case "lang":
		showLanguages()

	default:
		fmt.Println("Unknown command:", args[0])
		printHelp()
	}
}

func handleConfig(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage:")
		fmt.Println("	keyflip config show")
		fmt.Println("	keyflip config set from=xx to=yy")

		return
	}

	switch args[0] {
	case "show":
		cfg, err := windows.LoadConfig()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("from:		%s\n", cfg.From)
		fmt.Printf("to:		%s\n", cfg.To)
		fmt.Printf("hotkey:	%s\n", cfg.Hotkey)

	case "set":
		cfg, err := windows.LoadConfig()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		layoutsPath, err := windows.LayoutsPath()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to find layouts.json", err)
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

		if err := windows.SaveConfig(cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("Config saved.")
	}
}

func printHelp() {
	fmt.Println(`
KeyFlip — keyboard layout fixer (Windows)

USAGE:
	keyflip					Run daemon
	keyflip install			Install autostart
	keyflip uninstall		Remove autostart
	keyflip restart			Reinstall autostart
	keyflip config			Configure behavior
	keyflip lang			Show available languages
	keyflip help			Show this message`)
}

func showLanguages() {
	path, err := windows.LayoutsPath()

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
		fmt.Println(" - ", lang)
	}
}