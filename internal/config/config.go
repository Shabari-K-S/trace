package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const TraceDir = ".trace"
const ConfigFile = ".trace/config.json"

// Config defines the trace configuration.
type Config struct {
	TrackedFiles    []string `json:"tracked_files"`
	DefaultBranch   string   `json:"default_branch,omitempty"`
	BackupOnRestore bool     `json:"backup_on_restore,omitempty"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		TrackedFiles:    []string{".env"},
		DefaultBranch:   "main",
		BackupOnRestore: true,
	}
}

// InitConfig creates a default config file if it does not exist.
func InitConfig() error {
	if err := os.MkdirAll(TraceDir, 0755); err != nil {
		return err
	}

	// Do not overwrite existing config
	if _, err := os.Stat(ConfigFile); err == nil {
		return nil
	}

	cfg := DefaultConfig()
	return Save(cfg)
}

// Save writes the config to disk.
func Save(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, data, 0644)
}

// Load reads .trace/config.json, returns default if missing.
func Load() (Config, error) {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	// Apply defaults for missing fields
	if cfg.DefaultBranch == "" {
		cfg.DefaultBranch = "main"
	}

	return cfg, nil
}

// GetTrackedFiles returns the list of files to track.
func GetTrackedFiles() ([]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	return cfg.TrackedFiles, nil
}

// AddTrackedFile adds a file to the tracking list.
func AddTrackedFile(file string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	// Check if already tracked
	for _, f := range cfg.TrackedFiles {
		if filepath.Clean(f) == filepath.Clean(file) {
			return nil // Already tracked
		}
	}

	cfg.TrackedFiles = append(cfg.TrackedFiles, file)
	return Save(cfg)
}
