---
name: wmti
description: Author a wmti Walkthrough — an ordered JSON tour of a change set or codebase that the wmti Neovim viewer steps a reader through.
disable-model-invocation: true
---

# Authoring a wmti Walkthrough

A **Walkthrough** is one JSON file in `.wmti/` that guides a **Reader** through a slice
of code: an ordered list of **Steps**, each a `file` + line range with a `title` and a
**Why**. The wmti viewer is a *dumb reader* — it does no generation and never fixes a
bad Walkthrough — so the JSON you write here is the whole product. Get the `why` and the
`anchor` right.

You are the **Author**. Whether this Walkthrough is a throwaway change-overview or a
durable onboarding tour is just lifespan; the schema is identical.

Full contract: [`walkthrough.schema.json`](./walkthrough.schema.json).

**Style — lean, not bare.** Favour many small Steps over a few big ones: one idea per
Step, a focused range. Each `why` is one or two sentences that carry real intent — why it
matters, or how it connects — without padding into a paragraph or restating the code.
Substance, not sprawl. When a topic is large, split it across several short Walkthroughs
(see [a narrative across walkthroughs](#a-narrative-across-walkthroughs)) rather than
letting one sprawl.

## Process

### 1. Scope the tour

Decide what the Reader should be walked through:

- A **change set** — `git diff`, a branch's commits, or the edits you just made.
- An **area to onboard** — an entry point and the flow beneath it.

Completion: you can state in one sentence what the Reader will understand at the end.

### 2. Explore before writing

Find the entry points, trace the data flow, and read enough to know *why* each part
matters. Break the flow into small, single-idea Steps; skip boilerplate (imports,
getters, generated code).

Completion: you have a list of candidate Steps, each tied to a specific `file` + lines.

### 3. Order the Steps

Sequence them as a logical progression — entry points → flow → complexity, simple to
complex. Array order *is* Step order; there are no ids.

Completion: reading the titles top-to-bottom tells a coherent story.

### 4. Write each Step

For every Step:

- `file` — path relative to the project root, forward slashes. Never absolute, never `..`.
- `line_start` / `line_end` — **1-based**, `line_end >= line_start`. One idea per Step — a
  tight range (often ~5–25 lines), not a whole function dumped in.
- `title` — ~5–7 words.
- `why` — the **Why**, one or two sentences: the *intent* (why it matters / why it
  changed) plus, where it earns it, how it connects to the next Step or the bigger
  picture. Trim padding and anything the `title` or code already says — but leave enough
  to actually inform.
  - Too long: "This function validates the user's credentials against the database using bcrypt hashing, and if they are correct it then proceeds to generate and return a signed JWT that expires after 24 hours, which is later consumed by the API layer."
  - Too thin: "Checks the password, then mints the JWT."
  - Good: "Hashes and verifies the password, then mints a 24h JWT — the token step 4 trusts. This is the whole auth boundary."

Completion: every Step has all five required fields and a `why` of one or two informative sentences.

### 5. Set the anchor (this is what survives drift)

Set `anchor` to the **verbatim text of the `line_start` line** — copy it exactly,
whitespace and all. This is mandatory in practice: it's how the viewer relocates a Step
when code moves after authoring (see [Why anchors matter](#why-anchors-matter)).

Completion: each Step's `anchor` is byte-for-byte the current text of its start line.

### 6. Write the file

Write the Walkthrough to `.wmti/<kebab-name>.json` (e.g. `.wmti/auth-refactor.json`).
Set `title` and `steps`; add `version`, `description`, and `base_sha` (the commit you
authored against — enables the viewer's staleness note) when useful.

Completion: the JSON parses and conforms to the schema.

### 7. Verify

Before finishing, check every Step:

- the `file` exists,
- `line_start`/`line_end` are accurate and 1-based,
- the `anchor` equals the start line **exactly**.

Completion: every Step verified against the real files — no guesses.

## A narrative across walkthroughs

One Walkthrough should stay tight. When a topic is too big for that, split it into a
series of short Walkthroughs that read in order rather than one sprawling file:

- Name them with an ordering prefix — `01-entry.json`, `02-flow.json`, `03-internals.json`.
  The picker lists files sorted, so the prefix sets the reading order.
- Give each a focused `title` and `description` for its slice.
- Let the last Step's `why` hand off to the next Walkthrough ("continued in 02-flow").

Each piece stays glanceable, and the Reader chooses how deep to go.

## Why anchors matter

The viewer resolves each Step's location *at view time*, against the file as it exists
now, so it lands the Reader on the right code despite **Drift** (a Step's lines having
moved since authoring):

1. Search the buffer for the `anchor` — **exact** match first, then **whitespace-trimmed**.
2. Among matches, pick the one **nearest** the recorded `line_start`; re-span by the range length.
3. No match → jump to the recorded `line_start` and flag the Step **"may have moved"**.
4. No `anchor` → use the raw line range as-is.

So a wrong or missing `anchor` doesn't error — it silently degrades to a line-number
guess that may send the Reader to the wrong place. That's why step 5 is non-negotiable.

`base_sha` is separate: it's per-Walkthrough **Stale** detection. If set and the repo's
`HEAD` differs, the viewer shows a non-blocking "authored against older code" note. A
Walkthrough can be stale yet still anchor every Step cleanly.

## Minimal example

```json
{
  "version": "1.0.0",
  "title": "Recent changes: auth refactor",
  "base_sha": "a1b2c3d",
  "steps": [
    {
      "file": "src/auth/login.go",
      "line_start": 10,
      "line_end": 24,
      "anchor": "func issueToken(u *User) (string, error) {",
      "title": "Issue the session token",
      "why": "Checks the password, then mints the 24h JWT step 3 trusts."
    }
  ]
}
```

## Common mistakes

- Absolute paths or `..` in `file` (use a path relative to the project root).
- Zero-based line numbers (the first line is 1).
- Ranges too large to take in (split them).
- A `why` that restates the code instead of its intent.
- An `anchor` that doesn't match the start line verbatim — or omitting it.
- Missing top-level `title` or `steps`.

## Installing this skill elsewhere

This folder is self-contained (the schema travels with it). To add it to another agent
or project, see [`INSTALL.md`](./INSTALL.md).
