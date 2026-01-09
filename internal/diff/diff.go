package diff

import (
	"fmt"

	"trace/internal/core"
)

// EnvDiff holds the differences between environment keys.
type EnvDiff struct {
	Added   []string
	Removed []string
	Changed []string // Keys that exist in both but have different value hashes
}

// FileDiff holds the differences between tracked files.
type FileDiff struct {
	Added    []string
	Removed  []string
	Modified []string
}

// CompareEnv compares environment keys between two snapshots.
func CompareEnv(oldKeys, newKeys map[string]string) EnvDiff {
	var added, removed, changed []string

	// Find added and changed
	for key, newHash := range newKeys {
		oldHash, exists := oldKeys[key]
		if !exists {
			added = append(added, key)
		} else if oldHash != newHash {
			changed = append(changed, key)
		}
	}

	// Find removed
	for key := range oldKeys {
		if _, exists := newKeys[key]; !exists {
			removed = append(removed, key)
		}
	}

	return EnvDiff{
		Added:   added,
		Removed: removed,
		Changed: changed,
	}
}

// CompareFiles compares tracked files between two snapshots.
func CompareFiles(oldFiles, newFiles map[string]string) FileDiff {
	var added, removed, modified []string

	// Find added and modified
	for path, newHash := range newFiles {
		oldHash, exists := oldFiles[path]
		if !exists {
			added = append(added, path)
		} else if oldHash != newHash {
			modified = append(modified, path)
		}
	}

	// Find removed
	for path := range oldFiles {
		if _, exists := newFiles[path]; !exists {
			removed = append(removed, path)
		}
	}

	return FileDiff{
		Added:    added,
		Removed:  removed,
		Modified: modified,
	}
}

// CompareSnapshots compares two snapshots and returns both diffs.
func CompareSnapshots(old, new *core.Snapshot) (EnvDiff, FileDiff) {
	oldEnv := make(map[string]string)
	oldFiles := make(map[string]string)
	if old != nil {
		oldEnv = old.EnvKeys
		oldFiles = old.Files
	}

	newEnv := make(map[string]string)
	newFiles := make(map[string]string)
	if new != nil {
		newEnv = new.EnvKeys
		newFiles = new.Files
	}

	return CompareEnv(oldEnv, newEnv), CompareFiles(oldFiles, newFiles)
}

// RenderEnvDiff prints environment differences.
func RenderEnvDiff(d EnvDiff) {
	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
		return
	}

	for _, k := range d.Added {
		fmt.Printf("  \033[32m+ [ENV ADDED]\033[0m   %s\n", k)
	}
	for _, k := range d.Removed {
		fmt.Printf("  \033[31m- [ENV REMOVED]\033[0m %s\n", k)
	}
	for _, k := range d.Changed {
		fmt.Printf("  \033[33m* [ENV CHANGED]\033[0m %s\n", k)
	}
}

// RenderFileDiff prints file differences.
func RenderFileDiff(d FileDiff) {
	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Modified) == 0 {
		return
	}

	for _, p := range d.Added {
		fmt.Printf("  \033[32m+ [FILE ADDED]\033[0m    %s\n", p)
	}
	for _, p := range d.Removed {
		fmt.Printf("  \033[31m- [FILE REMOVED]\033[0m  %s\n", p)
	}
	for _, p := range d.Modified {
		fmt.Printf("  \033[33m* [FILE MODIFIED]\033[0m %s\n", p)
	}
}

// IsEmpty returns true if there are no differences.
func (d EnvDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// IsEmpty returns true if there are no differences.
func (d FileDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Modified) == 0
}
