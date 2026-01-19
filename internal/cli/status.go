package cli

import (
	"fmt"
	"os"
	"strings"

	"trace/internal/core"
	"trace/internal/diff"
	"trace/internal/monitor"
	"trace/internal/store"
)

// Status shows the current environment drift from HEAD.
func Status() error {
	// Get current branch
	branch, _ := core.GetCurrentBranch()
	if branch != "" {
		fmt.Printf("On branch \033[36m%s\033[0m\n", branch)
	} else {
		head, _ := core.GetHEAD()
		if head != "" {
			fmt.Printf("HEAD detached at \033[33m%s\033[0m\n", core.ShortHash(head))
		}
	}

	envDiff, fileDiff, procs, err := GetStatus()
	if err != nil {
		if err.Error() == "no commits" {
			fmt.Println("\nNo commits yet.")
			fmt.Println("  (use \"trace snap <message>\" to create your first commit)")
			return nil
		}
		return err
	}

	if envDiff.IsEmpty() && fileDiff.IsEmpty() {
		fmt.Println("\n✨ Nothing to commit, working environment clean")
	} else {
		fmt.Println("\nChanges not committed:")
		fmt.Println("  (use \"trace snap <message>\" to commit)")
		fmt.Println()

		if !fileDiff.IsEmpty() {
			diff.RenderFileDiff(os.Stdout, fileDiff)
		}

		if !envDiff.IsEmpty() {
			diff.RenderEnvDiff(os.Stdout, envDiff)
		}
	}

	if len(procs) > 0 {
		fmt.Println("\n🚀 Active Processes:")
		for _, p := range procs {
			portStr := ""
			if len(p.Ports) > 0 {
				strPorts := make([]string, len(p.Ports))
				for i, port := range p.Ports {
					strPorts[i] = fmt.Sprintf("%d", port)
				}
				portStr = fmt.Sprintf(" (ports: %s)", strings.Join(strPorts, ", "))
			}
			fmt.Printf("  • [%d] %s%s\n", p.PID, p.Name, portStr)
		}
	}

	return nil
}

// GetStatus returns the current drift and active processes.
func GetStatus() (diff.EnvDiff, diff.FileDiff, []monitor.ProcessInfo, error) {
	// Get HEAD commit
	head, err := core.GetHEAD()
	if err != nil {
		return diff.EnvDiff{}, diff.FileDiff{}, nil, fmt.Errorf("get HEAD: %w", err)
	}

	if head == "" {
		return diff.EnvDiff{}, diff.FileDiff{}, nil, fmt.Errorf("no commits")
	}

	// Load HEAD commit
	headCommit, err := store.LoadCommit(head)
	if err != nil {
		return diff.EnvDiff{}, diff.FileDiff{}, nil, fmt.Errorf("load HEAD: %w", err)
	}

	// Collect current state
	current, err := collectSnapshot()
	if err != nil {
		return diff.EnvDiff{}, diff.FileDiff{}, nil, fmt.Errorf("collect snapshot: %w", err)
	}

	// Compare
	envDiff, fileDiff := diff.CompareSnapshots(&headCommit.Snapshot, &current)

	// Phase 2: Process Detection
	cwd, _ := os.Getwd()
	procs, err := monitor.GetProjectProcesses(cwd)
	if err != nil {
		// Log error but don't fail?
		// For GetStatus, getting processes is part of the status.
		// If it fails, maybe return nil procs.
	}

	return envDiff, fileDiff, procs, nil
}
