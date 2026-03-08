package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

		listView := m.list.View()

		pageInfo := ""
		if m.list.Paginator.TotalPages > 1 {
			pageInfo = fmt.Sprintf("Page %d of %d  |  ", m.list.Paginator.Page+1, m.list.Paginator.TotalPages)
		}

		footerBorder := m.styles.FooterBorderStyle.Render(strings.Repeat("─", m.width))

		if m.deleteConfirmPath != "" {
			filename := filepath.Base(m.deleteConfirmPath)
			confirmText := m.styles.StepDescriptionStyle.Render("Delete " + filename + "?")
			footerContent := lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.styles.RenderKeybinding("y", "yes"),
				" ",
				m.styles.RenderKeybinding("n", "no"),
			)
			footerInner := m.styles.FooterStyle.Render(confirmText)
			footerInner2 := m.styles.FooterStyle.Render(footerContent)
			footer := lipgloss.JoinVertical(
				lipgloss.Left,
				footerBorder,
				footerInner,
				footerInner2,
			)

			listLines := strings.Split(listView, "\n")
			var content strings.Builder
			content.WriteString(listView)
			for i := len(listLines); i < m.height-lipgloss.Height(footer); i++ {
				content.WriteString("\n")
			}
			content.WriteString(footer)

			return content.String()
		}

		footerContent := lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.styles.StepDescriptionStyle.Render(pageInfo),
			m.styles.RenderKeybinding("↑↓", "select"),
			" ",
			m.styles.RenderKeybinding("Enter", "open"),
			" ",
			m.styles.RenderKeybinding("d", "delete"),
			" ",
			m.styles.RenderKeybinding("q", "quit"),
			" ",
			m.styles.RenderKeybinding("m", "model"),
		)
		footerInner := m.styles.FooterStyle.Render(footerContent)
		footer := lipgloss.JoinVertical(
			lipgloss.Left,
			footerBorder,
			footerInner,
		)

		listLines := strings.Split(listView, "\n")
		var content strings.Builder
		content.WriteString(listView)
		for i := len(listLines); i < m.height-lipgloss.Height(footer); i++ {
			content.WriteString("\n")
		}
		content.WriteString(footer)

		return content.String()

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

	case StateProviderSelect:
		var content strings.Builder

		title := m.styles.HeaderFileStyle.Render("Select Provider")
		content.WriteString(title)
		content.WriteString("\n\n")

		for i, provider := range m.providerList {
			name := provider
			if i == m.providerSelectedIndex {
				highlightStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#85DCFF")).
					Background(lipgloss.Color("#2D2D3D")).
					Padding(0, 1)
				content.WriteString(highlightStyle.Render(" " + name + " "))
			} else {
				item := m.styles.StepDescriptionStyle.Render(name)
				content.WriteString(" ")
				content.WriteString(item)
			}
			content.WriteString("\n")
		}

		if m.providerAskingFor == "key" {
			content.WriteString("\n")
			prompt := "Enter API key for " + m.providerList[m.providerSelectedIndex] + ":"
			content.WriteString(m.styles.StepDescriptionStyle.Render(prompt))
			content.WriteString("\n")
			content.WriteString(m.modelTextInput.View())
		}

		if m.connectionTestResult != "" {
			content.WriteString("\n\n")
			if strings.HasPrefix(m.connectionTestResult, "✓") {
				content.WriteString(m.styles.StepTitleStyle.Render(m.connectionTestResult))
			} else {
				content.WriteString(m.styles.ErrorTextStyle.Render(m.connectionTestResult))
			}
		}

		content.WriteString("\n\n")

		footerBorder := m.styles.FooterBorderStyle.Render(strings.Repeat("─", m.width))
		if m.providerAskingFor == "key" {
			footerContent := lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.styles.RenderKeybinding("Enter", "save"),
				" ",
				m.styles.RenderKeybinding("Esc", "cancel"),
			)
			footerInner := m.styles.FooterStyle.Render(footerContent)
			footer := lipgloss.JoinVertical(
				lipgloss.Left,
				footerBorder,
				footerInner,
			)
			content.WriteString(footer)
		} else {
			footerContent := lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.styles.RenderKeybinding("Enter", "select"),
				" ",
				m.styles.RenderKeybinding("t", "test"),
				" ",
				m.styles.RenderKeybinding("q", "back"),
			)
			footerInner := m.styles.FooterStyle.Render(footerContent)
			footer := lipgloss.JoinVertical(
				lipgloss.Left,
				footerBorder,
				footerInner,
			)
			content.WriteString(footer)
		}

		listLines := strings.Split(content.String(), "\n")
		var result strings.Builder
		for i := len(listLines); i < m.height; i++ {
			result.WriteString("\n")
		}
		result.WriteString(content.String())

		return result.String()

	case StateConfirmDelete:
		return m.renderCentered(fmt.Sprintf(
			"Deleting %s...",
			filepath.Base(m.deleteConfirmPath),
		))

	case StateModelSelect:
		var content strings.Builder

		title := m.styles.HeaderFileStyle.Render("Select Model")
		content.WriteString(title)
		content.WriteString("\n\n")

		for i, model := range m.modelList {
			name := model
			if i == m.modelSelectedIndex {
				highlightStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#85DCFF")).
					Background(lipgloss.Color("#2D2D3D")).
					Padding(0, 1)
				content.WriteString(highlightStyle.Render(" " + name + " "))
			} else {
				item := m.styles.StepDescriptionStyle.Render(name)
				content.WriteString(" ")
				content.WriteString(item)
			}
			content.WriteString("\n")
		}

		if m.modelIsAdding {
			content.WriteString("\n")
			prompt := "Enter model name:"
			if m.modelAskingFor == "key" {
				prompt = "Enter API key:"
			}
			content.WriteString(m.styles.StepDescriptionStyle.Render(prompt))
			content.WriteString("\n")
			content.WriteString(m.modelTextInput.View())
			if m.modelError != "" {
				content.WriteString("\n")
				content.WriteString(m.styles.ErrorTextStyle.Render(m.modelError))
			}
		}

		if m.connectionTestResult != "" {
			content.WriteString("\n\n")
			if strings.HasPrefix(m.connectionTestResult, "✓") {
				content.WriteString(m.styles.StepTitleStyle.Render(m.connectionTestResult))
			} else {
				content.WriteString(m.styles.ErrorTextStyle.Render(m.connectionTestResult))
			}
		}

		content.WriteString("\n\n")

		footerBorder := m.styles.FooterBorderStyle.Render(strings.Repeat("─", m.width))
		if m.modelIsAdding {
			footerContent := lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.styles.RenderKeybinding("Enter", "save"),
				" ",
				m.styles.RenderKeybinding("Esc", "cancel"),
			)
			footerInner := m.styles.FooterStyle.Render(footerContent)
			footer := lipgloss.JoinVertical(
				lipgloss.Left,
				footerBorder,
				footerInner,
			)
			content.WriteString(footer)
		} else if m.modelDeleteConfirm {
			modelName := m.modelList[m.modelSelectedIndex]
			confirmText := m.styles.StepDescriptionStyle.Render("Delete " + modelName + "?")
			footerContent := lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.styles.RenderKeybinding("y", "yes"),
				" ",
				m.styles.RenderKeybinding("n", "no"),
			)
			footerInner := m.styles.FooterStyle.Render(confirmText)
			footerInner2 := m.styles.FooterStyle.Render(footerContent)
			footer := lipgloss.JoinVertical(
				lipgloss.Left,
				footerBorder,
				footerInner,
				footerInner2,
			)
			content.WriteString(footer)
		} else {
			footerContent := lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.styles.RenderKeybinding("↑↓", "select"),
				" ",
				m.styles.RenderKeybinding("Enter", "confirm"),
				" ",
				m.styles.RenderKeybinding("t", "test"),
				" ",
				m.styles.RenderKeybinding("d", "delete"),
				" ",
				m.styles.RenderKeybinding("q", "back"),
			)
			footerInner := m.styles.FooterStyle.Render(footerContent)
			footer := lipgloss.JoinVertical(
				lipgloss.Left,
				footerBorder,
				footerInner,
			)
			content.WriteString(footer)
		}

		listLines := strings.Split(content.String(), "\n")
		var result strings.Builder
		for i := len(listLines); i < m.height; i++ {
			result.WriteString("\n")
		}
		result.WriteString(content.String())

		return result.String()

	case StateViewing:
		if m.err != nil {
			return m.renderCentered(fmt.Sprintf("Error: %s\n\nPress Tab to try next step, or q to quit", m.err.Error()))
		}
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

		sidebarWidth := int(float64(m.width) * 0.30)
		if sidebarWidth < 25 {
			sidebarWidth = 25
		}
		codeWidth := m.width - sidebarWidth - 2

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
		header := lipgloss.JoinVertical(
			lipgloss.Left,
			headerInner,
			m.styles.HeaderBorderStyle.Render(strings.Repeat("─", m.width)),
		)

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
		footer := lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.FooterBorderStyle.Render(strings.Repeat("─", m.width)),
			footerInner,
		)

		headerHeight := lipgloss.Height(header)
		footerHeight := lipgloss.Height(footer)
		contentHeight := m.height - headerHeight - footerHeight - 2
		if contentHeight < 3 {
			contentHeight = 3
		}

		m.styles.SetDimensions(m.width, codeWidth, sidebarWidth, contentHeight)

		linePrefixWidth := 7
		maxCodeLineWidth := codeWidth - linePrefixWidth - 1
		if maxCodeLineWidth < 1 {
			maxCodeLineWidth = 1
		}

		isHighlighted := func(lineNum int) bool {
			return lineNum >= m.highlightStart && lineNum <= m.highlightEnd
		}

		var codeBlock string
		if len(m.fileContent) > 0 {
			highlightedLines, err := m.highlighter.HighlightFile(step.File, m.fileContent, m.displayStart)
			if err != nil || len(highlightedLines) == 0 {
				var codeContent strings.Builder
				for i, line := range m.fileContent {
					if i >= contentHeight {
						break
					}
					lineNum := m.displayStart + i
					truncatedLine := truncate(line, maxCodeLineWidth)
					lineContent := fmt.Sprintf("%4d │ %s", lineNum, truncatedLine)

					// Apply highlight background to the line
					if isHighlighted(lineNum) {
						codeContent.WriteString(m.styles.HighlightLineStyle.Render(lineContent))
					} else {
						codeContent.WriteString(lineContent)
					}
					codeContent.WriteString("\n")
				}
				codeBlock = m.styles.CodeStyle.Render(codeContent.String())
			} else {
				var codeContent strings.Builder
				for i := range m.fileContent {
					if i >= contentHeight {
						break
					}
					lineNum := m.displayStart + i

					// Check if we have highlighted content for this line
					hasHighlightContent := i < len(highlightedLines)

					// If highlighted, keep syntax colors but add a subtle background
					// by using a lighter background color that lets colors show through
					if isHighlighted(lineNum) {
						lineNumStr := fmt.Sprintf("%4d │ ", lineNum)
						codeContent.WriteString(m.styles.HighlightLineStyle.Render(lineNumStr))

						remainingWidth := maxCodeLineWidth
						if hasHighlightContent {
							for _, token := range highlightedLines[i].Tokens {
								tokenWidth := displayWidth(token.Text)
								if tokenWidth > remainingWidth {
									if remainingWidth > 0 {
										truncatedText := truncateClean(token.Text, remainingWidth)
										// Use token's color plus lighter background
										codeContent.WriteString(token.Style.Background(lipgloss.Color("#4a4a5e")).Render(truncatedText))
									}
									break
								}
								// Add background to the token's existing style
								codeContent.WriteString(token.Style.Background(lipgloss.Color("#4a4a5e")).Render(token.Text))
								remainingWidth -= tokenWidth
								if remainingWidth <= 0 {
									break
								}
							}
							// Fill remaining width with highlight background
							if remainingWidth > 0 {
								codeContent.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("#4a4a5e")).Render(strings.Repeat(" ", remainingWidth)))
							}
						}
						codeContent.WriteString("\n")
					} else {
						// Use syntax highlighting for non-highlighted lines
						lineNumStr := fmt.Sprintf("%4d │ ", lineNum)
						codeContent.WriteString(lineNumStr)

						remainingWidth := maxCodeLineWidth
						if hasHighlightContent {
							for _, token := range highlightedLines[i].Tokens {
								tokenWidth := displayWidth(token.Text)
								if tokenWidth > remainingWidth {
									if remainingWidth > 0 {
										truncatedText := truncateClean(token.Text, remainingWidth)
										codeContent.WriteString(token.Style.Render(truncatedText))
									}
									break
								}
								codeContent.WriteString(token.Style.Render(token.Text))
								remainingWidth -= tokenWidth
								if remainingWidth <= 0 {
									break
								}
							}
						}
						codeContent.WriteString("\n")
					}
				}
				codeBlock = m.styles.CodeStyle.Render(codeContent.String())
			}
		} else {
			codeBlock = m.styles.CodeStyle.Render("Loading file content...\n")
		}

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

		separator := m.styles.SeparatorStyle.Height(contentHeight).Render(strings.Repeat("│\n", contentHeight))

		contentRow := lipgloss.JoinHorizontal(
			lipgloss.Top,
			codeBlock,
			separator,
			sidebarBlock,
		)

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

	for i := len(lines); i < m.height; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
}
