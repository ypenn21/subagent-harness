# appteam — Progress Log

## Session 23 — 2026-03-26

### What was accomplished
- **F-0122–F-0126: YAML frontmatter for all agent and skill templates** — Added structured YAML frontmatter blocks to all 8 agent templates and all 21 skill templates, enabling Claude Code to parse metadata (name, description, model, tool restrictions) at the platform level
- **Agent frontmatter** (F-0122) — Added YAML frontmatter to PM, TPM, SWE, SWE-Test, SWE-QA, Reviewer, Platform, and CustomAgent templates. Each includes `name`, `description`, and `model` (using `{{.ModelName}}`). Non-coding agents (PM, TPM, Reviewer) additionally get `disallowedTools` to enforce no-code-writing at the platform level
- **Skill frontmatter** (F-0123) — Added YAML frontmatter to all 21 skill templates with `name`, `description`, and `user-invocable: true`. Five workflow skills (brainstorm, cuj-list, cuj-test, roadmap, standup) additionally get `disable-model-invocation: true` to prevent automatic triggering
- **Test updates** (F-0124) — Updated tests to account for frontmatter prefix in template output
- **Version bump** (F-0125) — Version bumped to 0.18.0 in main.go
- **Test verification** (F-0126) — 97 tests pass, 0 failures across all packages
- Spec: `docs/specs/F-0122-agent-skill-frontmatter.md`

### Key decisions
- Only included essential frontmatter fields (name, description, model, disallowedTools, user-invocable, disable-model-invocation) — advanced fields like hooks, mcpServers, isolation deferred to future enhancement
- PM, TPM, Reviewer get `disallowedTools` to enforce no-code-writing at platform level, reinforcing the agent-only execution rule
- 5 workflow skills (brainstorm, cuj-list, cuj-test, roadmap, standup) get `disable-model-invocation: true` to prevent accidental automatic triggering

### Next steps
- Consider adding more frontmatter fields (tools, permissionMode, skills preloading) in future versions
- Remaining backlog: F-0027 dry-run, F-0081 retrospective, F-0082 multi-project

---

## Session 22 — 2026-03-19

### What was accomplished
- **F-0121: CLAUDE.md branching consistency** — Discovered that 7 CLAUDE.md pipeline/version-control sections still said "feature branch" even for Scion projects, contradicting the branching convention and agent template instructions
- Made all 7 references framework-conditional: Scion → "worktree branch", Claude Code → "feature branch"
- **Full E2E Scion pipeline test successful** — bookmark-cli app with lean team (PM + SWE-1 + SWE-Test):
  - PM created spec and backlog entry
  - SWE-1 implemented bookmark CRUD (main.go, storage.go, tests) — committed to `swe-1` branch, NOT feature/* branches
  - SWE-Test merged `swe-1` branch successfully, ran 99 tests with zero failures
  - Orchestrator relay pattern worked smoothly
- Version bumped to 0.17.1

### Files changed (4 files)
| File | Change |
|------|--------|
| `internal/templates/claude_md.go` | 7 "feature branch" refs made conditional on framework (Scion = worktree, CC = feature) |
| `main.go` | Version `"0.17.0"` → `"0.17.1"` |
| `docs/BACKLOG.md` | Added F-0121, marked v0.17.0 and v0.17.1 complete |
| `docs/RELEASENOTES.md` | Added v0.17.1 entry |

### Key decisions
- Every single "feature branch" reference in CLAUDE.md must be conditional — agents read CLAUDE.md as their primary instructions, and contradictory branching guidance caused agents to create feature branches in Scion (breaking the `git merge <agent>` handoff model)

### Next steps
- Push, tag v0.17.1
- Update SCION_LESSONS.md with this finding

---

## Session 21 — 2026-03-19

### What was accomplished
- **F-0117–F-0120: Scion worktree merge instructions** — Added "Worktree Isolation" section to all 8 Scion agent templates with role-specific `git merge` commands
- SWE-2 implemented merge instructions with team-size conditionals (lean vs standard/full)
- 2 new tests for conditional merge logic (Reviewer, Platform)
- Version bumped to 0.17.0
- Spec: `docs/specs/F-0117-scion-merge-instructions.md`

### Files changed (4 files)
| File | Change |
|------|--------|
| `internal/templates/scion_agents_md.go` | Added "Worktree Isolation" section to all 8 templates |
| `internal/templates/templates_test.go` | Updated existing tests + 2 new conditional merge tests |
| `main.go` | Version `"0.16.2"` → `"0.17.0"` |
| `main_test.go` | Version assertion updated |

### Key decisions
- Merge instructions placed before "Signaling Protocol" for logical flow (merge first, then work, then signal)
- SWEs always merge `pm` in lean teams (PM assigns directly) vs `tpm` in standard/full

### Next steps
- E2E Scion test to verify merge instructions work in practice
- Push, tag v0.17.0

---

## Session 20 — 2026-03-19

### What was accomplished
- **F-0114–F-0116: Scion grove-id generation fix** — E2E Scion testing revealed appteam creates `.scion/templates/` but not `.scion/grove-id`, causing scion to fail to recognize the grove
- **Generator fix** (F-0114) — SWE-2 added `generateUUID()` using `crypto/rand` (stdlib, zero deps). Generator now creates `.scion/grove-id` with UUID v4 and `.scion/agents/` directory during Scion project generation
- **Version bump + tests** (F-0115) — SWE-2 bumped version to 0.16.2, added 4 new generator tests
- **Test verification** (F-0116) — SWE-Test confirmed all tests pass with clean static analysis
- Product spec: `docs/specs/F-0114-scion-grove-id.md`

### Files changed (4 files)
| File | Change |
|------|--------|
| `internal/generator/generator.go` | Added `generateUUID()`, grove-id + agents dir creation in `generateScionFiles()` |
| `internal/generator/generator_test.go` | 4 new tests: grove-id exists, valid UUID, agents dir, no grove-id for claude-code, uniqueness |
| `main.go` | Version bump `"0.16.1"` → `"0.16.2"` |
| `main_test.go` | Version assertion updated |

### Key decisions
- UUID v4 generated via `crypto/rand` — no external dependencies needed
- `.scion/agents/` directory created alongside grove-id (matches what `scion init` creates)

### Next steps
- Push commits, tag v0.16.2
- Redeploy docs site

---

## Session 19 — 2026-03-19

### What was accomplished
- **F-0109–F-0113: Scion inter-agent messaging fix — orchestrator relay model** — E2E testing of v0.16.0 revealed that `scion message` fails inside Docker containers (Docker CLI unavailable). Fixed by switching all agent templates from direct `scion message` commands to `sciontool status task_completed` signaling with orchestrator relay from the host
- **Agent template signaling protocol** (F-0109) — SWE-2 replaced all "Messaging Protocol" sections in `internal/templates/scion_agents_md.go` with "Signaling Protocol" sections using `sciontool status task_completed`. Added Docker unavailability warning to each agent template
- **Pipeline + CLAUDE.md orchestrator relay** (F-0110) — SWE-2 updated `/pipeline` skill in `internal/templates/skill_pipeline.go` and all three Scion orchestration sections (lean/standard/full) in `internal/templates/claude_md.go` to describe the orchestrator relay model: agents signal completion, orchestrator relays information between agents from the host
- **Gitignore fix** (F-0111) — SWE-2 fixed `.gitignore` entry in `internal/generator/generator.go` from `.scion/` to `.scion/agents/` so that `.scion/templates/` remains tracked in git
- **Version bump + tests** (F-0112) — SWE-2 bumped version to `0.16.1` in `main.go`. Updated all Scion-related test assertions in `internal/templates/templates_test.go` to verify signaling protocol content
- **Tests verified** (F-0113) — SWE-Test confirmed all 180 tests pass across all 5 packages with clean static analysis. No regressions in non-Scion tests
- Product spec: `docs/specs/F-0109-scion-signaling-fix.md`
- Backlog items F-0109–F-0113 — All 5 items completed, tested, and verified

### Files changed (6 files)
| File | Change |
|------|--------|
| `internal/templates/scion_agents_md.go` | Replaced "Messaging Protocol" with "Signaling Protocol" using `sciontool status task_completed`; added Docker unavailability warning |
| `internal/templates/skill_pipeline.go` | Updated Scion section for orchestrator relay model |
| `internal/templates/claude_md.go` | Updated lean/standard/full Scion orchestration sections for orchestrator relay |
| `internal/generator/generator.go` | Fixed `.gitignore` entry: `.scion/` → `.scion/agents/` |
| `internal/templates/templates_test.go` | Updated Scion test assertions for signaling protocol content |
| `main.go` | Version bump `"0.16.0"` → `"0.16.1"` |

### Key decisions
- Orchestrator relay model chosen over Docker-in-Docker or socket mounting — simpler, more secure, no container runtime changes needed
- Agents use `sciontool status task_completed` (in-container tool, no Docker CLI required) instead of `scion message` (requires Docker CLI)
- Orchestrator on host monitors agent status and relays information between agents using `scion message` from the host side
- `.gitignore` scoped to `.scion/agents/` (runtime-generated) so `.scion/templates/` (user config) remains tracked

### Test summary
- **180 tests** across 5 packages — all pass, 0 failures
- All Scion signaling protocol tests updated and passing
- No regressions in non-Scion tests
- Clean static analysis (go vet)

### Next steps
- Push, tag v0.16.1, redeploy docs
- Remaining backlog: F-0027 dry-run, F-0081 retrospective, F-0082 multi-project

---

## Session 18 — 2026-03-18

### What was accomplished
- **F-0103–F-0108: Scion parallel agent launch with inter-agent messaging** — Replaced sequential orchestrator-managed agent launch with simultaneous parallel launch and peer-to-peer messaging via `scion message`
- **Parallel `/pipeline` skill** (F-0103) — SWE-1 rewrote the Scion section of `SkillPipelineTemplate` in `internal/templates/skill_pipeline.go` to launch ALL agents simultaneously in a single block of `scion start` commands. Orchestrator role changed to monitor-only (via `scion list` and `scion attach`), stopping all agents with `scion stop` when PM reports milestone complete
- **Inter-agent messaging protocol** (F-0104) — SWE-1 added "Messaging Protocol" sections to all 8 agent templates in `internal/templates/scion_agents_md.go`. Each role specifies who it waits for, who it messages, and the exact `scion message` commands to use for handoffs
- **CLAUDE.md parallel orchestration** (F-0105) — SWE-1 updated Scion orchestration sections in `internal/templates/claude_md.go` for all three team sizes (lean/standard/full). Orchestrator instructions now describe parallel launch with self-coordinating agents
- **Team size messaging variations** (F-0106) — Lean: PM→SWE-1 direct (no TPM); Standard: PM→TPM→SWEs→SWE-Test→Reviewer→PM; Full: includes SWE-QA between SWE-Test and Reviewer
- **Version bump + tests** (F-0107) — Version bumped to 0.16.0 in `main.go`. 13 new test functions with 20+ subtests in `templates_test.go` verify messaging protocol content in rendered templates
- **Tests verified** (F-0108) — All tests pass, 0 failures. Messaging protocol content verified in template output for all team sizes and agent roles
- **TeamSize field plumbing** — Added `TeamSize` field to SWE and custom agent template data structs, passed from generator for conditional template rendering
- Product spec: `docs/specs/F-0103-scion-parallel-launch.md`
- Backlog items F-0103–F-0108 — All 6 items completed, tested, and verified

### Files changed (9 files, +468 / -29 lines)
| File | Change |
|------|--------|
| `internal/templates/skill_pipeline.go` | Scion section: parallel launch, monitor-only orchestrator, `scion stop` on completion |
| `internal/templates/scion_agents_md.go` | Added "Messaging Protocol" section to all 8 agent templates with `scion message` commands |
| `internal/templates/claude_md.go` | Scion orchestration sections updated for parallel launch with self-coordination |
| `internal/templates/swe.go` | Added `TeamSize` field to `SWETemplateData` for conditional rendering |
| `internal/templates/custom_agent.go` | Added `TeamSize` field to `CustomAgentTemplateData` for conditional rendering |
| `internal/generator/generator.go` | Pass `TeamSize` to SWE and custom agent template data |
| `internal/templates/templates_test.go` | 13 new test functions with 20+ subtests for messaging protocol verification |
| `main.go` | Version bump `"0.15.0"` → `"0.16.0"` |
| `main_test.go` | Updated version string check |

### Key decisions
- All agents launched simultaneously — no sequential start/wait cycles. Pipeline order enforced through messaging dependencies, not orchestrator sequencing
- Orchestrator is monitor-only after launch — only intervenes via `scion message` if needed
- Each agent's messaging protocol is role-specific with exact `scion message` commands, not generic instructions
- Lean teams skip TPM in messaging chain (PM→SWE-1 direct), matching the lean pipeline philosophy
- Claude Code (`claude-code`) framework templates are completely unaffected — changes scoped to Scion only

### Test summary
- All tests pass across all 5 packages — 0 failures
- 13 new test functions with 20+ subtests verify messaging protocol content
- All existing v0.15.0 tests continue to pass (no regressions)
- Messaging protocol verified for PM, TPM, SWE, SWE-Test, SWE-QA, Reviewer across team sizes

### Next steps
- Tag v0.16.0 after PO approval, push, install
- Remaining backlog: F-0027 dry-run, F-0081 retrospective, F-0082 multi-project

---

## Session 17 — 2026-03-18

### What was accomplished
- **F-0096–F-0102: Scion framework support** — Full implementation of Scion (Google Cloud Platform multi-agent orchestration) as an alternative team framework option alongside Claude Code Agent Teams
- **Config fields** (F-0096) — SWE-1 added `Framework string` and `DefaultHarness string` fields to `ProjectConfig` in `internal/config/config.go`. `backfillDefaults` in `internal/config/settings.go` defaults `Framework` to `"claude-code"` and `DefaultHarness` to `"claude"` for backward compatibility with all pre-v0.15.0 settings
- **Wizard framework selection** (F-0097) — SWE-1 added framework selection as new Step 2 in the wizard (after Project Basics, before Git Repository), renumbering subsequent steps from 8 to 9 total. Presents "Claude Code Agent Teams" vs "Scion" choice with descriptions. Scion selection triggers harness prompt (claude/gemini/opencode/codex)
- **Scion template constants** (F-0098) — SWE-2 created three new template files in `internal/templates/`:
  - `scion_agent_yaml.go` — `ScionAgentYAMLTemplate` for `scion-agent.yaml` config files (name, description, harness config)
  - `scion_system_prompt.go` — `ScionSystemPromptTemplate` for role persona definitions per agent
  - `scion_agents_md.go` — `ScionAgentsMDTemplate` for operational instructions per agent (PM, TPM, SWE, SWE-Test, SWE-QA, Reviewer, Platform)
- **Generator Scion output** (F-0099) — SWE-2 updated `internal/generator/generator.go` with conditional output: when `Framework == "scion"`, generates `.scion/templates/<role>/` directories (each containing `scion-agent.yaml`, `agents.md`, `system-prompt.md`) instead of `.claude/agents/`. Skills embedded into agent `agents.md` files. Team size rules (lean/standard/full) apply to both frameworks
- **CLAUDE.md Scion variant** (F-0100) — SWE-2 updated `internal/templates/claude_md.go` with conditional agent management section: Scion variant references `scion start/list/attach/message/stop` commands instead of tmux/TeamCreate/TeamDelete
- **Version bump** (F-0101) — SWE-1 bumped version to `0.15.0` in `main.go`
- **Tests verified** (F-0102) — SWE-Test confirmed 148 tests pass across all 5 packages, 0 failures. All existing tests continue to pass with backward compatibility intact
- **Product spec** — PM created `docs/specs/F-0096-scion-framework.md` with 17 requirements and 17 acceptance criteria
- **Backlog items F-0096–F-0102** — All 7 items completed, tested, and verified

### Files changed (9 files, +917 / -68 lines)
| File | Change |
|------|--------|
| `internal/config/config.go` | Added `Framework string` and `DefaultHarness string` fields to `ProjectConfig` |
| `internal/config/settings.go` | Backfill `Framework` to `"claude-code"` and `DefaultHarness` to `"claude"` |
| `internal/wizard/wizard.go` | New Step 2 (Framework selection), harness prompt, renumbered steps to 9 total |
| `internal/wizard/wizard_test.go` | Updated for framework selection and step renumbering |
| `internal/templates/scion_agent_yaml.go` | New — `ScionAgentYAMLTemplate` for scion-agent.yaml |
| `internal/templates/scion_system_prompt.go` | New — `ScionSystemPromptTemplate` for role personas |
| `internal/templates/scion_agents_md.go` | New — `ScionAgentsMDTemplate` for agent operational instructions |
| `internal/templates/claude_md.go` | Conditional agent management: tmux vs scion commands |
| `internal/generator/generator.go` | Conditional output: `.claude/agents/` vs `.scion/templates/<role>/`; skills embedding for Scion |
| `main.go` | Version bump `"0.14.0"` to `"0.15.0"` |
| `main_test.go` | Updated version string check |
| `docs/specs/F-0096-scion-framework.md` | New — product spec for Scion framework support |

### Key decisions
- Framework selection is a wizard step (Step 2), not a CLI flag — keeps the interactive wizard pattern consistent
- Scion output structure follows Scion conventions: `.scion/templates/<role>/` with `scion-agent.yaml`, `agents.md`, `system-prompt.md`
- Skills are embedded into Scion agent `agents.md` files (PM skills in PM's file, dev skills in SWE files, etc.) since Scion doesn't have a separate skills directory
- CLAUDE.md is generated for both frameworks (Scion with claude harness reads CLAUDE.md)
- Default harness is "claude" — users can change to gemini, opencode, or codex
- Backward compatible: pre-v0.15.0 settings default to "claude-code" framework, producing identical output

### Test summary
- **148 tests** across 5 packages — all pass, 0 failures
- All existing v0.14.0 tests continue to pass
- Backward compatibility verified for pre-v0.15.0 settings.json files

### Next steps
- Tag v0.15.0 after PO approval, push, install
- Remaining backlog: F-0027 dry-run, F-0081 retrospective, F-0082 multi-project

---

## Session 16 — 2026-03-18

### What was accomplished
- **F-0095: GitHub Pages site with Docsify** — Created a public-facing website for appteam using Docsify, hosted on GitHub Pages from the `docs/` directory
- **Docsify SPA** — `docs/index.html` with dark theme, custom CSS, sidebar navigation, and mobile-responsive layout
- **Presentation pages** — Four content pages covering the full appteam story:
  - `docs/home.md` — Landing page with hero section, value prop, feature highlights, and quick start
  - `docs/why.md` — Problem statement: pain points of managing AI agent teams manually
  - `docs/what.md` — Feature overview: generated files, team structure, skills system, CLI flags
  - `docs/how.md` — How it works: installation, wizard walkthrough, regeneration, team sizing
  - `docs/project-roadmap.md` — Completed milestones (v0.1.0–v0.13.1) and upcoming features
- **Reveal.js slide deck** — `docs/slides.html` presentation for live demos and team onboarding
- **Supporting files** — `docs/.nojekyll` (prevent Jekyll processing), `docs/_sidebar.md` (navigation)
- **Existing docs unaffected** — BACKLOG.md, PROGRESS.md, RELEASENOTES.md, specs/, guide/ all unchanged
- Product spec: `docs/specs/F-0095-github-pages-site.md`
- Backlog item F-0095 completed

### Key decisions
- Docsify chosen over alternatives (Jekyll, Hugo, MkDocs) for zero-build simplicity — single `index.html` + markdown files
- Dark theme with custom CSS for modern look — no external theme dependency
- Reveal.js slides embedded as standalone HTML for presentation use case
- Site served from `docs/` on main branch — no separate `gh-pages` branch needed
- No version bump — this is a documentation milestone, not a code release

### Next steps
- Enable GitHub Pages in repo settings (Settings > Pages > Source: main, /docs)
- Remaining backlog: F-0027 dry-run, F-0081 retrospective, F-0082 multi-project

---

## Session 15 — 2026-03-17

### What was accomplished
- **F-0094: CUJ cleanup phase** — Added Phase 6 (Cleanup) to `/cuj-test` skill template: browser close, artifact management, test data cleanup, state reset, cleanup report
- SWE-QA agent template updated to include cleanup in testing workflow
- Version bumped to 0.13.1
- All 148 tests pass (0 failures)

### Next steps
- Tag v0.13.1 after PO approval
- Remaining backlog: F-0027 dry-run, F-0081 retrospective, F-0082 multi-project

---

## Session 14 — 2026-03-18

### What was accomplished
- **SWE-QA CUJ testing + new skills** (F-0088–F-0093) — Full implementation of CUJ-oriented QA agent template rewrite and two new QA skills (`/cuj-list`, `/cuj-test`)
- **SWE-QA agent template rewrite** (F-0088) — SWE-1 rewrote `internal/templates/swe_qa.go` (`SWEQATemplate`) with comprehensive CUJ-oriented testing philosophy, detailed Puppeteer/Chromium instructions (navigation, selectors, form interaction, screenshots, waiting strategies, console monitoring, multi-page flows, viewport configuration), CUJ inventory management (`docs/CUJ.md` with CUJ-NNN IDs, priorities, steps, outcomes), structured test result reporting, and Lighthouse integration for performance/accessibility/SEO audits
- **`/cuj-list` skill** (F-0089) — SWE-1 created `internal/templates/skill_cuj_list.go` with `SkillCUJListTemplate` constant. Skill reads existing `docs/CUJ.md`, auto-increments CUJ-NNN IDs, collects journey name/description/priority/steps/outcomes, and writes structured CUJ inventory
- **`/cuj-test` skill** (F-0090) — SWE-1 created `internal/templates/skill_cuj_testing.go` with `SkillCUJTestTemplate` constant. Skill loads CUJ inventory, allows selection by ID or priority, executes headless Chromium tests via Puppeteer, captures screenshots, monitors console errors, generates pass/fail reports, and updates CUJ file with results
- **Skill registration** (F-0091) — SWE-1 registered both skills in wizard Step 7 under "Development Skills — Tier 2" (default N), added to `skillFiles` in generator, added to `backfillSkillDefaults()` with default `false`, added to `allSkills` in `skillsSummary()`
- **Version bump** (F-0092) — SWE-1 updated version constant to `"0.13.0"` in main.go
- **Tests verified** (F-0093) — SWE-Test confirmed 148 tests pass (80 top-level + subtests) across all 5 packages, 0 failures. New tests: SWE-QA template rendering, `/cuj-list` template rendering, `/cuj-test` template rendering, backfill defaults for cuj-list and cuj-test keys
- **Product spec** — PM created `docs/specs/F-0088-swe-qa-cuj-testing.md` with detailed requirements and acceptance criteria
- **Backlog items F-0088–F-0093** — All 6 items completed, tested, and verified
- **Agent team pipeline** — Used TeamCreate with PM → TPM → SWE-1 → SWE-Test for implementation

### Files changed (11 files, +714 / -41 lines)
| File | Change |
|------|--------|
| `internal/templates/swe_qa.go` | Rewritten — CUJ-oriented testing, Puppeteer/Chromium instructions, CUJ inventory, test reporting, Lighthouse |
| `internal/templates/skill_cuj_list.go` | New — `SkillCUJListTemplate` for `/cuj-list` skill |
| `internal/templates/skill_cuj_testing.go` | New — `SkillCUJTestTemplate` for `/cuj-test` skill |
| `internal/config/settings.go` | Added `cuj-list` and `cuj-test` to `backfillSkillDefaults()` with default `false` |
| `internal/wizard/wizard.go` | Added both skills to Tier 2 display and `allSkills` in `skillsSummary()` |
| `internal/generator/generator.go` | Added both skills to `skillFiles` map |
| `main.go` | Version bump `"0.12.0"` → `"0.13.0"` |
| `main_test.go` | Updated version string check |
| `internal/templates/templates_test.go` | New SWE-QA, `/cuj-list`, `/cuj-test` template rendering tests |
| `internal/config/settings_test.go` | Updated backfill test for cuj-list and cuj-test keys |
| `internal/wizard/wizard_test.go` | Updated for cuj-list and cuj-test skills in wizard |

### Key decisions
- Both CUJ skills placed in Tier 2 (default N) — not all projects have frontends requiring browser testing
- SWE-QA template follows CUJ-oriented philosophy — tests organized as end-to-end user journeys, not isolated page checks
- CUJ inventory uses `docs/CUJ.md` with structured format: CUJ-NNN IDs, P0/P1/P2 priority, numbered steps, expected outcomes, last tested date/result
- `/cuj-test` skill generates structured pass/fail reports with screenshot evidence and console error capture
- Lighthouse integration included for automated performance, accessibility, SEO, and best practices auditing
- Backward compatible: `backfillSkillDefaults()` adds both skills with `false` for pre-v0.13.0 settings.json files

### Test summary
- **148 tests** (80 top-level + subtests) across 5 packages — all pass, 0 failures
- New tests cover SWE-QA template rendering, both CUJ skill template rendering, and backfill defaults
- All existing v0.12.0 tests continue to pass

### Next steps
- Tag v0.13.0 after PO approval, push, install
- Explore remaining backlog items (F-0027 dry-run, F-0081 retrospective, F-0082 multi-project)

---

## Session 13 — 2026-03-17

### What was accomplished
- **Team sizing — lean / standard / full** (F-0079, F-0080, F-0083–F-0087) — Full implementation of configurable team sizes that adapt the generated CLAUDE.md pipeline, roles, and agent files to the project's scope
- **TeamSize config field** (F-0079) — SWE-1 added `TeamSize string` field to `ProjectConfig` in `internal/config/config.go`. `backfillDefaults` in `internal/config/settings.go` defaults empty/missing `TeamSize` to `"standard"` for backward compatibility with pre-v0.12.0 settings.json files
- **Wizard team size selection** (F-0083) — SWE-1 added team size prompt as the first question in wizard Step 5 (Agent Team), with descriptions for each option. Lean skips SWE count and optional agent prompts (hardcodes 1 SWE + SWE-Test). Full defaults to 5 SWEs with all optional agents auto-enabled. Standard remains unchanged
- **CLAUDE.md conditional pipeline** (F-0084) — SWE-2 updated `internal/templates/claude_md.go` with conditional pipeline rules per team size. Lean pipeline: PM handles backlog/progress/release notes directly (no TPM), simplified 10-step flow, roles limited to PO/PM/SWE-1/SWE-Test. Standard: current behavior unchanged. Full: all agents always documented in roles section
- **Backlog enforcement** (F-0080) — SWE-2 added "Backlog Tracking (Non-Negotiable)" section to CLAUDE.md template, present in all three team sizes — every piece of work must have a backlog entry regardless of team configuration
- **Generator conditional generation** (F-0085) — SWE-2 updated `internal/generator/generator.go` for conditional agent file generation. Lean: generates only PM, SWE-1, SWE-Test (no TPM, no SWE-2+, no optional agents). Standard: uses config flags. Full: always generates all agent files. Custom agents generated in all team sizes
- **Version bump** (F-0086) — SWE-1 updated version constant to `"0.12.0"` in main.go
- **Tests verified** (F-0087) — SWE-Test confirmed 143 tests pass (72 top-level + subtests) across all 5 packages, 0 failures. New tests cover config backfill, wizard flows for all three team sizes, template rendering per team size, generator output per team size, backlog enforcement presence
- **Product spec** — PM created `docs/specs/F-0079-team-sizing.md` with detailed requirements and acceptance criteria
- **Backlog items F-0079, F-0080, F-0083–F-0087** — All 7 items completed, tested, and verified
- **Agent team pipeline** — Used TeamCreate with PM → TPM → SWE-1/SWE-2 → SWE-Test for parallel implementation

### Files changed (13 files, +893 / -28 lines)
| File | Change |
|------|--------|
| `internal/config/config.go` | Added `TeamSize string` field to `ProjectConfig` |
| `internal/config/settings.go` | Backfill `TeamSize` to `"standard"` for backward compat |
| `internal/config/settings_test.go` | Tests for TeamSize backfill default and round-trip |
| `internal/wizard/wizard.go` | Team size selection in Step 5; conditional agent prompts per size |
| `internal/wizard/wizard_test.go` | Tests for lean/standard/full wizard flows |
| `internal/templates/claude_md.go` | Conditional pipeline rules per team size; backlog enforcement section |
| `internal/templates/templates_test.go` | Tests for CLAUDE.md rendering per team size + backlog enforcement |
| `internal/generator/generator.go` | Conditional agent file generation based on TeamSize |
| `internal/generator/generator_test.go` | Tests for lean/standard/full agent file generation |
| `main.go` | Version bump `"0.11.0"` → `"0.12.0"` |
| `main_test.go` | Updated version string check |
| `docs/specs/F-0079-team-sizing.md` | New — product spec for team sizing feature |
| `docs/BACKLOG.md` | v0.12.0 milestone items (F-0079, F-0080, F-0083–F-0087) |

### Key decisions
- Three team sizes: lean (PM + 1 SWE + SWE-Test), standard (current default), full (all agents auto-enabled, 5 SWEs)
- Lean mode removes TPM — PM handles backlog, progress, and release notes directly to reduce coordination overhead
- Backlog enforcement is mandatory in all team sizes — lean teams don't get to skip backlog tracking
- Generator overrides config values that conflict with team size (e.g., lean forces SWECount=1)
- Custom agents still allowed in all team sizes — lean teams can add specialized agents as needed
- Backward compatible: existing settings.json files default to "standard" team size via backfill

### Test summary
- **143 tests** (72 top-level + subtests) across 5 packages — all pass, 0 failures
- New tests cover all three team sizes across config, wizard, templates, and generator
- All existing v0.11.0 tests continue to pass

### Next steps
- Tag v0.12.0 after PO approval, push, install
- Explore remaining backlog items (F-0027 dry-run, F-0081 retrospective, F-0082 multi-project)

---

## Session 12 — 2026-03-17

### What was accomplished
- **`/brainstorm` PM skill** (F-0074–F-0078) — Full implementation of a new `/brainstorm` slash command for structured product ideation sessions between the PM agent and PO
- **SkillBrainstormTemplate** (F-0074) — SWE-1 created `internal/templates/skill_brainstorm.go` with `SkillBrainstormTemplate` constant. Template implements a 5-phase brainstorm flow: Product Vision Check-in, Feature Ideation, Competitive Analysis, Prioritization Discussion, and Action Items. Includes `{{.ProjectName}}`, `{{.OwnerName}}`, and `{{.OwnerEmail}}` template variables. Instructs PM to offer capturing outcomes via `/roadmap` and `/spec` at session end
- **Wizard registration** (F-0075) — SWE-1 added `/brainstorm` to wizard Step 7 under "Project Management Skills" with default Y, and added "brainstorm" to `allSkills` list in `skillsSummary()`
- **Generator + backfill** (F-0076) — SWE-1 added `brainstorm` entry to `skillFiles` slice in generator and added `"brainstorm": true` to `backfillSkillDefaults()` for backward compatibility with pre-v0.11.0 settings.json files
- **Version bump** (F-0077) — SWE-1 updated version constant to `"0.11.0"` in main.go, added new template rendering tests
- **Tests verified** (F-0078) — SWE-Test confirmed 125 tests pass (61 top-level + 64 subtests) across all 5 packages, 0 failures. New tests: brainstorm template rendering with variable substitution, backward-compat backfill for brainstorm key
- **Product spec** — PM created `docs/specs/F-0074-brainstorm-skill.md` with detailed requirements and acceptance criteria
- **Backlog items F-0074–F-0078** — All 5 items completed, tested, and verified
- **Agent team pipeline** — Used TeamCreate with PM → TPM → SWE-1 → SWE-Test for implementation

### Files changed (9 files, +120 / -12 lines)
| File | Change |
|------|--------|
| `internal/templates/skill_brainstorm.go` | New — `SkillBrainstormTemplate` constant with 5-phase brainstorm flow |
| `internal/wizard/wizard.go` | Added `/brainstorm` to PM skills group in Step 7; added to `allSkills` in `skillsSummary()` |
| `internal/generator/generator.go` | Added `brainstorm` entry to `skillFiles` slice |
| `internal/config/settings.go` | Added `"brainstorm": true` to `backfillSkillDefaults()` |
| `main.go` | Version bump `"0.10.0"` → `"0.11.0"` |
| `internal/templates/templates_test.go` | New brainstorm template rendering tests (variable substitution, structural elements) |
| `internal/config/settings_test.go` | Updated backfill test for brainstorm key |
| `internal/wizard/wizard_test.go` | Updated for brainstorm skill in wizard |
| `main_test.go` | Updated version string check |

### Key decisions
- Brainstorm skill is a PM-tier skill (default Y) — consistent with all other PM skills
- Template uses 5 structured phases to guide productive ideation without being overly rigid
- Skill suggests using `/roadmap` and `/spec` for follow-up but does not create backlog items directly — keeps concerns separated
- No session persistence — each `/brainstorm` invocation is a fresh session

### Test summary
- **125 tests** (61 top-level + 64 subtests) across 5 packages — all pass, 0 failures
- New tests added for v0.11.0 brainstorm template rendering and backward-compat backfill
- All existing v0.10.0 tests continue to pass

### Next steps
- Tag v0.11.0 after PO approval, push, install
- Explore remaining backlog items (F-0027 dry-run)

---

## Session 11 — 2026-03-17

### What was accomplished
- **`-d`/`--dir` target directory flag** (F-0070–F-0073) — Full implementation of CLI flag to specify target directory, with auto-creation, wizard integration, and `-r` support
- **Flag parsing overhaul** (F-0070) — SWE-1 refactored `main.go` flag parsing from simple `os.Args[1]` switch to loop-based multi-flag parser supporting `-d <folder>` / `--dir <folder>`. Flag validates non-empty argument; errors if missing or starts with `-`
- **Directory auto-creation** (F-0070) — `os.MkdirAll` creates target directory (including parents) before wizard or regenerate logic runs
- **`--help` updated** (F-0070) — Added `-d, --dir <folder>` line to help output
- **Wizard skip** (F-0071) — `wizard.Run()` now accepts `targetDir string` parameter. When non-empty, skips Step 1 target directory prompt and displays `"(from -d flag)"` indicator
- **`-r -d` combo** (F-0071) — `regenerate()` accepts `targetDir` and `dirFlagSet` params. When `-d` is set, uses specified directory instead of CWD for loading settings.json and regenerating
- **Existing settings detection** (F-0071) — Settings.json lookup uses `-d` directory when provided, falling back to CWD
- **Version bump** (F-0072) — SWE-1 updated version constant to `"0.10.0"` in main.go
- **Tests verified** (F-0073) — SWE-Test confirmed 64 top-level tests pass across all 5 packages (up from 117 subtests in v0.9.0; subtests reorganized), 6 new `-d` flag tests
- **Product spec** — PM created `docs/specs/F-0070-target-dir-flag.md` with detailed requirements and acceptance criteria
- **Backlog items F-0070–F-0073** — All 4 items completed, tested, and verified
- **Agent team pipeline** — Used TeamCreate with PM → TPM → SWE-1 → SWE-Test for implementation

### Files changed (4 files, +235 / -40 lines)
| File | Change |
|------|--------|
| `main.go` | Refactored flag parsing to loop-based multi-flag parser; added `-d`/`--dir` flag; `os.MkdirAll` for dir creation; updated `--help`; `regenerate()` accepts targetDir; version bump to 0.10.0 |
| `main_test.go` | New tests for `-d` flag parsing, `--dir` long form, `-d` with missing arg, `-d` with empty string, `-r -d` combo, `--help` includes `-d`, version string update |
| `internal/wizard/wizard.go` | `Run()` accepts `targetDir string` param; skips target dir prompt when non-empty; displays `"(from -d flag)"` |
| `internal/wizard/wizard_test.go` | Updated `Run()` calls with new `targetDir` parameter |

### Key decisions
- `-d` flag creates directory before any other logic (wizard or regenerate)
- Missing value and flag-as-value (`-d -r`) detected and rejected
- Wizard `Run()` accepts `targetDir string` parameter — empty means prompt as usual
- Flag parsing refactored from `os.Args[1]` switch to loop-based parser — supports combining multiple flags in any order
- No `flag` package migration — extended existing manual parsing pattern per spec guidance

### Test summary
- **64 top-level tests** across 5 packages — all PASS
- 6 new tests added for v0.10.0 `-d` flag functionality
- All existing v0.9.0 tests continue to pass

### Next steps
- Tag v0.10.0 after PO approval, push, install
- Explore remaining backlog items (F-0027 dry-run)

---

## Session 10 — 2026-03-17

### What was accomplished
- **Custom agent definitions** (F-0064–F-0069) — Full implementation of user-defined agent roles during wizard, generated as `.claude/agents/<name>.md` files
- **CustomAgentConfig struct** (F-0064) — SWE-1 added `CustomAgentConfig{Name, Title, Description, Instructions}` struct and `CustomAgents []CustomAgentConfig` field to `ProjectConfig` in `internal/config/config.go`
- **Wizard custom agent input** (F-0065) — SWE-1 added custom agent section to wizard Step 5 (Agent Team), after optional agents and before model selection. Loop accepts kebab-case name, title, description, and multi-line instructions for each custom agent. Blank name or declining "Add another?" exits the loop
- **CustomAgentTemplate** (F-0066) — SWE-2 created `internal/templates/custom_agent.go` with `CustomAgentTemplate` constant and `CustomAgentTemplateData` struct. Template renders role title, description, project context, user-provided instructions, and git conventions
- **Generator custom agents** (F-0067) — SWE-2 updated `internal/generator/generator.go` to loop through `cfg.CustomAgents` and generate `.claude/agents/<name>.md` for each, using the `fileSpec` pattern consistent with built-in agents
- **Version bump** (F-0068) — SWE-1 updated version constant to `"0.9.0"` in main.go, updated main_test.go version check
- **Tests verified** (F-0069) — SWE-Test confirmed 117 tests pass across all 5 packages (up from 55 in v0.8.0). New tests: custom agent config round-trip, custom agent template rendering, wizard custom agent input, generator custom agent file creation
- **Product spec** — PM created `docs/specs/F-0064-custom-agents.md` with detailed requirements and acceptance criteria
- **Backlog items F-0064–F-0069** — All 6 items completed, tested, and verified
- **Agent team pipeline** — Used TeamCreate with PM → TPM → SWE-1/SWE-2 → SWE-Test for parallel implementation

### Files changed (10 files, +304 / -3 lines)
| File | Change |
|------|--------|
| `internal/config/config.go` | Added `CustomAgentConfig` struct + `CustomAgents []CustomAgentConfig` field |
| `internal/wizard/wizard.go` | Custom agent input loop in Step 5 (name/title/description/instructions) |
| `internal/templates/custom_agent.go` | New — `CustomAgentTemplate` constant + `CustomAgentTemplateData` struct |
| `internal/generator/generator.go` | Custom agent file generation loop using `fileSpec` pattern |
| `main.go` | Version bump `"0.8.0"` → `"0.9.0"` |
| `internal/config/settings_test.go` | Custom agent config round-trip test |
| `internal/templates/templates_test.go` | Custom agent template rendering tests |
| `internal/generator/generator_test.go` | Custom agent generation + empty custom agents tests |
| `internal/wizard/wizard_test.go` | Custom agent wizard input test |
| `main_test.go` | Updated version string check |

### Key decisions
- Custom agents follow the same `fileSpec` generation pattern as built-in agents — consistent architecture
- Wizard places custom agent input after optional agents (Platform, Reviewer, SWE-Test, SWE-QA) and before model selection — logical flow within Step 5
- Multi-line instructions use blank-line termination — consistent with custom skills input pattern from v0.8.0
- Custom agents persist in `settings.json` via existing `SaveSettings`/`LoadSettings` — backward compatible (nil slice is zero value)
- Template includes project context (name, owner, email, model) plus user-provided instructions — follows SWE template structure

### Test summary
- **117 tests** across 5 packages — all pass
- 62 new tests added for v0.9.0 functionality (including subtests)
- All existing v0.8.0 tests continue to pass

### Next steps
- Tag v0.9.0 after PO approval
- Explore remaining backlog items (F-0027 dry-run)

---

## Session 9 — 2026-03-17

### What was accomplished
- **12 new dev skill templates** (F-0059, F-0060) — SWE-2 created template files in `internal/templates/` for all 12 development skills:
  - Tier 1 (default Y): `/debug`, `/test`, `/review`, `/docs`, `/refactor`, `/hotfix`
  - Tier 2 (default N): `/api-design`, `/schema`, `/deploy`, `/security`, `/adr`, `/standup`
- **CustomSkillConfig struct** (F-0057) — SWE-1 added `CustomSkillConfig{Name, Description}` struct and `CustomSkills []CustomSkillConfig` field to `ProjectConfig` in `internal/config/config.go`
- **Wizard grouped skills + custom skills** (F-0058) — SWE-1 updated wizard Step 7 to display skills in two groups ("Project Management Skills:" and "Development Skills:" headers). Added custom skills input loop accepting `name: description` format (blank line to finish)
- **Generator updates** (F-0061) — SWE-2 expanded `skillFiles` slice from 6 to 18 predefined skills. Added custom skill generation loop creating `.claude/skills/<name>.md` for each user-defined custom skill
- **Backward-compat skill defaults** (F-0057) — `LoadSettings` now backfills missing skill keys: PM + Tier 1 skills default true, Tier 2 skills default false. Pre-v0.8.0 settings.json files with only 6 keys work seamlessly
- **Version bump** (F-0062) — SWE-1 updated version constant to `"0.8.0"` in main.go
- **Tests verified** (F-0063) — SWE-Test confirmed 55 tests pass across all 5 packages (up from 51 in v0.7.0). 4 new tests: custom skill config round-trip, backward-compat skill backfill, custom skill generation, dev skill template rendering
- **Product spec** — PM created `docs/specs/F-0057-universal-skills.md` with detailed requirements and acceptance criteria
- **Backlog items F-0057–F-0063** — All 7 items completed, tested, and verified
- **Agent team pipeline** — Used TeamCreate with PM → TPM → SWE-1/SWE-2 → SWE-Test for parallel implementation

### Files changed (22 files, +946 / -31 lines)
| File | Change |
|------|--------|
| `internal/config/config.go` | Added `CustomSkillConfig` struct + `CustomSkills` field |
| `internal/config/settings.go` | Backward-compat backfill for missing skill keys in `LoadSettings` |
| `internal/wizard/wizard.go` | Grouped skill display (PM + Dev headers), custom skills input loop |
| `internal/templates/skill_debug.go` | New — `/debug` skill template |
| `internal/templates/skill_test_write.go` | New — `/test` skill template |
| `internal/templates/skill_review.go` | New — `/review` skill template |
| `internal/templates/skill_docs.go` | New — `/docs` skill template |
| `internal/templates/skill_refactor.go` | New — `/refactor` skill template |
| `internal/templates/skill_hotfix.go` | New — `/hotfix` skill template |
| `internal/templates/skill_api_design.go` | New — `/api-design` skill template |
| `internal/templates/skill_schema.go` | New — `/schema` skill template |
| `internal/templates/skill_deploy.go` | New — `/deploy` skill template |
| `internal/templates/skill_security.go` | New — `/security` skill template |
| `internal/templates/skill_adr.go` | New — `/adr` skill template |
| `internal/templates/skill_standup.go` | New — `/standup` skill template |
| `internal/generator/generator.go` | Expanded skillFiles (6→18), custom skill generation loop |
| `main.go` | Version bump `"0.7.0"` → `"0.8.0"` |
| `internal/config/settings_test.go` | Custom skill round-trip + backward-compat backfill tests |
| `internal/templates/templates_test.go` | Dev skill template rendering tests |
| `internal/generator/generator_test.go` | Custom skill generation test |
| `internal/wizard/wizard_test.go` | Updated for grouped skills display |
| `main_test.go` | Updated version string check |

### Key decisions
- Skills displayed in two groups in wizard: "Project Management Skills:" (6 existing) and "Development Skills:" (12 new) — clear visual separation
- Tier 1 dev skills default to Y (common workflows every project benefits from), Tier 2 default to N (specialized, opt-in)
- Custom skills use simple `name: description` input format with blank line to finish — no complex schema needed
- `LoadSettings` backfills missing skill keys using tier-based defaults — ensures seamless upgrade from v0.7.0 settings.json files
- Template file `skill_test_write.go` avoids `_test.go` suffix to prevent Go's test runner from treating it as a test file

### Test summary
- **55 tests** across 5 packages — all pass
- 4 new tests added for v0.8.0 functionality
- All existing v0.7.0 tests continue to pass

### Next steps
- Tag v0.8.0 after PO approval
- Explore remaining backlog items (F-0027 dry-run, F-0028 custom agents)

---

## Session 8 — 2026-03-17

### What was accomplished
- **/roadmap skill template** (F-0053) — SWE-2 created `internal/templates/skill_roadmap.go` with `SkillRoadmapTemplate` constant. Instructs Claude to read BACKLOG.md, find highest F-NNNN ID, prompt for feature name/priority/description, and append to Future Items table
- **Selectable skills in wizard** (F-0051, F-0052) — SWE-1 added `SelectedSkills map[string]bool` field to `ProjectConfig`, new Step 7 (Skills) with Y/n prompts for all 6 skills (spec, release, pipeline, status, regenerate, roadmap), bumped `totalSteps` to 8, updated summary display to show selected skills
- **Conditional skill generation** (F-0054) — SWE-2 updated `internal/generator/generator.go` to only generate skill files for selected skills. Added backward-compatibility: nil/empty `SelectedSkills` map defaults to all skills enabled (pre-v0.7.0 settings.json support)
- **Tests updated** (F-0055) — SWE-Test verified 51 tests pass across all 5 packages (up from 45 in v0.6.0). 6 new tests: roadmap template rendering, skills wizard step, conditional skill generation, backward-compat nil-map fallback, settings round-trip with SelectedSkills
- **Version bump** (F-0056) — SWE-1 updated version constant to `"0.7.0"` in main.go
- **Product spec** — PM created `docs/specs/F-0051-roadmap-selectable-skills.md` with detailed requirements and acceptance criteria for both features
- **Backlog items F-0051–F-0056** — All 6 items completed, tested, and verified
- **Agent team pipeline** — Used TeamCreate with PM → TPM → SWE-1/SWE-2 → SWE-Test for parallel implementation

### Files changed (12 files, +503 / -20 lines)
| File | Change |
|------|--------|
| `internal/config/config.go` | Added `SelectedSkills map[string]bool` field with `AllSkillNames()` helper |
| `internal/wizard/wizard.go` | New Step 7 (Skills selection), totalSteps 7→8, updated summary |
| `internal/templates/skill_roadmap.go` | New — `SkillRoadmapTemplate` constant |
| `internal/generator/generator.go` | Conditional skill generation with backward-compat nil-map fallback |
| `main.go` | Version bump `"0.6.0"` → `"0.7.0"` |
| `internal/config/settings_test.go` | SelectedSkills round-trip test |
| `internal/templates/templates_test.go` | Roadmap template rendering test |
| `internal/generator/generator_test.go` | Conditional skill generation + backward-compat tests |
| `internal/wizard/wizard_test.go` | Skills wizard step test |
| `main_test.go` | Updated version string check |
| `docs/specs/F-0051-roadmap-selectable-skills.md` | New — product spec |
| `.gitignore` | Added `coverage.out` |

### Key decisions
- Skills default to all selected (Y/n prompts) — backward compatible with existing workflows
- Nil/empty `SelectedSkills` map treated as "all skills enabled" in generator — ensures pre-v0.7.0 settings.json files continue to work without migration
- /roadmap skill scans all backlog sections (not just Future Items) for highest F-NNNN ID to ensure sequential numbering

### Test summary
- **51 tests** across 5 packages — all pass
- 6 new tests added for v0.7.0 functionality
- All existing v0.6.0 tests continue to pass

### Next steps
- Tag v0.7.0 after PO approval
- Explore remaining backlog items (F-0027 dry-run, F-0028 custom agents)

---

## Session 7 — 2026-03-17

### What was accomplished
- **Comprehensive test suite** — 45 tests across all 5 packages, from zero automated tests to full coverage
- **Config tests** (7 tests, 85% coverage) — `settings_test.go`: SettingsPath, SettingsExist, SaveAndLoad round-trip, directory auto-creation, LoadNotExist, InvalidJSON, JSONFormat
- **Template tests** (18 tests + 40 subtests, 100% coverage) — `templates_test.go`: all 4 funcMap helpers (add, sub, seq, linkRange), Render, ClaudeMD with/without GCP, SWE variations (1/3/5), optional agents, conventions, all 18 template constants smoke test
- **Generator tests** (5 tests, 58% coverage) — `generator_test.go`: directory creation, file completeness, minimal config, settings persistence, content correctness
- **Wizard tests** (11 tests, 82% coverage) — `style_test.go` + `wizard_test.go`: Styler color on/off, PadBold, Banner, StepHeader, Divider, ask/askBool/askInt helpers, full Run integration, cancel flow
- **Main tests** (4 tests) — `main_test.go`: --help, --version, unknown flag, -r without settings (via exec.Command)
- **SWE-Test verified** — all 45 tests pass, all acceptance criteria met, APPROVED
- **Product spec** — PM created `docs/specs/F-0045-test-coverage.md` with detailed requirements and acceptance criteria
- **Backlog items F-0045–F-0050** — All 6 items completed and verified
- **Merged** feature/test-coverage to main
- **Agent team pipeline** — Used TeamCreate with PM → TPM → SWE-1/SWE-2 → SWE-Test for parallel implementation

### Key decisions
- Tests use `t.TempDir()` for filesystem operations (automatic cleanup, no leftover artifacts)
- No external dependencies in tests (no git, no network, no gh CLI)
- Generator tests use `InitGit=false`, `CreateRepo=false` to avoid git operations
- Main package tests use compiled binary via `exec.Command` (coverage reads 0% but functionality is tested)
- `setupGitRepo` and `IsTTY` excluded from coverage targets by design

### Coverage summary
| Package | Coverage | Test File |
|---------|----------|-----------|
| `internal/config` | 85.0% | `settings_test.go` |
| `internal/templates` | 100.0% | `templates_test.go` |
| `internal/generator` | 57.7% | `generator_test.go` |
| `internal/wizard` | 81.7% | `style_test.go`, `wizard_test.go` |
| `main` | 0.0%* | `main_test.go` |

\* Main package uses `exec.Command` to test the compiled binary, so `go test -cover` doesn't capture coverage, but all CLI behaviors are tested.

### Next steps
- Tag v0.6.0 after PO approval
- Explore remaining backlog items (F-0027 dry-run, F-0028 custom agents)
- CI/CD pipeline setup (future milestone)

---

## Session 6 — 2026-03-16

### What was accomplished
- **`.appteam/settings.json`** — New `internal/config/settings.go` with `SaveSettings`, `LoadSettings`, `SettingsExist`, `SettingsPath` functions. Generator saves config as JSON after file generation; startup detects existing settings and offers "Use saved config? (Y/n)" prompt
- **`-r` / `--regenerate` flag** — Regenerate all files from saved settings without running the wizard
- **TargetDir portability fix** — Override `cfg.TargetDir` with CWD when loading from settings.json (Reviewer-flagged bug)
- **Bootstrap skills generation** — 5 Claude Code skills created in `.claude/skills/`: `/spec`, `/release`, `/pipeline`, `/status`, `/regenerate`
- **5 new template files** — `skill_spec.go`, `skill_release.go`, `skill_pipeline.go`, `skill_status.go`, `skill_regenerate.go`
- **PM specs** — Created `docs/specs/F-0040-settings-json.md` and `docs/specs/F-0041-bootstrap-skills.md`
- **Agent team pipeline** — Used TeamCreate with PM → SWE-1/SWE-2 → Reviewer for implementation
- **Generator output** — Now creates 17 files (was 12), including `.claude/skills/` directory
- **Tagged v0.5.0**

### Key decisions
- Settings stored as `.appteam/settings.json` (not `--config file.json`) — simpler UX, auto-detected on startup
- Skills always generated (not conditional) — they bootstrap new apps with useful slash commands
- F-0026 and F-0029 superseded by the settings.json approach
- TargetDir always overridden with CWD to ensure portability across machines

### Next steps
- Add user-facing documentation to the appteam app
- Test coverage (still zero test files)
- Explore remaining backlog items (F-0027 dry-run, F-0028 custom agents)

---

## Session 5 — 2026-03-16

### What was accomplished
- **Context reorganization** — Moved all tracking files (BACKLOG.md, PROGRESS.md, RELEASENOTES.md, PIPELINE.md) from project root to `docs/` directory
- **PM spec system** — Created `docs/specs/` directory with `TEMPLATE.md` for structured product requirement specs (one per feature)
- **PM agent overhaul** — Rewrote PM template with spec creation workflow, file naming conventions (`F-NNNN-slug.md`), and `docs/specs/` ownership
- **All agent templates updated** — PM, TPM, SWE, SWE-Test, SWE-QA, Reviewer, Platform all reference `docs/` paths and `docs/specs/` for context
- **CLAUDE.md template updated** — Key References now includes `docs/specs/`, `docs/RELEASENOTES.md`, `docs/PIPELINE.md`; pipeline steps reference spec creation
- **Backlog template** — Added Spec column with links to `docs/specs/F-NNNN-slug.md`
- **Generator updated** — Creates `docs/` and `docs/specs/` directories, writes tracking files to `docs/`, generates `docs/specs/TEMPLATE.md`
- **Tagged v0.4.0**

### Key decisions
- `docs/` directory for all project management context (BACKLOG, PROGRESS, RELEASENOTES, PIPELINE, specs)
- Root keeps only CLAUDE.md and README.md (Claude Code convention and GitHub convention)
- PM creates a spec file for every feature/bug/enhancement before it enters the backlog
- Spec files named `F-NNNN-short-slug.md` matching backlog IDs
- Agent-to-file ownership matrix embedded in templates (PM owns specs, TPM owns backlog/progress)

### Next steps
- Test coverage (still zero test files)
- CLI flags for non-interactive mode (F-0026)
- Explore future backlog items

---

## Session 4 — 2026-03-16

### What was accomplished
- **MermaidJS pipeline diagram** — SWE-2 created `internal/templates/pipeline.go` with dynamic flowchart (PO → PM → TPM → SWEs → Test/QA → Review → TPM → PM → PO). Renders based on config: SWE count, optional agents
- **Git repo initialization** — SWE-1 added wizard Step 2 with `.git` detection, `git init`, `gh repo create`, and manual remote URL support via `os/exec`
- **Model selection** — SWE-1 added agent model picker (Opus 4.6, Sonnet 4.6, Haiku 4.5) and parameterized `{{.ModelName}}` / `{{.ModelID}}` across all templates
- **Backlog format update** — Changed to `F-0001` sequential IDs across all milestones with Dependencies column
- **Bug fixes from review** — Target dir ordering (moved to Step 1), dead code removal, `setupGitRepo` error handling, hardcoded defaults removed
- **Reviewer approved** — All changes verified, no blocking issues
- **Tagged v0.3.0** — Committed, tagged, and pushed to origin

### Key decisions
- Target directory asked in Step 1 (Project Basics) so git detection in Step 2 uses the correct path
- `setupGitRepo` returns error — `git init` failure is fatal, `gh`/remote failures are non-fatal warnings
- GitHub org default is empty (not hardcoded) since it's asked before Product Owner step
- Model parameterization covers Co-Authored-By lines, CLAUDE.md agent rules, and TeamCreate instructions
- PIPELINE.md always generated (not conditional)

### Next steps
- Add public/private repo choice for `gh repo create` (F-0026 candidate)
- Test coverage (reviewer flagged zero test files)
- Explore future backlog items (CLI flags, dry-run, custom agents)

---

## Session 3 — 2026-03-16

### What was accomplished
- **TUI bug fixes** — Fixed banner box misalignment (`boxLine()` helper) and summary `%-14s` ANSI alignment (`PadBold()` method)
- **New templates** — Added BACKLOG.md, PROGRESS.md, and RELEASENOTES.md generation to the tool
- **README update** — Comprehensive README with global/local install, usage, generated files tables, NO_COLOR docs
- **Project tracking** — Created BACKLOG.md, PROGRESS.md, and RELEASENOTES.md for appteam itself using semver `vMAJOR.MINOR.PATCH (YYYY-MM-DD)` format
- **Code review passed** — Reviewer approved all changes, no blocking issues
- **Merged and tagged** — `feature/tui-modernization` merged to `main`, tagged `v0.2.0`
- **Regenerated** — Re-ran appteam on itself to regenerate CLAUDE.md and all agent files

### Key decisions
- Versioning uses semver with date: `v0.1.0 (2026-03-16)`
- BACKLOG.md, PROGRESS.md, and RELEASENOTES.md are always generated (not conditional)
- Release notes follow Keep a Changelog style (Added, Changed, Fixed sections)

### Next steps
- Explore future backlog items (CLI flags, dry-run, custom agents)

---

## Session 2 — 2026-03-16

### What was accomplished
- **TUI Modernization** — Created agent team (`tui-modernize`) with PM, SWE-1, and Reviewer
- **PM defined requirements** — 8 detailed requirements with acceptance criteria (REQ-1 through REQ-8)
- **SWE-1 implemented TUI** — New `internal/wizard/style.go` with ANSI color helpers, TTY detection, `NO_COLOR` support
- **Reviewer identified 2 bugs** — Banner box misalignment and summary column alignment with ANSI codes
- **Files changed:**
  - `internal/wizard/style.go` (new) — Styler struct with color methods, Banner(), StepHeader(), Divider(), IsTTY()
  - `internal/wizard/wizard.go` — Updated to accept `color bool`, styled prompts, box-drawing summary
  - `internal/generator/generator.go` — Accepts `color bool`, styled output
  - `main.go` — TTY detection, color flag passthrough
- **Cleaned up stale agent files** — Removed leftover platform.md, swe-3.md, swe-qa.md from test run

### Key decisions
- Color controlled by `bool` flag from `main.go`, not embedded in the writer
- TTY detection uses `os.Stdout.Stat()` + `ModeCharDevice` (stdlib only)
- Unicode box-drawing chars always emitted — only ANSI escape codes suppressed when not TTY
- `NO_COLOR` env var respected per https://no-color.org/

### Next steps
- Fix alignment bugs from review
- Add BACKLOG/PROGRESS/RELEASENOTES generation

---

## Session 1 — 2026-03-16

### What was accomplished
- **Project created from scratch** — git init, go mod init, remote set to github.com/ahafin/appteam
- **Full implementation of appteam CLI tool:**
  - `internal/config/config.go` — ProjectConfig, SWEConfig, GCPConfig structs
  - `internal/templates/` — 8 template files + templates.go with Render() and funcMap
  - `internal/wizard/wizard.go` — 6-step interactive wizard with ask/askBool/askInt helpers
  - `internal/generator/generator.go` — File generation orchestration
  - `main.go` — Entry point wiring wizard → generator
- **Dogfooded** — Ran appteam on itself to generate CLAUDE.md and 6 agent definition files
- **README.md** — Install, usage, and tech stack documentation
- **Initial commit and push** to github.com/ahafin/appteam

### Key decisions
- Zero external dependencies (stdlib only — bufio, text/template, os, fmt)
- Templates as Go string constants in separate .go files
- Single parametric SWE template handles SWE-1 through SWE-5
- GCP sections conditional via `{{if .GCP.Enabled}}`
- All templates match swole project reference structure and tone

### Next steps
- TUI modernization (colors, box-drawing, styled prompts)
