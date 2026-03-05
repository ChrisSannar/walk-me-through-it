# Development Session Summary

## Project: Walk Me Through It (wmti)

**Current State**: A functional TUI application that displays code walkthroughs with file content and step-by-step navigation.

---

## What We've Built

### Core Functionality
- TUI startup with walkthrough file discovery (pattern: `wmti_*.json`)
- File selection list using Bubble Tea's list component
- Split-pane view showing:
  - **Top**: File content with line numbers (lines X-Y from the specified file)
  - **Bottom**: Step instructions, progress (Step X of Y), and keyboard controls
- Cyclic navigation: Tab moves forward, Shift+Tab moves backward, wrapping at ends
- Proper screen clearing to prevent rendering artifacts

### Naming Convention
- Walkthrough files must follow `wmti_<title>.json` pattern
- Existing files renamed: `wmti_project_setup.json`, `.wmti/wmti_security_measures.json`

---

## Key Files & Recent Changes

### `internal/tui/model.go` (Main TUI logic)
- `Model` struct with state management (StateSelecting, StateEnteringPath, StateLoading, StateViewing)
- `walkthroughItem` with Title/Description methods for list display
- `loadCurrentStepFile()` - Async file content loading
- `View()` - Renders split-pane layout with proper screen filling
- Keyboard handling: Tab/Shift+Tab for cyclic navigation, q/Ctrl+C to quit

### `internal/navigator/navigator.go` (Navigation logic)
- `NextCycle()` / `PreviousCycle()` - Modular navigation that wraps around
- Standard `Next()` / `Previous()` with bounds checking still available
- `CurrentStep()`, `Progress()`, `HasNext()`, `HasPrevious()`

### `internal/walker/finder.go` (File discovery)
- `FindWalkthroughFiles()` - Searches for `wmti_*.json` in root and `.wmti/` subdirectory
- Validates files are valid JSON objects
- `GetDefaultWalkthroughPath()` for default location

### `internal/walker/walker.go` (File reading)
- `ReadFileLines(path, start, end)` - Secure file reading with line range
- Security features: path sanitization, extension allowlisting, 10MB size limit, audit logging

---

## Current User Experience

1. Run `wmti` → Shows list of available walkthrough files
2. Select walkthrough → Loads and displays first step
3. View shows:
   ```
   📄 cmd/wmti/main.go (lines 1-15)
   ──────────────────────────────────
      1 │ package main
      2 │ 
      3 │ import (
   ...
   ──────────────────────────────────
   📋 Step 1 of 2: Main Entry Point
   📝 The application starts here...
   ⌨️  Tab: Next | Shift+Tab: Previous | q: Quit
   ```
4. Tab/Shift+Tab cycle through steps infinitely

---

## Potential Next Steps

### Immediate Improvements
- Syntax highlighting for the file content (using Glamour or Chroma)
- Better visual styling with Lipgloss (colors, borders, padding)
- Scrollable file content for large line ranges
- Show file tree or breadcrumb navigation
- Add step numbers to the file content display

### Features from README Roadmap
- Vim-style navigation keys (hjkl)
- Interactive file tree browser
- Export walkthroughs to Markdown
- AI-powered walkthrough generation (API integration in `internal/api/api.go`)

### Technical Debt
- Add tests for TUI components
- Add tests for navigator cyclic methods
- Error handling improvements in file loading
- Configuration management (the `init` command exists but may need updates)

### Build Commands
```bash
make build    # Build the application
make test     # Run tests
make lint     # Run linter
```

---

## Status

The application is at a solid MVP state with core navigation and file display working. The next logical step would likely be syntax highlighting or improved visual styling to make the code more readable.

## Session 2026-02-05: Layout and Text Wrapping Fixes

### Problem Solved
Fixed text wrapping issues in the code view. Long lines were wrapping and breaking the TUI layout, causing visual artifacts and misalignment.

### Solution Implemented
Implemented proper text truncation based on visual width (not character count) to handle:
- Unicode characters (multi-byte)
- Tab characters (displayed as 4 spaces but count as 1 rune)
- ANSI escape codes in syntax highlighted output

### Layout Architecture (Three-Section Fixed Layout)
```
┌─────────────────────────────────────────────────────┐
│ Header (2 lines): File info with underline          │
├──────────────────┬──────────────────────────────────┤
│ Code Area        │ Sidebar                          │
│ (70% width)      │ (30% width, min 25 chars)        │
│ • Line numbers   │ • Step title                     │
│ • Syntax         │ • Step description               │
│   highlighted    │ • Progress                       │
│ • Truncated at   │ • Controls                       │
│   edge           │                                  │
├──────────────────┴──────────────────────────────────┤
│ Footer (2 lines): Controls help                     │
└─────────────────────────────────────────────────────┘
```

### Key Implementation Details

**Width Calculations in `internal/tui/model.go`:**
```go
headerHeight := 2
footerHeight := 2
contentHeight := m.height - headerHeight - footerHeight
sidebarWidth := int(float64(m.width) * 0.30)  // 30%, min 25 chars
codeWidth := m.width - sidebarWidth - 2        // -2 for separator
linePrefixWidth := 7                           // "999 │ " format
maxCodeLineWidth := codeWidth - linePrefixWidth - 1  // 1-char safety buffer
```

**Helper Functions:**
- `displayWidth(s string) int` - Calculates visual width (tabs = 4 spaces)
- `truncate(s string, maxWidth int) string` - Truncates at visual width, no ellipsis

### Layout Construction
Using `lipgloss.JoinHorizontal()` and `lipgloss.JoinVertical()` for precise control:
- Left panel: Code with line numbers
- Vertical separator: `│` character in #666 color
- Right panel: Step info with word wrapping enabled
- No wrapping in code area (clean truncation at edge)

### Files Modified
- `internal/tui/model.go` - Updated `View()` function and added helper functions

### Build Command
```bash
go build ./cmd/wmti/
```

---


# Session Summary

## What We've Done

### 1. Added Audit Log to Footer
**Files**: `internal/walker/walker.go`, `internal/tui/model.go`
- Modified `Walker` struct to store `lastAuditMessage`
- Added `GetLastAuditMessage()` method
- Footer now displays `[AUDIT] Reading file: ...` message before the controls

### 2. Styled the Sidebar (Step Instructions Panel)
**File**: `internal/tui/model.go`
- Added padding: `Padding(1, 2)` inside the sidebar
- **Step header** (Step X of Y): Green (#25A065), bold
- **Step title**: White (#FFFDF5), bold, with margins
- **Step description**: Light gray (#BBBBBB)
- Used `lipgloss.JoinVertical()` for proper spacing

### 3. Attempted Header Fix + Debug
**File**: `internal/tui/model.go`
- Removed lipgloss height constraints from header/footer styles
- Added debug logging to `/tmp/wmti_debug.log` and `/tmp/wmti_render.log`
- **Key Finding**: Header IS being generated correctly (proven by debug output showing "📄 cmd/wmti/main.go (lines 1-15)" at the top of the render)
- Cleaned up debug code
- Kept the newline between audit message and instructions as requested

## Current Issue
**The header is generated correctly (visible in debug logs) but doesn't appear visually in the terminal.** This suggests a display/rendering issue rather than a code logic problem.

## Files Being Modified
- `internal/walker/walker.go` - Audit message tracking
- `internal/tui/model.go` - Main TUI styling and layout (header, footer, sidebar)

## Next Steps
1. Investigate why the generated header isn't displaying (terminal height issue? lipgloss JoinVertical behavior? ANSI codes?)
2. Continue styling improvements once header works
3. Consider adding syntax highlighting for the code area (mentioned in roadmap)
4. Test on different terminal sizes/emulators

## Status
The build currently compiles successfully and tests pass. The header content is confirmed to be generated correctly in the debug logs.

---

## Session 2026-02-06: Syntax Highlighting and Theming

### What Was Accomplished

#### 1. Dynamic Layout System (`internal/tui/model.go`)
- Fixed header/footer display issues using `lipgloss.Height()` for dynamic height calculation
- Header and footer now have "placement priority" - they render first, content fills remaining space
- Content area has minimum 3-line protection

#### 2. Navigation Improvements (`internal/tui/model.go`)
- 'q' key now returns to file selection page (instead of quitting app)
- Ctrl+C still quits the entire application
- No rescanning when going back - uses previously found files

#### 3. UI Theming System (`internal/tui/styles.go`)
- Extracted all styling to dedicated file
- **3 UI Themes**: Dark, Light, Midnight
- **ThemeColors struct** for easy theme swapping
- Methods: `SetUITheme()`, `GetUITheme()`, `applyTheme()`

#### 4. Syntax Highlighting (NEW: `internal/highlighter/highlighter.go`)
- Integrated **Chroma** library (`github.com/alecthomas/chroma/v2`)
- **6 Code Themes**: Monokai, Dracula, GitHub, OneDark, Solarized, Vim
- **20+ Language Support**: Go, JS/TS/JSX/TSX, Python, Ruby, PHP, Java/Kotlin/Scala, C/C++, Rust, Swift, Shell, HTML/CSS/SCSS, Markdown, JSON, YAML
- **Caching**: LRU cache (100 files max) with thread-safe RWMutex
- Cache key: `filename:startLine:lineCount`

#### 5. Text Truncation with Highlighting (`internal/tui/model.go`)
- Fixed wrapping issues when code exceeds terminal width
- Tracks `remainingWidth` as tokens render
- Truncates individual tokens when they exceed available space
- Uses `displayWidth()` function to handle Unicode/tabs properly

### File Structure
```
internal/
├── tui/
│   ├── model.go          # Main TUI - MODIFIED (dynamic heights, highlighting, navigation)
│   └── styles.go         # NEW - UI theming system
├── highlighter/
│   └── highlighter.go    # NEW - Chroma-based syntax highlighting
├── walker/
│   └── walker.go         # File reading (allowedExtensions reference)
└── navigator/
    └── navigator.go      # Walkthrough navigation (unchanged)
```

### Working Features
- File selection with list UI
- Step navigation (Tab/Shift+Tab)
- Styled header/footer with keybinding boxes
- Full syntax highlighting for supported languages
- Dynamic layout that adapts to terminal size
- Proper text truncation for long lines

### Default Configuration
- UI Theme: Dark (`ThemeDark`)
- Code Theme: Dracula (`CodeThemeDracula`)

### Build Status
✅ Tests passing
✅ Build successful

---

## Session 2026-03-01: Fixed Broken Example Bug and Updated Example Walkthrough

### What Was Done

#### 1. Fixed the Broken Example Walkthrough Bug
The `example.wmti.json` was pointing to a non-existent file (`path/to/file.go`), causing the TUI to get stuck on "Loading step..." when navigating between steps.

**Root Cause**: When loading a step's file fails (file doesn't exist), the app remained stuck in `StateLoadingStep` state with no way to recover.

**Fixes Implemented**:

- **Fixed in `internal/tui/model.go`**: Added logic to revert state from `StateLoadingStep` back to `StateViewing` when an error occurs during file loading
- **Fixed in `internal/tui/view.go`**: Added error display when in `StateViewing` state so users can see what went wrong
- **Fixed `example.wmti.json`**: Updated to reference an actual existing file (`internal/tui/types.go`)

#### 2. Updated the Example Walkthrough as a Meta Tutorial
Changed `example.wmti.json` to be a 4-step self-referential tour explaining how wmti works:

1. **"JSON Walkthrough Files"** - Shows itself explaining JSON structure
2. **"Each Step Defines a View"** - Shows the steps array within the JSON
3. **"Cyclic Navigation"** - Shows the Tab/Shift+Tab navigation code in `model.go`
4. **"Creating Walkthroughs"** - Shows `navigation.wmti.json` as another example

#### 3. Updated the Init Template
Updated `cmd/init/init.go` to generate the same new 4-step example when users run `wmti init`

### Files Modified

- `internal/tui/model.go` - Bug fix for stuck loading state
- `internal/tui/view.go` - Added error message display
- `.wmti/example.wmti.json` - Updated content
- `cmd/init/init.go` - Updated template content

### Build Status
✅ Tests passing
✅ Build successful

### What's Next
The user mentioned they want to add a "create a walkthrough" option. Before proceeding with that feature, we should:
- Build and test the app to verify the fixes work
- Then discuss the design for the "create walkthrough" feature (CLI command, interactive wizard, AI-assisted, etc.)

### Open Questions to Address Next Session
1. **Audit message positioning**: Footer layout may need adjustment based on terminal width
2. **Very long filenames**: Header truncation not yet implemented for file paths
3. **Theme synchronization**: UI theme and code theme are independent - should they be linked?
4. **Cache size**: Currently limited to 100 entries - is this appropriate for typical usage?
5. **Theme switching UI**: Add keybindings to cycle through themes at runtime?
6. **Configuration file**: Persist user's theme preferences?
7. **Line number styling**: Make line numbers themable (currently hardcoded)?

---

## Session 2026-02-08: Fixed Truncation Bug

### Problem
The code view had a truncation bug where one line wasn't truncating properly, causing visual overflow and layout issues.

### Root Causes Found
1. **Line 446**: `if remainingWidth >= 0` should be `if remainingWidth > 0` - when remainingWidth was exactly 0, it still tried to truncate, outputting "..." which caused overflow
2. **Ellipsis in code view**: The `truncate()` function always adds "..." but in the code display area, we want clean truncation at the edge without ellipsis

### Solution Implemented
**File**: `internal/tui/model.go`

1. Changed condition at line 446 from `>= 0` to `> 0`
2. Created new `truncateClean()` function that truncates without adding ellipsis
3. Updated highlighted code path to use `truncateClean()` for clean edge truncation
4. Kept `truncate()` with ellipsis for other uses (like sidebar text)

### Code Changes
```go
// In highlighted code path (line ~446):
if remainingWidth > 0 {  // Changed from >= 0
    truncatedText := truncateClean(token.Text, remainingWidth)  // Changed from truncate
    codeContent.WriteString(token.Style.Render(truncatedText))
}

// New function added:
func truncateClean(s string, maxWidth int) string {
    // Truncates without adding "..."
    // Used for code display where we want clean truncation at the edge
}
```

### Build Status
✅ Tests passing
✅ Build successful

---

## Session 2026-02-08: Changed Walkthrough File Naming Convention

### Changes Made

#### 1. New Naming Convention
- **Old pattern**: `wmti_*.json` (anywhere in project)
- **New pattern**: `*.wmti.json` (only in `.wmti/` directory)
- **Examples**: `auth-flow.wmti.json`, `tutorial.wmti.json`, `self.wmti.json`

#### 2. Auto-Creation of .wmti/ Directory
**File**: `internal/tui/model.go`
- When `wmti` runs, it automatically creates `.wmti/` directory if it doesn't exist
- Creates `self.wmti.json` template file - a self-referential walkthrough that explains the walkthrough format using itself as an example
- Template includes 3 steps showing: structure, steps array, and step properties

#### 3. Updated Init Command
**File**: `cmd/init/init.go`
- Creates `.wmti/` directory
- Creates `example.wmti.json` with starter template
- Provides next steps guidance

#### 4. Updated File Discovery
**File**: `internal/walker/finder.go`
- Changed search pattern to `.wmti/*.wmti.json`
- Removed root-level search (clean separation - all walkthroughs in `.wmti/`)

#### 5. Renamed Existing Files
- `wmti_project_setup.json` → `.wmti/project_setup.wmti.json`
- `wmti_security_measures.json` → `.wmti/security_measures.wmti.json`

#### 6. Updated Documentation
- `AGENTS.md` - Updated naming convention references
- `docs/NOTES.md` - Removed completed task, updated all references
- `docs/SESSIONS.md` - This entry

### Benefits
1. **Cleaner project root** - All walkthroughs in hidden `.wmti/` directory
2. **Clearer file association** - `.wmti.json` extension makes ownership obvious
3. **Better organization** - Single location for all walkthrough files
4. **Self-documenting** - `self.wmti.json` teaches users the format

### Migration Path
Old `wmti_*.json` files in root will no longer be detected. Users should:
1. Move files to `.wmti/` directory
2. Rename from `wmti_name.json` to `name.wmti.json`

### Build Status
✅ Tests passing
✅ Build successful

---

## Session 2026-02-28: Model Refactoring, Testing, and Line Highlighting

### What Was Accomplished

#### 1. Refactored model.go into Multiple Files
Split the 694-line `internal/tui/model.go` into 5 focused files for better maintainability:

- **types.go** - AppState enum, walkthroughItem struct, Model struct, message types
- **helpers.go** - displayWidth, padRight, truncate, truncateClean utility functions
- **loaders.go** - NewModel, Init, file loading functions (async content loading)
- **view.go** - View rendering logic (split-pane layout, code display, sidebar)
- **model.go** - Update state machine only (message handling)

#### 2. Created Comprehensive Test Suite (76 total tests)

| File | Tests | Coverage |
|------|-------|----------|
| `internal/navigator/navigator_test.go` | 20 | Navigation logic |
| `internal/walker/walker_test.go` | 19 | File reading, path security |
| `pkg/models/walkthrough_test.go` | 8 | Data structures |
| `internal/highlighter/highlighter_test.go` | 16 | Syntax highlighting |
| `internal/tui/helpers_test.go` | 13 | Utility functions |

#### 3. Implemented Line Highlighting Feature
Modified code rendering to visually highlight specific lines from the step's line range:

**Files Modified**: `internal/tui/styles.go`, `internal/tui/view.go`
- Added `HighlightLineStyle` with background color (`#1e1e2e`)
- Modified View rendering to apply background highlighting only to lines within `step.LineStart` to `step.LineEnd`
- Preserves syntax highlighting colors while adding background

#### 4. Added Context Lines and Centering
Extended functionality to show more code context around highlighted lines:

- Added new fields to `Model` and `fileContentMsg`:
  - `highlightStart`, `highlightEnd` - The lines to highlight
  - `displayStart`, `displayEnd` - The range of lines to actually display
- Modified `loaders.go` to read a wider range of lines to fill available screen space
- Displayed range is dynamically calculated to center the highlighted section

#### 5. Fixed and Refined Highlighting
- Added `hasHighlightContent` check to prevent panic from out-of-bounds access
- Applied background highlighting to full line width using space padding
- Lightened highlight color from `#2a2a3e` to `#1e1e2e` for better visibility

#### 6. Updated Walkthrough Files
Shortened step definitions in `.wmti/` files to be more focused (1-16 lines each):
- `.wmti/model-organization.wmti.json` - Shows the new file structure
- `.wmti/self.wmti.json` - Self-referential walkthrough

### Files Modified

```
internal/
├── tui/
│   ├── types.go           # NEW - Type definitions
│   ├── helpers.go         # NEW - Utility functions
│   ├── loaders.go         # NEW - File loading logic
│   ├── view.go            # NEW - View rendering
│   ├── model.go           # MODIFIED - Reduced to Update only
│   └── styles.go          # MODIFIED - Added HighlightLineStyle
├── navigator/
│   └── navigator_test.go  # NEW - 20 tests
├── walker/
│   ├── walker_test.go     # NEW - 19 tests
│   └── finder_test.go     # NEW
├── highlighter/
│   └── highlighter_test.go # NEW - 16 tests
└── pkg/
    └── models/
        └── walkthrough_test.go # NEW - 8 tests

internal/tui/
└── helpers_test.go        # NEW - 13 tests

.wmti/
├── model-organization.wmti.json
└── self.wmti.json
```

### Key Technical Details

**Line Highlighting Logic** (in `view.go`):
```go
// Apply background only to highlighted lines
isHighlightedLine := lineNum >= m.highlightStart && lineNum <= m.highlightEnd
if isHighlightedLine {
    lineContent = HighlightLineStyle.Render(lineContent + strings.Repeat(" ", paddingNeeded))
}
```

**Centering Calculation** (in `loaders.go`):
- Reads extra lines above and below the highlight range
- Centers the highlighted region when it doesn't fill the screen
- Falls back gracefully when file is smaller than display area

### Build Status
✅ Tests passing (76 tests)
✅ Build successful

---

## Session 2026-02-28 (Continued): Line Highlighting Improvements, Screen Flickering Fix, Delete Feature, and Walkthrough Creation

### What Was Accomplished

#### 1. Line Highlighting Improvements
- Made line highlighting extend to the end of the line (not just the immediate code)
- Changed highlight color from `#2a2a3e` to lighter `#4a4a5e`
- Added code to fill remaining width with highlight background after processing tokens

#### 2. Screen Flickering Fix
- Initially tried `tea.ClearScreen` which caused flickering
- Tried rendering footer in StateLoadingStep - caused stacking issue
- Tried `tea.Sequence(tea.ClearScreen, ...)` 
- User eventually fixed it themselves with a different approach

#### 3. Delete Walkthrough Feature
- Added `StateConfirmDelete` to AppState enum in `types.go`
- Added `deleteConfirmPath` field to Model in `types.go`
- Added delete key handling (`d` and `Del` keys) in `model.go`
- Added confirmation modal rendering in `view.go`
- Fixed bug where deletion would call `findWalkthroughFiles()` which recreated `self.wmti.json`
- Fixed state handling so app doesn't quit when pressing `q` in delete confirmation

#### 4. Created 5 Walkthrough Files
- Removed existing walkthroughs from `.wmti/` directory
- Created new walkthroughs covering project topics:
  - `file-walker.wmti.json` - File Walker Security
  - `navigation.wmti.json` - Walkthrough Navigation
  - `syntax-highlighting.wmti.json` - Syntax Highlighting
  - `configuration.wmti.json` - Configuration Management
  - `tui-model.wmti.json` - TUI State Management

#### 5. Removed self.wmti.json Creator
- Removed the `createSelfTemplate` function from `loaders.go`
- Removed the auto-creation logic from `findWalkthroughFiles()`

### Files Modified

- `internal/tui/view.go` - Highlighting and delete modal rendering
- `internal/tui/styles.go` - Highlight color change
- `internal/tui/types.go` - Added StateConfirmDelete and deleteConfirmPath
- `internal/tui/model.go` - Delete key handling and deletion logic
- `internal/tui/loaders.go` - Removed self.wmti.json template creation

### Build Status
✅ Tests passing
✅ Build successful

---

## Session 2026-03-03: Auto-Init, Model Selection UI, Footer Improvements, and Delete Confirmation Refactor

### What Was Accomplished

#### 1. Auto-Init for `wmti` Command
**Files**: `cmd/wmti/root.go`, `cmd/init/init.go`, `internal/config/config.go`

- Modified `cmd/wmti/root.go` to automatically run init if not initialized
- Added `isInitialized()` function that checks for `.wmti.json` files in `.wmti/` directory
- Exported `RunInit()` in `cmd/init/init.go` to be called from root.go
- Updated `internal/config/config.go` to include `Models` and `Selected` fields for future model management

#### 2. Model Selection Page (New TUI State)
**Files**: `internal/tui/types.go`, `internal/tui/view.go`, `internal/tui/model.go`, `internal/tui/loaders.go`

- Added `StateModelSelect` to `AppState` enum in `types.go`
- Added model-related fields to Model struct:
  - `modelList` - slice of available models
  - `modelSelectedIndex` - current selection
  - `modelTextInput` - for new model name input
  - `modelIsAdding` - flag for adding mode
  - `modelAddingName` - storing model name during two-step input
  - `modelAskingFor` - tracks "name" or "api_key" input phase
  - `modelDeleteConfirm` - delete confirmation state
- Created menu view in `view.go` with:
  - Highlighted selection (cyan on dark gray background)
  - Green arrow replaced with highlight
  - "Add new model" option at bottom
  - Two-step input: first model name, then API key
- Added key handling in `model.go`:
  - Arrow keys (↑↓ or j/k) for navigation
  - Enter to confirm selection or add new model
  - 'q' types into input field, only 'Esc' cancels
  - 'd' to delete model with confirmation in footer

#### 3. Walkthrough Selection Footer Improvements
**File**: `internal/tui/view.go`

- Replaced Bubble Tea's built-in status bar with custom footer
- Added pagination info ("Page X of Y") to footer
- Added 'd' delete and 'm' model selection options to footer
- Made footer colors lighter/more visible

#### 4. Delete Confirmation Refactor
**Files**: `internal/tui/types.go`, `internal/tui/view.go`, `internal/tui/model.go`

- Unified delete confirmation to use footer instead of centered message
- Shows "Delete <filename>?" in footer with "Enter confirm" and "Esc cancel"
- Works for both walkthrough files and model deletion

### Files Modified

```
cmd/
├── wmti/root.go           # Auto-init logic
└── init/init.go           # Exported RunInit()

internal/
├── config/config.go      # Added Models and Selected fields
└── tui/
    ├── types.go           # Added StateModelSelect, model fields
    ├── view.go            # Model selection view, updated footer
    ├── model.go           # Key handling for model selection
    └── loaders.go         # Model list initialization with dummy data
```

### Next Steps

- Actually storing API keys locally (user mentioned "eventually want to store the API keys on a local root folder")
- Building LLM integration to generate walkthroughs
- Other model management features

### Build Status
✅ Tests passing
✅ Build successful

---

## Session 2026-03-05: API Key Storage, UI Improvements, and Security Fixes

### What Was Accomplished

#### 1. API Key Storage Investigation and Fix
**Issue**: API keys were stored in plaintext in `~/.config/wmti/config.yaml` but weren't actually being saved when entered in the TUI.

**Investigation**: Found that opencode stores API keys in `~/.local/share/opencode/auth.json` (also plaintext).

**Solution**: Implemented opencode-style auth storage:
- Created `internal/config/auth.go` - new auth package storing keys in `~/.local/share/wmti/auth.json`
- Removed API key field from main config (`internal/config/config.go`)
- Fixed bug in `internal/tui/model.go` to actually save API keys when entered
- Also deletes API key from storage when model is deleted

#### 2. UI Improvements
**Files**: `internal/tui/model.go`, `internal/tui/view.go`, `internal/tui/styles.go`

- Changed model delete confirmation from Enter key to 'y'/'n' keys (consistent with walkthrough deletion)
- Added duplicate model name validation with error display
- Added `ErrorTextStyle` in `internal/tui/styles.go` (red color)
- Added `modelError` field to Model struct for tracking/displaying errors

#### 3. Security Fixes
**Files**: `internal/walker/walker.go`, `internal/navigator/navigator.go`, `pkg/models/walkthrough.go`

Fixed 4 security vulnerabilities:

1. **Manual path entry not sanitized** - Fixed by validating paths in `StateEnteringPath` state
2. **Walkthrough file references could escape root** - Fixed by validating each step's path in navigator
3. **No line range validation** - Fixed by adding `LineStart > 0` and `LineEnd >= LineStart` checks
4. **No input length limits** - Fixed by adding `CharLimit` to text inputs (512 for paths, 64 for model names)

### Files Modified

```
cmd/wmti/root.go                     # Auto-init logic
cmd/init/init.go                      # Exported RunInit()

internal/
├── config/
│   ├── auth.go                      # NEW - Auth storage (~/.local/share/wmti/auth.json)
│   └── config.go                    # Removed API key, added Models/Selected

internal/
├── tui/
│   ├── types.go                     # Added modelError field
│   ├── view.go                      # Error display, y/n delete confirmation
│   ├── model.go                     # API key saving, duplicate validation
│   ├── loaders.go                   # Model list initialization
│   └── styles.go                    # Added ErrorTextStyle

internal/
├── walker/
│   ├── walker.go                    # Path sanitization
│   └── walker_test.go               # Updated tests

internal/
├── navigator/
│   └── navigator.go                 # Step path validation

pkg/
└── models/
    ├── walkthrough.go               # Line range validation
    ├── errors.go                    # Error types
    └── walkthrough_test.go          # Updated tests
```

### Current State

All tests pass, build succeeds. The application now has:
- Separate auth storage for API keys (matching opencode's approach)
- y/n confirmation for model deletion
- Duplicate model name prevention with error display
- Security hardening against path traversal and DoS attacks

### Potential Next Steps

- Add environment variable support for API keys (e.g., `WMTI_API_KEY`)
- Implement actual API integration for the `internal/api` package
- Add keychain support for extra security
- Consider encrypting the auth.json file

### Build Status
✅ Tests passing
✅ Build successful
