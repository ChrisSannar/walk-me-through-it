# Plugin is a git-agnostic viewer; walkthrough generation is offloaded to an agent skill

`wmti` is a Neovim plugin whose only job is to read a JSON walkthrough and step the cursor through it. It does not compute diffs, run git, or call an LLM. Walkthroughs are authored by a coding agent, guided by a portable Agent Skill (`SKILL.md`) that carries the schema and authoring directions.

## Why

The valuable, hard part — deciding what matters and *why* — is judgement the coding agent already has in context while it writes code. Keeping that out of the plugin makes the viewer deterministic, testable, dependency-free, and reusable for any walkthrough (change overview, onboarding tour, PR review), not just "recent changes".

## Considered and rejected

- **In-plugin `git diff` parsing** + a Claude Code **PostToolUse hook** that auto-records changed hunks. Rejected: heavy, fragile plumbing (hook + git + agent/git reconciliation), and it hard-wires the tool to git and to one agent's hook system.

## Consequences

- No automatic safety net: if no agent (or human) authors a walkthrough, there is none. Accepted deliberately.
- The authoring skill must be portable across agents (Agent Skills / `SKILL.md` open standard), not Claude-only.
