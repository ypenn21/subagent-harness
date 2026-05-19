# F-0057: Universal Development Skills + Custom Skills

**Type:** Feature
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-17

## Problem

The current skills system (v0.7.0) only includes 6 project management skills (`/spec`, `/release`, `/pipeline`, `/status`, `/regenerate`, `/roadmap`). These are useful for appteam-style agent workflows, but any generated project also needs universal development skills that apply regardless of domain — debugging, testing, code review, documentation, refactoring, and more. Additionally, users need the ability to define custom skills specific to their project, without modifying the tool itself.

## Requirements

### Two Skill Categories

1. **Project Management Skills** (existing, tied to appteam agent workflow):
   - `/spec` — Create product specs (default Y)
   - `/release` — Release notes, commit, tag, push (default Y)
   - `/roadmap` — Add items to backlog (default Y)
   - `/status` — Milestone status summary (default Y)
   - `/pipeline` — Spin up agent team (default Y)
   - `/regenerate` — Regenerate from settings (default Y)

2. **Development Skills** (NEW, universal for any project):

   **Tier 1** (default Y — common workflows every project benefits from):
   - `/debug` — Systematic debugging: reproduce, isolate root cause, fix, verify with test
   - `/test` — Write tests for a module/function. Analyze code, generate table-driven tests, run and verify
   - `/review` — Code review checklist: security (OWASP top 10), performance, correctness, error handling, style
   - `/docs` — Generate/update documentation for a module, API, or feature
   - `/refactor` — Safe refactoring: run tests before, make changes, run tests after, verify no regressions
   - `/hotfix` — Emergency fix: create hotfix branch, fix, test, merge, tag patch release

   **Tier 2** (default N — specialized workflows users opt into):
   - `/api-design` — Design REST/GraphQL endpoints: routes, request/response schemas, error codes, auth
   - `/schema` — Database schema design: tables, migrations, indexes, relationships
   - `/deploy` — Deployment checklist: pre-deploy checks, deploy, post-deploy verification, rollback plan
   - `/security` — Security audit: injection, XSS, auth issues, secrets, dependency vulnerabilities
   - `/adr` — Architecture Decision Record: context, decision, consequences
   - `/standup` — Standup summary: done (git log), next (backlog), blockers

### Custom Skills

3. The system must support user-defined custom skills via a new `CustomSkills` field in `ProjectConfig`
4. Each custom skill has a `Name` (string) and `Description` (string)
5. The wizard must prompt for custom skills after predefined skills, accepting `name: description` input (blank line to finish)
6. The generator must create `.claude/skills/<name>.md` for each custom skill using a basic template
7. Custom skills must be persisted in `settings.json` alongside predefined skill selections

### Config Changes

8. `SelectedSkills map[string]bool` grows from 6 keys to 18 (6 PM + 12 dev skills)
9. New `CustomSkills []CustomSkillConfig` field where `CustomSkillConfig` has `Name` and `Description` fields
10. New `CustomSkillConfig` struct added to `config/config.go`

### Wizard UX

11. Step 7 must display skills in two groups with headers:
    - "Project Management Skills:" (6 existing skills)
    - "Development Skills:" (12 new skills, Tier 1 defaults Y, Tier 2 defaults N)
12. Custom skills input appears after predefined skills with prompt "Custom skills (blank line to finish)"
13. Custom skill input format: `name: description` per line

### Implementation Files

14. **12 new template files** in `internal/templates/`:
    - `skill_debug.go` — `/debug` skill template
    - `skill_test_write.go` — `/test` skill template (named `_test_write` to avoid collision with Go's `_test.go` convention)
    - `skill_review.go` — `/review` skill template
    - `skill_docs.go` — `/docs` skill template
    - `skill_refactor.go` — `/refactor` skill template
    - `skill_hotfix.go` — `/hotfix` skill template
    - `skill_api_design.go` — `/api-design` skill template
    - `skill_schema.go` — `/schema` skill template
    - `skill_deploy.go` — `/deploy` skill template
    - `skill_security.go` — `/security` skill template
    - `skill_adr.go` — `/adr` skill template
    - `skill_standup.go` — `/standup` skill template

15. **Config changes** in `internal/config/config.go`:
    - Add `CustomSkillConfig` struct with `Name` and `Description` fields
    - Add `CustomSkills []CustomSkillConfig` to `ProjectConfig`

16. **Wizard changes** in `internal/wizard/wizard.go`:
    - Grouped skill display with "Project Management Skills:" and "Development Skills:" headers
    - Custom skill input loop after predefined skills

17. **Generator changes** in `internal/generator/generator.go`:
    - Expanded `skillFiles` slice with all 18 predefined skills
    - Custom skill generation loop creating `.claude/skills/<name>.md` for each custom skill

18. **Version bump** in `main.go`: version set to `0.8.0`

### Backward Compatibility

19. Pre-v0.8.0 `settings.json` files with only 6 keys in `SelectedSkills` must continue to work
20. Missing keys default based on tier: PM skills default to true, Tier 1 dev skills default to true, Tier 2 dev skills default to false
21. Missing `CustomSkills` field defaults to empty slice (no custom skills)

## Acceptance Criteria

- [ ] All 18 predefined skills can be individually toggled Y/N in the wizard
- [ ] Wizard displays skills in two groups with "Project Management Skills:" and "Development Skills:" headers
- [ ] Tier 1 dev skills default to Y, Tier 2 dev skills default to N
- [ ] Custom skills can be entered during wizard with `name: description` format
- [ ] Custom skills generate `.claude/skills/<name>.md` files with skill content
- [ ] `settings.json` persists both predefined skill selections and custom skill definitions
- [ ] Pre-v0.8.0 settings files with 6-key `SelectedSkills` load correctly (missing keys use tier-based defaults)
- [ ] Each of the 12 new skill templates contains actionable, well-structured instructions (matching the quality/style of existing skill templates like `/spec` and `/roadmap`)
- [ ] All existing tests pass
- [ ] New tests cover: custom skill config serialization, wizard skill grouping, generator custom skill output
- [ ] Version reads `0.8.0`

## Out of Scope

- Skill dependency management (e.g., `/deploy` requires `/test`)
- Skill template customization beyond the initial content
- Importing skills from external sources or registries
- Skill versioning or update mechanism
- Interactive skill editing after initial creation

## Dependencies

- F-0051 (v0.7.0 — Roadmap + Selectable Skills) — completed, provides the `SelectedSkills` map foundation

## Open Questions

- None — PO has approved this design
