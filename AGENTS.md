# AI Agent Instructions for wmti

This document provides instructions for AI agents working on the Walk Me Through It (wmti) project.

## Project Overview

**wmti** is a TUI (Terminal User Interface) application written in Go that helps developers navigate codebases through interactive, JSON-based walkthroughs. It uses:

- **Bubble Tea** - TUI framework
- **Lipgloss** - Styling and layout
- **Chroma** - Syntax highlighting
- **Cobra** - CLI framework
- **Viper** - Configuration management

## Build & Development

### Commands
```bash
make build    # Build the application (output: build/wmti)
make test     # Run tests
make lint     # Run linter (golangci-lint)
make fmt      # Format code
make clean    # Clean build artifacts
go run ./cmd/wmti/  # Run without building
```

### Do Not Run Build/Test Unless Asked
**Unless the user specifically asks you to run a build or test, do not execute these commands.** The user will run them when ready.

## Project Structure

```
cmd/
├── wmti/              # Main CLI entry point (main.go, root.go)
└── init/              # Init command implementation

internal/
├── tui/               # Bubble Tea TUI models and views
│   ├── model.go       # Main TUI logic
│   └── styles.go      # UI theming system
├── highlighter/       # Syntax highlighting (Chroma integration)
│   └── highlighter.go
├── navigator/         # Walkthrough navigation logic
├── walker/            # File system traversal and reading
├── api/               # AI API integration (future)
└── config/            # Configuration management

pkg/
└── models/            # Shared data structures

docs/
├── SESSIONS.md        # Development session summaries
└── NOTES.md           # Feature ideas and goals

testdata/              # Sample walkthrough files
```

## Coding Conventions

### Go Code Style
- Follow standard Go conventions (gofmt)
- Use `goimports` for imports
- Error handling: return errors, don't panic
- Context usage: pass context.Context as first param when needed

### TUI Patterns
- **States**: Use `AppState` enum (StateSelecting, StateEnteringPath, StateLoading, StateViewing, StateLoadingStep)
- **Messages**: Define custom tea.Msg types for async operations
- **Styling**: Use Lipgloss styles from `internal/tui/styles.go`
- **Dimensions**: Use `SetDimensions()` for responsive layouts

### File Naming
- Walkthrough files: `wmti_*.json` pattern
- Internal packages: lowercase, descriptive
- Test files: `*_test.go`

## Walkthrough Format

Walkthrough files are JSON with this structure:
```json
{
  "title": "Walkthrough Title",
  "description": "Description",
  "version": "1.0.0",
  "steps": [
    {
      "id": 1,
      "title": "Step Title",
      "description": "Step description",
      "file": "path/to/file.go",
      "line_start": 1,
      "line_end": 10,
      "action": "read"
    }
  ]
}
```

## UI Themes

### Available UI Themes (styles.go)
- `ThemeDark` (default)
- `ThemeLight`
- `ThemeMidnight`

### Available Code Themes (highlighter.go)
- `ThemeDracula` (default)
- `ThemeMonokai`
- `ThemeGitHub`
- `ThemeOneDark`
- `ThemeSolarized`
- `ThemeVim`

## Key Features to Maintain

1. **Cyclic Navigation**: Tab/Shift+Tab cycle through steps infinitely
2. **File Selection**: List of `wmti_*.json` files with selection
3. **Syntax Highlighting**: 20+ languages supported via Chroma
4. **Security**: Path validation, extension allowlisting, audit logging
5. **Responsive Layout**: Adapts to terminal size changes

## Navigation Keys

- **Tab**: Next step (cyclic)
- **Shift+Tab**: Previous step (cyclic)
- **q**: Return to file selection (when viewing), or quit
- **Ctrl+C**: Quit application
- **Enter**: Select item in list
- **j/k**: Navigate list (vim-style, loops at ends)

## Session Documentation

When making significant changes, update `docs/SESSIONS.md` with:
- Date
- What was accomplished
- Files modified
- Key implementation details
- Build status

## Common Gotchas

1. **Truncation**: Code view uses `truncateClean()` (no ellipsis) while other areas use `truncate()` (with ellipsis)
2. **Highlighting Background**: Don't apply background colors from Chroma themes - keep code area transparent
3. **List Looping**: Arrow keys and j/k both need custom handling to loop from first↔last
4. **Lipgloss Height**: Use `lipgloss.Height()` to calculate actual rendered dimensions

## Testing

- Unit tests in `*_test.go` files
- Run with: `go test ./...` or `make test`
- Walker package has tests (`internal/walker/finder_test.go`)

## Future Roadmap (from NOTES.md)

- Line highlighting for current step
- Chat box in content view
- Scan for new walkthrough files
- Auto-detect AI keys during init
- Store walkthroughs in local storage
