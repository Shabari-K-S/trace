package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Snapshot represents the state of the project at a point in time
type Snapshot struct {
	Timestamp string   `json:"timestamp"`
	Project   string   `json:"project_path"`
	EnvKeys   []string `json:"env_keys"`
}

func main() {
	fmt.Println("🛰️ Trace: Initializing Project Snap...")

	// 1. Get current working directory
	dir, _ := os.Getwd()

	// 2. Look for .env file
	envPath := dir + "/.env"
	keys, err := parseEnvKeys(envPath)
	
	if err != nil {
		fmt.Printf("⚠️ No .env found at %s. Scanning for project structure instead.\n", envPath)
	}

	// 3. Create Snapshot object
	snap := Snapshot{
		Timestamp: time.Now().Format(time.RFC3339),
		Project:   dir,
		EnvKeys:   keys,
	}

	// 4. Save to JSON
	file, _ := json.MarshalIndent(snap, "", "  ")
	err = os.WriteFile("snap.json", file, 0644)

	if err == nil {
		fmt.Println("✅ Snap created! Check 'snap.json' to see your project state.")
		fmt.Printf("Captured %d environment keys.\n", len(keys))
	}
}

// parseEnvKeys reads a file and returns only the keys (not the values!)
func parseEnvKeys(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Ignore comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Split by '=' and take the first part (the key)
		parts := strings.SplitN(line, "=", 2)
		if len(parts) > 0 {
			keys = append(keys, strings.TrimSpace(parts[0]))
		}
	}
	return keys, scanner.Err()
}