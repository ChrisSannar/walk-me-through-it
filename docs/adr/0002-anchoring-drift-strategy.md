# Drift is resolved at view time by start-line-text anchoring, not extmarks

A step stores its line range plus an optional `anchor` (the text of its start line). When opening a walkthrough, the viewer locates each step by trying the original `line_start`, then searching the buffer for the anchor (exact match first, then whitespace-trimmed), choosing the occurrence nearest the original `line_start`. On a miss it jumps to the original line and flags the step "may have moved". Range length is derived from `line_start`/`line_end`.

## Why

Walkthroughs are authored against a snapshot but viewed later against live, edited buffers, so line numbers drift. Anchoring relocates a step deterministically at load time without any persistent state in the file.

## Considered and rejected

- **Neovim extmarks.** The obvious tool, but extmarks "cannot be preserved on writing and loading a buffer to file" — they don't survive reload or a new session, so they cannot relocate a step authored earlier. (Extmarks are still fine for keeping a highlight stable *during* a live viewing session.)
- **Full-range snippet text.** Bulky in JSON and brittle — any change to an interior line breaks the match. Start-line text + derived length is smaller and more tolerant.

## Consequences

- An optional top-level `base_sha` gives a cheap, opportunistic "authored against older code" warning when the repo is at a different HEAD — a global signal anchoring alone can't provide.
