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
- Walkthrough files: `*.wmti.json` pattern in `.wmti/` directory
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
      "action": "read"
    },
    {
      "id": 2,
      "title": "Credential Validation",
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
      "action": "read"
    }
  ]
}
```

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
