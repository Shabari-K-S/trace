package config

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

const TraceDir = ".trace"
const ConfigFile = "config.json"

type Config struct {
    TrackedFiles []string `json:"tracked_files"`
}

// InitConfig creates a default config file if it does not exist.
func InitConfig() error {
    if err := os.MkdirAll(TraceDir, 0755); err != nil {
        return err
    }

    path := filepath.Join(TraceDir, ConfigFile)

    // Do not overwrite existing config
    if _, err := os.Stat(path); err == nil {
        return nil
    }

    cfg := Config{
        TrackedFiles: []string{
            ".env",
        },
    }

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0644)
}

// Load reads .trace/config.json, returns default if missing.
func Load() (Config, error) {
    path := filepath.Join(TraceDir, ConfigFile)

    data, err := os.ReadFile(path)
    if err != nil {
        // No config yet: return empty config (no tracked files)
        if os.IsNotExist(err) {
            return Config{}, nil
        }
        return Config{}, fmt.Errorf("load config: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return Config{}, fmt.Errorf("parse config: %w", err)
    }

    return cfg, nil
}
