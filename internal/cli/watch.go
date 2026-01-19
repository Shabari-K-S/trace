package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"trace/internal/diff"
)

// Watch monitors the environment for changes and updates the display.
func Watch(interval time.Duration) error {
	fmt.Println("👀 Watching for changes... (Ctrl+C to stop)")

	// Create a buffer to render output to check for changes
	var lastOutput string

	// Initial check
	output, err := checkAndRender()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	} else {
		// Clear screen and print initial state
		printUpdate(output)
		lastOutput = output
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		output, err := checkAndRender()
		if err != nil {
			// If transient error, maybe just log it?
			// For now, let's keep retrying
			continue
		}

		if output != lastOutput {
			printUpdate(output)
			lastOutput = output
		}
	}

	return nil
}

func checkAndRender() (string, error) {
	envDiff, fileDiff, procs, err := GetStatus()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	// We only show output if there are "changes" OR active processes.
	// Status() shows "Nothing to commit" if clean.
	// Watch mode usually shows only "relevant" info.
	// Let's mimic Status() output but capturing it.

	clean := envDiff.IsEmpty() && fileDiff.IsEmpty()

	if clean {
		// If clean, and no processes, maybe show just "Watching..." or nothing?
		// User might want to see active processes even if clean.
		// Let's just say "Environment clean"
		fmt.Fprintln(&buf, "✨ Environment clean")
	} else {
		fmt.Fprintln(&buf, "Changes not committed:")
		fmt.Fprintln(&buf)

		if !fileDiff.IsEmpty() {
			diff.RenderFileDiff(&buf, fileDiff)
		}

		if !envDiff.IsEmpty() {
			diff.RenderEnvDiff(&buf, envDiff)
		}
	}

	if len(procs) > 0 {
		fmt.Fprintln(&buf, "\n🚀 Active Processes:")
		for _, p := range procs {
			portStr := ""
			if len(p.Ports) > 0 {
				strPorts := make([]string, len(p.Ports))
				for i, port := range p.Ports {
					strPorts[i] = fmt.Sprintf("%d", port)
				}
				portStr = fmt.Sprintf(" (ports: %s)", strings.Join(strPorts, ", "))
			}
			fmt.Fprintf(&buf, "  • [%d] %s%s\n", p.PID, p.Name, portStr)
		}
	} else if clean {
		// If clean and no processes, status is fully clean
		// We already printed "Environment clean"
	}

	return buf.String(), nil
}

func printUpdate(output string) {
	// Clear screen and move cursor to top left
	fmt.Print("\033[H\033[2J")
	fmt.Println("👀 Watching for changes... (Ctrl+C to stop)")
	fmt.Println()
	fmt.Print(output)
}
