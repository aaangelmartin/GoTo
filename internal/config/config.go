// Package config loads and saves user configuration (TOML) and resolves XDG paths.
package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config is the user-facing configuration.
type Config struct {
	Browser        string  `toml:"browser"`
	SearchEngine   string  `toml:"search_engine"`
	Theme          string  `toml:"theme"`
	ConfirmDelete  bool    `toml:"confirm_delete"`
	FuzzyThreshold float64 `toml:"fuzzy_threshold"`
}

// Default returns the zero-value-safe defaults.
func Default() Config {
	return Config{
		Browser:        "default",
		SearchEngine:   "https://www.google.com/search?q={q}",
		Theme:          "default",
		ConfirmDelete:  true,
		FuzzyThreshold: 0.4,
	}
}

// Load reads the config file, returning defaults if missing. Keys present in
// the file override the defaults; absent keys keep their default values.
func Load(path string) (Config, error) {
	cfg := Default()
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if _, err := toml.Decode(string(b), &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Save writes the config as TOML.
func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

// ConfigPath returns the path to config.toml, honoring XDG_CONFIG_HOME.
func ConfigPath() (string, error) {
	if x := os.Getenv("GOTO_CONFIG"); x != "" {
		return x, nil
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "goto", "config.toml"), nil
}

// AliasesPath returns the path to aliases.json, honoring XDG_DATA_HOME.
func AliasesPath() (string, error) {
	if x := os.Getenv("GOTO_ALIASES"); x != "" {
		return x, nil
	}
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "goto", "aliases.json"), nil
}
