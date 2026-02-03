# Project: walk-me-through-it (wmti)

An interactive TUI application that guides developers through codebases using step-by-step walkthroughs.

## Overview

`wmti` is a terminal-based IDE-like experience that helps developers understand unfamiliar codebases by following structured walkthrough documents. Users can navigate through code explanations using simple keyboard shortcuts.

## Tech Stack

- **Language**: Go 1.21+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Elm-style architecture for interactive terminal apps
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) - CSS-like styling for terminals
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles) - Pre-built TUI components
- **CLI**: [Cobra](https://github.com/spf13/cobra) - Command structure and flag parsing
- **Config**: [Viper](https://github.com/spf13/viper) - Configuration management and API key storage
- **File Navigation**: Standard `os` and `filepath` packages

## Project Structure

```
walk-me-through-it/
├── cmd/
│   ├── wmti/              # Main entry point (root command)
│   └── init/              # Init command implementation
├── internal/
│   ├── config/            # Configuration management (Viper)
│   ├── tui/               # Bubble Tea models, views, and update logic
│   ├── navigator/         # JSON walkthrough document parser
│   ├── walker/            # File system traversal utilities
│   └── api/               # Optional AI API integration
├── pkg/
│   └── models/            # Shared data structures (JSON schema types)
├── testdata/              # Sample JSON walkthrough documents for testing
├── go.mod
├── go.sum
├── .golangci.yml          # Linting configuration
├── Makefile               # Build, test, and development commands
└── AGENTS.md              # This file
```

## Development Commands

```bash
# Build the application
make build

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Run the application locally
go run ./cmd/wmti

# Install dependencies
go mod tidy
```

## Architecture

### Command Structure (Cobra)

1. **Root Command (`wmti`)**: Launches the TUI IDE
   - Loads configuration
   - Opens walkthrough navigator
   - Starts interactive session

2. **Init Command (`wmti init`)**: Sets up initial configuration
   - Prompts for optional API key
   - Creates config directory and file
   - Validates setup

### TUI Architecture (Bubble Tea)

The application follows the Bubble Tea Model-Update-View pattern:

- **Model**: Application state including current step, file content, navigation history
- **Update**: Handles keyboard input (Tab/Shift+Tab for navigation, q to quit)
- **View**: Renders the IDE interface with file viewer and description panel

### JSON Walkthrough Schema

Walkthrough documents define the navigation path:

```json
{
  "title": "Understanding the Authentication Flow",
  "description": "A walkthrough of how user authentication works in this codebase",
  "version": "1.0.0",
  "steps": [
    {
      "id": 1,
      "title": "Entry Point",
      "description": "The authentication process starts here when a user submits login credentials.",
      "file": "src/auth/login.go",
      "line_start": 15,
      "line_end": 45,
      "action": "read"
    },
    {
      "id": 2,
      "title": "Credential Validation",
      "description": "This function validates the provided credentials against the database.",
      "file": "src/auth/validator.go",
      "line_start": 32,
      "line_end": 78,
      "action": "read"
    }
  ]
}
```

**Schema Fields:**
- `title`: Walkthrough title displayed in header
- `description`: Overview of what this walkthrough covers
- `version`: Document version for compatibility
- `steps[]`: Array of navigation steps
  - `id`: Unique step identifier (sequential)
  - `title`: Step title shown in navigation
  - `description`: Detailed explanation of this step
  - `file`: Relative path to file from project root
  - `line_start`: Starting line number (1-based)
  - `line_end`: Ending line number (1-based)
  - `action`: Action type (currently only "read")

### Navigation Flow

1. User launches `wmti` in a project directory
2. Application searches for `.wmti/walkthrough.json` or accepts path argument
3. TUI displays first step with file content and description
4. User presses `Tab` to advance to next step
5. User presses `Shift+Tab` to go back to previous step
6. User presses `q` or `Ctrl+C` to exit

## Key Patterns

### Code Organization

- **Separation of Concerns**: Business logic in `internal/`, CLI in `cmd/`, shared types in `pkg/`
- **Interface-Based Design**: Define interfaces for components to enable testing
- **Error Handling**: Use Go's explicit error returns; wrap errors with context
- **Configuration**: Use Viper for all config; support env vars and config files

### Testing Approach

- Unit tests for all packages in `internal/`
- Table-driven tests for multiple scenarios
- Mock external dependencies (file system, API calls)
- Test data in `testdata/` directory
- Integration tests for TUI components using Bubble Tea's testing utilities

### Code Style

- Follow standard Go conventions (gofmt, goimports)
- Use `golangci-lint` for comprehensive linting
- Prefer explicit over implicit
- Document all exported functions and types
- Keep functions small and focused (single responsibility)

### TUI Patterns

- Use Bubble Tea's message-based architecture
- Separate view logic from state management
- Support both mouse and keyboard navigation where applicable
- Provide clear visual feedback for user actions
- Handle terminal resize events gracefully

## AI Development Guidelines

### Before Making Changes

1. **Check Current Branch**: Always verify which branch you're on
2. **Create Feature Branch**: Branch from the current working branch (NOT main)
   ```bash
   git checkout -b feature/descriptive-name
   ```
3. **Understand Context**: Read relevant files before modifying
4. **Follow Existing Patterns**: Match code style and architecture of surrounding code

### Branch Workflow

```
main (prod) ←── dev ←── feature/xyz ←── AI works here
     ↑            ↑
   releases    your work
```

**Rules:**
- **main**: Production releases only - NEVER commit directly
- **dev**: Integration branch - merge feature branches here
- **feature/*** : AI and feature work branches
- Always branch from the current working branch
- Use descriptive branch names: `feature/add-json-parser`, `fix/navigation-bug`

### Making Changes

1. **Small, Focused Commits**: One logical change per commit
2. **Test Changes**: Run `make test` and `make lint` before finishing
3. **Update Documentation**: If you change behavior, update comments and AGENTS.md
4. **Handle Errors**: Never ignore errors; always return or handle them
5. **Type Safety**: Use strong typing; avoid `interface{}` unless necessary

### Code Review Checklist

Before considering work complete:

- [ ] Code follows Go conventions (gofmt -l . shows no issues)
- [ ] All tests pass (`go test ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] No hardcoded secrets or API keys
- [ ] Error messages are clear and actionable
- [ ] New code has corresponding tests
- [ ] Documentation updated if needed
- [ ] Commit messages are descriptive

### Common Pitfalls to Avoid

1. **Don't modify go.mod directly** - Use `go get` or `go mod tidy`
2. **Don't ignore errors** - Always check error returns
3. **Don't use global state** - Pass dependencies explicitly
4. **Don't block the TUI** - Use Bubble Tea's Cmd pattern for async operations
5. **Don't assume terminal capabilities** - Use Lipgloss for safe styling

## External Resources

- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper#readme)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Contact & Issues

- Project: walk-me-through-it
- License: MIT
- Author: Christopher Sannar
