package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var userConfigDirFn = os.UserConfigDir

// Config is gx's user configuration.
type Config struct {
	UseNerdFontIcons bool `json:"use-nerdfont-icons"`
}

// Default returns the default configuration.
func Default() Config {
	return Config{
		UseNerdFontIcons: false,
	}
}

// FilePath returns the config file path, typically ~/.config/gx/config.json.
func FilePath() (string, error) {
	base, err := userConfigDirFn()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(base, "gx", "config.json"), nil
}

// Load reads user config from disk. Missing file returns defaults.
func Load() (Config, error) {
	cfg := Default()
	path, err := FilePath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config %s: %w", path, err)
	}

	// Support both kebab-case and snake_case key variants.
	var raw struct {
		UseNerdFontIconsKebab *bool `json:"use-nerdfont-icons"`
		UseNerdFontIconsSnake *bool `json:"use_nerdfont_icons"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}
	if raw.UseNerdFontIconsKebab != nil {
		cfg.UseNerdFontIcons = *raw.UseNerdFontIconsKebab
	} else if raw.UseNerdFontIconsSnake != nil {
		cfg.UseNerdFontIcons = *raw.UseNerdFontIconsSnake
	}

	return cfg, nil
}

// Init writes the default config file and returns its path.
// It returns an error if the file already exists.
func Init() (string, error) {
	path, err := FilePath()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("config already exists at %s", path)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat config %s: %w", path, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("create config dir: %w", err)
	}

	b, err := json.MarshalIndent(Default(), "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode default config: %w", err)
	}
	b = append(b, '\n')

	if err := os.WriteFile(path, b, 0644); err != nil {
		return "", fmt.Errorf("write config %s: %w", path, err)
	}
	return path, nil
}
