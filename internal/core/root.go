package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindProjectRoot looks for the project root by walking up from the current directory.
// It searches for .trace, .git, or go.mod.
func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get cwd: %w", err)
	}

	dir := cwd
	for {
		// Check for markers
		if isRoot(dir) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root of filesystem
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("project root not found (no .trace, .git, or go.mod found)")
}

func isRoot(dir string) bool {
	// Priority 1: .trace directory
	if _, err := os.Stat(filepath.Join(dir, ".trace")); err == nil {
		return true
	}

	// Priority 2: .git directory
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		return true
	}

	// Priority 3: go.mod file
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return true
	}

	return false
}
