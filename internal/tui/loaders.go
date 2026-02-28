package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065")).
		Padding(0, 1)

	ti := textinput.New()
	ti.Placeholder = "Enter path to walkthrough file..."
	ti.Focus()

	return Model{
		state:       StateSelecting,
		rootPath:    cwd,
		list:        l,
		textInput:   ti,
		walker:      walker.NewWalker(cwd),
		styles:      NewStyles(),
		highlighter: highlighter.NewHighlighter(highlighter.ThemeDracula),
	}
}

// Init initializes the TUI model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.findWalkthroughFiles(),
		tea.EnterAltScreen,
	)
}

// findWalkthroughFiles searches for walkthrough files and creates directory/template if needed
func (m Model) findWalkthroughFiles() tea.Cmd {
	return func() tea.Msg {
		wmtiDir := filepath.Join(m.rootPath, ".wmti")
		if err := os.MkdirAll(wmtiDir, 0755); err != nil {
			return errMsg{fmt.Errorf("failed to create .wmti directory: %w", err)}
		}

		selfTemplatePath := filepath.Join(wmtiDir, "self.wmti.json")
		if _, err := os.Stat(selfTemplatePath); os.IsNotExist(err) {
			if err := m.createSelfTemplate(selfTemplatePath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not create self template: %v\n", err)
			}
		}

		files, err := walker.FindWalkthroughFiles(m.rootPath)
		if err != nil {
			return errMsg{err}
		}
		return walkthroughFilesMsg{files}
	}
}

// createSelfTemplate creates a self-referential walkthrough template
func (m Model) createSelfTemplate(path string) error {
	template := `{
  "title": "Understanding Walkthrough Files",
  "description": "A walkthrough that explains the structure of .wmti.json files using this file as an example",
  "version": "0",
  "1.0.steps": [
    {
      "id": 1,
      "title": "Walkthrough Structure",
      "description": "Every .wmti.json file starts with metadata: title, description, and version. The title appears in the file selection list.",
      "file": ".wmti/self.wmti.json",
      "line_start": 1,
      "line_end": 5,
      "action": "read"
    },
    {
      "id": 2,
      "title": "Steps Array",
      "description": "The 'steps' array contains all navigation points. Each step has an ID, title, description, file path, and line range.",
      "file": ".wmti/self.wmti.json",
      "line_start": 6,
      "line_end": 18,
      "action": "read"
    },
    {
      "id": 3,
      "title": "Step Properties",
      "description": "Each step specifies which file to open and which lines to highlight. The 'action' field determines what to do (currently only 'read' is supported).",
      "file": ".wmti/self.wmti.json",
      "line_start": 10,
      "line_end": 17,
      "action": "read"
    }
  ]
}
`
	return os.WriteFile(path, []byte(template), 0644)
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
