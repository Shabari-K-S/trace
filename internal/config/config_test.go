package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLifecycle(t *testing.T) {
	// Setup temporary directory
	tmpDir, err := os.MkdirTemp("", "trace_config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Mock project root by creating .trace folder
	traceDir := filepath.Join(tmpDir, ".trace")
	err = os.Mkdir(traceDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	configFile := filepath.Join(traceDir, "config.json")

	// We need to carefully mock FindingProjectRoot if the config package relies on it.
	// Looking at config.go, `getConfigPath` uses `core.FindProjectRoot`.
	// Since we can't easily mock `core.FindProjectRoot` without changing code structure (dependency injection),
	// this integration test might be flaky if run outside a project.
	// However, `InitConfig` calls `getConfigPath`.

	// Workaround: We will test Save and Load with a specific config if we refactor,
	// but seemingly `Save` calls `getConfigPath`.
	// Ideally `config` package shouldn't depend on `core` for path finding if we want isolated unit tests,
	// or we should be able to pass the path.

	// Let's assume for this test we are testing the logic of Load/Save if we could force the path,
	// checking `config.go` again... `Save(cfg Config)` calls `getConfigPath`.

	// Wait, `getConfigPath` fails if no root found.
	// So we must run this test in a way that FindProjectRoot works.
	// FindProjectRoot looks for `.trace` or `.git`.
	// We created `.trace` in `tmpDir`.
	// So if we Chdir to `tmpDir`, it should work.

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(tmpDir)

	// 1. Test Init (implicitly tests Save)
	err = InitConfig()
	if err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("Config file not created at %s", configFile)
	}

	// 2. Test Load
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.TrackedFiles) == 0 {
		t.Error("Expected default tracked files, got empty")
	}

	// 3. Test AddTrackedFile
	err = AddTrackedFile("newfile.txt")
	if err != nil {
		t.Fatalf("AddTrackedFile failed: %v", err)
	}

	cfg, err = Load()
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, f := range cfg.TrackedFiles {
		if f == "newfile.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Error("newfile.txt was not added locally")
	}

	// 4. Test Duplicate Add
	err = AddTrackedFile("newfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	cfg, _ = Load()
	count := 0
	for _, f := range cfg.TrackedFiles {
		if f == "newfile.txt" {
			count++
		}
	}
	if count > 1 {
		t.Error("Duplicate file added")
	}
}
