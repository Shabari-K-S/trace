package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"trace/internal/config"
	"trace/internal/core"
	"trace/internal/store"
)

// Snap creates a new commit with the current environment state.
func Snap(message string) error {
	if message == "" {
		return fmt.Errorf("commit message required: trace snap \"your message\"")
	}

	// Collect current state
	snapshot, err := collectSnapshot()
	if err != nil {
		return fmt.Errorf("collect snapshot: %w", err)
	}

	// Get parent commit (current HEAD)
	parent, err := core.GetHEAD()
	if err != nil {
		return fmt.Errorf("get HEAD: %w", err)
	}

	// Create commit
	commit := core.NewCommit(parent, message, snapshot)

	// Save commit
	if err := store.SaveCommit(commit); err != nil {
		return fmt.Errorf("save commit: %w", err)
	}

	// Update HEAD
	if err := core.SetHEAD(commit.Hash); err != nil {
		return fmt.Errorf("update HEAD: %w", err)
	}

	// Print summary
	fmt.Printf("📸 Committed: %s\n", commit.ShortHash())
	fmt.Printf("   Message: %s\n", message)
	if len(snapshot.Files) > 0 {
		fmt.Printf("   Files: %d tracked\n", len(snapshot.Files))
	}
	if len(snapshot.EnvKeys) > 0 {
		fmt.Printf("   Env keys: %d captured\n", len(snapshot.EnvKeys))
	}

	return nil
}

// collectSnapshot captures the current environment state.
func collectSnapshot() (core.Snapshot, error) {
	cfg, err := config.Load()
	if err != nil {
		return core.Snapshot{}, err
	}

	snapshot := core.Snapshot{
		EnvKeys: make(map[string]string),
		Files:   make(map[string]string),
	}

	for _, file := range cfg.TrackedFiles {
		path := filepath.Clean(file)

		ignored, _ := core.ShouldIgnore(path)
		if ignored {
			continue
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Skip missing files
			}
			return core.Snapshot{}, fmt.Errorf("read %s: %w", path, err)
		}

		// Store blob for restore capability
		hash, err := store.SaveBlob(content)
		if err != nil {
			return core.Snapshot{}, fmt.Errorf("save blob for %s: %w", path, err)
		}

		snapshot.Files[path] = hash

		// If it's a .env file, also capture keys
		if strings.HasSuffix(path, ".env") || strings.Contains(path, ".env.") {
			keys := parseEnvKeys(content)
			for key, value := range keys {
				// Store hash of value, not the value itself
				snapshot.EnvKeys[key] = core.HashString(value)
			}
		}
	}

	return snapshot, nil
}

// parseEnvKeys extracts key-value pairs from .env file content.
func parseEnvKeys(content []byte) map[string]string {
	keys := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(content)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		if idx := strings.Index(line, "="); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			keys[key] = value
		}
	}

	return keys
}
