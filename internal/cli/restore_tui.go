package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type restoreModel struct {
	files    []string
	selected map[int]struct{}
	cursor   int
}

func initialRestoreModel(files []string) restoreModel {
	return restoreModel{
		files:    files,
		selected: make(map[int]struct{}),
		cursor:   0,
	}
}

func (m restoreModel) Init() tea.Cmd {
	return nil
}

func (m restoreModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m restoreModel) View() string {
	s := strings.Builder{}
	s.WriteString("Select files to restore:\n\n")

	for i, file := range m.files {
		res := " "
		if _, ok := m.selected[i]; ok {
			res = "x"
		}

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		style := lipgloss.NewStyle()
		if m.cursor == i {
			style = selectedStyle
		}

		s.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, res, style.Render(file)))
	}

	s.WriteString("\n(Press [space] to select, [enter] to confirm)\n")
	return s.String()
}

// InteractiveRestoreSelection launches the TUI to select files.
func InteractiveRestoreSelection(files []string) ([]string, error) {
	p := tea.NewProgram(initialRestoreModel(files))
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	model := m.(restoreModel)
	var selected []string
	for i := range model.selected {
		selected = append(selected, model.files[i])
	}

	return selected, nil
}
