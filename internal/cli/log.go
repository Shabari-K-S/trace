package cli

import (
	"fmt"
	"strings"
	"time"

	"trace/internal/core"
	"trace/internal/store"
)

// Log displays the commit history.
func Log(count int) error {
	head, err := core.GetHEAD()
	if err != nil {
		return fmt.Errorf("get HEAD: %w", err)
	}

	if head == "" {
		fmt.Println("No commits yet. Run 'trace snap \"message\"' to create your first commit.")
		return nil
	}

	history, err := store.GetCommitHistory(head)
	if err != nil {
		return fmt.Errorf("get history: %w", err)
	}

	if count > 0 && len(history) > count {
		history = history[:count]
	}

	branch, _ := core.GetCurrentBranch()

	for i, c := range history {
		// Header
		fmt.Printf("\033[33mcommit %s\033[0m", c.Hash)
		if i == 0 {
			fmt.Printf(" \033[36m(HEAD")
			if branch != "" {
				fmt.Printf(" -> %s", branch)
			}
			fmt.Print(")\033[0m")
		}
		fmt.Println()

		// Date
		t, _ := time.Parse(time.RFC3339, c.Timestamp)
		fmt.Printf("Date:   %s\n", t.Format("Mon Jan 2 15:04:05 2006 -0700"))

		// Message
		fmt.Printf("\n    %s\n", c.Message)

		// Summary
		var summary []string
		if len(c.Snapshot.Files) > 0 {
			var files []string
			for path := range c.Snapshot.Files {
				files = append(files, path)
			}
			summary = append(summary, fmt.Sprintf("Files: %s", strings.Join(files, ", ")))
		}
		if len(c.Snapshot.EnvKeys) > 0 {
			summary = append(summary, fmt.Sprintf("Env: %d keys", len(c.Snapshot.EnvKeys)))
		}
		if len(summary) > 0 {
			fmt.Printf("\n    %s\n", strings.Join(summary, " | "))
		}

		fmt.Println()
	}

	return nil
}
