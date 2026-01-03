package snapshot

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type State struct {
	Timestamp string   `json:"timestamp"`
	EnvKeys   []string `json:"env_keys"`
}

// Collect captures the current project environment
func Collect() (State, error) {
	keys := []string{}
	
	// Scan for .env keys
	file, err := os.Open(".env")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" && !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
				key := strings.Split(line, "=")[0]
				keys = append(keys, strings.TrimSpace(key))
			}
		}
	}

	return State{
		Timestamp: time.Now().Format(time.RFC3339),
		EnvKeys:   keys,
	}, nil
}