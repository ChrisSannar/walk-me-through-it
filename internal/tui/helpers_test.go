package tui

import "testing"

func TestDisplayWidth_ASCII(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"", 0},
		{"a", 1},
		{"Hello World", 11},
		{"1234567890", 10},
	}

	for _, tt := range tests {
		result := displayWidth(tt.input)
		if result != tt.expected {
			t.Errorf("displayWidth(%q) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestDisplayWidth_Unicode(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"日本語", 3},
		{"🎉", 1},
		{"hello日本語", 8},
		{"a🎉b", 3},
		{"é", 1},
		{"café", 4},
	}

	for _, tt := range tests {
		result := displayWidth(tt.input)
		if result != tt.expected {
			t.Errorf("displayWidth(%q) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestDisplayWidth_Tab(t *testing.T) {
	result := displayWidth("\t")
	if result != 4 {
		t.Errorf("displayWidth(tab) = %d; want 4", result)
	}

	result = displayWidth("a\tb")
	if result != 6 { // a(1) + tab(4) + b(1)
		t.Errorf("displayWidth('a\\tb') = %d; want 6", result)
	}
}

func TestDisplayWidth_ControlCharacters(t *testing.T) {
	// Control characters should be ignored
	result := displayWidth("hello\x00world")
	if result != 10 {
		t.Errorf("displayWidth with control chars = %d; want 10", result)
	}

	result = displayWidth("\x01\x02\x03")
	if result != 0 {
		t.Errorf("displayWidth with only control chars = %d; want 0", result)
	}
}

func TestTruncate_WithEllipsis(t *testing.T) {
	tests := []struct {
		input    string
		maxWidth int
		expected string
	}{
		{"hello", 3, "..."},   // can't fit any char + ellipsis
		{"hello", 4, "h..."},  // 1 char + ellipsis fits exactly
		{"hello", 5, "hello"}, // exact fit, no ellipsis needed
		{"hello", 6, "hello"}, // fits, no ellipsis
		{"hi", 10, "hi"},      // smaller than max
		{"", 5, ""},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.maxWidth)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q; want %q", tt.input, tt.maxWidth, result, tt.expected)
		}
	}
}

func TestTruncate_NoEllipsisWhenFits(t *testing.T) {
	result := truncate("hello", 5)
	if result != "hello" {
		t.Errorf("truncate('hello', 5) = %q; want 'hello'", result)
	}

	result = truncate("hello", 10)
	if result != "hello" {
		t.Errorf("truncate('hello', 10) = %q; want 'hello'", result)
	}
}

func TestTruncate_Tab(t *testing.T) {
	// tab is 4 characters wide
	result := truncate("a\tb", 5) // a(1) + tab(4) = 5, but adds ellipsis
	if result == "" {
		t.Error("truncate with tab should not be empty")
	}
}

func TestTruncateClean_NoEllipsis(t *testing.T) {
	tests := []struct {
		input    string
		maxWidth int
		expected string
	}{
		{"hello", 3, "hel"},
		{"hello", 4, "hell"},
		{"hello", 5, "hello"},
		{"hello", 6, "hello"},
		{"", 5, ""},
		{"helloworld", 4, "hell"},
	}

	for _, tt := range tests {
		result := truncateClean(tt.input, tt.maxWidth)
		if result != tt.expected {
			t.Errorf("truncateClean(%q, %d) = %q; want %q", tt.input, tt.maxWidth, result, tt.expected)
		}
	}
}

func TestTruncateClean_Tab(t *testing.T) {
	result := truncateClean("\t\t", 5) // each tab is 4
	if result != "\t" {
		t.Errorf("truncateClean('\\t\\t', 5) = %q; want '\\t'", result)
	}
}

func TestPadRight_AddsSpaces(t *testing.T) {
	tests := []struct {
		input    string
		width    int
		expected string
	}{
		{"hello", 10, "hello     "},
		{"hello", 5, "hello"},
		{"", 5, "     "},
		{"hi", 2, "hi"},
		{"test", 8, "test    "},
	}

	for _, tt := range tests {
		result := padRight(tt.input, tt.width)
		if result != tt.expected {
			t.Errorf("padRight(%q, %d) = %q; want %q", tt.input, tt.width, result, tt.expected)
		}
	}
}

func TestPadRight_Unicode(t *testing.T) {
	result := padRight("日本語", 8)
	// Each CJK character counts as 1 visual width in our simple implementation
	// "日本語" is 3 chars, so we add 5 spaces
	expected := "日本語     "
	if result != expected {
		t.Errorf("padRight('日本語', 8) = %q; want %q", result, expected)
	}
}

func TestPadRight_AlreadyLonger(t *testing.T) {
	result := padRight("hello", 3)
	if result != "hello" {
		t.Errorf("padRight('hello', 3) = %q; want 'hello'", result)
	}
}

func TestTruncateVsTruncateClean(t *testing.T) {
	input := "function"

	// truncate adds ellipsis
	r1 := truncate(input, 5)
	// truncateClean does not
	r2 := truncateClean(input, 5)

	if r1 == r2 {
		t.Error("truncate and truncateClean should produce different results")
	}

	// Both should not exceed maxWidth
	if len(r1) > 5 && r1 != "f..." {
		t.Errorf("truncate exceeded width: %s", r1)
	}
	if len(r2) > 5 {
		t.Errorf("truncateClean exceeded width: %s", r2)
	}
}
