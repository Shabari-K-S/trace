package cli

import (
	"fmt"

	"trace/internal/core"
	"trace/internal/diff"
	"trace/internal/store"
)

// Status shows the current environment drift from HEAD.
func Status() error {
	// Get current branch
	branch, _ := core.GetCurrentBranch()
	if branch != "" {
		fmt.Printf("On branch \033[36m%s\033[0m\n", branch)
	} else {
		head, _ := core.GetHEAD()
		if head != "" {
			fmt.Printf("HEAD detached at \033[33m%s\033[0m\n", core.ShortHash(head))
		}
	}

	// Get HEAD commit
	head, err := core.GetHEAD()
	if err != nil {
		return fmt.Errorf("get HEAD: %w", err)
	}

	if head == "" {
		fmt.Println("\nNo commits yet.")
		fmt.Println("  (use \"trace snap <message>\" to create your first commit)")
		return nil
	}

	// Load HEAD commit
	headCommit, err := store.LoadCommit(head)
	if err != nil {
		return fmt.Errorf("load HEAD: %w", err)
	}

	// Collect current state
	current, err := collectSnapshot()
	if err != nil {
		return fmt.Errorf("collect snapshot: %w", err)
	}

	// Compare
	envDiff, fileDiff := diff.CompareSnapshots(&headCommit.Snapshot, &current)

	if envDiff.IsEmpty() && fileDiff.IsEmpty() {
		fmt.Println("\n✨ Nothing to commit, working environment clean")
		return nil
	}

	fmt.Println("\nChanges not committed:")
	fmt.Println("  (use \"trace snap <message>\" to commit)")
	fmt.Println()

	if !fileDiff.IsEmpty() {
		diff.RenderFileDiff(fileDiff)
	}

	if !envDiff.IsEmpty() {
		diff.RenderEnvDiff(envDiff)
	}

	return nil
}
