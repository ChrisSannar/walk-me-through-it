# wmti — walk me through it

A dumb, dependency-free Neovim plugin that walks you through a codebase. Point it
at a **Walkthrough** — an ordered list of steps stored as JSON in `.wmti/` — and it
steps your editor through each spot, jumping the code window to the right lines and
showing a short explanation of *why* they matter in a side panel.

It is intentionally a **viewer only**: it reads a Walkthrough and navigates it. It
never computes diffs, calls an LLM, or generates content. Walkthroughs are written
separately (by hand, or by a coding agent) against the [schema](#walkthrough-format)
below. Because the viewer is git-agnostic, the same tool serves recent-change
overviews, onboarding tours, and PR reviews — the difference is just what you put in
the JSON.

Walkthroughs are written either by hand or by a coding agent guided by the bundled
[authoring skill](#authoring) — see [`.wmti/sample.json`](.wmti/sample.json) for a worked
example.

## Requirements

- Neovim **0.10+** (uses `vim.json`, `vim.system`, `vim.ui.select`).
- **Zero** third-party plugins.

## Install

### lazy.nvim — from a local checkout (development)

```lua
{
  dir = "/path/to/walk-me-through-it",
  name = "wmti",
  cmd = { "Wmti" },
  -- Optional: lets you type lowercase `wmti` / `wmti next` instead of `:Wmti`.
  init = function()
    vim.cmd([[cnoreabbrev <expr> wmti (getcmdtype() == ':' && getcmdline() ==# 'wmti') ? 'Wmti' : 'wmti']])
  end,
  -- Optional: no default keymaps ship; opt into your own.
  keys = {
    { "<leader>wo", "<cmd>Wmti<cr>",       desc = "wmti: open" },
    { "<leader>wn", "<cmd>Wmti next<cr>",  desc = "wmti: next" },
    { "<leader>wp", "<cmd>Wmti prev<cr>",  desc = "wmti: prev" },
    { "<leader>wq", "<cmd>Wmti close<cr>", desc = "wmti: close" },
  },
  opts = {},
}
```

`opts = {}` calls `require("wmti").setup({})`. See [Configuration](#configuration).

## Usage

Walkthroughs are discovered in `.wmti/*.json` **relative to your current working
directory**, so launch Neovim from the project root.

One command, `:Wmti`, with subcommands (Vim command names must be uppercase +
alphanumeric, so there is no literal lowercase `wmti` or `wmti-next` command — the
abbreviation above bridges that):

| Command | What it does |
|---------|--------------|
| `:Wmti` / `:Wmti open` | Open: 0 found → notify · 1 → auto-open · many → picker |
| `:Wmti select` | Always show the picker, even with one Walkthrough |
| `:Wmti next` / `:Wmti prev` | Step forward / back (stops at the ends) |
| `:Wmti dive` | Move focus into the code window at the current step |
| `:Wmti close` | Close the sidebar and end the session |

Subcommands Tab-complete.

### In the sidebar

Focus stays in the sidebar so you can keep stepping (sidebar-driven preview); the
code window follows without stealing focus.

| Key | Action |
|-----|--------|
| `n` / `<Tab>` | Next step |
| `p` / `<S-Tab>` | Previous step |
| `<CR>` | Dive into the code at the current step |
| `q` | Close |

## Configuration

```lua
require("wmti").setup({
  side = "right", -- "right" | "left" — which side the sidebar opens on
  width = 44,     -- sidebar width in columns
  dir = ".wmti",  -- directory (relative to cwd) scanned for *.json Walkthroughs
})
```

The sidebar is color-coded with highlight groups that link to standard groups with
`default = true`, so your colorscheme wins. Override any of them after `setup`:
`WmtiTitle`, `WmtiRule`, `WmtiStale`, `WmtiStepCurrent`, `WmtiStepIdle`,
`WmtiLocation`, `WmtiFlag`, `WmtiWhy`, `WmtiHelp`, `WmtiCurrentBg` (sidebar) and
`WmtiStep` (the highlighted range in the code).

## Walkthrough format

A Walkthrough is JSON in `.wmti/`. Steps are shown in array order — there are no ids.

```jsonc
{
  "title": "Recent changes: auth refactor",   // required
  "steps": [                                    // required, >= 1
    {
      "file": "src/auth/login.go",  // required — relative to project root, forward slashes
      "line_start": 10,             // required — 1-based
      "line_end": 24,               // required — 1-based, >= line_start
      "title": "Issue the session", // required — short
      "why": "Validates credentials, then issues a 24h JWT.", // required — intent, not a restatement
      "anchor": "func issueToken(u *User) (string, error) {"  // optional — text of the START line
    }
  ],
  "version": "1.0.0",       // optional, semver
  "description": "What changed and why", // optional
  "base_sha": "a1b2c3d"     // optional — enables the staleness note
}
```

### How locations survive drift

Code moves after a Walkthrough is written. At view time each step is resolved
against the file as it exists now:

1. If a step has an `anchor`, search the buffer for that line — **exact** match
   first, then **whitespace-trimmed** — and pick the occurrence **nearest** the
   recorded `line_start`, re-spanning by the original length.
2. If the anchor can't be found, jump to the recorded `line_start` and flag the
   step **may have moved**.
3. If the file is gone, flag the step **file missing** (the rest of the tour still
   works).
4. With no `anchor`, the raw `line_start`/`line_end` range is used as-is.

If `base_sha` is set and `HEAD` differs (in a git repo), the sidebar shows a
non-blocking **authored against older code** note. Outside git, or with no
`base_sha`, staleness checks stay silent.

## Authoring

Walkthroughs can be written by hand, but the intended path is to let the coding agent
that made a change also write the tour — it already holds the *why*. A portable,
self-contained Agent Skill ships in this repo for that:

```
skills/wmti/
  SKILL.md                  # authoring directions (user-invoked)
  walkthrough.schema.json   # the JSON Schema (draft 2020-12) contract
  INSTALL.md                # how to add the skill to another agent / project
```

It's plain markdown + JSON, so it works with any agent that supports the Agent Skills /
`SKILL.md` standard (Claude Code, Codex CLI, Gemini CLI, Cursor, Copilot, …). In Claude
Code, invoke it with `/wmti`. To add it to another project or agent, see
[`skills/wmti/INSTALL.md`](skills/wmti/INSTALL.md). The schema doubles as the validation
contract.

## Development

The logic lives in a single pure module, `lua/wmti/core.lua` (load / resolve /
staleness), with everything else a thin Neovim shell around it:

```
plugin/wmti.lua          :Wmti command + subcommand dispatch (lazy entry)
lua/wmti/init.lua        setup() + open/select orchestration
lua/wmti/config.lua      defaults + setup merge
lua/wmti/walkthrough.lua discover .wmti/*.json, read → core.load (the I/O boundary)
lua/wmti/core.lua        pure load / resolve / staleness — no Neovim, git, or fs
lua/wmti/nav.lua         session: jump the code window, drive the sidebar, git staleness
lua/wmti/sidebar.lua     the side panel: rendering + buffer-local keymaps
```

### Tests

Zero-dependency, run under headless Neovim from the repo root:

```bash
nvim --headless -l tests/core_spec.lua    # pure core unit tests
nvim --headless -l tests/shell_smoke.lua  # end-to-end shell smoke test
```

Both print `ok`/`FAIL` lines and exit non-zero on failure.
