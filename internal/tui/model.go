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
	"github.com/chrissannar/walk-me-through-it/internal/highlighter"
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
	StateLoadingStep
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
	lastAuditMsg string
	err          error
	styles       *Styles
	highlighter  *highlighter.Highlighter
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
		// Ensure .wmti/ directory exists
		wmtiDir := filepath.Join(m.rootPath, ".wmti")
		if err := os.MkdirAll(wmtiDir, 0755); err != nil {
			return errMsg{fmt.Errorf("failed to create .wmti directory: %w", err)}
		}

		// Create self-referential template if it doesn't exist
		selfTemplatePath := filepath.Join(wmtiDir, "self.wmti.json")
		if _, err := os.Stat(selfTemplatePath); os.IsNotExist(err) {
			if err := m.createSelfTemplate(selfTemplatePath); err != nil {
				// Don't fail if template creation fails, just log it
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
  "version": "1.0.0",
  "steps": [
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
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state == StateViewing {
				// Go back to file selection without rescanning
				m.navigator = nil
				m.fileContent = []string{}
				m.selectedFile = ""
				if len(m.list.Items()) == 0 {
					m.state = StateEnteringPath
				} else {
					m.state = StateSelecting
				}
				return m, nil
			}
			return m, tea.Quit
		case "tab":
			if m.state == StateSelecting && len(m.list.Items()) == 0 {
				m.state = StateEnteringPath
				return m, nil
			}
			if m.state == StateViewing && m.navigator != nil {
				_, err := m.navigator.NextCycle()
				if err == nil {
					m.state = StateLoadingStep
					m.fileContent = []string{}
					return m, tea.Batch(tea.ClearScreen, m.loadCurrentStepFile())
				}
			}
			return m, nil
		case "shift+tab":
			if m.state == StateViewing && m.navigator != nil {
				_, err := m.navigator.PreviousCycle()
				if err == nil {
					m.state = StateLoadingStep
					m.fileContent = []string{}
					return m, tea.Batch(tea.ClearScreen, m.loadCurrentStepFile())
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

			// Handle looping navigation
			switch msg.String() {
			case "down", "j":
				if m.list.Index() == len(m.list.Items())-1 {
					// At last item, loop to first
					m.list.Select(0)
					return m, nil
				}
			case "up", "k":
				if m.list.Index() == 0 {
					// At first item, loop to last
					m.list.Select(len(m.list.Items()) - 1)
					return m, nil
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
		m.lastAuditMsg = m.walker.GetLastAuditMessage()
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

	case StateLoadingStep:
		if m.navigator == nil {
			return m.renderCentered("Error: Navigator not initialized")
		}
		step, err := m.navigator.CurrentStep()
		if err != nil {
			return m.renderCentered(fmt.Sprintf("Error: %v", err))
		}
		return m.renderCentered(fmt.Sprintf("Loading step...\n\n📄 %s\n📝 %s", step.Title, step.Description))

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

		// Calculate widths first (needed for content)
		sidebarWidth := int(float64(m.width) * 0.30)
		if sidebarWidth < 25 {
			sidebarWidth = 25
		}
		codeWidth := m.width - sidebarWidth - 2 // -2 for separator

		// STEP 1: Build header and footer first to get their actual heights
		// HEADER - styled with background and borders
		headerContent := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				"📄 ",
				m.styles.HeaderFileStyle.Render(step.File),
			),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.styles.HeaderLineLabelStyle.Render("lines "),
				m.styles.HeaderLineStyle.Render(fmt.Sprintf("%d-%d", step.LineStart, step.LineEnd)),
			),
		)
		headerInner := m.styles.HeaderStyle.Render(headerContent)
		// Wrap with full-width background using border
		header := lipgloss.JoinVertical(
			lipgloss.Left,
			headerInner,
			m.styles.HeaderBorderStyle.Render(strings.Repeat("─", m.width)),
		)

		// FOOTER - styled with keybinding boxes and color-coded content
		keybindings := lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.styles.RenderKeybinding("Tab", "Next"),
			" ",
			m.styles.RenderKeybinding("Shift+Tab", "Previous"),
			" ",
			m.styles.RenderKeybinding("q", "Quit"),
		)

		var footerContent string
		if m.lastAuditMsg != "" {
			// Truncate audit message to fit
			auditText := m.lastAuditMsg
			if len(auditText) > m.width-4 {
				auditText = auditText[:m.width-7] + "..."
			}
			auditStyled := m.styles.AuditMsgStyle.Render("🔒 " + auditText)
			footerContent = lipgloss.JoinVertical(
				lipgloss.Left,
				auditStyled,
				" ",
				keybindings,
			)
		} else {
			footerContent = keybindings
		}
		footerInner := m.styles.FooterStyle.Render(footerContent)
		// Wrap with full-width border
		footer := lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.FooterBorderStyle.Render(strings.Repeat("─", m.width)),
			footerInner,
		)

		// STEP 2: Calculate content height based on actual header/footer heights
		headerHeight := lipgloss.Height(header)
		footerHeight := lipgloss.Height(footer)
		contentHeight := m.height - headerHeight - footerHeight - 2
		if contentHeight < 3 {
			contentHeight = 3 // Minimum content area
		}

		// STEP 3: Update styles and render content with calculated height
		m.styles.SetDimensions(m.width, codeWidth, sidebarWidth, contentHeight)

		// CONTENT AREA: Two columns side by side
		// Left: Code window with syntax highlighting
		// Prefix is: "999 │ " = 4 digits + space + │ + space = 7 chars
		linePrefixWidth := 7
		maxCodeLineWidth := codeWidth - linePrefixWidth - 1 // -1 for safety buffer

		// Highlight code content
		var codeBlock string
		if len(m.fileContent) > 0 {
			highlightedLines, err := m.highlighter.HighlightFile(step.File, m.fileContent, step.LineStart)
			if err != nil || len(highlightedLines) == 0 {
				// Fallback to plain text if highlighting fails
				var codeContent strings.Builder
				for i, line := range m.fileContent {
					if i >= contentHeight {
						break
					}
					lineNum := step.LineStart + i
					truncatedLine := truncate(line, maxCodeLineWidth)
					codeContent.WriteString(fmt.Sprintf("%4d │ %s\n", lineNum, truncatedLine))
				}
				codeBlock = m.styles.CodeStyle.Render(codeContent.String())
			} else {
				// Render highlighted lines with proper width tracking
				var codeContent strings.Builder
				for i, hlLine := range highlightedLines {
					if i >= contentHeight {
						break
					}
					// Render line number
					lineNumStr := fmt.Sprintf("%4d │ ", hlLine.LineNumber)
					codeContent.WriteString(lineNumStr)

					// Track remaining width for this line
					remainingWidth := maxCodeLineWidth

					// Render highlighted tokens with width tracking
					for _, token := range hlLine.Tokens {
						tokenWidth := displayWidth(token.Text)

						if tokenWidth > remainingWidth {
							// Token doesn't fit - truncate it cleanly (no ellipsis in code view)
							if remainingWidth > 0 {
								truncatedText := truncateClean(token.Text, remainingWidth)
								codeContent.WriteString(token.Style.Render(truncatedText))
							}
							break // Stop rendering more tokens
						}

						// Token fits - render it
						codeContent.WriteString(token.Style.Render(token.Text))
						remainingWidth -= tokenWidth

						// If no width left, stop
						if remainingWidth <= 0 {
							break
						}
					}
					codeContent.WriteString("\n")
				}
				codeBlock = m.styles.CodeStyle.Render(codeContent.String())
			}
		} else {
			codeBlock = m.styles.CodeStyle.Render("Loading file content...\n")
		}

		// Right: Sidebar with styled content
		stepHeader := m.styles.StepHeaderStyle.Render(fmt.Sprintf("📋 Step %d of %d", current, total))
		stepTitle := m.styles.StepTitleStyle.Render(step.Title)
		stepDesc := m.styles.StepDescriptionStyle.Render(step.Description)
		sidebarContent := lipgloss.JoinVertical(
			lipgloss.Left,
			stepHeader,
			stepTitle,
			stepDesc,
		)
		sidebarBlock := m.styles.SidebarStyle.Render(sidebarContent)

		// Vertical separator
		separator := m.styles.SeparatorStyle.Height(contentHeight).Render(strings.Repeat("│\n", contentHeight))

		// Join content horizontally
		contentRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			codeBlock,
			separator,
			sidebarBlock,
		)

		// STEP 4: Join all sections vertically
		return lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			contentRow,
			footer,
		)
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

// displayWidth returns the visual width of a string (handles Unicode properly)
func displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r == '\t' {
			width += 4 // Tab is typically 4 spaces
		} else if r < 32 {
			// Control characters - ignore or count as 0
		} else {
			width++
		}
	}
	return width
}

// padRight pads a string with spaces to the right to reach the desired visual width
func padRight(s string, width int) string {
	currentWidth := displayWidth(s)
	for currentWidth < width {
		s += " "
		currentWidth++
	}
	return s
}

// truncate truncates a string to fit within the maximum visual width (with ellipsis)
func truncate(s string, maxWidth int) string {
	if displayWidth(s) <= maxWidth {
		return s
	}

	// Need to truncate
	result := ""
	width := 0
	for _, r := range s {
		runeWidth := 1
		if r == '\t' {
			runeWidth = 4
		}

		if width+runeWidth > maxWidth-3 {
			// Don't have room for this rune, add ellipsis instead
			return result + "..."
		}

		result += string(r)
		width += runeWidth
	}
	return result
}

// truncateClean truncates a string to fit within the maximum visual width (no ellipsis)
// Used for code display where we want clean truncation at the edge
func truncateClean(s string, maxWidth int) string {
	if displayWidth(s) <= maxWidth {
		return s
	}

	// Need to truncate without ellipsis
	result := ""
	width := 0
	for _, r := range s {
		runeWidth := 1
		if r == '\t' {
			runeWidth = 4
		}

		if width+runeWidth > maxWidth {
			// Don't have room for this rune, stop here
			return result
		}

		result += string(r)
		width += runeWidth
	}
	return result
}
