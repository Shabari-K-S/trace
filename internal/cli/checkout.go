package cli

import (
	"fmt"
	"os"

	"trace/internal/core"
	"trace/internal/store"
)

// Checkout moves HEAD to a specific commit (detached HEAD state).
func Checkout(ref string) error {
	// Try to resolve as branch first
	branches, err := core.ListBranches()
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if branch == ref {
			// Checkout branch
			if err := core.SetHEADToBranch(branch); err != nil {
				return fmt.Errorf("set HEAD: %w", err)
			}

			hash, _ := core.GetBranch(branch)
			if hash != "" {
				commit, err := store.LoadCommit(hash)
				if err == nil {
					fmt.Printf("Switched to branch '%s'\n", branch)
					fmt.Printf("   Latest: %s - %s\n", commit.ShortHash(), commit.Message)
				}
			} else {
				fmt.Printf("Switched to branch '%s' (no commits)\n", branch)
			}
			return nil
		}
	}

	// Resolve as commit hash
	hash, err := store.ResolveCommit(ref)
	if err != nil {
		return err
	}

	commit, err := store.LoadCommit(hash)
	if err != nil {
		return fmt.Errorf("load commit: %w", err)
	}

	// Set HEAD directly to commit (detached)
	if err := os.WriteFile(core.HeadFile, []byte(hash+"\n"), 0644); err != nil {
		return fmt.Errorf("set HEAD: %w", err)
	}

	fmt.Printf("HEAD is now at %s %s\n", commit.ShortHash(), commit.Message)
	fmt.Println("\n⚠️  You are in 'detached HEAD' state.")
	fmt.Println("   To return to a branch: trace checkout main")

	return nil
}
