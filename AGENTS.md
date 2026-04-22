<<<<<<< HEAD
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
=======
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
>>>>>>> main

## Project Structure

```
<<<<<<< HEAD
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
- Walkthrough files: `*.wmti.json` pattern in `.wmti/` directory
- Internal packages: lowercase, descriptive
- Test files: `*_test.go`

## Walkthrough Format

Walkthrough files are JSON with this structure:
```json
{
  "title": "Walkthrough Title",
  "description": "Description",
=======
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
>>>>>>> main
  "version": "1.0.0",
  "steps": [
    {
      "id": 1,
<<<<<<< HEAD
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

---

## LLM Instructions for Creating Walkthrough Files

### QUICK REFERENCE
- **Location**: Create files in `.wmti/` directory (at project root)
- **Naming**: Use `*.wmti.json` pattern (e.g., `auth-flow.wmti.json`)
- **Format**: Valid JSON with specific structure
- **Action**: Always use `"read"` (only supported action currently)

### File Location & Naming Convention
```
project-root/
├── .wmti/                    # ALL walkthrough files go here
│   ├── getting-started.wmti.json
│   ├── auth-flow.wmti.json
│   └── api-integration.wmti.json
└── src/
    └── ...
```

**Naming Rules**:
- Use descriptive, kebab-case names
- Must end with `.wmti.json`
- Avoid: `wmti_*.json` (old pattern, deprecated)

### JSON Structure Template

```json
{
  "title": "Descriptive Walkthrough Title",
  "description": "One or two sentences explaining what this walkthrough covers and who it's for",
  "version": "1.0.0",
  "steps": [
    {
      "id": 1,
      "title": "Brief Step Title (5-7 words)",
      "description": "Explain what the user should understand here. Include context, purpose, and how this relates to previous/next steps. Be specific about what code concepts or patterns are being demonstrated.",
      "file": "relative/path/from/project/root.go",
      "line_start": 10,
      "line_end": 35,
      "action": "read"
    }
  ]
}
```

### Field Reference

#### Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Appears in file selection list. Be descriptive but concise (3-6 words). |
| `description` | string | Yes | Shown before walkthrough starts. 1-2 sentences explaining scope and audience. |
| `version` | string | Yes | Semantic versioning (e.g., "1.0.0"). Update when steps change. |
| `steps` | array | Yes | Array of step objects. Must have at least 1 step. |

#### Step Object Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | integer | Yes | Sequential starting at 1. Used for ordering. |
| `title` | string | Yes | Short title shown in sidebar (5-7 words max). |
| `description` | string | Yes | Detailed explanation shown in sidebar. Explain WHY this code matters, not just WHAT it does. |
| `file` | string | Yes | Relative path from project root. Use forward slashes. |
| `line_start` | integer | Yes | First line to display (1-based indexing). |
| `line_end` | integer | Yes | Last line to display. Keep ranges focused (10-50 lines ideal). |
| `action` | string | Yes | Currently only `"read"` is supported. |

### Best Practices for AI-Generated Walkthroughs

1. **Analyze Code First**
   - Identify the entry point and main flow
   - Group related functionality into logical steps
   - Skip boilerplate, focus on business logic

2. **Write Descriptive Step Descriptions**
   - ❌ Bad: "This is the login function"
   - ✅ Good: "The login function validates credentials against the database using bcrypt for password hashing. On success, it generates a JWT token with a 24-hour expiration."

3. **Choose Appropriate Line Ranges**
   - Include enough context (imports, function signatures)
   - But keep focused: 10-50 lines per step is ideal
   - Don't show entire files unless necessary

4. **Create Logical Progression**
   - Start with entry points (main, handlers, controllers)
   - Follow the data flow through the system
   - Build complexity gradually

5. **Use Relative Paths**
   - ✅ Good: `src/auth/login.go`
   - ❌ Bad: `/home/user/project/src/auth/login.go`
   - ❌ Bad: `../src/auth/login.go`

6. **Version Your Walkthroughs**
   - Start with "1.0.0"
   - Increment when adding/removing steps
   - Document major changes in description

### Example Complete Walkthrough

```json
{
  "title": "Authentication Flow",
  "description": "Understanding how user authentication works in this codebase, from login request to session creation.",
  "version": "1.0.0",
  "steps": [
    {
      "id": 1,
      "title": "Login Handler Entry",
      "description": "The authentication flow begins here. The HTTP handler receives login credentials, validates the request format, and extracts the email and password from the JSON body. This is the API entry point that clients call.",
      "file": "src/handlers/auth.go",
      "line_start": 45,
      "line_end": 72,
=======
      "title": "Entry Point",
      "description": "The authentication process starts here when a user submits login credentials.",
      "file": "src/auth/login.go",
      "line_start": 15,
      "line_end": 45,
>>>>>>> main
      "action": "read"
    },
    {
      "id": 2,
      "title": "Credential Validation",
<<<<<<< HEAD
      "description": "The ValidateCredentials function queries the database for the user record by email. It uses bcrypt.CompareHashAndPassword to securely check the password without storing plaintext. Returns a User struct on success or authentication error.",
      "file": "src/auth/validator.go",
      "line_start": 28,
      "line_end": 55,
      "action": "read"
    },
    {
      "id": 3,
      "title": "JWT Token Generation",
      "description": "Upon successful validation, GenerateToken creates a JWT containing the user ID and role. The token is signed with a secret key from environment variables and expires after 24 hours. This token will be used for subsequent authenticated requests.",
      "file": "src/auth/token.go",
      "line_start": 15,
      "line_end": 42,
      "action": "read"
    },
    {
      "id": 4,
      "title": "Session Creation",
      "description": "The CreateSession function stores the active session in Redis with the token as the key. This enables session invalidation and tracking. The session includes user metadata and expiration time matching the JWT.",
      "file": "src/auth/session.go",
      "line_start": 33,
      "line_end": 58,
=======
      "description": "This function validates the provided credentials against the database.",
      "file": "src/auth/validator.go",
      "line_start": 32,
      "line_end": 78,
>>>>>>> main
      "action": "read"
    }
  ]
}
```

<<<<<<< HEAD
### Step Creation Workflow

When asked to create a walkthrough:

1. **Understand the Goal**
   - What concept/feature should the walkthrough explain?
   - Who is the target audience (new dev, reviewer, etc.)?

2. **Explore the Codebase**
   - Find the main entry points
   - Trace the execution flow
   - Identify key functions and their relationships

3. **Plan the Steps**
   - Create a rough outline of the logical flow
   - Each step should have a clear purpose
   - Ensure progressive disclosure (simple → complex)

4. **Verify File Paths**
   - Use relative paths from project root
   - Confirm files exist and are readable
   - Check line numbers are accurate

5. **Write Descriptions**
   - Explain WHY the code is structured this way
   - Mention important patterns or decisions
   - Connect to previous/next steps for flow

6. **Review and Test**
   - Validate JSON syntax
   - Check all required fields are present
   - Ensure IDs are sequential starting at 1

### Common Mistakes to Avoid

1. **Using Absolute Paths**
   - ❌ `"file": "/home/user/project/src/main.go"`
   - ✅ `"file": "src/main.go"`

2. **Skipping Step IDs**
   - ❌ IDs: 1, 2, 4, 5 (missing 3)
   - ✅ IDs: 1, 2, 3, 4

3. **Zero-Based Line Numbers**
   - ❌ `"line_start": 0` (first line is 1)
   - ✅ `"line_start": 1`

4. **Too Many Lines Per Step**
   - ❌ 200+ lines (user loses context)
   - ✅ 10-50 lines (focused and readable)

5. **Vague Descriptions**
   - ❌ "This function does stuff"
   - ✅ "This function validates user input using regex patterns and returns sanitized data"

6. **Forgetting Version**
   - Always include version field
   - Use semantic versioning

---

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
2. **File Selection**: List of `*.wmti.json` files with selection from `.wmti/` directory
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
=======
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
>>>>>>> main
