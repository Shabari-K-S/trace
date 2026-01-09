package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// Commit represents a point-in-time snapshot of the project environment.
// Similar to a Git commit, it has a parent, message, and captured state.
type Commit struct {
	Hash      string   `json:"hash"`
	Parent    string   `json:"parent,omitempty"`
	Timestamp string   `json:"timestamp"`
	Message   string   `json:"message"`
	Snapshot  Snapshot `json:"snapshot"`
}

// Snapshot holds the environment state at commit time.
type Snapshot struct {
	EnvKeys map[string]string `json:"env_keys"`  // key -> hash of value
	Files   map[string]string `json:"files"`     // path -> content hash
}

// NewCommit creates a new commit with the given parent, message, and snapshot.
// It computes the hash based on the commit content.
func NewCommit(parent, message string, snapshot Snapshot) *Commit {
	c := &Commit{
		Parent:    parent,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   message,
		Snapshot:  snapshot,
	}
	c.Hash = c.computeHash()
	return c
}

// computeHash generates a SHA256 hash of the commit content.
func (c *Commit) computeHash() string {
	data, _ := json.Marshal(struct {
		Parent    string   `json:"parent"`
		Timestamp string   `json:"timestamp"`
		Message   string   `json:"message"`
		Snapshot  Snapshot `json:"snapshot"`
	}{
		Parent:    c.Parent,
		Timestamp: c.Timestamp,
		Message:   c.Message,
		Snapshot:  c.Snapshot,
	})
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// ShortHash returns the first 7 characters of the commit hash.
func (c *Commit) ShortHash() string {
	if len(c.Hash) >= 7 {
		return c.Hash[:7]
	}
	return c.Hash
}
