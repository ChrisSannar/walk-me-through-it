package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the TUI state
type Model struct {
	width  int
	height int
}

// NewModel creates a new TUI model
func NewModel() Model {
	return Model{}
}

// Init initializes the TUI model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.height == 0 {
		return "Walk Me Through It - Press 'q' to quit"
	}
	// Display text at the top, fill rest with empty lines
	content := "Walk Me Through It - Press 'q' to quit"
	emptyLines := m.height - 1
	return content + strings.Repeat("\n", emptyLines)
}
