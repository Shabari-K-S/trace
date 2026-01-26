package core

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ShouldIgnore checks if a file should be ignored based on .traceignore.
func ShouldIgnore(path string) (bool, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return false, err
	}

	ignoreFile := filepath.Join(root, ".traceignore")

	file, err := os.Open(ignoreFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	// Naive implementation: simple string match or glob
	// For now, let's support exact relative path or simple glob
	// path should be relative to root for comparison
	absPath, _ := filepath.Abs(path)
	relPath, err := filepath.Rel(root, absPath)
	if err != nil {
		return false, nil
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Use filepath.Match for glob support
		matched, _ := filepath.Match(line, relPath)
		if matched {
			return true, nil
		}

		// Exact match
		if line == relPath {
			return true, nil
		}

		// Directory prefix match (e.g. "temp/" should match "temp/file" and "temp")
		if strings.HasSuffix(line, string(os.PathSeparator)) {
			dirPrefix := line
			// If relPath is inside dirPrefix
			if strings.HasPrefix(relPath, dirPrefix) {
				return true, nil
			}
			// If relPath IS the directory (ignoring trailing slash)
			if relPath == strings.TrimSuffix(dirPrefix, string(os.PathSeparator)) {
				return true, nil
			}
		}
	}

	return false, nil
}
