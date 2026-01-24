//go:build windows

package windows

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	From string `json:"from"`
	To string `json:"to"`
	Hotkey string `json:"hotkey"`
}

func defaultConfig() Config {
	return Config{
		From: "en",
		To: "ua",
		Hotkey: "ctrl+shift+k",
	}
}

func configPath() (string, error) {
	dir := os.Getenv("APPDATA")

	return filepath.Join(dir, "KeyFlip", "config.json"), nil
}

func LoadConfig() (Config, error) {
	path, _ := configPath()

	b, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		cfg := defaultConfig()
		_ = SaveConfig(cfg)

		return cfg, nil
	}

	var cfg Config
	_ = json.Unmarshal(b, &cfg)

	return cfg, nil
}

func SaveConfig(cfg Config) error {
	path, _ := configPath()
	_ = os.MkdirAll(filepath.Dir(path), 0755)

	b, _ := json.MarshalIndent(cfg, "", " ")

	return os.WriteFile(path, b, 0644)
}