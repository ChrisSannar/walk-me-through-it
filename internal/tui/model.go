package tui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chrissannar/walk-me-through-it/internal/navigator"
	"github.com/chrissannar/walk-me-through-it/internal/walker"
)

// AppState represents the current state of the application
type AppState int

const (
	StateSelecting AppState = iota
	StateEnteringPath
	StateLoading
	StateViewing
)

// walkthroughItem represents a walkthrough file for the list
type walkthroughItem struct {
	path string
	name string
}

func (w walkthroughItem) FilterValue() string { return w.name }
func (w walkthroughItem) Title() string       { return w.name }
func (w walkthroughItem) Description() string { return w.path }

// Model represents the TUI state
type Model struct {
	state        AppState
	width        int
	height       int
	list         list.Model
	textInput    textinput.Model
	navigator    *navigator.Navigator
	walker       *walker.Walker
	rootPath     string
	selectedFile string
	err          error
}

// NewModel creates a new TUI model
func NewModel() Model {
	// Get current working directory
	cwd, _ := os.Getwd()

	// Create list for walkthrough selection
	listItems := []list.Item{}
	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a walkthrough file"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)

	// Create text input for manual path entry
	ti := textinput.New()
	ti.Placeholder = "Enter path to walkthrough file..."
	ti.Focus()

	return Model{
		state:     StateSelecting,
		rootPath:  cwd,
		list:      l,
		textInput: ti,
		walker:    walker.NewWalker(cwd),
	}
}

// Init initializes the TUI model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.findWalkthroughFiles(),
		tea.EnterAltScreen,
	)
}

// findWalkthroughFiles searches for walkthrough files
func (m Model) findWalkthroughFiles() tea.Cmd {
	return func() tea.Msg {
		files, err := walker.FindWalkthroughFiles(m.rootPath)
		if err != nil {
			return errMsg{err}
		}
		return walkthroughFilesMsg{files}
	}
}

// Messages
type walkthroughFilesMsg struct {
	files []string
}

type errMsg struct {
	err error
}

type walkthroughLoadedMsg struct {
	nav *navigator.Navigator
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.state == StateSelecting && len(m.list.Items()) == 0 {
				m.state = StateEnteringPath
				return m, nil
			}
		}

		// Handle state-specific key events
		switch m.state {
		case StateSelecting:
			if msg.String() == "enter" {
				if item, ok := m.list.SelectedItem().(walkthroughItem); ok {
					m.selectedFile = item.path
					return m, m.loadWalkthrough(item.path)
				}
			}
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd

		case StateEnteringPath:
			if msg.String() == "enter" {
				path := m.textInput.Value()
				if path != "" {
					m.selectedFile = path
					return m, m.loadWalkthrough(path)
				}
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case walkthroughFilesMsg:
		items := make([]list.Item, len(msg.files))
		for i, file := range msg.files {
			items[i] = walkthroughItem{
				path: file,
				name: filepath.Base(file),
			}
		}
		m.list.SetItems(items)
		if len(items) == 0 {
			m.state = StateEnteringPath
		}
		return m, nil

	case walkthroughLoadedMsg:
		m.navigator = msg.nav
		m.state = StateViewing
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

// loadWalkthrough loads a walkthrough file
func (m Model) loadWalkthrough(path string) tea.Cmd {
	return func() tea.Msg {
		nav := navigator.NewNavigator()
		if err := nav.LoadWalkthrough(path); err != nil {
			return errMsg{err}
		}
		return walkthroughLoadedMsg{nav}
	}
}

// View renders the TUI
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	switch m.state {
	case StateSelecting:
		if len(m.list.Items()) == 0 {
			return m.renderCentered("Searching for walkthrough files...\n\nPress Tab to enter a path manually, or q to quit")
		}
		return m.list.View()

	case StateEnteringPath:
		return m.renderCentered(fmt.Sprintf(
			"No walkthrough files found in current directory.\n\n%s\n\nPress Enter to load, or q to quit",
			m.textInput.View(),
		))

	case StateLoading:
		return m.renderCentered("Loading walkthrough...")

	case StateViewing:
		if m.navigator == nil {
			return m.renderCentered("Error: Navigator not initialized")
		}

		wt := m.navigator.Walkthrough()
		if wt == nil {
			return m.renderCentered("Error: No walkthrough loaded")
		}

		current, total := m.navigator.Progress()
		step, err := m.navigator.CurrentStep()
		if err != nil {
			return m.renderCentered(fmt.Sprintf("Error: %v", err))
		}

		content := fmt.Sprintf(
			"Walkthrough: %s\nStep %d of %d: %s\n\nFile: %s (lines %d-%d)\n\n%s\n\nPress Tab for next, Shift+Tab for previous, q to quit",
			wt.Title,
			current,
			total,
			step.Title,
			step.File,
			step.LineStart,
			step.LineEnd,
			step.Description,
		)
		return m.renderTop(content)
	}

	return ""
}

// renderCentered centers content on screen
func (m Model) renderCentered(content string) string {
	lines := strings.Split(content, "\n")
	paddingTop := (m.height - len(lines)) / 2
	if paddingTop < 0 {
		paddingTop = 0
	}

	var sb strings.Builder
	for i := 0; i < paddingTop; i++ {
		sb.WriteString("\n")
	}
	sb.WriteString(content)

	// Fill remaining space
	currentLines := paddingTop + len(lines)
	for i := currentLines; i < m.height; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderTop renders content at the top of the screen
func (m Model) renderTop(content string) string {
	lines := strings.Split(content, "\n")
	var sb strings.Builder
	sb.WriteString(content)

	// Fill remaining space
	for i := len(lines); i < m.height; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
}

// Helper to satisfy interface
var _ io.Writer = (*strings.Builder)(nil)
