package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type diffModel struct {
	content  string
	ready    bool
	viewport viewport.Model
}

func initialDiffModel(content string) diffModel {
	return diffModel{
		content: content,
	}
}

func (m diffModel) Init() tea.Cmd {
	return nil
}

func (m diffModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m diffModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	title := titleStyle.Render("Trace Diff Viewer")
	info := dimStyle.Render("(scroll with j/k/arrows, q to quit)")
	header := fmt.Sprintf("%s\n%s\n", title, info)
	footer := fmt.Sprintf("\n%s", dimStyle.Render("End of diff"))

	return fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer)
}

// ShowDiffTUI displays the diff content in a scrollable viewport.
func ShowDiffTUI(content string) error {
	p := tea.NewProgram(
		initialDiffModel(content),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
