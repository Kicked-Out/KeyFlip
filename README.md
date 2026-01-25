
```markdown
# KeyFlip — Keyboard Layout Fixer for macOS

KeyFlip automatically converts text in your clipboard between different keyboard layouts.  
Especially useful when you accidentally type text in the wrong layout.

---

## 📥 Download & Installation

### 1. Clone the repository
```bash
git clone https://github.com/YourUsername/KeyFlip.git
cd KeyFlip
```

### 2. Build for macOS

```bash
go build -o keyflip ./package/macos
```

### 3. Run

* **Daemon mode (runs in the background):**

```bash
./keyflip
```

* **CLI commands:**

```bash
./keyflip help
```

---

## 🖥 CLI Commands

| Command                              | Description                                  |
| ------------------------------------ | -------------------------------------------- |
| `./keyflip`                          | Run KeyFlip daemon                           |
| `./keyflip install`                  | Install autostart (launchd agent)            |
| `./keyflip uninstall`                | Remove autostart                             |
| `./keyflip restart`                  | Restart the autostart agent                  |
| `./keyflip config show`              | Show current configuration (from/to layouts) |
| `./keyflip config set from=xx to=yy` | Change default conversion layouts            |
| `./keyflip lang`                     | List available languages                     |

> **Note:** Changing the `hotkey` in the config has no effect. The actual hotkey is hardcoded to **Cmd+Shift+K** on macOS.

---

## 🎹 Default Layouts

KeyFlip includes two default layouts: **English (en)** and **Ukrainian (ua)**.

### Example `layouts.json` structure:

```json
{
  "layouts": {
    "en": {
      "q": "q", "w": "w", "e": "e", "r": "r", "t": "t", "y": "y", "u": "u", "i": "i", "o": "o", "p": "p",
      "a": "a", "s": "s", "d": "d", "f": "f", "g": "g", "h": "h", "j": "j", "k": "k", "l": "l",
      "z": "z", "x": "x", "c": "c", "v": "v", "b": "b", "n": "n", "m": "m"
    },
    "ua": {
      "q": "й", "w": "ц", "e": "у", "r": "к", "t": "е", "y": "н", "u": "г", "i": "ш", "o": "щ", "p": "з",
      "a": "ф", "s": "і", "d": "в", "f": "а", "g": "п", "h": "р", "j": "о", "k": "л", "l": "д",
      "z": "я", "x": "ч", "c": "с", "v": "м", "b": "и", "n": "т", "m": "ь"
    }
  }
}
```

> You can add new layouts by following the same structure: map each character in the "from" layout to its corresponding character in the "to" layout.

---

## ⌨️ Usage

1. Copy some text in the "from" layout (e.g., English).
2. Press **Cmd+Shift+K** (hotkey).
3. KeyFlip automatically converts the clipboard text to the "to" layout (e.g., Ukrainian) and pastes it wherever your cursor is.

---

## 🔧 Notes

* Only **macOS** is currently supported.
* The default hotkey is **Cmd+Shift+K** and cannot be changed via CLI.
* Config file is stored at:

```
~/Library/Application Support/KeyFlip/config.json
```

* Layouts file (`layouts.json`) is stored alongside the executable.

---

## 📄 Documentation

Full documentation and examples are available in the `docs/` folder. You can open `index.html` in a browser to see a static site with usage instructions and download links.

---

## ⚡ Contribution

1. Fork the repository.
2. Create a feature branch:

```bash
git checkout -b feature/my-new-layout
```

3. Make your changes.
4. Push and submit a pull request.

---

## 🏷 License

This project is open-source. See `LICENSE` for details.


