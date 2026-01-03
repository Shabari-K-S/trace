package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const TraceDir = ".trace"
const SnapsDir = ".trace/snaps"

func Init() error {
	return os.MkdirAll(SnapsDir, 0755)
}

func SaveSnap(data interface{}) error {
	// Create filename based on Unix timestamp
	filename := fmt.Sprintf("%d.json", time.Now().Unix())
	path := filepath.Join(SnapsDir, filename)

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Update history log
	logPath := filepath.Join(TraceDir, "history.log")
	logEntry := fmt.Sprintf("[%s] Saved snap to %s\n", time.Now().Format(time.Kitchen), filename)
	
	f, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(logEntry)

	return os.WriteFile(path, file, 0644)
}

// GetLatestTwoSnaps returns the data from the two most recent snapshot files
func GetLatestTwoSnaps() (prev, current []string, err error) {
	files, err := os.ReadDir(SnapsDir)
	if err != nil || len(files) < 2 {
		return nil, nil, fmt.Errorf("not enough snapshots to compare (need at least 2)")
	}

	// ReadDir returns files sorted by name. Since we use Unix timestamps, 
	// the last two files are the most recent.
	f1 := filepath.Join(SnapsDir, files[len(files)-2].Name())
	f2 := filepath.Join(SnapsDir, files[len(files)-1].Name())

	prevData, _ := os.ReadFile(f1)
	currData, _ := os.ReadFile(f2)

	// Simple helper to unmarshal specific format
	var s1, s2 struct { EnvKeys []string `json:"env_keys"` }
	json.Unmarshal(prevData, &s1)
	json.Unmarshal(currData, &s2)

	return s1.EnvKeys, s2.EnvKeys, nil
}