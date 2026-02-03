# Walk Me Through It (wmti)

An interactive terminal-based IDE that guides developers through unfamiliar codebases using step-by-step walkthrough documents.

## Overview

`wmti` is a TUI (Terminal User Interface) application designed to help developers understand complex codebases by following structured, JSON-based walkthroughs. Whether you're onboarding to a new project, reviewing code, or documenting architecture, wmti provides an interactive, step-by-step navigation experience directly in your terminal.

## Features

- **Interactive TUI**: Navigate through codebases with an IDE-like terminal interface
- **Step-by-Step Walkthroughs**: Follow structured paths through code with detailed explanations
- **Sequential Navigation**: Use `Tab` to advance and `Shift+Tab` to go back between steps
- **Configuration Management**: Simple setup with `wmti init` for API keys and preferences
- **JSON-Based Walkthroughs**: Human-readable walkthrough documents that can be hand-written or AI-generated
- **File Range Highlighting**: View specific line ranges within files for focused learning

## Installation

### Prerequisites

- Go 1.21 or later

### From Source

```bash
# Clone the repository
git clone https://github.com/chrissannar/walk-me-through-it.git
cd walk-me-through-it

# Build the application
make build

# Install to $GOPATH/bin
make install
```

### Using Go Install

```bash
go install github.com/chrissannar/walk-me-through-it/cmd/wmti@latest
```

## Quick Start

### 1. Initialize Configuration

```bash
wmti init
```

This creates a configuration file at `~/.config/wmti/config.yaml`. You can optionally add an API key for AI-generated walkthroughs (feature coming soon).

### 2. Create a Walkthrough

Create a JSON walkthrough document (e.g., `.wmti/walkthrough.json`):

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

### 3. Launch the TUI

```bash
# Automatically finds .wmti/walkthrough.json
wmti

# Or specify a custom path
wmti --walkthrough path/to/walkthrough.json
```

### 4. Navigate the Walkthrough

Once the TUI opens:
- **Tab**: Advance to the next step
- **Shift+Tab**: Go back to the previous step
- **q** or **Ctrl+C**: Exit the application

## Walkthrough Document Format

Walkthrough documents are JSON files that define your navigation path:

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Walkthrough title displayed in header |
| `description` | string | Overview of what this walkthrough covers |
| `version` | string | Document version for compatibility |
| `steps` | array | Array of navigation steps |
| `steps[].id` | int | Unique step identifier (sequential, starting at 1) |
| `steps[].title` | string | Step title shown in navigation |
| `steps[].description` | string | Detailed explanation of this step |
| `steps[].file` | string | Relative path to file from project root |
| `steps[].line_start` | int | Starting line number (1-based) |
| `steps[].line_end` | int | Ending line number (1-based) |
| `steps[].action` | string | Action type (currently only "read") |

## Development

### Building from Source

```bash
# Build the application
make build

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Clean build artifacts
make clean
```

### Project Structure

```
walk-me-through-it/
├── cmd/
│   ├── wmti/              # Main CLI entry point
│   └── init/              # Init command implementation
├── internal/
│   ├── config/            # Configuration management (Viper)
│   ├── tui/               # Bubble Tea TUI models and views
│   ├── navigator/         # JSON walkthrough document parser
│   ├── walker/            # File system traversal utilities
│   └── api/               # AI API integration (future)
├── pkg/
│   └── models/            # Shared data structures
├── testdata/              # Sample JSON walkthroughs
└── Makefile               # Build automation
```

### Tech Stack

- **Language**: Go 1.21+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **CLI**: [Cobra](https://github.com/spf13/cobra)
- **Config**: [Viper](https://github.com/spf13/viper)

## Roadmap

- [x] Basic CLI structure with Cobra
- [x] JSON walkthrough document parsing
- [x] File system navigation
- [ ] Full TUI implementation with Bubble Tea
- [ ] Syntax highlighting for code display
- [ ] AI-powered walkthrough generation
- [ ] Vim-style navigation keys
- [ ] Interactive file tree browser
- [ ] Export walkthroughs to Markdown

## Contributing

1. Check the current branch (`git branch`)
2. Create a feature branch from the working branch: `git checkout -b feature/your-feature`
3. Make your changes following Go conventions
4. Run tests and linter: `make test && make lint`
5. Submit a pull request to the `dev` branch

See [AGENTS.md](AGENTS.md) for detailed development guidelines.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Christopher Sannar
