package macos

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Hotkey string `json:"hotkey"`
}

// Returns the default configuration, hotkey is just a text here the real hotkey handling in listener
func defaultConfig() Config {
	return Config{
		From:   "en",
		To:     "ua",
		Hotkey: "cmd+shift+k",
	}
}

// Returns the path to the config file
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "Library", "Application Support", "KeyFlip")
	return filepath.Join(dir, "config.json"), nil
}

// Loads the configuration from file, if not exists creates default
func LoadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := defaultConfig()
			if saveErr := SaveConfig(cfg); saveErr != nil {
				
				return cfg, fmt.Errorf("config not found, failed to create default: %w", saveErr)
			}
			return cfg, nil
		}
		return Config{}, err
	}

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}

	// Basic validation
	if cfg.From == "" || cfg.To == "" {
		return Config{}, errors.New("invalid config: from/to must not be empty")
	}

	return cfg, nil
}

// Saves the configuration to file
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
