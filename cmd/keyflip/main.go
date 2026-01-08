package main

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
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
	fmt.Println("Press your configured hotkey to transform selected text.")

	macos.StartHotkeyListener(macos.Process)

	select {} // daemon живе тут
}

func runCLI(args []string) {
	switch args[0] {

	case "help", "--help", "-h":
		printHelp()

	case "status":
		fmt.Println("KeyFlip status: (stub)")

	case "install":
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		exe, _ = filepath.EvalSymlinks(exe)

		if err := macos.InstallLaunchd(exe); err != nil {
			fmt.Fprintln(os.Stderr, "Install failed:", err)
			os.Exit(1)
		}

		fmt.Println("KeyFlip installed as launchd agent.")
		fmt.Println("  IMPORTANT:")
		fmt.Println("Grant Accessibility permissions:")
		fmt.Println("System Settings → Privacy & Security → Accessibility → enable KeyFlip")

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
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
		fmt.Println("Usage:")
		fmt.Println("  keyflip config show")
		fmt.Println("  keyflip config set from=en to=ua")
		fmt.Println("  keyflip config set hotkey=cmd+shift+k")
		return
	}

	switch args[0] {
	case "show":
		cfg, err := macos.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Printf("from:   %s\n", cfg.From)
		fmt.Printf("to:     %s\n", cfg.To)
		fmt.Printf("hotkey: %s\n", cfg.Hotkey)

	case "set":
		cfg, err := macos.LoadConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		// парсимо аргументи формату key=value
		for _, a := range args[1:] {
			kv := strings.SplitN(a, "=", 2)
			if len(kv) != 2 {
				continue
			}
			k, v := kv[0], kv[1]
			switch k {
			case "from":
				cfg.From = v
			case "to":
				cfg.To = v
			case "hotkey":
				cfg.Hotkey = v
			}
		}

		if err := macos.SaveConfig(cfg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println("Config saved.")

	default:
		fmt.Println("Unknown config command:", args[0])
	}
}


func printHelp() {
	fmt.Println(`
KeyFlip — keyboard layout fixer

USAGE:
  keyflip              Run background daemon
  keyflip help         Show this help
  keyflip status       Show daemon status
  keyflip install      Install autostart (launchd)
  keyflip uninstall    Remove autostart
  keyflip config       Configure language & hotkey

EXAMPLES:
  keyflip
  keyflip help
`)
}
