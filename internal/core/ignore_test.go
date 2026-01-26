package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldIgnore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "trace_ignore_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create project root
	rootDir := filepath.Join(tmpDir, "proj")
	os.MkdirAll(filepath.Join(rootDir, ".trace"), 0755)

	// Create .traceignore
	ignoreContent := `
# Comment
secret.txt
*.log
temp/
`
	os.WriteFile(filepath.Join(rootDir, ".traceignore"), []byte(ignoreContent), 0644)

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(rootDir)

	tests := []struct {
		path   string
		expect bool
	}{
		{"secret.txt", true},
		{"app.log", true},
		{"temp/cache.tmp", true}, // Glob matching logic depends on implementation
		{"main.go", false},
		{"temp", true},
	}

	// Note: The current glob implementation in `ignore.go` uses `filepath.Match`.
	// `filepath.Match` does NOT support `**` or recursive directory matching like git.
	// It only matches the pattern against the name.
	// But our `ignore.go` implementation:
	// `matched, _ := filepath.Match(line, relPath)`
	// `if matched || line == relPath`
	// So `temp/` in ignore file might not match `temp/cache.tmp` with `filepath.Match("temp/", "temp/cache.tmp")` -> False.
	// Let's check `ignore.go` logic again.

	// Logic in ignore.go:
	// line := strings.TrimSpace(scanner.Text())
	// matched, _ := filepath.Match(line, relPath)
	// if matched || line == relPath { return true }

	// So `temp/` implies `line` is `temp/`.
	// `filepath.Match("temp/", "temp/cache.tmp")` is likely False.
	// So `temp/cache.tmp` will NOT be ignored unless we fix the logic or the test expectation matches the current simple logic.
	// The current logic is "Exact match relative path OR Glob match relative path".
	// So `*.log` matches `app.log`.
	// `secret.txt` matches `secret.txt`.
	// `temp/` matches `temp/` (directory).
	// But `temp` does not match `temp/` strictly string-wise if trailing slash is trimmed?
	// `strings.TrimSpace` removes spaces, not internal slashes.

	// Let's adjust expectations to what the CURRENT code likely does, then we can decide if we improve it.
	// Current code is simple.

	for _, tt := range tests {
		// Create dummy file for resolve? No, ShouldIgnore only checks path string.
		// But ShouldIgnore calls `FindProjectRoot` -> `filepath.Rel`.

		// For `temp/cache.tmp`:
		// relPath = "temp/cache.tmp"
		// line = "temp/"
		// Match("temp/", "temp/cache.tmp") -> bad pattern? or just no match.

		// Let's stick to simple tests that should pass with current logic.
		// If I trigger a failure, I know I need to fix the logic.

		ignored, err := ShouldIgnore(tt.path)
		if err != nil {
			t.Errorf("ShouldIgnore(%s) error: %v", tt.path, err)
			continue
		}

		// For `temp/cache.tmp` with `temp/` pattern, current logic probably returns False.
		// I'll comment out the recursive case expectation being True unless I fix the code.
		// But let's set it to what we EXPECT a user wants (True) and see if it fails.
		// If it fails, I'll fix the code.

		if ignored != tt.expect {
			t.Errorf("ShouldIgnore(%s) = %v, want %v", tt.path, ignored, tt.expect)
		}
	}
}
