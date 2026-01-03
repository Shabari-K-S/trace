package diff

import (
	"fmt"
)

// Compare returns the added and removed keys between two snapshots
func Compare(oldKeys, newKeys []string) (added []string, removed []string) {
	keyMap := make(map[string]bool)
	for _, k := range oldKeys {
		keyMap[k] = true
	}

	// Check for added keys
	newKeyMap := make(map[string]bool)
	for _, k := range newKeys {
		newKeyMap[k] = true
		if !keyMap[k] {
			added = append(added, k)
		}
	}

	// Check for removed keys
	for _, k := range oldKeys {
		if !newKeyMap[k] {
			removed = append(removed, k)
		}
	}

	return added, removed
}

// RenderDiff prints a pretty output to the terminal
func RenderDiff(added, removed []string) {
	if len(added) == 0 && len(removed) == 0 {
		fmt.Println("✨ No environmental drift detected. Everything matches!")
		return
	}

	for _, k := range added {
		fmt.Printf("➕ [ADDED]   %s\n", k)
	}
	for _, k := range removed {
		fmt.Printf("➖ [REMOVED] %s\n", k)
	}
}