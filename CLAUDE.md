# Agent Instructions for wmti

> This file is the single source of truth for AI agents. `CLAUDE.md` is a hardlink to this file — edit one, both change. Do not split them.

## Project

**Walk Me Through It (`wmti`)** is a terminal-based, IDE-like TUI written in Go that guides developers through unfamiliar codebases via structured, JSON-based walkthroughs. It reads source files in read-only mode and shows focused, syntax-highlighted line ranges alongside step descriptions. It can also call an LLM provider to generate walkthroughs.

**Tech stack:** Go 1.25.6 · Bubble Tea (TUI, Elm-style model/update/view) · Bubbles (components) · Lipgloss (styling) · Chroma (syntax highlighting) · Cobra (CLI) · Viper (config).

## Commands

```bash
make build      # builds to ./build/wmti
make test       # go test -v ./...
make lint       # golangci-lint run (config in .golangci.yml)
make fmt        # go fmt ./... then gofumpt -w .  (gofumpt required)
make run        # go run ./cmd/wmti
make run-init   # go run ./cmd/wmti init
make dev        # auto-reload (requires air)

go test ./internal/walker/                       # one package
go test -run TestSanitizePath ./internal/walker/ # one test
```

`gofumpt` and `air` are installed separately via `go install`. **Do not run build/test unless asked** — the user runs them when ready. `.golangci.yml` is strict (gosec, gocyclo min-complexity 15, lll 120, gomnd, etc.).

## Architecture

Layering is strict: `cmd/` (Cobra commands) → `internal/` (logic) → `pkg/models` (shared JSON schema types). `internal` never imports `cmd`.

- **`cmd/wmti`** — `root.go` is the root Cobra command. Running `wmti` auto-runs init when `.wmti/` has no `*.wmti.json`, then launches the Bubble Tea program. `cmd/init` creates `.wmti/` plus a self-referential `example.wmti.json` tutorial.
- **`internal/tui`** — the Bubble Tea app, split by concern (not one model file):
  - `types.go` — the `Model` struct, the `AppState` enum (state machine: `StateSelecting`, `StateEnteringPath`, `StateLoading`, `StateViewing`, `StateLoadingStep`, `StateConfirmDelete`, `StateModelSelect`, `StateProviderSelect`, `StateNewWalkthrough`), and all `tea.Msg` types.
  - `loaders.go` — `NewModel()` and the `tea.Cmd` closures that do async I/O. **All file/network work happens inside these commands, never in `Update`** (don't block the TUI).
  - `model.go` — the `Update(msg)` switch; key handling is keyed off `m.state`.
  - `view.go` / `styles.go` / `helpers.go` — rendering, Lipgloss styling, helpers.
- **`internal/navigator`** — parses walkthrough JSON into `models.Walkthrough`, tracks the current step, exposes `NextCycle`/`PreviousCycle` (Tab/Shift+Tab wrap around). On load it validates the schema and that each step `file` stays within root.
- **`internal/walker`** — read-only file access and the security boundary. `ReadFileLines`: `SanitizePath` blocks traversal/symlink escape, extensions are allowlisted, files >10MB rejected, opens `O_RDONLY`, records an audit string (`GetLastAuditMessage`). `finder.go` globs `.wmti/*.wmti.json`. Route new file access through `Walker`, not raw `os`.
- **`internal/highlighter`** — Chroma highlighting with a per-line cache; emits tokens as Lipgloss styles. Do **not** apply Chroma background colors — the code pane stays transparent.
- **`internal/api`** — LLM integration. `providers.go` is a static registry (`Providers`): google, openai, anthropic, groq, ollama, each with base URL, auth scheme, default models. `client.go`'s `Chat` dispatches on `provider.ID` to per-provider request/response shapes. Note: `Client` currently sends `provider.DefaultModels[0]`, not the selected model.
- **`internal/config`** — two separate stores. `config.go` (Viper, `~/.config/wmti/config.yaml`) holds the model list and selection. `auth.go` (`~/.local/share/wmti/auth.json`, mode `0600`) holds API keys keyed by model name. **API keys live only in `auth.json`**, never in the Viper config.
- **`pkg/models`** — `Walkthrough`/`Step` structs and `Walkthrough.Validate()` (step IDs sequential from 1; non-empty `file`; `line_start > 0` and `line_end >= line_start`).

## Themes

- UI themes (`styles.go`): `ThemeDark` (default), `ThemeLight`, `ThemeMidnight`.
- Code themes (`highlighter.go`): `Dracula` (default), `Monokai`, `GitHub`, `OneDark`, `Solarized`, `Vim`.

## Navigation keys

- **Tab** — next step (cyclic) · **Shift+Tab** — previous step (cyclic)
- **q** — back to file selection when viewing, otherwise quit · **Ctrl+C** — quit
- **Enter** — select item · **j/k** — vim-style list nav (loops at ends)

## Walkthrough documents

Live in `.wmti/` and **must** be named `*.wmti.json` (kebab-case, e.g. `auth-flow.wmti.json`; the old `wmti_*.json` pattern is deprecated). They are valid JSON validated against the schema below.

| Top-level | Type | Notes |
|-----------|------|-------|
| `title` | string | shown in selection list / header |
| `description` | string | 1–2 sentences: scope and audience |
| `version` | string | semver, e.g. `"1.0.0"` |
| `steps` | array | ≥ 1 step object |

| Step field | Type | Notes |
|------------|------|-------|
| `id` | int | sequential from 1, no gaps |
| `title` | string | 5–7 words |
| `description` | string | explain *why* the code matters, not just what |
| `file` | string | relative path from project root, forward slashes |
| `line_start` | int | 1-based |
| `line_end` | int | 1-based, `>= line_start`; keep ranges focused (~10–50 lines) |
| `action` | string | only `"read"` is supported |

### Authoring guidance (for LLM-generated walkthroughs)

1. Explore first: find entry points, trace the data flow, group related logic into steps; skip boilerplate.
2. Order steps as a logical progression (entry points → flow → complexity), simple to complex.
3. Write descriptions that explain intent and connect to adjacent steps. Bad: "This is the login function." Good: "Validates credentials with bcrypt, then issues a 24h JWT."
4. Use relative paths from project root (`src/auth/login.go`), never absolute or `../`.
5. Verify files exist and line numbers are accurate before finishing.

Common mistakes: absolute paths; gaps in step IDs; zero-based line numbers (first line is 1); ranges too large to follow; vague descriptions; missing `version`.

## Conventions & gotchas

- Standard Go style (`gofmt`/`goimports`, local prefix `github.com/chrissannar/walk-me-through-it`). Return errors with context; don't panic. Pass `context.Context` first when needed. Avoid global state; pass dependencies explicitly. Don't edit `go.mod` by hand — use `go get`/`go mod tidy`.
- `truncateClean()` truncates without an ellipsis (used for the code view); `truncate()` adds an ellipsis (used elsewhere).
- List arrow keys and j/k need custom handling to loop first↔last.
- Use `lipgloss.Height()` for actual rendered dimensions.
- Tests are table-driven `*_test.go`; sample data in `testdata/`.

## Branch workflow

`main` is production — **never commit directly**. `dev` is the integration branch and the PR target. Do feature work on `feature/*` branches cut from the current working branch. Small, focused commits.

## Docs

`docs/SESSIONS.md` (development session history — update it after significant changes) and `docs/NOTES.md` (feature backlog: line highlighting for the current step, chat box in content view, rescan for new walkthroughs, auto-detect AI keys during init, local storage of walkthroughs).
