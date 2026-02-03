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
	fileContent  []string
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

type fileContentMsg struct {
	lines []string
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
			if m.state == StateViewing && m.navigator != nil {
				_, err := m.navigator.NextCycle()
				if err == nil {
					return m, m.loadCurrentStepFile()
				}
			}
			return m, nil
		case "shift+tab":
			if m.state == StateViewing && m.navigator != nil {
				_, err := m.navigator.PreviousCycle()
				if err == nil {
					return m, m.loadCurrentStepFile()
				}
			}
			return m, nil
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
		// Load file content for the first step
		return m, m.loadCurrentStepFile()

	case fileContentMsg:
		m.fileContent = msg.lines
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

// loadCurrentStepFile loads the file content for the current step
func (m Model) loadCurrentStepFile() tea.Cmd {
	return func() tea.Msg {
		if m.navigator == nil {
			return errMsg{fmt.Errorf("navigator not initialized")}
		}

		step, err := m.navigator.CurrentStep()
		if err != nil {
			return errMsg{err}
		}

		lines, err := m.walker.ReadFileLines(step.File, step.LineStart, step.LineEnd)
		if err != nil {
			return errMsg{fmt.Errorf("failed to read file %s: %w", step.File, err)}
		}

		return fileContentMsg{lines}
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
		// Ensure list fills the entire screen to prevent artifacts
		listView := m.list.View()
		lines := strings.Split(listView, "\n")
		var result strings.Builder
		result.WriteString(listView)
		for i := len(lines); i < m.height; i++ {
			result.WriteString("\n")
		}
		return result.String()

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

		// Build the file content section (top portion)
		var fileSection strings.Builder
		fileSection.WriteString(fmt.Sprintf("📄 %s (lines %d-%d)\n", step.File, step.LineStart, step.LineEnd))
		fileSection.WriteString(strings.Repeat("─", m.width) + "\n")

		if len(m.fileContent) > 0 {
			for i, line := range m.fileContent {
				lineNum := step.LineStart + i
				fileSection.WriteString(fmt.Sprintf("%4d │ %s\n", lineNum, line))
			}
		} else {
			fileSection.WriteString("Loading file content...\n")
		}

		// Build the instructions section (bottom window)
		var instructionSection strings.Builder
		instructionSection.WriteString(strings.Repeat("─", m.width) + "\n")
		instructionSection.WriteString(fmt.Sprintf("📋 Step %d of %d: %s\n", current, total, step.Title))
		instructionSection.WriteString(fmt.Sprintf("📝 %s\n", step.Description))
		instructionSection.WriteString(fmt.Sprintf("⌨️  Tab: Next | Shift+Tab: Previous | q: Quit"))

		// Combine sections with proper spacing
		fileLines := strings.Split(fileSection.String(), "\n")
		instructionLines := strings.Split(instructionSection.String(), "\n")

		// Calculate available space for file content
		instructionHeight := len(instructionLines)
		availableHeight := m.height - instructionHeight - 1 // -1 for spacing

		// Truncate file content if needed
		if len(fileLines) > availableHeight {
			fileLines = fileLines[:availableHeight-1]
			fileLines = append(fileLines, "... (content truncated)")
		}

		// Combine everything and ensure we fill the entire screen
		var result strings.Builder
		result.WriteString(strings.Join(fileLines, "\n"))

		// Fill remaining space between file content and instructions
		currentLineCount := len(fileLines)
		for i := currentLineCount; i < availableHeight; i++ {
			result.WriteString("\n")
		}

		// Add instruction section
		result.WriteString(strings.Join(instructionLines, "\n"))

		// Ensure we fill to the bottom of the screen
		totalLines := len(fileLines) + (availableHeight - len(fileLines)) + len(instructionLines)
		for i := totalLines; i < m.height; i++ {
			result.WriteString("\n")
		}

		return result.String()
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
