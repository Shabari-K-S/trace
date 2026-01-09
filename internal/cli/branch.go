package cli

import (
	"fmt"
	"os"

	"trace/internal/core"
)

// Branch manages branches.
func Branch(name string, delete bool) error {
	if name == "" {
		return listBranches()
	}

	if delete {
		return deleteBranch(name)
	}

	return createBranch(name)
}

// listBranches shows all branches.
func listBranches() error {
	branches, err := core.ListBranches()
	if err != nil {
		return err
	}

	currentBranch, _ := core.GetCurrentBranch()

	if len(branches) == 0 {
		fmt.Println("No branches yet.")
		return nil
	}

	for _, branch := range branches {
		if branch == currentBranch {
			fmt.Printf("* \033[32m%s\033[0m\n", branch)
		} else {
			fmt.Printf("  %s\n", branch)
		}
	}

	return nil
}

// createBranch creates a new branch at HEAD.
func createBranch(name string) error {
	// Check if branch already exists
	hash, _ := core.GetBranch(name)
	if hash != "" {
		return fmt.Errorf("branch '%s' already exists", name)
	}

	// Get current HEAD
	head, err := core.GetHEAD()
	if err != nil {
		return fmt.Errorf("get HEAD: %w", err)
	}

	if head == "" {
		// No commits yet, just create empty branch and switch to it
		if err := core.SetHEADToBranch(name); err != nil {
			return fmt.Errorf("create branch: %w", err)
		}
		fmt.Printf("Created and switched to new branch '%s'\n", name)
		return nil
	}

	// Create branch pointing to HEAD
	if err := core.SetBranch(name, head); err != nil {
		return fmt.Errorf("create branch: %w", err)
	}

	// Switch to new branch
	if err := core.SetHEADToBranch(name); err != nil {
		return fmt.Errorf("switch branch: %w", err)
	}

	fmt.Printf("Created and switched to new branch '%s'\n", name)
	return nil
}

// deleteBranch deletes a branch.
func deleteBranch(name string) error {
	currentBranch, _ := core.GetCurrentBranch()
	if name == currentBranch {
		return fmt.Errorf("cannot delete current branch '%s'", name)
	}

	// Check if branch exists
	hash, _ := core.GetBranch(name)
	if hash == "" {
		return fmt.Errorf("branch '%s' not found", name)
	}

	// Delete branch file
	branchPath := core.HeadsDir + "/" + name
	if err := os.Remove(branchPath); err != nil {
		return fmt.Errorf("delete branch: %w", err)
	}

	fmt.Printf("Deleted branch '%s'\n", name)
	return nil
}
