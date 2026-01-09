package cli

import (
	"fmt"

	"trace/internal/core"
	"trace/internal/diff"
	"trace/internal/store"
)

// Diff compares the current state with a specific commit or between two commits.
func Diff(target string) error {
	head, err := core.GetHEAD()
	if err != nil {
		return fmt.Errorf("get HEAD: %w", err)
	}

	if head == "" {
		fmt.Println("No commits yet.")
		return nil
	}

	// Collect current state
	current, err := collectSnapshot()
	if err != nil {
		return fmt.Errorf("collect snapshot: %w", err)
	}

	var targetCommit *core.Commit

	if target == "" {
		// Compare with HEAD
		targetCommit, err = store.LoadCommit(head)
		if err != nil {
			return fmt.Errorf("load HEAD: %w", err)
		}
		fmt.Printf("🔍 Comparing working environment with HEAD (%s)\n\n", targetCommit.ShortHash())
	} else {
		// Resolve target reference
		targetHash, err := store.ResolveCommit(target)
		if err != nil {
			return err
		}

		targetCommit, err = store.LoadCommit(targetHash)
		if err != nil {
			return fmt.Errorf("load commit: %w", err)
		}
		fmt.Printf("🔍 Comparing working environment with %s (%s)\n\n", target, targetCommit.ShortHash())
	}

	// Compare
	envDiff, fileDiff := diff.CompareSnapshots(&targetCommit.Snapshot, &current)

	if envDiff.IsEmpty() && fileDiff.IsEmpty() {
		fmt.Println("✨ No differences found.")
		return nil
	}

	if !fileDiff.IsEmpty() {
		diff.RenderFileDiff(fileDiff)
	}

	if !envDiff.IsEmpty() {
		diff.RenderEnvDiff(envDiff)
	}

	return nil
}

// DiffCommits compares two specific commits.
func DiffCommits(from, to string) error {
	fromHash, err := store.ResolveCommit(from)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", from, err)
	}

	toHash, err := store.ResolveCommit(to)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", to, err)
	}

	fromCommit, err := store.LoadCommit(fromHash)
	if err != nil {
		return fmt.Errorf("load %s: %w", from, err)
	}

	toCommit, err := store.LoadCommit(toHash)
	if err != nil {
		return fmt.Errorf("load %s: %w", to, err)
	}

	fmt.Printf("🔍 Comparing %s..%s\n\n", fromCommit.ShortHash(), toCommit.ShortHash())

	envDiff, fileDiff := diff.CompareSnapshots(&fromCommit.Snapshot, &toCommit.Snapshot)

	if envDiff.IsEmpty() && fileDiff.IsEmpty() {
		fmt.Println("✨ No differences found.")
		return nil
	}

	if !fileDiff.IsEmpty() {
		diff.RenderFileDiff(fileDiff)
	}

	if !envDiff.IsEmpty() {
		diff.RenderEnvDiff(envDiff)
	}

	return nil
}
