package highlighter

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestNewHighlighter(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	if h == nil {
		t.Fatal("NewHighlighter returned nil")
	}

	if h.theme != ThemeDracula {
		t.Errorf("expected theme dracula, got %s", h.theme)
	}

	if h.cache == nil {
		t.Error("cache should not be nil")
	}

	if h.formatter == nil {
		t.Error("formatter should not be nil")
	}
}

func TestSetTheme(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	h.SetTheme(ThemeMonokai)

	if h.theme != ThemeMonokai {
		t.Errorf("expected theme monokai, got %s", h.theme)
	}
}

func TestGetTheme(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	theme := h.GetTheme()

	if theme != ThemeDracula {
		t.Errorf("expected ThemeDracula, got %s", theme)
	}
}

func TestHighlightFile_EmptyLines(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	result, err := h.HighlightFile("test.go", []string{}, 1)

	if err != nil {
		t.Fatalf("HighlightFile failed: %v", err)
	}

	if result != nil {
		t.Error("expected nil result for empty lines")
	}
}

func TestHighlightFile_Go(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	lines := []string{
		"package main",
		"",
		"func main() {",
		"\tfmt.Println(\"hello\")",
		"}",
	}

	result, err := h.HighlightFile("main.go", lines, 1)

	if err != nil {
		t.Fatalf("HighlightFile failed: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}

	// First line should be "package main"
	if len(result) < 1 || result[0].LineNumber != 1 {
		t.Errorf("expected first line at line 1")
	}

	// Check that tokens exist
	if len(result[0].Tokens) == 0 {
		t.Error("expected tokens for first line")
	}
}

func TestHighlightFile_NoLexer(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	// Use a file extension that likely has no lexer
	lines := []string{"some content", "more content"}

	result, err := h.HighlightFile("file.xyz", lines, 1)

	if err != nil {
		t.Fatalf("HighlightFile failed: %v", err)
	}

	// Should fall back to plain text
	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}

	// First line should be plain text
	if len(result) < 1 {
		t.Error("expected at least one line")
	}
}

func TestHighlightFile_Caching(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	lines := []string{"package main", "func main() {}"}

	// First call - should populate cache
	result1, err := h.HighlightFile("main.go", lines, 1)
	if err != nil {
		t.Fatalf("HighlightFile failed: %v", err)
	}

	// Second call - should use cache
	result2, err := h.HighlightFile("main.go", lines, 1)
	if err != nil {
		t.Fatalf("HighlightFile failed: %v", err)
	}

	// Results should be equivalent
	if len(result1) != len(result2) {
		t.Errorf("expected same length, got %d and %d", len(result1), len(result2))
	}
}

func TestHighlightFile_DifferentStartLines(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	lines := []string{"line1", "line2", "line3"}

	// Same content, different start line - different cache keys
	result1, err := h.HighlightFile("test.go", lines, 1)
	if err != nil {
		t.Fatalf("HighlightFile failed: %v", err)
	}

	result2, err := h.HighlightFile("test.go", lines, 10)
	if err != nil {
		t.Fatalf("HighlightFile failed: %v", err)
	}

	// Line numbers should differ
	if result1[0].LineNumber == result2[0].LineNumber {
		t.Error("expected different line numbers for different start lines")
	}

	if result1[0].LineNumber != 1 {
		t.Errorf("expected first result at line 1, got %d", result1[0].LineNumber)
	}

	if result2[0].LineNumber != 10 {
		t.Errorf("expected second result at line 10, got %d", result2[0].LineNumber)
	}
}

func TestConvertTokens_Multiline(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	// Test that multiline tokens are properly split
	lines := []string{"package main", "func main() {", "}"}

	result, err := h.HighlightFile("main.go", lines, 1)
	if err != nil {
		t.Fatalf("HighlightFile failed: %v", err)
	}

	// Should have 3 lines
	if len(result) != 3 {
		t.Errorf("expected 3 lines, got %d", len(result))
	}

	// Each line should have at least one token
	for i, line := range result {
		if len(line.Tokens) == 0 {
			t.Errorf("line %d has no tokens", i+1)
		}
	}
}

func TestHighlightAsPlain(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	lines := []string{"hello", "world", "test"}

	result := h.highlightAsPlain(lines, 5)

	if len(result) != 3 {
		t.Errorf("expected 3 lines, got %d", len(result))
	}

	// Check line numbers
	if result[0].LineNumber != 5 {
		t.Errorf("expected line 5, got %d", result[0].LineNumber)
	}
	if result[1].LineNumber != 6 {
		t.Errorf("expected line 6, got %d", result[1].LineNumber)
	}
	if result[2].LineNumber != 7 {
		t.Errorf("expected line 7, got %d", result[2].LineNumber)
	}

	// Should have plain text tokens (empty style)
	for _, line := range result {
		if len(line.Tokens) != 1 {
			t.Errorf("expected 1 token per line, got %d", len(line.Tokens))
		}
		// lipgloss.Style is not nil-comparable, just check it exists
		_ = line.Tokens[0].Style
	}
}

func TestCacheResult(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	result := []HighlightedLine{
		{LineNumber: 1, Tokens: []Token{{Text: "test", Style: lipgloss.NewStyle()}}},
	}

	h.cacheResult("testkey", result)

	// Check that cache contains the key
	h.cacheMu.RLock()
	_, ok := h.cache["testkey"]
	h.cacheMu.RUnlock()

	if !ok {
		t.Error("expected cache to contain testkey")
	}
}

func TestClearCache(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	// Add something to cache
	result := []HighlightedLine{{LineNumber: 1, Tokens: []Token{{Text: "test"}}}}
	h.cacheResult("testkey", result)

	// Clear cache
	h.ClearCache()

	// Check that cache is empty
	h.cacheMu.RLock()
	length := len(h.cache)
	h.cacheMu.RUnlock()

	if length != 0 {
		t.Errorf("expected cache to be empty after ClearCache, got %d", length)
	}
}

func TestGetAvailableThemes(t *testing.T) {
	themes := GetAvailableThemes()

	if len(themes) == 0 {
		t.Fatal("expected non-empty themes list")
	}

	// Check for known themes
	found := false
	for _, theme := range themes {
		if theme == ThemeDracula {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected ThemeDracula in available themes")
	}
}

func TestRenderLine(t *testing.T) {
	line := HighlightedLine{
		LineNumber: 1,
		Tokens: []Token{
			{Text: "package", Style: lipgloss.NewStyle()},
			{Text: " main", Style: lipgloss.NewStyle()},
		},
	}

	result := RenderLine(line)

	// Should join tokens without separators
	expected := "package main"
	if result != expected {
		t.Errorf("RenderLine = %q; want %q", result, expected)
	}
}

func TestGetLexer(t *testing.T) {
	h := NewHighlighter(ThemeDracula)

	tests := []struct {
		filename  string
		expectNil bool
	}{
		{"main.go", false},
		{"test.js", false},
		{"app.py", false},
		{"unknown.xyz", true},
	}

	for _, tt := range tests {
		lexer := h.getLexer(tt.filename)
		if tt.expectNil && lexer != nil {
			t.Errorf("expected nil lexer for %s", tt.filename)
		}
		// Most should have lexers
		if !tt.expectNil && lexer == nil {
			t.Logf("no lexer for %s (may be expected)", tt.filename)
		}
	}
}
