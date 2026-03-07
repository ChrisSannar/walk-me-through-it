package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chrissannar/walk-me-through-it/internal/api"
	"github.com/chrissannar/walk-me-through-it/internal/config"
	"github.com/chrissannar/walk-me-through-it/internal/highlighter"
	"github.com/chrissannar/walk-me-through-it/internal/navigator"
	"github.com/chrissannar/walk-me-through-it/internal/walker"
)

// NewModel creates a new TUI model
func NewModel() Model {
	cwd, _ := os.Getwd()

	listItems := []list.Item{}
	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a walkthrough file"
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)
	l.Styles.StatusBar = lipgloss.NewStyle().Foreground(lipgloss.Color("transparent"))

	ti := textinput.New()
	ti.Placeholder = "Enter path to walkthrough file..."
	ti.CharLimit = 512
	ti.Focus()

	modelTi := textinput.New()
	modelTi.Placeholder = "Enter model name..."
	modelTi.CharLimit = 64
	modelTi.Focus()

	keyTi := textinput.New()
	keyTi.Placeholder = "Enter API key..."
	keyTi.CharLimit = 512
	keyTi.EchoMode = textinput.EchoPassword
	keyTi.Focus()

	providerTi := textinput.New()
	providerTi.Placeholder = "Enter API key..."
	providerTi.CharLimit = 512
	providerTi.EchoMode = textinput.EchoPassword
	providerTi.Focus()

	providerList := []string{}
	for _, p := range api.Providers {
		providerList = append(providerList, p.DisplayName)
	}

	cfg, _ := loadConfig()
	modelNames := []string{}
	for _, m := range cfg.Models {
		modelNames = append(modelNames, m.Name)
	}
	modelNames = append(modelNames, "+ Add new model")

	return Model{
		state:                 StateProviderSelect,
		rootPath:              cwd,
		list:                  l,
		textInput:             ti,
		walker:                walker.NewWalker(cwd),
		styles:                NewStyles(),
		highlighter:           highlighter.NewHighlighter(highlighter.ThemeDracula),
		modelList:             modelNames,
		modelSelectedIndex:    0,
		modelTextInput:        modelTi,
		modelIsAdding:         false,
		providerList:          providerList,
		providerSelectedIndex: 0,
		providerAskingFor:     "",
	}
}

func loadConfig() (*config.Config, error) {
	return config.LoadOrCreate()
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
		wmtiDir := filepath.Join(m.rootPath, ".wmti")
		if err := os.MkdirAll(wmtiDir, 0755); err != nil {
			return errMsg{fmt.Errorf("failed to create .wmti directory: %w", err)}
		}

		files, err := walker.FindWalkthroughFiles(m.rootPath)
		if err != nil {
			return errMsg{err}
		}
		return walkthroughFilesMsg{files}
	}
}

// loadWalkthrough loads a walkthrough file
func (m Model) loadWalkthrough(path string) tea.Cmd {
	return func() tea.Msg {
		nav := navigator.NewNavigatorWithRoot(m.rootPath)
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

		highlightStart := step.LineStart
		highlightEnd := step.LineEnd
		highlightLines := highlightEnd - highlightStart + 1

		// Calculate how many lines to read to fill the screen
		// Account for header (~3 lines) and footer (~3 lines)
		targetLines := m.height - 10
		if targetLines < 10 {
			targetLines = 20 // minimum
		}

		// Expand range to center the highlight
		displayStart := highlightStart
		displayEnd := highlightEnd

		if targetLines > highlightLines {
			extra := targetLines - highlightLines
			before := extra / 2
			after := extra - before

			displayStart = highlightStart - before
			displayEnd = highlightEnd + after

			// Don't go below line 1
			if displayStart < 1 {
				displayStart = 1
				// If we can't go back, show more after
				after = (highlightEnd - displayStart + 1) + extra - highlightLines
				displayEnd = highlightEnd + after
			}
		}

		lines, err := m.walker.ReadFileLines(step.File, displayStart, displayEnd)
		if err != nil {
			return errMsg{fmt.Errorf("failed to read file %s: %w", step.File, err)}
		}

		return fileContentMsg{
			lines:          lines,
			highlightStart: highlightStart,
			highlightEnd:   highlightEnd,
			displayStart:   displayStart,
			displayEnd:     displayEnd,
		}
	}
}
