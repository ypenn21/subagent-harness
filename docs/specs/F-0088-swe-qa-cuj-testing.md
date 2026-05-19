# F-0088: SWE-QA CUJ Testing — Critical User Journey Testing with Headless Chromium

**Type:** Feature
**Priority:** P1 (important)
**Status:** Draft
**Requested by:** PO
**Date:** 2026-03-17

## Problem

The SWE-QA agent currently has a minimal template that mentions Puppeteer/Chromium but lacks actionable structure for real QA work. It doesn't define *what* to test, *how* to organize test journeys, or *how* to report results. Teams using the "full" team size get a SWE-QA agent that needs significant manual prompting to be useful.

Two key gaps exist:

1. **No CUJ inventory** — There is no standard way to define, maintain, and prioritize the Critical User Journeys (CUJs) that the QA agent should test. Without a structured CUJ list, testing is ad hoc and coverage is unpredictable.

2. **No dedicated skills** — The QA agent has no slash commands for its core workflows. Other agents have dedicated skills (e.g., `/spec`, `/debug`, `/test`), but QA has none. A `/cuj-list` skill for managing the CUJ inventory and a `/cuj-test` skill for executing tests would make the QA agent self-sufficient and consistent.

## Requirements

### 1. Enhanced SWE-QA Agent Template

1. Rewrite `internal/templates/swe_qa.go` (`SWEQATemplate`) with a comprehensive agent definition covering:
   - **CUJ-oriented testing philosophy** — the agent thinks in terms of Critical User Journeys, not isolated page checks
   - **Headless Chromium / Puppeteer detailed instructions:**
     - Page navigation (`page.goto()`, waiting for network idle)
     - Element selection (`page.waitForSelector()`, `page.$()`, `page.$$()`)
     - Form interaction (typing, clicking, selecting dropdowns)
     - Screenshot capture at key checkpoints (`page.screenshot()`)
     - Waiting strategies (network idle, element visibility, custom conditions)
     - Console error monitoring (`page.on('console')`, `page.on('pageerror')`)
     - Multi-page flows (login -> dashboard -> feature -> logout)
     - Viewport configuration for responsive testing
   - **CUJ inventory management** — maintain a `docs/CUJ.md` file with:
     - CUJ ID format: `CUJ-001`, `CUJ-002`, etc.
     - Journey name and description
     - Steps (numbered list of user actions)
     - Expected outcomes per step
     - Priority: P0 (critical path), P1 (important), P2 (nice to have)
     - Last tested date and result (PASS/FAIL/SKIP)
   - **Test result reporting** — structured pass/fail reports with:
     - Per-CUJ status
     - Screenshot evidence
     - Console errors captured
     - Timing information
     - Failure details with reproduction steps
   - **Lighthouse integration** — performance, accessibility, SEO, and best practices audits
   - **Rules section** — commit conventions, handoff protocol to TPM/SWEs

2. The template must use the same template variables as the existing SWE-QA template:
   - `{{.ProjectName}}`, `{{.OwnerName}}`, `{{.OwnerEmail}}`, `{{.ModelName}}`

### 2. New Skill: `/cuj-list`

3. Create `internal/templates/skill_cuj_list.go` with a `SkillCUJListTemplate` constant
4. The skill template defines a `/cuj-list` slash command with the following flow:
   - Read existing `docs/CUJ.md` if it exists; find the highest `CUJ-NNN` ID to auto-increment
   - Ask the user about the application's key user journeys
   - For each journey: collect name, description, priority (P0/P1/P2), step-by-step user actions, expected outcomes
   - Write/update `docs/CUJ.md` with the complete inventory in a structured markdown format
   - Commit the change with standard commit conventions
5. Template variables: `{{.ProjectName}}`, `{{.OwnerName}}`, `{{.OwnerEmail}}`

### 3. New Skill: `/cuj-test`

6. Create `internal/templates/skill_cuj_test.go` with a `SkillCUJTestTemplate` constant
7. The skill template defines a `/cuj-test` slash command with the following flow:
   - Read `docs/CUJ.md` to load the CUJ inventory
   - Ask the user which CUJs to test: all, specific IDs (e.g., `CUJ-001,CUJ-003`), or by priority (e.g., `P0 only`)
   - For each selected CUJ:
     a. Launch headless Chromium via Puppeteer (`npx puppeteer` or project-installed)
     b. Execute each step in the journey (navigate, click, fill, wait)
     c. Capture screenshots at key checkpoints (saved to a `screenshots/` directory)
     d. Monitor browser console for errors and warnings
     e. Verify expected outcomes (element presence, text content, URL changes)
   - Generate a test report summarizing:
     - Overall pass/fail count
     - Per-CUJ results with details
     - Screenshots taken
     - Console errors found
     - Recommendations for failures
   - Update `docs/CUJ.md` with last tested date and result per CUJ
   - Commit the updated CUJ file
8. Template variables: `{{.ProjectName}}`, `{{.OwnerName}}`, `{{.OwnerEmail}}`

### 4. Registration

9. Add both skills to the predefined skill list in `internal/config/settings.go`:
   - `"cuj-list"` — default `false` (Tier 2, not all apps have frontends)
   - `"cuj-test"` — default `false` (Tier 2, not all apps have frontends)
10. Add both skills to the wizard Step 7 display under "Development Skills — Tier 2" group in `internal/wizard/wizard.go`
11. Add both skills to the `skillFiles` map in `internal/generator/generator.go` mapping to their template constants
12. Add both skills to the `allSkills` list in `skillsSummary()` in `internal/wizard/wizard.go`
13. Add both skills to `backfillSkillDefaults()` in `internal/config/settings.go` with default `false`

### 5. Version

14. Bump version to `0.13.0` in `main.go`

### 6. Tests

15. New tests must cover:
    - Template rendering: `SWEQATemplate` renders correctly with sample config
    - Template rendering: `SkillCUJListTemplate` renders correctly
    - Template rendering: `SkillCUJTestTemplate` renders correctly
    - Config: `backfillSkillDefaults()` adds `cuj-list` and `cuj-test` with `false` default
    - Generator: `cuj-list` and `cuj-test` appear in `skillFiles` and generate correctly when selected
    - Wizard: both skills appear in Tier 2 skills display and in `skillsSummary()`

## Implementation Files

| File | Change |
|------|--------|
| `internal/templates/swe_qa.go` | Rewrite `SWEQATemplate` with comprehensive CUJ-oriented QA agent definition |
| `internal/templates/skill_cuj_list.go` | New file: `SkillCUJListTemplate` for `/cuj-list` skill |
| `internal/templates/skill_cuj_test.go` | New file: `SkillCUJTestTemplate` for `/cuj-test` skill |
| `internal/config/settings.go` | Add `cuj-list` and `cuj-test` to `backfillSkillDefaults()` with default `false` |
| `internal/wizard/wizard.go` | Add both skills to Tier 2 display and `skillsSummary()` |
| `internal/generator/generator.go` | Add both skills to `skillFiles` map |
| `main.go` | Version bump to `0.13.0` |
| `internal/templates/swe_qa_test.go` | Tests for enhanced SWE-QA template rendering |
| `internal/templates/skill_cuj_list_test.go` | Tests for `/cuj-list` skill template rendering |
| `internal/templates/skill_cuj_test_test.go` | Tests for `/cuj-test` skill template rendering |
| `internal/config/config_test.go` | Tests for skill backfill defaults |
| `internal/generator/generator_test.go` | Tests for skill file generation |

## Acceptance Criteria

- [ ] `SWEQATemplate` rewritten with CUJ-oriented testing philosophy, Puppeteer instructions, CUJ inventory management, test reporting, and Lighthouse integration
- [ ] `SkillCUJListTemplate` created in `skill_cuj_list.go` with complete `/cuj-list` skill flow (read existing CUJs, collect new journeys, write `docs/CUJ.md`, commit)
- [ ] `SkillCUJTestTemplate` created in `skill_cuj_test.go` with complete `/cuj-test` skill flow (load CUJs, select targets, run Puppeteer tests, capture screenshots, report results, update CUJ file)
- [ ] Both skills registered in `backfillSkillDefaults()` with default `false`
- [ ] Both skills appear in wizard Tier 2 group display
- [ ] Both skills included in `skillsSummary()` output
- [ ] Both skills mapped in generator `skillFiles` to their template constants
- [ ] Version bumped to `0.13.0` in `main.go`
- [ ] All new templates render without errors using standard test config
- [ ] All existing tests continue to pass
- [ ] New tests cover template rendering, config backfill, generator mapping, and wizard display

## Out of Scope

- Actual Puppeteer/Chromium installation or dependency management (the skill instructs the agent; the agent handles tooling)
- Visual regression testing with baseline image diffing (future enhancement)
- Integration with CI/CD pipelines for automated CUJ testing
- Non-browser testing (API testing, load testing)
- CUJ test parallelization across multiple browser instances

## Dependencies

- F-0079 (v0.12.0 — Team Sizing) — SWE-QA agent is generated in "full" team size and optionally in "standard"
- F-0057 (v0.8.0 — Universal Skills) — provides skill registration patterns (`backfillSkillDefaults`, `skillFiles`, Tier 2 grouping)
- F-0074 (v0.11.0 — Brainstorm Skill) — reference for multi-phase skill template structure

## Open Questions

- None — PO has described the CUJ testing approach and skill flows. The implementation follows established patterns from existing skills and agent templates.
