# F-0051: /roadmap Skill + Selectable Skills in Wizard

**Type:** Feature
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-17

## Problem

Currently, all 5 skills (/spec, /release, /pipeline, /status, /regenerate) are always generated unconditionally during project setup. Users have no way to control which skills are included — some projects may not need all of them, and there is no mechanism to add new skills without regenerating everything.

Additionally, there is no `/roadmap` skill for quickly adding future work items to the backlog. Today, adding a roadmap item requires manually editing `docs/BACKLOG.md`, finding the next available F-NNNN ID, and formatting the table row correctly. A skill would automate this.

## Feature A: /roadmap Skill

### Requirements

1. The system must generate a new skill file at `.claude/skills/roadmap.md` that implements a `/roadmap` slash command
2. When invoked, the skill must instruct Claude to prompt the user for:
   - **Feature name** — short title for the roadmap item
   - **Priority** — P0 (critical), P1 (important), or P2 (nice to have), defaulting to P2
   - **Description** — brief notes about the feature (optional)
3. The skill must instruct Claude to auto-determine the next sequential F-NNNN ID by reading `docs/BACKLOG.md` and scanning for the highest existing ID
4. The skill must instruct Claude to append a new row to the "Future Items" table in `docs/BACKLOG.md` with the determined ID, feature name, priority, status `TODO`, no dependencies, and the description as notes
5. The skill template must be a Go string constant in `internal/templates/skill_roadmap.go`, following the same pattern as existing skill templates (e.g., `skill_spec.go`)
6. The template must include `{{.ProjectName}}` and `{{.OwnerName}}` / `{{.OwnerEmail}}` template variables in the Project Context section, consistent with other skills

### Template Content

The `/roadmap` skill must instruct the agent to:
1. Read `docs/BACKLOG.md` and find the highest F-NNNN ID across all sections (not just Future Items)
2. Compute the next sequential ID (e.g., if highest is F-0051, next is F-0052)
3. Prompt the user for feature name, priority (default P2), and optional description/notes
4. Append a properly formatted table row to the "Future Items" section
5. Report back with the assigned ID and a summary of what was added

## Feature B: Selectable Skills in Wizard

### Requirements

1. The system must add a new field to `ProjectConfig` to track which skills are selected:
   ```go
   SelectedSkills map[string]bool
   ```
   Keys are skill names: `"spec"`, `"release"`, `"pipeline"`, `"status"`, `"regenerate"`, `"roadmap"`

2. The wizard must add a new sub-step within the existing Step 6 (Conventions) or as a new dedicated step presenting all 6 available skills as a checklist:
   ```
   ━━ Step N of M ━━ Skills
     Select which skills to include:
     ▸ /spec       — Create product specs          (Y/n):
     ▸ /release    — Generate release notes         (Y/n):
     ▸ /pipeline   — Show agent pipeline diagram    (Y/n):
     ▸ /status     — Milestone status summary       (Y/n):
     ▸ /regenerate — Regenerate all files           (Y/n):
     ▸ /roadmap    — Add items to backlog roadmap   (Y/n):
   ```
   Each skill defaults to "Yes" (all selected by default).

3. The wizard `totalSteps` must increase from 7 to 8 to accommodate the new Skills step. The new step should be inserted as Step 7 (before the existing Confirm step, which becomes Step 8).

4. The generator must only generate skill files for skills where `cfg.SelectedSkills[name]` is `true`. The current unconditional append block in `generator.go` (lines 84-90) must be replaced with conditional generation:
   ```go
   skillFiles := map[string]fileSpec{
       "spec":       {filepath.Join(skillsDir, "spec.md"), templates.SkillSpecTemplate, cfg},
       "release":    {filepath.Join(skillsDir, "release.md"), templates.SkillReleaseTemplate, cfg},
       "pipeline":   {filepath.Join(skillsDir, "pipeline.md"), templates.SkillPipelineTemplate, cfg},
       "status":     {filepath.Join(skillsDir, "status.md"), templates.SkillStatusTemplate, cfg},
       "regenerate": {filepath.Join(skillsDir, "regenerate.md"), templates.SkillRegenerateTemplate, cfg},
       "roadmap":    {filepath.Join(skillsDir, "roadmap.md"), templates.SkillRoadmapTemplate, cfg},
   }
   for name, sf := range skillFiles {
       if cfg.SelectedSkills[name] {
           files = append(files, sf)
       }
   }
   ```

5. The `SelectedSkills` field must be persisted in `.appteam/settings.json` so that `-r` regeneration respects the user's skill selection. Since `ProjectConfig` is serialized via `json.MarshalIndent`, the field will be included automatically with the JSON tag `SelectedSkills`.

6. When loading settings from a pre-v0.7.0 `settings.json` that lacks the `SelectedSkills` field, the system must treat a nil/empty map as "all skills selected" for backward compatibility. This logic belongs in the generator, not in LoadSettings.

7. The configuration summary (Step 8 / Confirm) must display the selected skills:
   ```
   │  Skills:        /spec, /release, /pipeline, /status, /regenerate, /roadmap
   ```
   Or if some are deselected:
   ```
   │  Skills:        /spec, /release, /status (3 of 6)
   ```

## Implementation Files

| File | Change |
|------|--------|
| `internal/config/config.go` | Add `SelectedSkills map[string]bool` field to `ProjectConfig` |
| `internal/wizard/wizard.go` | Add skills selection step (new Step 7), bump totalSteps to 8, update Confirm to Step 8 |
| `internal/wizard/wizard.go` | Add selected skills to `printSummary` |
| `internal/templates/skill_roadmap.go` | New file — `SkillRoadmapTemplate` constant |
| `internal/templates/templates.go` | No changes needed (Render is generic) |
| `internal/generator/generator.go` | Replace unconditional skill append with conditional generation based on `SelectedSkills`; add backward-compat nil-map fallback |
| `main.go` | Update version constant to `"0.7.0"` |

## Acceptance Criteria

- [ ] New file `internal/templates/skill_roadmap.go` exists with `SkillRoadmapTemplate` constant
- [ ] `/roadmap` skill template instructs Claude to read BACKLOG.md, find next F-NNNN ID, prompt for details, and append to Future Items table
- [ ] `ProjectConfig` has a `SelectedSkills map[string]bool` field
- [ ] Wizard presents a skills selection checklist with all 6 skills, each defaulting to Yes
- [ ] Wizard `totalSteps` is 8; skills step is Step 7, confirm is Step 8
- [ ] Generator only generates skill files for selected skills
- [ ] When `SelectedSkills` is nil or empty (pre-v0.7.0 settings), all skills are generated (backward compatibility)
- [ ] `.appteam/settings.json` includes `SelectedSkills` after generation
- [ ] `-r` regeneration respects saved skill selections
- [ ] Configuration summary displays selected skills
- [ ] All existing tests continue to pass
- [ ] New tests added for: skill_roadmap template rendering, skills selection wizard step, conditional skill generation in generator

## Out of Scope

- Interactive toggle UI (arrow-key checkboxes) — the simple `(Y/n)` per skill is sufficient for v0.7.0
- Skill dependency management (e.g., /release depends on /spec) — all skills are independent
- Custom user-defined skills — only the 6 built-in skills are supported
- Removing previously generated skill files when deselected during regeneration — only controls what gets generated, not cleanup

## Dependencies

- None — this builds on the existing skill infrastructure from F-0041

## Open Questions

- None — all design decisions resolved in this spec
