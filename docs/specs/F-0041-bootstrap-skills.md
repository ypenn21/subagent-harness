# F-0041: Bootstrap Skills — Generate Claude Code Skills

**Type:** Feature
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-16

## Problem

Claude Code supports user-invocable skills (slash commands) defined as markdown files in `.claude/skills/`. Currently, appteam generates agent definitions but no skills, leaving teams without pre-built automation for common workflows like creating specs, releasing versions, running the pipeline, checking status, or regenerating configs.

Generating a starter set of skills that are parameterized with the project's config would immediately bootstrap teams with useful automation.

## Requirements

1. **Generate `.claude/skills/` directory** — The generator must create `.claude/skills/` alongside `.claude/agents/` in the target directory
2. **Five skill files** — Generate the following five skill markdown files, each as a Go template rendered with `ProjectConfig`:

   ### `/spec` — `spec.md`
   - Instructs the PM agent to create a new product spec
   - Reads `docs/BACKLOG.md` to determine the next sequential `F-NNNN` ID
   - Copies `docs/specs/TEMPLATE.md` and fills in the feature title provided by the user
   - Opens the spec for editing
   - References project-specific paths (project name, owner)

   ### `/release` — `release.md`
   - Updates `docs/RELEASENOTES.md` with a new version entry
   - Bumps the version constant in the codebase (if applicable)
   - Commits changes with a release commit message
   - Creates an annotated git tag (`vX.Y.Z`)
   - Pushes the commit and tag to origin
   - Uses configured git author (`{{.OwnerName}}` / `{{.OwnerEmail}}`)

   ### `/pipeline` — `pipeline.md`
   - Spins up the full agent team using `TeamCreate`
   - Creates tasks for the specified feature/milestone
   - Spawns PM, TPM, SWE, and SWE-Test agents in tmux panes
   - Follows the mandatory pipeline workflow defined in CLAUDE.md
   - Uses configured model (`{{.ModelID}}`) for all agents

   ### `/status` — `status.md`
   - Reads `docs/BACKLOG.md` and `docs/PROGRESS.md`
   - Summarizes the current milestone: what's done, what's in progress, what's blocked
   - Reports team utilization (which SWEs are assigned to what)
   - Provides a concise status update suitable for PO review

   ### `/regenerate` — `regenerate.md`
   - Runs `appteam -r` to regenerate all files from saved `.appteam/settings.json`
   - Verifies the settings file exists before running
   - Reports what files were regenerated

3. **Go template rendering** — Each skill must be a Go string constant template in `internal/templates/` (following the existing pattern), rendered with `ProjectConfig` data so all project-specific values (project name, owner, email, model, paths) are baked in
4. **Correct Claude Code skill format** — Each skill file must follow the Claude Code skill markdown format:
   - Clear title and description
   - Step-by-step instructions for Claude to execute
   - No hardcoded project-specific values — all parameterized via templates
5. **Generator integration** — Add the five skill files to the `files` slice in `generator.Generate()`, so they are created alongside agents and docs
6. **Skills directory creation** — Create `.claude/skills/` directory in the generator (similar to how `.claude/agents/` is created)

## Acceptance Criteria

- [ ] Running `appteam` (wizard or `-r`) generates `.claude/skills/` directory with 5 files: `spec.md`, `release.md`, `pipeline.md`, `status.md`, `regenerate.md`
- [ ] Each skill file references the correct project name (`{{.ProjectName}}`)
- [ ] Each skill file references the correct owner name and email (`{{.OwnerName}}`, `{{.OwnerEmail}}`)
- [ ] Each skill file references the correct model ID (`{{.ModelID}}`)
- [ ] `/spec` skill correctly references `docs/BACKLOG.md` and `docs/specs/TEMPLATE.md` paths
- [ ] `/release` skill uses the configured git author for commits and tags
- [ ] `/pipeline` skill references the correct model and agent structure from config
- [ ] `/status` skill reads from `docs/BACKLOG.md` and `docs/PROGRESS.md`
- [ ] `/regenerate` skill depends on F-0040 (`-r` flag) being implemented
- [ ] No hardcoded project-specific values in any skill file — all parameterized
- [ ] Template constants follow existing naming pattern (`SkillSpecTemplate`, `SkillReleaseTemplate`, etc.)
- [ ] Generator output lists all 5 skill files with green checkmarks

## Out of Scope

- Custom user-defined skills (users adding their own beyond the 5)
- Skill argument parsing or parameterization at invocation time
- Integration testing of skills (they are instructions for Claude, not executable code)

## Dependencies

- F-0040 (settings.json) — The `/regenerate` skill depends on the `-r` flag from F-0040

## Open Questions

- None
