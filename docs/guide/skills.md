# Skills

appteam generates 5 Claude Code slash commands in `.claude/skills/`. These are markdown files that define instructions for Claude Code to follow when a user invokes the command.

## `/spec` — Create a new product spec

**File:** `.claude/skills/spec.md`

Creates a new product spec in `docs/specs/` with the next sequential `F-NNNN` ID.

**Usage:**
```
/spec Add dark mode support
```

**What it does:**
1. Scans `docs/BACKLOG.md` for the highest existing `F-NNNN` ID
2. Creates `docs/specs/F-NNNN-short-slug.md` from the template
3. Fills in the spec with title, type, priority, requirements, and acceptance criteria
4. Reports the file path and summary

## `/release` — Cut a new release

**File:** `.claude/skills/release.md`

Updates release notes, commits, tags, and pushes a new version.

**Usage:**
```
/release v0.6.0
```

If no version is provided, it reads `docs/RELEASENOTES.md` and suggests the next version bump.

**What it does:**
1. Adds a new version entry to `docs/RELEASENOTES.md` (Added/Changed/Fixed)
2. Commits all pending changes
3. Creates an annotated git tag
4. Pushes the commit and tag to origin

## `/pipeline` — Spin up the full agent team

**File:** `.claude/skills/pipeline.md`

Creates a full agent team via tmux for a feature or milestone.

**Usage:**
```
/pipeline Add user authentication
```

**What it does:**
1. Creates a team with `TeamCreate`
2. Creates tasks for PM, TPM, SWEs, testing, and review
3. Spawns all agents in tmux panes
4. Assigns tasks and follows the mandatory pipeline
5. Shuts down agents and cleans up when complete

## `/status` — Milestone status summary

**File:** `.claude/skills/status.md`

Provides a quick status report of the current milestone.

**Usage:**
```
/status
```

**What it does:**
1. Reads `docs/BACKLOG.md` and `docs/PROGRESS.md`
2. Summarizes completed, in-progress, blocked, and pending items
3. Reports team utilization and key decisions
4. Lists immediate next actions

## `/regenerate` — Regenerate from saved settings

**File:** `.claude/skills/regenerate.md`

Regenerates all project files from `.appteam/settings.json` without running the wizard.

**Usage:**
```
/regenerate
```

**What it does:**
1. Verifies `.appteam/settings.json` exists
2. Runs `appteam -r`
3. Restores tracking files (`BACKLOG.md`, `PROGRESS.md`, `RELEASENOTES.md`) from git so they aren't overwritten
4. Reports which files were regenerated

**Note:** Tracking files contain project-specific history that would be lost if overwritten by the template. The skill handles this automatically.
