package tui

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
		m.highlightStart = msg.highlightStart
		m.highlightEnd = msg.highlightEnd
		m.displayStart = msg.displayStart
		m.displayEnd = msg.displayEnd
		m.lastAuditMsg = m.walker.GetLastAuditMessage()
		m.state = StateViewing
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	return m, nil
}
