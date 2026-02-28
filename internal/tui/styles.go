package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/chrissannar/walk-me-through-it/internal/highlighter"
)

// UITheme represents the UI color theme type
type UITheme string

const (
	ThemeDark     UITheme = "dark"
	ThemeLight    UITheme = "light"
	ThemeMidnight UITheme = "midnight"
)

// ThemeColors holds all colors for a theme
type ThemeColors struct {
	Background  string
	Primary     string
	Text        string
	TextMuted   string
	TextSidebar string
	Border      string
	BorderLight string
}

// Predefined themes
var themes = map[UITheme]ThemeColors{
	ThemeDark: {
		Background:  "#1a1a2e",
		Primary:     "#25A065",
		Text:        "#FFFDF5",
		TextMuted:   "#888888",
		TextSidebar: "#BBBBBB",
		Border:      "#333333",
		BorderLight: "#444444",
	},
	ThemeLight: {
		Background:  "#f5f5f5",
		Primary:     "#2d8659",
		Text:        "#1a1a1a",
		TextMuted:   "#666666",
		TextSidebar: "#444444",
		Border:      "#cccccc",
		BorderLight: "#dddddd",
	},
	ThemeMidnight: {
		Background:  "#0d1117",
		Primary:     "#58a6ff",
		Text:        "#c9d1d9",
		TextMuted:   "#8b949e",
		TextSidebar: "#aeb4ba",
		Border:      "#30363d",
		BorderLight: "#21262d",
	},
}

// CodeTheme represents the syntax highlighting theme
type CodeTheme = highlighter.Theme

// Available code themes
const (
	CodeThemeMonokai   CodeTheme = highlighter.ThemeMonokai
	CodeThemeDracula   CodeTheme = highlighter.ThemeDracula
	CodeThemeGitHub    CodeTheme = highlighter.ThemeGitHub
	CodeThemeOneDark   CodeTheme = highlighter.ThemeOneDark
	CodeThemeSolarized CodeTheme = highlighter.ThemeSolarized
	CodeThemeVim       CodeTheme = highlighter.ThemeVim
)

// Styles holds all UI styling definitions
type Styles struct {
	uiTheme   UITheme
	codeTheme CodeTheme

	// Layout styles
	CodeStyle      lipgloss.Style
	SeparatorStyle lipgloss.Style
	SidebarStyle   lipgloss.Style

	// Header styles
	HeaderStyle          lipgloss.Style
	HeaderFileStyle      lipgloss.Style
	HeaderLineStyle      lipgloss.Style
	HeaderLineLabelStyle lipgloss.Style
	HeaderBorderStyle    lipgloss.Style

	// Footer styles
	FooterStyle         lipgloss.Style
	FooterBorderStyle   lipgloss.Style
	KeybindingBoxStyle  lipgloss.Style
	KeybindingTextStyle lipgloss.Style
	KeybindingDescStyle lipgloss.Style
	FooterDividerStyle  lipgloss.Style
	AuditMsgStyle       lipgloss.Style

	// Sidebar content styles
	StepHeaderStyle      lipgloss.Style
	StepTitleStyle       lipgloss.Style
	StepDescriptionStyle lipgloss.Style

	// Code highlight style
	HighlightLineStyle lipgloss.Style
}

// NewStyles creates a new Styles instance with default theme
func NewStyles() *Styles {
	s := &Styles{
		uiTheme:   ThemeDark,
		codeTheme: CodeThemeDracula,
	}
	s.applyTheme()
	return s
}

// SetUITheme changes the UI theme
func (s *Styles) SetUITheme(theme UITheme) {
	if _, ok := themes[theme]; ok {
		s.uiTheme = theme
		s.applyTheme()
	}
}

// SetCodeTheme changes the code highlighting theme
func (s *Styles) SetCodeTheme(theme CodeTheme) {
	s.codeTheme = theme
}

// GetUITheme returns the current UI theme
func (s *Styles) GetUITheme() UITheme {
	return s.uiTheme
}

// GetCodeTheme returns the current code theme
func (s *Styles) GetCodeTheme() CodeTheme {
	return s.codeTheme
}

// applyTheme applies the current UI theme colors
func (s *Styles) applyTheme() {
	colors, ok := themes[s.uiTheme]
	if !ok {
		colors = themes[ThemeDark]
	}

	// Layout styles
	s.CodeStyle = lipgloss.NewStyle().
		PaddingTop(1)

	s.SeparatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.BorderLight))

	s.SidebarStyle = lipgloss.NewStyle().
		Padding(1, 2)

	// Header styles
	s.HeaderStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colors.Background)).
		Padding(0, 2)

	s.HeaderFileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Bold(true)

	s.HeaderLineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Primary)).
		Bold(true)

	s.HeaderLineLabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.TextMuted))

	s.HeaderBorderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Border))

	// Footer styles
	s.FooterStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colors.Background)).
		Padding(0, 2)

	s.FooterBorderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Border))

	s.KeybindingBoxStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colors.Primary)).
		Foreground(lipgloss.Color(colors.Text)).
		Bold(true).
		Padding(0, 1)

	s.KeybindingDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.TextSidebar))

	s.FooterDividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.BorderLight))

	s.AuditMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.TextMuted)).
		Italic(true)

	// Sidebar content styles
	s.StepHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Primary)).
		Bold(true).
		MarginBottom(1)

	s.StepTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	s.StepDescriptionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.TextSidebar))

	// Code highlight - light background to highlight the relevant lines
	s.HighlightLineStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#2a2a3e"))
}

// RenderKeybinding creates a styled keybinding like [Tab] Next
func (s *Styles) RenderKeybinding(key, description string) string {
	keyBox := s.KeybindingBoxStyle.Render(key)
	desc := s.KeybindingDescStyle.Render(description)
	return lipgloss.JoinHorizontal(lipgloss.Center, keyBox, " ", desc)
}

// SetDimensions updates styles that depend on screen dimensions
func (s *Styles) SetDimensions(width, codeWidth, sidebarWidth, contentHeight int) {
	s.CodeStyle = s.CodeStyle.Width(codeWidth).Height(contentHeight)
	s.SidebarStyle = s.SidebarStyle.Width(sidebarWidth).Height(contentHeight)
}

// GetAvailableUIThemes returns a list of available UI themes
func GetAvailableUIThemes() []UITheme {
	return []UITheme{ThemeDark, ThemeLight, ThemeMidnight}
}

// GetAvailableCodeThemes returns a list of available code highlighting themes
func GetAvailableCodeThemes() []CodeTheme {
	return highlighter.GetAvailableThemes()
}
