package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	macos "github.com/Kicked-Out/KeyFlip/platforms/macos"
)

func main() {
	if len(os.Args) > 1 {
		runCLI(os.Args[1:])
		return
	}
	runDaemon()
}

func runDaemon() {
	fmt.Println("KeyFlip daemon starting...")

	macos.StartHotkeyListener(func() {
		cfg, err := macos.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			return
		}
		macos.ProcessWithConfig(cfg)
	})

	select {}
}

// =======================
// CLI
// =======================

func runCLI(args []string) {
	switch args[0] {

	case "help", "--help", "-h":
		printHelp()

	case "install":
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}

		if err := macos.InstallLaunchd(exe); err != nil {
			fmt.Fprintln(os.Stderr, "Install failed:", err)
			os.Exit(1)
		}

		fmt.Println("KeyFlip installed as launchd agent.")
		fmt.Println("Grant Accessibility permissions:")
		fmt.Println("System Settings → Privacy & Security → Accessibility → KeyFlip")

	case "uninstall":
		if err := macos.UninstallLaunchd(); err != nil {
			fmt.Fprintln(os.Stderr, "Uninstall failed:", err)
			os.Exit(1)
		}
		fmt.Println("KeyFlip uninstalled.")

	case "config":
		handleConfig(args[1:])

	default:
		fmt.Println("Unknown command:", args[0])
		printHelp()
	}
}

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

		for _, a := range args[1:] {
			kv := strings.SplitN(a, "=", 2)
			if len(kv) != 2 {
				fmt.Println("Invalid config:", a)
				continue
			}
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

		if err := macos.SaveConfig(cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Config saved.")
	}
}

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
