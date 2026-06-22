# Walk Me Through It (wmti)

A Neovim plugin that walks a developer through code by stepping the cursor to relevant locations, each with a short explanation. The plugin is a viewer; the content is authored by a coding agent.

## Language

### Artifacts

**Walkthrough**:
An ordered sequence of steps, authored as one JSON file, that guides a reader through a slice of the codebase. Whether it is a throwaway change-overview or a durable onboarding tour is a matter of lifespan, not a separate kind of thing.
_Avoid_: Tour, overview, guide, tutorial, changes file.

**Step**:
A single stop in a walkthrough: a file and line range to jump to, with a `title` and a `why`.
_Avoid_: Stop, entry, item, node.

**Why**:
A step's prose explaining *why* the code matters or why it changed — intent, not a restatement of what the code does. Covers both modes (changed-because / matters-because).
_Avoid_: Description, body, comment, note.

### Roles

**Author**:
Whoever produces a walkthrough — by default a coding agent guided by the authoring skill, occasionally a human writing one by hand.
_Avoid_: User, creator, generator.

**Reader**:
The developer stepping through a walkthrough in Neovim. The viewer serves the reader; it never authors.
_Avoid_: User, viewer (the viewer is the plugin, not the person).

### Resolution

**Anchor**:
The text of a step's start line, used to relocate the step when its line numbers have moved since authoring. The span is taken from the step's line range, not stored on the anchor.
_Avoid_: Marker, fingerprint, signature.

**Drift**:
A single step's authored line numbers no longer matching where its code currently sits. Resolved at view time via the anchor.
_Avoid_: Skew, slippage.

**Stale**:
A whole walkthrough that was authored against an earlier commit than the one being viewed (its base differs from current HEAD). Distinct from drift: a walkthrough can be stale yet have every step still anchor cleanly, or be current yet have a drifted step.
_Avoid_: Outdated, drifted (drift is per-step, staleness is per-walkthrough).
