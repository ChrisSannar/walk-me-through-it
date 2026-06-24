# Installing the `wmti` authoring skill

This folder is a self-contained, portable Agent Skill — `SKILL.md` plus its schema. It
teaches a coding agent to author a **Walkthrough** (`.wmti/*.json`) that the
[wmti Neovim viewer](https://github.com/chrissannar/walk-me-through-it) steps through.

It's plain markdown + JSON, so it works with any agent. The only per-agent difference is
**how you make the agent aware of it**. Copy the whole `wmti/` folder into the target
project first:

```bash
cp -r /path/to/walk-me-through-it/skills/wmti /path/to/other-project/skills/wmti
```

Then wire it up for your agent below.

> The skill only *writes* the JSON. To *view* a Walkthrough you also need the wmti
> Neovim plugin installed — the skill and viewer are decoupled by design.

## Claude Code

Claude Code auto-discovers skills under `.claude/skills/` and lets you invoke them by
name. Symlink (keeps the committed copy in `skills/` as the single source of truth):

```bash
cd /path/to/other-project
mkdir -p .claude/skills
ln -sfn ../../skills/wmti .claude/skills/wmti
```

(A plain `cp -r skills/wmti .claude/skills/wmti` works too, but then you maintain two
copies.) Invoke it with `/wmti`. Because the skill is user-invoked
(`disable-model-invocation: true`), Claude won't fire it on its own — you trigger it.

## Codex

Codex has no skills auto-loader; its instruction file is `AGENTS.md`. Make the agent
aware of the skill by referencing it from the project's root `AGENTS.md`:

```markdown
## Authoring walkthroughs

When asked to write or update a wmti walkthrough, follow `skills/wmti/SKILL.md` and
validate the result against `skills/wmti/walkthrough.schema.json`. Write it to
`.wmti/<kebab-name>.json`.
```

Then just ask: *"Author a wmti walkthrough of the changes on this branch, following the
wmti skill."*

Optional — a slash command. If your Codex version supports custom prompts
(`~/.codex/prompts/*.md`), add `~/.codex/prompts/wmti.md`:

```markdown
Read skills/wmti/SKILL.md in this repo and follow it to author a wmti walkthrough for:
$ARGUMENTS
Validate against skills/wmti/walkthrough.schema.json before finishing.
```

Then `/wmti recent auth changes`.

## Any other agent (Cursor, Copilot, Gemini CLI, …)

Same two-step pattern:

1. Copy the `wmti/` folder into the project (a path the agent can read).
2. Reference it from whatever instruction file that agent reads (`AGENTS.md` is the
   common one), telling it to follow `skills/wmti/SKILL.md` and validate against
   `skills/wmti/walkthrough.schema.json` when authoring a walkthrough.

The skill folder itself never changes per agent — it's the single source of truth.
