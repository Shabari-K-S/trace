package cli

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"trace/internal/diff"
	"trace/internal/monitor"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")) // Pink
	subTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))  // Blue
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))             // Green
	warnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))            // Orange
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))            // Grey
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))            // Pink
)

type model struct {
	interval time.Duration
	envDiff  diff.EnvDiff
	fileDiff diff.FileDiff
	procs    []monitor.ProcessInfo
	err      error
	cursor   int
	message  string // Status message
}

type tickMsg time.Time
type statusMsg struct {
	envDiff  diff.EnvDiff
	fileDiff diff.FileDiff
	procs    []monitor.ProcessInfo
	err      error
}

// clearMsg clears the status message
type clearMsg struct{}

func Watch(interval time.Duration) error {
	p := tea.NewProgram(initialModel(interval))
	_, err := p.Run()
	return err
}

func initialModel(interval time.Duration) model {
	return model{
		interval: interval,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		checkStatusCmd,
		tickCmd(m.interval),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.procs)-1 {
				m.cursor++
			}
		case "enter", "x", "delete":
			if len(m.procs) > 0 && m.cursor >= 0 && m.cursor < len(m.procs) {
				p := m.procs[m.cursor]
				err := KillPID(p.PID)
				if err != nil {
					m.message = fmt.Sprintf("Error killing %s: %v", p.Name, err)
				} else {
					m.message = fmt.Sprintf("Killed %s (PID %d)", p.Name, p.PID)
				}
				// Clear message after 3 seconds
				return m, tea.Batch(
					checkStatusCmd,
					tickCmd(100*time.Millisecond), // Rapid refresh
					clearMessageCmd(),
				)
			}
		}

	case tickMsg:
		return m, tea.Batch(
			checkStatusCmd,
			tickCmd(m.interval),
		)

	case statusMsg:
		m.envDiff = msg.envDiff
		m.fileDiff = msg.fileDiff
		m.procs = msg.procs
		m.err = msg.err

		if m.cursor >= len(m.procs) {
			m.cursor = len(m.procs) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}

	case clearMsg:
		m.message = ""
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("👀 Trace Watch"))
	s.WriteString(" " + dimStyle.Render("(Arrows to nav, Enter/'x' to kill, 'q' to quit)"))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(warnStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		return s.String()
	}

	// Message Bar
	if m.message != "" {
		if strings.HasPrefix(m.message, "Error") {
			s.WriteString(warnStyle.Render(m.message))
		} else {
			s.WriteString(successStyle.Render(m.message))
		}
		s.WriteString("\n\n")
	}

	clean := m.envDiff.IsEmpty() && m.fileDiff.IsEmpty()
	if clean {
		s.WriteString(successStyle.Render("✨ Environment Clean"))
	} else {
		s.WriteString(warnStyle.Render("⚠️  Changes Not Committed:"))
		s.WriteString("\n")
		renderDiffs(&s, m.fileDiff, m.envDiff)
	}

	s.WriteString("\n\n")

	if len(m.procs) > 0 {
		s.WriteString(subTitleStyle.Render(fmt.Sprintf("🚀 Active Processes (%d)", len(m.procs))))
		s.WriteString("\n")
		for i, p := range m.procs {
			cursor := "  "
			style := dimStyle
			if i == m.cursor {
				cursor = "> "
				style = selectedStyle
			}

			ports := ""
			if len(p.Ports) > 0 {
				strPorts := make([]string, len(p.Ports))
				for i, port := range p.Ports {
					strPorts[i] = fmt.Sprintf("%d", port)
				}
				ports = dimStyle.Render(fmt.Sprintf(":%s", strings.Join(strPorts, ",")))
			}

			line := fmt.Sprintf("%s• [%d] %s %s", cursor, p.PID, p.Name, ports)
			s.WriteString(style.Render(line) + "\n")
		}
	} else {
		s.WriteString(dimStyle.Render("No active processes detected."))
	}

	return s.String()
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func clearMessageCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return clearMsg{}
	})
}

func checkStatusCmd() tea.Msg {
	envDiff, fileDiff, procs, err := GetStatus()
	return statusMsg{
		envDiff:  envDiff,
		fileDiff: fileDiff,
		procs:    procs,
		err:      err,
	}
}

func renderDiffs(s *strings.Builder, files diff.FileDiff, env diff.EnvDiff) {
	for _, p := range files.Added {
		s.WriteString(fmt.Sprintf("  %s %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("+"), p))
	}
	for _, p := range files.Removed {
		s.WriteString(fmt.Sprintf("  %s %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render("-"), p))
	}
	for _, p := range files.Modified {
		s.WriteString(fmt.Sprintf("  %s %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("~"), p))
	}

	for _, k := range env.Added {
		s.WriteString(fmt.Sprintf("  %s %s (env)\n", lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("+"), k))
	}
	for _, k := range env.Removed {
		s.WriteString(fmt.Sprintf("  %s %s (env)\n", lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render("-"), k))
	}
	for _, k := range env.Changed {
		s.WriteString(fmt.Sprintf("  %s %s (env)\n", lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("~"), k))
	}
}
