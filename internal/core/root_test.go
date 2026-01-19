package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	// dedicated tmp directory for tests
	tmpDir, err := os.MkdirTemp("", "trace_root_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Case 1: .trace directory
	traceRoot := filepath.Join(tmpDir, "trace_proj")
	os.MkdirAll(filepath.Join(traceRoot, ".trace"), 0755)

	subDir := filepath.Join(traceRoot, "subdir", "deep")
	os.MkdirAll(subDir, 0755)

	// Change to subdir
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(subDir)

	root, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("Expected to find root, got error: %v", err)
	}

	// Resolve symlinks just in case
	root, _ = filepath.EvalSymlinks(root)
	expected, _ := filepath.EvalSymlinks(traceRoot)

	if root != expected {
		t.Errorf("Expected root %s, got %s", expected, root)
	}

	// Case 2: No root markers
	noRoot := filepath.Join(tmpDir, "noroot")
	os.MkdirAll(noRoot, 0755)
	os.Chdir(noRoot)

	_, err = FindProjectRoot()
	if err == nil {
		t.Error("Expected error when no root markers present, got nil")
	}
}
