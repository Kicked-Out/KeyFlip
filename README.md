# KeyFlip — Keyboard Layout Fixer

KeyFlip automatically converts text in your clipboard between different keyboard layouts.  
It is especially useful when text is accidentally typed using the wrong keyboard layout.

The application runs in the background and converts clipboard text using a global hotkey.

---

## 📥 Download & Installation

### Option 1: Download prebuilt binaries (recommended)

Prebuilt binaries are available on the **GitHub Releases** page.

1. Open the latest release
2. Download the archive for your operating system
3. Extract the archive to any directory

Each archive contains:

- the KeyFlip executable
- `layouts.json`

Both files must remain in the same directory.

---

### Option 2: Build from source

#### Requirements

- Go **1.22+**
- Git
- Supported OS: macOS, Windows 10 / 11

#### Clone repository

```bash
git clone https://github.com/Kicked-Out/KeyFlip
cd KeyFlip
```

---

### Build on macOS

```bash
go build -o keyflip ./cmd/keyflip
```

---

### Build on Windows

```powershell
go build -o keyflip.exe ./cmd/keyflip
```

Ensure that `layouts.json` is located next to the compiled binary.

---

## ▶ Running the application

### Background mode (default)

#### macOS

```bash
./keyflip
```

#### Windows

```powershell
keyflip.exe
```

The application runs in the background and listens for the global hotkey.

---

### CLI mode

```bash
keyflip help
```

or on Windows:

```powershell
keyflip.exe help
```

---

## 🖥 CLI Commands (macOS)

| Command                              | Description                     |
| ----------------------------------   | ------------------------------- |
| `./keyflip`                          | Run background daemon           |
| `./keyflip install`                  | Enable autostart                |
| `./keyflip uninstall`                | Disable autostart               |
| `./keyflip restart`                  | Restart background process      |
| `./keyflip config show`              | Show current configuration      |
| `./keyflip config set from=xx to=yy` | Set default conversion layouts  |
| `./keyflip lang`                     | List available keyboard layouts |

## 🖥 CLI Commands (Windows)

| Command                            | Description                     |
| ---------------------------------- | ------------------------------- |
| `keyflip`                          | Run background daemon           |
| `keyflip install`                  | Enable autostart                |
| `keyflip uninstall`                | Disable autostart               |
| `keyflip restart`                  | Restart background process      |
| `keyflip config show`              | Show current configuration      |
| `keyflip config set from=xx to=yy` | Set default conversion layouts  |
| `keyflip lang`                     | List available keyboard layouts |

---

## ⌨️ Hotkeys

| Operating System | Hotkey           |
| ---------------- | ---------------- |
| macOS            | Cmd + Shift + K  |
| Windows          | Ctrl + Shift + K |

> The hotkey is currently hardcoded and cannot be changed via configuration or CLI.

---

## 🎹 Keyboard Layouts

KeyFlip uses a JSON-based layout mapping system.

### Default layouts

- English (`en`)
- Ukrainian (`ua`)

### `layouts.json` example

```json
{
    "layouts": {
        "en": {
            "q": "q", "w": "w", "e": "e", "r": "r",
            "t": "t", "y": "y", "u": "u", "i": "i",
            "o": "o", "p": "p", "a": "a", "s": "s",
            "d": "d", "f": "f", "g": "g", "h": "h",
            "j": "j", "k": "k", "l": "l", "z": "z",
            "x": "x", "c": "c", "v": "v", "b": "b",
            "n": "n", "m": "m"
        },
        "ua": {
            "q": "й", "w": "ц", "e": "у", "r": "к",
            "t": "е", "y": "н", "u": "г", "i": "ш",
            "o": "щ", "p": "з", "a": "ф", "s": "і",
            "d": "в", "f": "а", "g": "п", "h": "р",
            "j": "о", "k": "л", "l": "д", "z": "я",
            "x": "ч", "c": "с", "v": "м", "b": "и",
            "n": "т", "m": "ь"
        }
    }
}
```

To add a new layout, define a new key and map each character from the source layout to its target character.

---

## 📂 Configuration

### macOS

```
~/Library/Application Support/KeyFlip/config.json
```

---

### Windows

```
%APPDATA%\KeyFlip\config.json
```

Example:

```
C:\Users\<username>\AppData\Roaming\KeyFlip\config.json
```

---

### Layouts file

For all operating systems:

- `layouts.json` must be stored **next to the executable**

---

## 🔁 Autostart

### Enable autostart

```bash
keyflip install
```

or on Windows:

```powershell
keyflip.exe install
```

- macOS: registers a `launchd` agent
- Windows: registers an entry in user startup

---

### Disable autostart

```bash
keyflip uninstall
```

---

### Restart background process

```bash
keyflip restart
```

---

## ⚠ Platform Notes

### macOS

- Requires accessibility permissions for keyboard events
- Uses `launchd` for background execution

### Windows

- Does not require administrator privileges
- Clipboard access uses native Windows APIs

---

## 📄 Documentation

Additional documentation and examples are available in the `docs/` directory.
Open `docs/index.html` in a browser to view the static documentation site.

---

## ⚡ Contribution

1. Fork the repository
2. Create a feature branch

```bash
git checkout -b feature/my-change
```

3. Commit your changes
4. Push and open a pull request

---

## 🏷 License

This project is open-source.
See the `LICENSE` file for details.
