package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chrissannar/walk-me-through-it/internal/api"
	"github.com/chrissannar/walk-me-through-it/internal/config"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state == StateModelSelect && m.modelIsAdding {
				// Type 'q' into the input field
				currentValue := m.modelTextInput.Value()
				m.modelTextInput.SetValue(currentValue + "q")
				return m, nil
			}
			if m.state == StateProviderSelect && m.providerAskingFor == "key" {
				// Type 'q' into the API key input
				currentValue := m.modelTextInput.Value()
				m.modelTextInput.SetValue(currentValue + "q")
				return m, nil
			}
			if m.state == StateModelSelect {
				// Go back to provider selection
				m.state = StateProviderSelect
				return m, nil
			}
			if m.state == StateProviderSelect {
				// Go back to walkthrough selection
				m.state = StateSelecting
				return m, nil
			}
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
			if m.state == StateModelSelect && m.modelIsAdding {
				m.modelTextInput.Reset()
				m.modelTextInput.Placeholder = "Enter model name..."
				m.modelIsAdding = false
				m.modelAddingName = ""
				m.modelAskingFor = ""
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
					m.lastAuditMsg = ""
					return m, m.loadCurrentStepFile()
				}
			}
			return m, nil
		case "shift+tab":
			if m.state == StateViewing && m.navigator != nil {
				_, err := m.navigator.PreviousCycle()
				if err == nil {
					m.state = StateLoadingStep
					m.fileContent = []string{}
					m.lastAuditMsg = ""
					return m, m.loadCurrentStepFile()
				}
			}
			return m, nil
		}

		// Handle state-specific key events
		switch m.state {
		case StateSelecting:
			// 'm' key to open provider selection
			if msg.String() == "m" {
				m.state = StateProviderSelect
				m.providerSelectedIndex = 0
				return m, nil
			}

			// 'n' key to create new walkthrough
			if msg.String() == "n" {
				m.state = StateNewWalkthrough
				m.newWalkthroughInput.Reset()
				m.newWalkthroughInput.Placeholder = "What can I walk you through?"
				return m, nil
			}

			if msg.String() == "enter" {
				if m.deleteConfirmPath != "" {
					// Confirm delete walkthrough
					err := os.Remove(m.deleteConfirmPath)
					if err != nil {
						m.err = err
						m.deleteConfirmPath = ""
						return m, nil
					}
					// Remove item from list directly instead of rescanning
					// (rescanning would recreate self.wmti.json template)
					var remainingItems []list.Item
					for _, item := range m.list.Items() {
						if wti, ok := item.(walkthroughItem); ok {
							if wti.path != m.deleteConfirmPath {
								remainingItems = append(remainingItems, item)
							}
						}
					}
					m.list.SetItems(remainingItems)
					m.deleteConfirmPath = ""
					if len(remainingItems) == 0 {
						m.state = StateEnteringPath
					}
					return m, nil
				}
				if item, ok := m.list.SelectedItem().(walkthroughItem); ok {
					m.selectedFile = item.path
					return m, m.loadWalkthrough(item.path)
				}
			}

			// Handle delete key
			if msg.String() == "d" || msg.String() == "Del" {
				if item, ok := m.list.SelectedItem().(walkthroughItem); ok {
					m.deleteConfirmPath = item.path
					return m, nil
				}
			}

			// Handle 'y' to confirm delete, 'n' to cancel
			if m.deleteConfirmPath != "" {
				if msg.String() == "y" {
					// Confirm delete walkthrough
					err := os.Remove(m.deleteConfirmPath)
					if err != nil {
						m.err = err
						m.deleteConfirmPath = ""
						return m, nil
					}
					var remainingItems []list.Item
					for _, item := range m.list.Items() {
						if wti, ok := item.(walkthroughItem); ok {
							if wti.path != m.deleteConfirmPath {
								remainingItems = append(remainingItems, item)
							}
						}
					}
					m.list.SetItems(remainingItems)
					m.deleteConfirmPath = ""
					if len(remainingItems) == 0 {
						m.state = StateEnteringPath
					}
					return m, nil
				}
				if msg.String() == "n" {
					m.deleteConfirmPath = ""
					return m, nil
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
					if _, err := m.walker.SanitizePath(path); err != nil {
						m.err = fmt.Errorf("invalid path: %w", err)
						return m, nil
					}
					m.selectedFile = path
					return m, m.loadWalkthrough(path)
				}
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd

		case StateProviderSelect:
			if m.providerAskingFor == "key" {
				if msg.String() == "enter" {
					value := m.modelTextInput.Value()
					if value != "" {
						providerID := m.selectedProvider
						if err := config.SetAPIKey(providerID, value, providerID); err != nil {
							m.err = err
						}
						m.providerAskingFor = ""
						m.state = StateModelSelect
						m.modelSelectedIndex = 0
						m.updateModelListForProvider(providerID)
					}
					return m, nil
				}
				if msg.String() == "Esc" {
					m.providerAskingFor = ""
					m.modelTextInput.Reset()
					m.modelTextInput.Placeholder = "Enter API key..."
					return m, nil
				}
				var cmd tea.Cmd
				m.modelTextInput, cmd = m.modelTextInput.Update(msg)
				return m, cmd
			}

			switch msg.String() {
			case "enter":
				selectedProviderName := m.providerList[m.providerSelectedIndex]
				m.selectedProvider = ""
				for _, p := range api.Providers {
					if p.DisplayName == selectedProviderName {
						m.selectedProvider = p.ID
						break
					}
				}
				if m.selectedProvider != "" {
					entry, err := config.GetAPIKeyEntry(m.selectedProvider)
					if err == nil && entry.APIKey != "" {
						m.state = StateModelSelect
						m.modelSelectedIndex = 0
						m.updateModelListForProvider(m.selectedProvider)
					} else {
						m.providerAskingFor = "key"
						m.modelTextInput.Reset()
						m.modelTextInput.Placeholder = "Enter API key for " + selectedProviderName + "..."
					}
				}
				return m, nil
			case "up", "k":
				if m.providerSelectedIndex > 0 {
					m.providerSelectedIndex--
				}
				return m, nil
			case "down", "j":
				if m.providerSelectedIndex < len(m.providerList)-1 {
					m.providerSelectedIndex++
				}
				return m, nil
			case "t":
				return m, m.testConnection()
			case "Esc":
				m.state = StateSelecting
				return m, nil
			}

		case StateModelSelect:
			if m.modelIsAdding {
				if msg.String() == "enter" {
					value := m.modelTextInput.Value()
					if value != "" {
						if m.modelAskingFor == "name" {
							for _, existing := range m.modelList {
								if existing == value {
									m.modelError = fmt.Sprintf("model %q already exists", value)
									return m, nil
								}
							}
							m.modelError = ""
							m.modelAddingName = value
							m.modelAskingFor = "key"
							m.modelTextInput.Reset()
							m.modelTextInput.Placeholder = "Enter API key..."
						} else {
							apiKey := value
							modelName := m.modelAddingName
							if err := config.SetAPIKey(modelName, apiKey, m.selectedProvider); err != nil {
								m.err = err
							}
							m.modelList = append(m.modelList[:len(m.modelList)-1], m.modelAddingName, "+ Add new model")
							m.modelSelectedIndex = len(m.modelList) - 2
							cfg, err := loadConfig()
							if err == nil {
								cfg.Models = append(cfg.Models, config.ModelConfig{
									Name:     m.modelAddingName,
									Provider: m.selectedProvider,
								})
								cfg.Save()
							}
							m.modelAddingName = ""
							m.modelAskingFor = ""
							m.modelTextInput.Reset()
							m.modelTextInput.Placeholder = "Enter model name..."
							m.modelIsAdding = false
						}
					}
					return m, nil
				}
				if msg.String() == "Esc" {
					m.modelTextInput.Reset()
					m.modelTextInput.Placeholder = "Enter model name..."
					m.modelIsAdding = false
					m.modelAddingName = ""
					m.modelAskingFor = ""
					m.modelError = ""
					return m, nil
				}
				var cmd tea.Cmd
				m.modelTextInput, cmd = m.modelTextInput.Update(msg)
				return m, cmd
			}

			switch msg.String() {
			case "enter":
				if m.modelSelectedIndex == len(m.modelList)-1 {
					m.modelIsAdding = true
					m.modelAskingFor = "name"
					m.modelError = ""
					m.modelTextInput.Placeholder = "Enter model name..."
				} else {
					m.selectedModel = m.modelList[m.modelSelectedIndex]
					m.state = StateSelecting
				}
				return m, nil
			case "d":
				if m.modelSelectedIndex < len(m.modelList)-1 {
					m.modelDeleteConfirm = true
				}
				return m, nil
			case "y":
				if m.modelDeleteConfirm {
					modelToDelete := m.modelList[m.modelSelectedIndex]
					m.modelList = append(m.modelList[:m.modelSelectedIndex], m.modelList[m.modelSelectedIndex+1:]...)
					cfg, err := loadConfig()
					if err == nil {
						var newModels []config.ModelConfig
						for _, mc := range cfg.Models {
							if mc.Name != modelToDelete {
								newModels = append(newModels, mc)
							}
						}
						cfg.Models = newModels
						cfg.Save()
					}
					if err := config.DeleteAPIKey(modelToDelete); err != nil {
						m.err = err
					}
					m.modelDeleteConfirm = false
				}
				return m, nil
			case "n":
				if m.modelDeleteConfirm {
					m.modelDeleteConfirm = false
				}
				return m, nil
			case "q":
				if m.modelDeleteConfirm {
					m.modelDeleteConfirm = false
					return m, nil
				}
				m.state = StateSelecting
				return m, nil
			case "up", "k":
				m.modelDeleteConfirm = false
				if m.modelSelectedIndex > 0 {
					m.modelSelectedIndex--
				}
				return m, nil
			case "down", "j":
				m.modelDeleteConfirm = false
				if m.modelSelectedIndex < len(m.modelList)-1 {
					m.modelSelectedIndex++
				}
				return m, nil
			case "t":
				return m, m.testConnection()
			case "Esc":
				if m.modelDeleteConfirm {
					m.modelDeleteConfirm = false
					return m, nil
				}
			}

		case StateNewWalkthrough:
			switch msg.String() {
			case "enter":
				prompt := m.newWalkthroughInput.Value()
				if prompt != "" {
					return m, m.generateWalkthrough(prompt)
				}
				return m, nil
			case "Esc":
				m.state = StateSelecting
				m.newWalkthroughInput.Reset()
				return m, nil
			}
			var cmd tea.Cmd
			m.newWalkthroughInput, cmd = m.newWalkthroughInput.Update(msg)
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
		} else {
			m.state = StateSelecting
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
		// If we were loading a step, go back to viewing state to show the error
		if m.state == StateLoadingStep {
			m.state = StateViewing
		}
		return m, nil

	case connectionTestMsg:
		if msg.success {
			m.connectionTestResult = "✓ " + msg.message
		} else {
			m.connectionTestResult = "✗ " + msg.message
		}
		return m, nil

	case walkthroughGeneratedMsg:
		m.state = StateSelecting
		m.newWalkthroughInput.Reset()
		return m, m.saveGeneratedWalkthrough(msg.content)

	case generationErrorMsg:
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

func (m *Model) updateModelListForProvider(providerID string) {
	provider := api.GetProvider(providerID)
	if provider == nil {
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		m.modelList = append([]string{}, provider.DefaultModels...)
		m.modelList = append(m.modelList, "+ Add new model")
		return
	}

	modelNames := []string{}
	for _, mc := range cfg.Models {
		if mc.Provider == providerID {
			modelNames = append(modelNames, mc.Name)
		}
	}

	if len(modelNames) == 0 {
		modelNames = append([]string{}, provider.DefaultModels...)
	}
	modelNames = append(modelNames, "+ Add new model")
	m.modelList = modelNames
}

func (m Model) testConnection() tea.Cmd {
	return func() tea.Msg {
		providerID := m.selectedProvider
		if providerID == "" {
			selectedProviderName := m.providerList[m.providerSelectedIndex]
			for _, p := range api.Providers {
				if p.DisplayName == selectedProviderName {
					providerID = p.ID
					break
				}
			}
		}

		if providerID == "" {
			return connectionTestMsg{success: false, message: "No provider selected"}
		}

		entry, err := config.GetAPIKeyEntry(providerID)
		if err != nil || entry.APIKey == "" {
			return connectionTestMsg{success: false, message: "No API key found for this provider"}
		}

		provider := api.GetProvider(providerID)
		if provider == nil {
			return connectionTestMsg{success: false, message: "Provider not found"}
		}

		client := api.NewClient(provider, entry.APIKey)
		err = client.TestConnection(context.Background())
		if err != nil {
			fullErr := err.Error()
			shortErr := truncateError(fullErr)
			writeErrorLog(providerID, fullErr)
			return connectionTestMsg{success: false, message: shortErr}
		}

		return connectionTestMsg{success: true, message: "Connection successful!"}
	}
}

func (m Model) generateWalkthrough(prompt string) tea.Cmd {
	return func() tea.Msg {
		providerID := m.selectedProvider
		if providerID == "" {
			cfg, _ := loadConfig()
			if cfg != nil {
				selectedConfig := cfg.GetSelectedModelConfig()
				if selectedConfig != nil {
					providerID = selectedConfig.Provider
					m.selectedProvider = providerID
				}
			}
		}

		if providerID == "" {
			return generationErrorMsg{err: fmt.Errorf("no provider or model selected. Press 'm' to select a model first.")}
		}

		entry, err := config.GetAPIKeyEntry(providerID)
		if err != nil || entry.APIKey == "" {
			return generationErrorMsg{err: fmt.Errorf("no API key found for this provider. Press 'm' to add one.")}
		}

		provider := api.GetProvider(providerID)
		if provider == nil {
			return generationErrorMsg{err: fmt.Errorf("provider not found: %s", providerID)}
		}

		client := api.NewClient(provider, entry.APIKey)

		systemPrompt := `You are an expert code walkthrough generator. Create a JSON walkthrough file that explains a codebase to developers.
Respond with ONLY valid JSON, no markdown formatting or explanation.
Use this exact schema:
{
  "title": "Descriptive Title",
  "description": "What this walkthrough covers (1-2 sentences)",
  "version": "1.0.0",
  "steps": [
    {
      "id": 1,
      "title": "Brief Step Title (5-7 words)",
      "description": "Detailed explanation of what this code does and why it matters.",
      "file": "relative/path/from/project/root.go",
      "line_start": 10,
      "line_end": 50,
      "action": "read"
    }
  ]
}
Include 3-5 steps. Use relative paths from project root. Keep line ranges focused (10-50 lines per step).`

		messages := []api.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		}

		response, err := client.Chat(context.Background(), messages)
		if err != nil {
			return generationErrorMsg{err: err}
		}

		return walkthroughGeneratedMsg{content: response}
	}
}

func (m Model) saveGeneratedWalkthrough(content string) tea.Cmd {
	return func() tea.Msg {
		parsed := api.ParseMarkdownCodeBlocks(content)

		walkthroughDir := filepath.Join(m.rootPath, ".wmti")
		if err := os.MkdirAll(walkthroughDir, 0755); err != nil {
			return errMsg{err: fmt.Errorf("failed to create walkthrough directory: %w", err)}
		}

		baseName := "generated-" + time.Now().Format("20060102-150405")
		filePath := filepath.Join(walkthroughDir, baseName+".wmti.json")

		if err := os.WriteFile(filePath, []byte(parsed), 0644); err != nil {
			return errMsg{err: fmt.Errorf("failed to save walkthrough: %w", err)}
		}

		return walkthroughLoadedMsg{}
	}
}

func truncateError(fullErr string) string {
	lines := strings.Split(fullErr, "\n")
	firstLine := lines[0]
	if len(firstLine) > 80 {
		return firstLine[:77] + "..."
	}
	return firstLine
}

func writeErrorLog(provider, fullErr string) {
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".local", "share", "wmti")
	os.MkdirAll(logDir, 0755)

	logPath := filepath.Join(logDir, "error.log")
	logMsg := fmt.Sprintf("[%s] Provider: %s\nError: %s\n\n",
		time.Now().Format("2006-01-02 15:04:05"), provider, fullErr)
	os.WriteFile(logPath, []byte(logMsg), 0644)
}
