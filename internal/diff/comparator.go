package diff

import (
    "fmt"

    "trace/internal/snapshot"
)

// CompareEnv returns the added and removed keys between two snapshots.
func CompareEnv(oldKeys, newKeys []string) (added []string, removed []string) {
    keyMap := make(map[string]bool)
    for _, k := range oldKeys {
        keyMap[k] = true
    }

    newKeyMap := make(map[string]bool)
    for _, k := range newKeys {
        newKeyMap[k] = true
        if !keyMap[k] {
            added = append(added, k)
        }
    }

    for _, k := range oldKeys {
        if !newKeyMap[k] {
            removed = append(removed, k)
        }
    }

    return added, removed
}

type FileDiff struct {
    Added    []string
    Removed  []string
    Modified []string
}

func CompareFiles(oldFiles, newFiles []snapshot.FileSnapshot) FileDiff {
    oldMap := make(map[string]string)
    newMap := make(map[string]string)

    for _, f := range oldFiles {
        oldMap[f.Path] = f.Hash
    }
    for _, f := range newFiles {
        newMap[f.Path] = f.Hash
    }

    var added, removed, modified []string

    // Added or modified
    for path, newHash := range newMap {
        oldHash, ok := oldMap[path]
        if !ok {
            added = append(added, path)
            continue
        }
        if oldHash != newHash {
            modified = append(modified, path)
        }
    }

    // Removed
    for path := range oldMap {
        if _, ok := newMap[path]; !ok {
            removed = append(removed, path)
        }
    }

    return FileDiff{
        Added:    added,
        Removed:  removed,
        Modified: modified,
    }
}

func RenderDiffEnv(added, removed []string) {
    if len(added) == 0 && len(removed) == 0 {
        fmt.Println("No env key drift detected.")
        return
    }

    for _, k := range added {
        fmt.Printf(" + [ENV ADDED]   %s\n", k)
    }
    for _, k := range removed {
        fmt.Printf(" - [ENV REMOVED] %s\n", k)
    }
}

func RenderDiffFiles(fd FileDiff) {
    if len(fd.Added) == 0 && len(fd.Removed) == 0 && len(fd.Modified) == 0 {
        fmt.Println("✨ No file drift detected.")
        return
    }

    for _, p := range fd.Added {
        fmt.Printf(" + [FILE ADDED]    %s\n", p)
    }
    for _, p := range fd.Removed {
        fmt.Printf(" - [FILE REMOVED]  %s\n", p)
    }
    for _, p := range fd.Modified {
        fmt.Printf(" * [FILE MODIFIED] %s\n", p)
    }
}
