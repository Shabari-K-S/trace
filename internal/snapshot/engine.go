package snapshot

import (
    "bufio"
    "crypto/sha256"
    "encoding/hex"
    "os"
    "path/filepath"
    "strings"
    "time"

    "trace/internal/config"
)

// Collect captures the current project environment
func Collect() (State, error) {
    keys := collectEnvKeys(".env")

    files, err := collectTrackedFiles()
    if err != nil {
        return State{}, err
    }

    return State{
        Timestamp: time.Now().Format(time.RFC3339),
        EnvKeys:   keys,
        Files:     files,
    }, nil
}

func collectEnvKeys(path string) []string {
    keys := []string{}

    file, err := os.Open(path)
    if err != nil {
        return keys
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
            continue
        }
        key := strings.SplitN(line, "=", 2)[0]
        keys = append(keys, strings.TrimSpace(key))
    }

    return keys
}

func collectTrackedFiles() ([]FileSnapshot, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }

    var files []FileSnapshot

    for _, p := range cfg.TrackedFiles {
        // Normalize path, but treat as project-relative
        clean := filepath.Clean(p)

        data, err := os.ReadFile(clean)
        if err != nil {
            // Skip missing files quietly for now
            continue
        }

        h := sha256.Sum256(data)
        files = append(files, FileSnapshot{
            Path: clean,
            Hash: hex.EncodeToString(h[:]),
        })
    }

    return files, nil
}
