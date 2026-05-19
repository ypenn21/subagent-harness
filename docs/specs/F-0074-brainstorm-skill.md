# F-0074: /brainstorm Skill for PM

**Type:** Feature
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-17

## Problem

The PO needs a way to have structured, conversational brainstorming sessions with the PM agent about product direction, feature ideas, roadmap items, and prioritization — without immediately committing to backlog entries or specs. Currently, brainstorming happens ad-hoc in conversation with no guided structure, making it easy to miss important considerations like competitive analysis, user impact, or technical feasibility. A dedicated `/brainstorm` skill gives the PM agent a repeatable framework for productive ideation sessions that can optionally feed into the existing `/roadmap` and `/spec` workflows.

## Requirements

### Template

1. The system must generate a new skill file at `.claude/skills/brainstorm.md` that implements a `/brainstorm` slash command
2. A new template constant `SkillBrainstormTemplate` must be defined in `internal/templates/skill_brainstorm.go`, following the same pattern as existing skill templates
3. The template must include `{{.ProjectName}}` and `{{.OwnerName}}` / `{{.OwnerEmail}}` template variables in the Project Context section, consistent with other skills

### Brainstorm Flow

4. The `/brainstorm` skill must guide the PM agent through a structured conversation with the PO covering these phases:
   - **Product Vision Check-in** — Confirm or refine the current product vision, target users, and core value proposition
   - **Feature Ideation** — Open-ended exploration of new feature ideas, improvements, and user pain points. The PM should ask probing questions, suggest adjacent possibilities, and help the PO think through ideas
   - **Competitive Analysis** — Discuss what competitors or similar tools do well, what gaps exist, and where differentiation opportunities lie
   - **Prioritization Discussion** — For the ideas surfaced, discuss impact vs. effort, dependencies, and sequencing. Help the PO think through what to build next
   - **Action Items** — Summarize the session outcomes and offer concrete next steps

5. At the end of the session, the PM agent must offer to capture outcomes by:
   - Adding promising items to `docs/BACKLOG.md` via the `/roadmap` skill
   - Creating product specs for high-priority ideas via the `/spec` skill
   - Documenting key decisions or strategic direction changes

6. The skill must instruct the PM to maintain a conversational, collaborative tone — this is a thinking session, not a requirements-gathering interview

### Registration

7. The skill must be registered in the wizard Step 7 under "Project Management Skills" with the entry:
   ```
   /brainstorm  — Brainstorm product ideas with PO    (Y/n):
   ```
   Default: Y (consistent with all PM skills defaulting to Y)

8. The skill must be added to the generator's `skillFiles` slice in `internal/generator/generator.go`:
   ```go
   {name: "brainstorm", path: filepath.Join(skillsDir, "brainstorm.md"), tmpl: templates.SkillBrainstormTemplate, data: cfg}
   ```

9. The skill must be added to the `allSkills` list in `skillsSummary()` in `internal/wizard/wizard.go` so it appears in the confirmation summary

### Backward Compatibility

10. `backfillSkillDefaults()` in `internal/config/settings.go` must be updated to include `"brainstorm"` with a default of `true` (PM skill tier), ensuring pre-v0.11.0 `settings.json` files that lack this key will have the skill enabled by default

### Version

11. Version must be bumped to `0.11.0` in `main.go`

### Tests

12. New tests must be added for the `SkillBrainstormTemplate` rendering, verifying:
    - Template renders without error
    - Output contains key structural elements (phase names, project context variables)
    - Template variables are correctly substituted

## Implementation Files

| File | Change |
|------|--------|
| `internal/templates/skill_brainstorm.go` | New file — `SkillBrainstormTemplate` constant |
| `internal/wizard/wizard.go` | Add `/brainstorm` to PM skills group in Step 7; add to `allSkills` in `skillsSummary()` |
| `internal/generator/generator.go` | Add `brainstorm` entry to `skillFiles` slice |
| `internal/config/settings.go` | Add `"brainstorm": true` to `backfillSkillDefaults()` |
| `main.go` | Update version to `"0.11.0"` |
| `internal/templates/skill_brainstorm_test.go` | New file — tests for template rendering |

## Acceptance Criteria

- [ ] New file `internal/templates/skill_brainstorm.go` exists with `SkillBrainstormTemplate` constant
- [ ] Template contains structured brainstorming flow: vision check-in, feature ideation, competitive analysis, prioritization, action items
- [ ] Template includes `{{.ProjectName}}`, `{{.OwnerName}}`, and `{{.OwnerEmail}}` variables
- [ ] Template instructs PM to offer capturing outcomes via `/roadmap` and `/spec` at session end
- [ ] Wizard Step 7 lists `/brainstorm` under "Project Management Skills" with default Y
- [ ] `skillsSummary()` includes "brainstorm" in its `allSkills` list
- [ ] Generator `skillFiles` includes brainstorm entry with correct path and template
- [ ] `backfillSkillDefaults()` includes `"brainstorm"` defaulting to `true`
- [ ] Pre-v0.11.0 `settings.json` files without `brainstorm` key load correctly with skill enabled
- [ ] Version reads `0.11.0`
- [ ] New tests cover template rendering and variable substitution
- [ ] All existing tests continue to pass

## Out of Scope

- The `/brainstorm` skill does not create backlog items or specs directly — it suggests using `/roadmap` and `/spec` for that
- No session persistence or brainstorm history — each invocation is a fresh session
- No integration with external brainstorming tools or whiteboards
- No AI-generated feature suggestions — the PM facilitates, the PO ideates

## Dependencies

- F-0057 (v0.8.0 — Universal Skills + Custom Skills) — provides the grouped skill wizard UX and `backfillSkillDefaults()` mechanism

## Open Questions

- None — PO has approved this design
