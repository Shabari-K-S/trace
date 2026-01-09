package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"trace/internal/core"
)

const (
	TraceDir   = ".trace"
	ObjectsDir = ".trace/objects"
	CommitsDir = ".trace/objects/commits"
	BlobsDir   = ".trace/objects/blobs"
	ConfigFile = ".trace/config.json"
	LogsDir    = ".trace/logs"
)

// Init creates the .trace directory structure.
func Init() error {
	dirs := []string{
		TraceDir,
		ObjectsDir,
		CommitsDir,
		BlobsDir,
		core.HeadsDir,
		LogsDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}
	return core.InitRefs()
}

// SaveCommit stores a commit object and returns its hash.
func SaveCommit(c *core.Commit) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal commit: %w", err)
	}

	path := filepath.Join(CommitsDir, c.Hash+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write commit: %w", err)
	}

	return nil
}

// LoadCommit reads a commit by its hash.
func LoadCommit(hash string) (*core.Commit, error) {
	path := filepath.Join(CommitsDir, hash+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read commit %s: %w", core.ShortHash(hash), err)
	}

	var c core.Commit
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse commit %s: %w", core.ShortHash(hash), err)
	}

	return &c, nil
}

// SaveBlob stores file content and returns its hash.
func SaveBlob(content []byte) (string, error) {
	hash := core.HashContent(content)
	path := filepath.Join(BlobsDir, hash)

	// Skip if blob already exists (content-addressable)
	if _, err := os.Stat(path); err == nil {
		return hash, nil
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return "", fmt.Errorf("write blob: %w", err)
	}

	return hash, nil
}

// LoadBlob reads blob content by its hash.
func LoadBlob(hash string) ([]byte, error) {
	path := filepath.Join(BlobsDir, hash)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read blob %s: %w", core.ShortHash(hash), err)
	}
	return data, nil
}

// CommitExists checks if a commit with the given hash exists.
func CommitExists(hash string) bool {
	path := filepath.Join(CommitsDir, hash+".json")
	_, err := os.Stat(path)
	return err == nil
}

// ResolveCommit resolves a partial hash or branch name to a full commit hash.
func ResolveCommit(ref string) (string, error) {
	// Try as exact commit hash first
	if CommitExists(ref) {
		return ref, nil
	}

	// Try as branch name
	hash, err := core.GetBranch(ref)
	if err == nil && hash != "" {
		return hash, nil
	}

	// Try as HEAD
	if ref == "HEAD" {
		return core.GetHEAD()
	}

	// Try as partial hash
	entries, err := os.ReadDir(CommitsDir)
	if err != nil {
		return "", fmt.Errorf("read commits: %w", err)
	}

	var matches []string
	for _, e := range entries {
		name := e.Name()
		hash := name[:len(name)-5] // Remove .json
		if len(ref) >= 4 && len(hash) >= len(ref) && hash[:len(ref)] == ref {
			matches = append(matches, hash)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("commit not found: %s", ref)
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("ambiguous commit reference: %s", ref)
	}

	return matches[0], nil
}

// GetCommitHistory returns commits from the given hash back to the root.
func GetCommitHistory(startHash string) ([]*core.Commit, error) {
	var history []*core.Commit
	hash := startHash

	for hash != "" {
		c, err := LoadCommit(hash)
		if err != nil {
			return nil, err
		}
		history = append(history, c)
		hash = c.Parent
	}

	return history, nil
}
