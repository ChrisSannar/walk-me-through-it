# Walk Me Through It (wmti)

A Neovim plugin that walks a developer through code by stepping the cursor to relevant locations, each with a short explanation. The plugin is a viewer; the content is authored by a coding agent.

## Language

**Walkthrough**:
An ordered sequence of steps, authored as one JSON file, that guides a reader through a slice of the codebase. Serves both quick change-overviews and durable onboarding tours.
_Avoid_: Tour, guide, tutorial, changes file.

**Step**:
A single stop in a walkthrough: a file and line range to jump to, with a `title` and a `why`.
_Avoid_: Stop, entry, item, node.

**Why**:
A step's prose explaining *why* the code matters or why it changed — intent, not a restatement of what the code does. Covers both modes (changed-because / matters-because).
_Avoid_: Description, body, comment, note.

**Anchor**:
The text of a step's start line plus the range length, used to relocate the step when line numbers have drifted since authoring.
_Avoid_: Marker, fingerprint, signature.

**Drift**:
The divergence between a step's authored line numbers and the file's current state. Resolved at view time via the anchor.
_Avoid_: Staleness, skew (reserve "stale" for the base_sha mismatch warning).
