package monitor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// ProcessInfo holds basic information about a process.
type ProcessInfo struct {
	PID   int32
	Name  string
	Ports []int
}

// GetProjectProcesses finds processes whose CWD matches the project root
// or is a subdirectory of the project root.
func GetProjectProcesses(root string) ([]ProcessInfo, error) {
	// Ensure root is cleaned and absolute if possible (caller should handle abs path)
	root = filepath.Clean(root)

	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("list processes: %w", err)
	}

	var projectProcs []ProcessInfo

	for _, p := range procs {
		cwd, err := p.Cwd()
		if err != nil {
			// Permission denied or process died - skip
			continue
		}

		// Clean the CWD
		cwd = filepath.Clean(cwd)

		// Check if CWD is within the project root
		// We use HasPrefix for simple containment check
		// To be precise, we need to ensure it's a directory match
		if strings.HasPrefix(cwd, root) {
			// Double check it's not a partial match like /foo/bar vs /foo/bar_baz
			rel, err := filepath.Rel(root, cwd)
			if err != nil {
				continue
			}
			if rel == "." || !strings.HasPrefix(rel, "..") {
				name, _ := p.Name()
				ports, _ := GetProcessPorts(p.Pid)

				projectProcs = append(projectProcs, ProcessInfo{
					PID:   p.Pid,
					Name:  name,
					Ports: ports,
				})
			}
		}
	}

	return projectProcs, nil
}

// GetProcessPorts returns the listening ports for a given PID.
func GetProcessPorts(pid int32) ([]int, error) {
	connections, err := net.ConnectionsPid("tcp", pid)
	if err != nil {
		return nil, err
	}

	var ports []int
	seen := make(map[int]bool)

	for _, conn := range connections {
		if conn.Status == "LISTEN" {
			port := int(conn.Laddr.Port)
			if !seen[port] {
				ports = append(ports, port)
				seen[port] = true
			}
		}
	}

	return ports, nil
}

// FindPidByPort returns the PID of the process listening on the given port.
func FindPidByPort(port int) (int32, error) {
	connections, err := net.Connections("tcp")
	if err != nil {
		return 0, fmt.Errorf("list connections: %w", err)
	}

	for _, conn := range connections {
		if conn.Status == "LISTEN" && int(conn.Laddr.Port) == port {
			return conn.Pid, nil
		}
	}

	return 0, fmt.Errorf("no process found listening on port %d", port)
}
