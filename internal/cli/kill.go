package cli

import (
	"fmt"
	"strconv"
	"syscall"
	"trace/internal/monitor"

	"github.com/shirou/gopsutil/v4/process"
)

// Kill stops a process by PID or Port.
func Kill(target string) error {
	pid, reason, err := ResolveKillTarget(target)
	if err != nil {
		return err
	}

	// Find process name for confirmation
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}
	name, _ := p.Name()

	fmt.Printf("Killing %s [%s]...\n", name, reason)

	if err := KillPID(pid); err != nil {
		return err
	}

	fmt.Println("✅ Process stopped.")
	return nil
}

// KillPID stops a process by PID (silent, for TUI).
func KillPID(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}

	if err := p.Kill(); err != nil {
		if err == syscall.EPERM {
			return fmt.Errorf("permission denied")
		}
		return fmt.Errorf("kill process: %w", err)
	}
	return nil
}

// ResolveKillTarget figures out the PID from a target string.
func ResolveKillTarget(target string) (int32, string, error) {
	val, err := strconv.Atoi(target)
	if err != nil {
		return 0, "", fmt.Errorf("invalid target '%s': must be a PID or Port number", target)
	}

	exists, _ := process.PidExists(int32(val))
	pidFromPort, errPort := monitor.FindPidByPort(val)
	isPort := errPort == nil

	if !exists && !isPort {
		return 0, "", fmt.Errorf("no process found with PID %d or listening on port %d", val, val)
	}

	if exists && isPort {
		if int32(val) == pidFromPort {
			return int32(val), fmt.Sprintf("PID %d (listening on port %d)", val, val), nil
		}
		return 0, "", fmt.Errorf("ambiguous target: %d is both a running PID and a monitored Port (PID %d). Use specific PID to be safe.", val, pidFromPort)
	} else if exists {
		return int32(val), fmt.Sprintf("PID %d", val), nil
	} else {
		return pidFromPort, fmt.Sprintf("PID %d (listening on port %d)", pidFromPort, val), nil
	}
}
