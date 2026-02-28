package tui

import (
	"io"
	"strings"
)

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
