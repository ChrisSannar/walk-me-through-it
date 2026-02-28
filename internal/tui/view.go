package tui

import (
	"fmt"
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

		var codeBlock string
		if len(m.fileContent) > 0 {
			highlightedLines, err := m.highlighter.HighlightFile(step.File, m.fileContent, step.LineStart)
			if err != nil || len(highlightedLines) == 0 {
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
				var codeContent strings.Builder
				for i, hlLine := range highlightedLines {
					if i >= contentHeight {
						break
					}
					lineNumStr := fmt.Sprintf("%4d │ ", hlLine.LineNumber)
					codeContent.WriteString(lineNumStr)

					remainingWidth := maxCodeLineWidth

					for _, token := range hlLine.Tokens {
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
					codeContent.WriteString("\n")
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
