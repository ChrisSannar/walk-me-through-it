// Package highlighter provides syntax highlighting functionality using Chroma
package highlighter

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a syntax highlighting theme
type Theme string

const (
	ThemeMonokai   Theme = "monokai"
	ThemeDracula   Theme = "dracula"
	ThemeGitHub    Theme = "github"
	ThemeOneDark   Theme = "onedark"
	ThemeSolarized Theme = "solarized-dark"
	ThemeVim       Theme = "vim"
)

// Highlighter provides syntax highlighting with caching
type Highlighter struct {
	cache     map[string][]HighlightedLine
	cacheMu   sync.RWMutex
	theme     Theme
	formatter chroma.Formatter
}

// HighlightedLine represents a single line of highlighted code
type HighlightedLine struct {
	LineNumber int
	Tokens     []Token
}

// Token represents a single highlighted token
type Token struct {
	Text  string
	Style lipgloss.Style
}

// NewHighlighter creates a new highlighter instance
func NewHighlighter(theme Theme) *Highlighter {
	return &Highlighter{
		cache:     make(map[string][]HighlightedLine),
		theme:     theme,
		formatter: formatters.TTY16m, // True color formatter
	}
}

// SetTheme changes the current theme and clears the cache
func (h *Highlighter) SetTheme(theme Theme) {
	h.cacheMu.Lock()
	defer h.cacheMu.Unlock()
	h.theme = theme
	h.cache = make(map[string][]HighlightedLine) // Clear cache
}

// GetTheme returns the current theme
func (h *Highlighter) GetTheme() Theme {
	return h.theme
}

// HighlightFile highlights code content from a file
func (h *Highlighter) HighlightFile(filename string, lines []string, startLine int) ([]HighlightedLine, error) {
	if len(lines) == 0 {
		return nil, nil
	}

	// Create cache key
	cacheKey := fmt.Sprintf("%s:%d:%d", filename, startLine, len(lines))

	// Check cache
	h.cacheMu.RLock()
	if cached, ok := h.cache[cacheKey]; ok {
		h.cacheMu.RUnlock()
		return cached, nil
	}
	h.cacheMu.RUnlock()

	// Get lexer for file type
	lexer := h.getLexer(filename)
	if lexer == nil {
		// No lexer found, return plain text
		result := h.highlightAsPlain(lines, startLine)
		h.cacheResult(cacheKey, result)
		return result, nil
	}

	// Join lines for lexing
	content := strings.Join(lines, "\n")

	// Tokenize
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize: %w", err)
	}

	// Get style
	style := h.getChromaStyle()

	// Convert to our format
	result := h.convertTokens(iterator, style, startLine)

	// Cache result
	h.cacheResult(cacheKey, result)

	return result, nil
}

// getLexer returns the appropriate lexer for a file
func (h *Highlighter) getLexer(filename string) chroma.Lexer {
	ext := strings.ToLower(filepath.Ext(filename))

	// Map extensions to Chroma lexer names
	lexerMap := map[string]string{
		".go":    "go",
		".js":    "javascript",
		".ts":    "typescript",
		".jsx":   "jsx",
		".tsx":   "tsx",
		".py":    "python",
		".rb":    "ruby",
		".php":   "php",
		".java":  "java",
		".kt":    "kotlin",
		".scala": "scala",
		".c":     "c",
		".cpp":   "cpp",
		".h":     "cpp",
		".hpp":   "cpp",
		".rs":    "rust",
		".swift": "swift",
		".md":    "markdown",
		".json":  "json",
		".yaml":  "yaml",
		".yml":   "yaml",
		".html":  "html",
		".css":   "css",
		".scss":  "scss",
		".sass":  "sass",
		".sh":    "bash",
		".bash":  "bash",
		".zsh":   "zsh",
	}

	if lexerName, ok := lexerMap[ext]; ok {
		return lexers.Get(lexerName)
	}

	// Try to match by filename
	return lexers.Match(filename)
}

// getChromaStyle returns the Chroma style for the current theme
func (h *Highlighter) getChromaStyle() *chroma.Style {
	styleName := string(h.theme)
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}
	return style
}

// convertTokens converts Chroma tokens to our HighlightedLine format
func (h *Highlighter) convertTokens(iterator chroma.Iterator, style *chroma.Style, startLine int) []HighlightedLine {
	var lines []HighlightedLine
	var currentLine HighlightedLine
	currentLine.LineNumber = startLine

	for token := iterator(); token != chroma.EOF; token = iterator() {
		text := token.Value
		linesInToken := strings.Split(text, "\n")

		for i, lineText := range linesInToken {
			if i > 0 {
				// New line, save current and start new
				if len(currentLine.Tokens) > 0 {
					lines = append(lines, currentLine)
				}
				currentLine = HighlightedLine{
					LineNumber: currentLine.LineNumber + 1,
					Tokens:     []Token{},
				}
			}

			if lineText != "" {
				// Convert Chroma style to lipgloss style
				lipglossStyle := h.chromaToLipgloss(token.Type, style)
				currentLine.Tokens = append(currentLine.Tokens, Token{
					Text:  lineText,
					Style: lipglossStyle,
				})
			}
		}
	}

	// Don't forget the last line
	if len(currentLine.Tokens) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// chromaToLipgloss converts a Chroma token type and style to lipgloss style
func (h *Highlighter) chromaToLipgloss(tokenType chroma.TokenType, style *chroma.Style) lipgloss.Style {
	entry := style.Get(tokenType)

	var s lipgloss.Style

	// Apply color
	if entry.Colour != 0 {
		color := entry.Colour.String()
		s = s.Foreground(lipgloss.Color(color))
	}

	// Note: Background colors are intentionally not applied to keep the code
	// area transparent and consistent with the terminal background

	// Apply styles
	if entry.Bold == chroma.Yes {
		s = s.Bold(true)
	}
	if entry.Italic == chroma.Yes {
		s = s.Italic(true)
	}
	if entry.Underline == chroma.Yes {
		s = s.Underline(true)
	}

	return s
}

// highlightAsPlain returns plain text without highlighting
func (h *Highlighter) highlightAsPlain(lines []string, startLine int) []HighlightedLine {
	var result []HighlightedLine
	for i, line := range lines {
		result = append(result, HighlightedLine{
			LineNumber: startLine + i,
			Tokens: []Token{
				{Text: line, Style: lipgloss.NewStyle()},
			},
		})
	}
	return result
}

// cacheResult stores highlighted result in cache
func (h *Highlighter) cacheResult(key string, result []HighlightedLine) {
	h.cacheMu.Lock()
	defer h.cacheMu.Unlock()

	// Limit cache size (simple LRU: clear if too big)
	if len(h.cache) > 100 {
		h.cache = make(map[string][]HighlightedLine)
	}

	h.cache[key] = result
}

// ClearCache clears the highlighting cache
func (h *Highlighter) ClearCache() {
	h.cacheMu.Lock()
	defer h.cacheMu.Unlock()
	h.cache = make(map[string][]HighlightedLine)
}

// GetAvailableThemes returns a list of available themes
func GetAvailableThemes() []Theme {
	return []Theme{
		ThemeMonokai,
		ThemeDracula,
		ThemeGitHub,
		ThemeOneDark,
		ThemeSolarized,
		ThemeVim,
	}
}

// RenderLine renders a highlighted line to a string
func RenderLine(line HighlightedLine) string {
	var result strings.Builder
	for _, token := range line.Tokens {
		result.WriteString(token.Style.Render(token.Text))
	}
	return result.String()
}
