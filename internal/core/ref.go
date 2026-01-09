package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	TraceDir   = ".trace"
	HeadFile   = ".trace/HEAD"
	RefsDir    = ".trace/refs"
	HeadsDir   = ".trace/refs/heads"
	DefaultRef = "main"
)

// GetHEAD returns the current commit hash that HEAD points to.
// HEAD can either point directly to a commit hash, or to a branch reference.
func GetHEAD() (string, error) {
	data, err := os.ReadFile(HeadFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No HEAD yet
		}
		return "", err
	}

	content := strings.TrimSpace(string(data))

	// Check if HEAD points to a ref (like "ref: refs/heads/main")
	if strings.HasPrefix(content, "ref: ") {
		refPath := strings.TrimPrefix(content, "ref: ")
		return readRef(filepath.Join(TraceDir, refPath))
	}

	// Otherwise, HEAD is a detached commit hash
	return content, nil
}

// SetHEAD updates HEAD to point to the given commit hash.
// If HEAD currently points to a branch, it updates the branch instead.
func SetHEAD(hash string) error {
	data, err := os.ReadFile(HeadFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	content := strings.TrimSpace(string(data))

	// If HEAD points to a ref, update that ref
	if strings.HasPrefix(content, "ref: ") {
		refPath := strings.TrimPrefix(content, "ref: ")
		return writeRef(filepath.Join(TraceDir, refPath), hash)
	}

	// Otherwise, update HEAD directly
	return os.WriteFile(HeadFile, []byte(hash+"\n"), 0644)
}

// GetCurrentBranch returns the name of the current branch, or empty if detached.
func GetCurrentBranch() (string, error) {
	data, err := os.ReadFile(HeadFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	content := strings.TrimSpace(string(data))
	if strings.HasPrefix(content, "ref: refs/heads/") {
		return strings.TrimPrefix(content, "ref: refs/heads/"), nil
	}

	return "", nil // Detached HEAD
}

// SetHEADToBranch sets HEAD to point to a branch reference.
func SetHEADToBranch(branch string) error {
	content := fmt.Sprintf("ref: refs/heads/%s\n", branch)
	return os.WriteFile(HeadFile, []byte(content), 0644)
}

// GetBranch returns the commit hash that a branch points to.
func GetBranch(name string) (string, error) {
	return readRef(filepath.Join(HeadsDir, name))
}

// SetBranch updates a branch to point to the given commit hash.
func SetBranch(name, hash string) error {
	return writeRef(filepath.Join(HeadsDir, name), hash)
}

// ListBranches returns all branch names.
func ListBranches() ([]string, error) {
	entries, err := os.ReadDir(HeadsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var branches []string
	for _, e := range entries {
		if !e.IsDir() {
			branches = append(branches, e.Name())
		}
	}
	return branches, nil
}

// InitRefs creates the refs directory structure and sets up default branch.
func InitRefs() error {
	if err := os.MkdirAll(HeadsDir, 0755); err != nil {
		return err
	}
	// Set HEAD to point to main branch (even though it doesn't exist yet)
	return SetHEADToBranch(DefaultRef)
}

func readRef(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func writeRef(path string, hash string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(hash+"\n"), 0644)
}
