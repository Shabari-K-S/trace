package cli

import (
	"fmt"

	"trace/internal/config"
	"trace/internal/core"
	"trace/internal/store"
)

// Init initializes a new trace repository in the current directory.
func Init() error {
	// Create directory structure
	if err := store.Init(); err != nil {
		return fmt.Errorf("init store: %w", err)
	}

	// Create default config
	if err := config.InitConfig(); err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	// Get current branch name for message
	branch, _ := core.GetCurrentBranch()
	if branch == "" {
		branch = "main"
	}

	fmt.Printf("✨ Initialized empty Trace repository in .trace/\n")
	fmt.Printf("   Branch: %s\n", branch)
	fmt.Printf("   Tracking: .env\n")
	fmt.Printf("\nRun 'trace snap \"initial state\"' to capture your first snapshot.\n")

	return nil
}
