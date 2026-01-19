package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"trace/internal/core"
)

const TraceDir = ".trace"

// Config defines the trace configuration.
type Config struct {
	TrackedFiles    []string `json:"tracked_files"`
	DefaultBranch   string   `json:"default_branch,omitempty"`
	BackupOnRestore bool     `json:"backup_on_restore,omitempty"`
	Hooks           Hooks    `json:"hooks,omitempty"`
}

// Hooks defines commands to run around lifecycle events.
type Hooks struct {
	PreRestore  string `json:"pre_restore,omitempty"`
	PostRestore string `json:"post_restore,omitempty"`
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
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	cfg := DefaultConfig()
	return Save(cfg)
}

func getConfigPath() (string, error) {
	root, err := core.FindProjectRoot()
	if err != nil {
		// If root not found, default to current directory but warn/error might be better?
		// For now, let's stick to CWD if root finding fails to avoid breaking non-project usage?
		// Actually, trace is project-scoped. If no root, we might want to default to CWD
		// so `init` works.
		// `init` creates .trace in CWD.
		// So if FindProjectRoot fails, we assume CWD for now, but `init` logic should be separate.
		// Let's assume ConfigFile is relative to root if found.
		// Let's assume ConfigFile is relative to root if found.
		return ".trace/config.json", nil
	}
	return filepath.Join(root, TraceDir, "config.json"), nil
}

// Save writes the config to disk.
func Save(cfg Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads .trace/config.json, returns default if missing.
func Load() (Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
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
