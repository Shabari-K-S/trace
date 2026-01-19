package cli

import (
	"fmt"
	"path/filepath"
	"trace/internal/config"
	"trace/internal/core"
)

// Track adds files to the tracking list in .trace/config.json.
func Track(files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("usage: trace track <file>...")
	}

	fmt.Println("Tracking new files:")
	for _, file := range files {
		clean := filepath.Clean(file)

		ignored, _ := core.ShouldIgnore(clean)
		if ignored {
			fmt.Printf("  ⚠️  Skipping %s (ignored by .traceignore)\n", clean)
			continue
		}

		err := config.AddTrackedFile(clean)
		if err != nil {
			return fmt.Errorf("track %s: %w", file, err)
		}
		fmt.Printf("  + %s\n", clean)
	}

	fmt.Println("\nUpdated .trace/config.json")
	return nil
}
