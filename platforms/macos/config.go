package macos

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Hotkey string `json:"hotkey"`
}

func defaultConfig() Config {
	return Config{
		From:   "en",
		To:     "ua",
		Hotkey: "cmd+shift+k",
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "Library", "Application Support", "KeyFlip")
	return filepath.Join(dir, "config.json"), nil
}

func LoadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := defaultConfig()
			_ = SaveConfig(cfg)
			return cfg, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}

	// мінімальна валідація
	if cfg.From == "" || cfg.To == "" {
		return Config{}, errors.New("invalid config: from/to must not be empty")
	}

	return cfg, nil
}

func SaveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
