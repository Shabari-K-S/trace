package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"trace/internal/config"
	"trace/internal/core"
	"trace/internal/store"
)

// RestoreOptions configures the restore behavior.
type RestoreOptions struct {
	CommitRef string   // Specific commit to restore from (empty = HEAD)
	Files     []string // Specific files to restore (empty = all)
	NoBackup  bool     // Skip creating backup files
}

// Restore restores tracked files to a specific commit state.
func Restore(opts RestoreOptions) error {
	// Determine target commit
	var targetHash string
	var err error

	if opts.CommitRef == "" {
		targetHash, err = core.GetHEAD()
		if err != nil {
			return fmt.Errorf("get HEAD: %w", err)
		}
		if targetHash == "" {
			return fmt.Errorf("no commits yet")
		}
	} else {
		targetHash, err = store.ResolveCommit(opts.CommitRef)
		if err != nil {
			return err
		}
	}

	// Load target commit
	commit, err := store.LoadCommit(targetHash)
	if err != nil {
		return fmt.Errorf("load commit: %w", err)
	}

	// Load config for backup setting
	cfg, _ := config.Load()
	createBackup := cfg.BackupOnRestore && !opts.NoBackup

	// Determine which files to restore
	filesToRestore := make(map[string]string)
	if len(opts.Files) > 0 {
		// Restore specific files
		for _, file := range opts.Files {
			path := filepath.Clean(file)
			hash, exists := commit.Snapshot.Files[path]
			if !exists {
				fmt.Printf("⚠️  File not found in commit: %s\n", path)
				continue
			}
			filesToRestore[path] = hash
		}
	} else {
		// Restore all tracked files
		filesToRestore = commit.Snapshot.Files
	}

	if len(filesToRestore) == 0 {
		fmt.Println("No files to restore.")
		return nil
	}

	fmt.Printf("🔄 Restoring from commit %s\n", commit.ShortHash())
	fmt.Printf("   Message: %s\n\n", commit.Message)

	restored := 0
	for path, hash := range filesToRestore {
		// Load blob content
		content, err := store.LoadBlob(hash)
		if err != nil {
			fmt.Printf("❌ Failed to load %s: %v\n", path, err)
			continue
		}

		// Create backup if file exists and backup is enabled
		if createBackup {
			if _, err := os.Stat(path); err == nil {
				backupPath := fmt.Sprintf("%s.backup.%d", path, time.Now().Unix())
				if err := copyFile(path, backupPath); err == nil {
					fmt.Printf("   📦 Backup: %s\n", backupPath)
				}
			}
		}

		// Ensure directory exists
		dir := filepath.Dir(path)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("❌ Failed to create directory %s: %v\n", dir, err)
				continue
			}
		}

		// Write restored content
		if err := os.WriteFile(path, content, 0644); err != nil {
			fmt.Printf("❌ Failed to restore %s: %v\n", path, err)
			continue
		}

		fmt.Printf("   ✅ Restored: %s\n", path)
		restored++
	}

	fmt.Printf("\n✨ Restored %d file(s)\n", restored)
	return nil
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, content, 0644)
}
