package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - can be swapped for different themes
const (
	ColorBackground  = "#1a1a2e" // Deep navy/slate
	ColorPrimary     = "#25A065" // Green accent
	ColorText        = "#FFFDF5" // Off-white
	ColorTextMuted   = "#888888" // Medium gray
	ColorTextSidebar = "#BBBBBB" // Light gray for sidebar
	ColorBorder      = "#333333" // Dark gray border
	ColorBorderLight = "#444444" // Slightly lighter border
)

// Styles holds all UI styling definitions
type Styles struct {
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
}

// NewStyles creates a new Styles instance with default theme
func NewStyles() *Styles {
	s := &Styles{}

	// Layout styles
	s.CodeStyle = lipgloss.NewStyle().
		PaddingTop(1)

	s.SeparatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBorderLight))

	s.SidebarStyle = lipgloss.NewStyle().
		Padding(1, 2)

	// Header styles
	s.HeaderStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorBackground)).
		Padding(0, 2)

	s.HeaderFileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true)

	s.HeaderLineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPrimary)).
		Bold(true)

	s.HeaderLineLabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextMuted))

	s.HeaderBorderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBorder))

	// Footer styles
	s.FooterStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorBackground)).
		Padding(0, 2)

	s.FooterBorderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBorder))

	s.KeybindingBoxStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorPrimary)).
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Padding(0, 1)

	s.KeybindingDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextSidebar))

	s.FooterDividerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBorderLight))

	s.AuditMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextMuted)).
		Italic(true)

	// Sidebar content styles (keeping existing colors)
	s.StepHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPrimary)).
		Bold(true).
		MarginBottom(1)

	s.StepTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	s.StepDescriptionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextSidebar))

	return s
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
