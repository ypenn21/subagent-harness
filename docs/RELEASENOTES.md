# appteam — Release Notes

## v0.18.0 — Agent & Skill YAML Frontmatter (2026-03-26)

### Added
- **YAML frontmatter in all 8 agent templates** — Every generated `.claude/agents/*.md` file now begins with structured YAML frontmatter between `---` markers, enabling Claude Code to configure agents at load time:
  - `name` — lowercase-hyphen identifier matching the filename (e.g., `pm`, `tpm`, `swe-1`, `swe-test`, `reviewer`, `platform`)
  - `description` — one-sentence summary of when Claude should delegate to this agent
  - `model` — model selection using the wizard's `{{.ModelName}}` template variable
  - `disallowedTools` — PM, TPM, and Reviewer agents restrict `Bash`, `Edit`, `Write` tools, enforcing the "never write application code directly" rule at the platform level (not just in prose)
- **YAML frontmatter in all 21 skill templates** — Every generated `.claude/skills/*.md` file now begins with structured YAML frontmatter:
  - `name` — slash-command trigger (e.g., `spec`, `deploy`, `brainstorm`)
  - `description` — one-sentence summary of what the skill does
  - `user-invocable: true` — all 21 skills are registered as slash commands
  - `disable-model-invocation: true` — 5 workflow skills that require interactive user participation (brainstorm, cuj-list, cuj-test, roadmap, standup) are prevented from auto-triggering
- **Custom agent frontmatter** — Custom agents defined during the wizard also receive frontmatter with dynamic `name` and `description` from user input
- **SWE parametric frontmatter** — The parametric SWE template produces per-instance names (`swe-1`, `swe-2`, etc.) via `swe-{{.Number}}`
- Product spec: `docs/specs/F-0122-agent-skill-frontmatter.md`

### Changed
- All 29 template string constants updated (8 agents + 21 skills) with frontmatter prepended before existing markdown content
- Existing markdown body content unchanged in all templates — frontmatter is purely additive
- Test assertions updated to account for frontmatter prefix in rendered template output
- Version bumped to 0.18.0

---

## v0.17.1 — Scion CLAUDE.md Branching Consistency (2026-03-19)

### Fixed
- **Contradictory branching instructions in CLAUDE.md** — 7 pipeline and version-control sections said "implement on a feature branch" even for Scion projects, contradicting the Scion convention ("commit to your worktree branch"). All references are now framework-conditional: Scion projects say "worktree branch", Claude Code projects say "feature branch"
- E2E verified: Scion pipeline (PM → SWE-1 → SWE-Test) with lean team completed successfully. SWE-1 committed to `swe-1` branch (no feature branches), SWE-Test merged `swe-1` and ran 99 tests with zero failures

---

## v0.17.0 — Scion Worktree Merge Instructions (2026-03-19)

### Added
- **Git merge instructions in all 8 Scion agent templates** — Each agent now has a "Worktree Isolation" section explaining which branch to merge before starting work:
  - PM: no merge needed (starts pipeline)
  - TPM: `git merge pm`
  - SWE: `git merge pm` (lean) / `git merge tpm` (standard/full)
  - SWE-Test: `git merge swe-1` (+ additional SWE branches)
  - SWE-QA: `git merge swe-test`
  - Reviewer: conditional on SWE-QA presence
  - Platform: `git merge pm` (lean) / `git merge tpm` (standard/full)
  - Custom: generic merge guidance
- 2 new tests: `TestScionReviewerMergeConditional`, `TestScionPlatformMergeConditional`

---

## v0.16.2 — Scion Grove-ID Generation (2026-03-19)

### Fixed
- **Missing `.scion/grove-id`** — appteam was generating `.scion/templates/` for Scion projects but not `.scion/grove-id`, which scion needs to recognize the grove. Without it, `scion init` refuses to run and `scion start`/`scion list` can't find the grove
- Generator now creates `.scion/grove-id` with a UUID v4 (using `crypto/rand`, stdlib only) and `.scion/agents/` directory during Scion project generation

---

## v0.16.1 — Scion Signaling Protocol Fix (2026-03-19)

### Fixed
- **Scion inter-agent messaging** — `scion message` fails inside Docker containers because Docker CLI is unavailable in Scion containers (`docker: executable file not found in $PATH`). Replaced with orchestrator relay model: agents signal completion via `sciontool status task_completed`, orchestrator relays information between agents from the host
- **`.gitignore` entry** — Changed from `.scion/` to `.scion/agents/` so `.scion/templates/` stays tracked in git

### Changed
- All 8 Scion agent templates: "Messaging Protocol" sections replaced with "Signaling Protocol" using `sciontool status task_completed` instead of `scion message`. Each template includes Docker unavailability warning
- `/pipeline` skill Scion section updated to describe orchestrator relay workflow
- CLAUDE.md template Scion orchestration sections updated for all 3 team sizes (lean/standard/full) to use orchestrator relay model
- Version bumped to 0.16.1

---

## v0.16.0 — Scion Parallel Launch & Inter-Agent Messaging (2026-03-18)

### Added
- **Inter-agent messaging protocol** — All Scion agent templates now include a "Messaging Protocol" section with exact `scion message` commands for peer-to-peer handoffs. Each agent knows who it waits for and who it messages when its work is done:
  - PM: messages TPM after spec creation, messages orchestrator on milestone completion
  - TPM: waits for PM, messages SWEs with work assignments
  - SWE: waits for TPM/PM, messages SWE-Test on implementation completion
  - SWE-Test: waits for SWEs, messages Reviewer (or PM) on test pass/fail
  - SWE-QA (full teams): waits for SWE-Test, messages Reviewer on QA verification
  - Reviewer: waits for SWE-Test/QA, messages PM on approval
- 13 new test functions with 20+ subtests for messaging protocol verification
- Product spec: `docs/specs/F-0103-scion-parallel-launch.md`

### Changed
- **`/pipeline` skill (Scion section)** — All agents now launched simultaneously in a single block of `scion start` commands instead of sequential start-and-wait. Orchestrator role changed to monitor-only: observe via `scion list`, intervene via `scion message`, and `scion stop` all agents when PM reports completion
- **CLAUDE.md Scion orchestration** — All three team sizes (lean/standard/full) updated for parallel launch with self-coordinating agents
- **Team size messaging variations**:
  - Lean: PM → SWE-1 → SWE-Test → PM (no TPM in messaging chain)
  - Standard: PM → TPM → SWEs → SWE-Test → Reviewer → PM
  - Full: PM → TPM → SWEs → SWE-Test → SWE-QA → Reviewer → PM
- SWE and custom agent template data structs now include `TeamSize` field for conditional rendering
- Generator passes `TeamSize` to SWE and custom agent templates
- Version bumped to 0.16.0

---

## v0.15.0 — Scion Framework Support (2026-03-18)

### Added
- **Scion as alternative team framework** — appteam now supports generating [Scion](https://github.com/GoogleCloudPlatform/scion) multi-agent configurations in addition to Claude Code Agent Teams
  - New wizard Step 2 (Team Framework) lets users choose between Claude Code Agent Teams and Scion
  - When Scion is selected, generates `.scion/templates/<role>/` directories (PM, TPM, SWE-1..N, SWE-Test, SWE-QA, Reviewer, Platform) instead of `.claude/agents/`
  - Each Scion template directory contains `scion-agent.yaml` (config), `agents.md` (operational instructions), and `system-prompt.md` (role persona)
- **Harness selection for Scion** — Choose default harness: `claude` (default), `gemini`, `opencode`, or `codex`
- **`Framework` field** in ProjectConfig — values: `"claude-code"` (default) or `"scion"`, persisted in `.appteam/settings.json`
- **`DefaultHarness` field** in ProjectConfig — values: `"claude"`, `"gemini"`, `"opencode"`, `"codex"`, persisted in settings
- New template files:
  - `internal/templates/scion_agent_yaml.go` — `ScionAgentYAMLTemplate` for agent YAML config
  - `internal/templates/scion_system_prompt.go` — `ScionSystemPromptTemplate` for role personas
  - `internal/templates/scion_agents_md.go` — `ScionAgentsMDTemplate` for agent operational instructions
- Product spec: `docs/specs/F-0096-scion-framework.md`

### Changed
- Wizard expanded to 9 steps (new Framework step is Step 2, subsequent steps renumbered)
- CLAUDE.md template has conditional agent management section — Scion variant uses `scion start/list/attach/message/stop` commands instead of tmux/TeamCreate/TeamDelete
- Generator conditionally outputs `.claude/agents/` or `.scion/templates/<role>/` based on Framework setting
- Skills embedded into Scion agent `agents.md` files (PM skills in PM, dev skills in SWEs, testing skills in SWE-Test/SWE-QA)
- Team size rules (lean/standard/full) apply to both frameworks
- `backfillDefaults` sets `Framework` to `"claude-code"` and `DefaultHarness` to `"claude"` for backward compatibility with pre-v0.15.0 settings
- Version bumped to 0.15.0

---

## v0.14.0 — GitHub Pages Website (2026-03-18)

### Added
- **Docsify-based public website** — Static site hosted on GitHub Pages from the `docs/` directory, no build step required
  - `docs/index.html` — Docsify SPA with dark theme, custom CSS, sidebar navigation
  - `docs/home.md` — Landing page with hero section, value prop, feature highlights, quick start guide
  - `docs/why.md` — Problem statement: pain points of managing AI agent teams without automation
  - `docs/what.md` — Feature overview: generated files, team structure, skills system, CLI flags
  - `docs/how.md` — How it works: installation, wizard walkthrough, regeneration, team sizing
  - `docs/project-roadmap.md` — Completed milestones (v0.1.0 through v0.13.1) and upcoming features
  - `docs/_sidebar.md` — Navigation sidebar for all presentation pages
  - `docs/.nojekyll` — Prevents GitHub Pages Jekyll processing
- **Reveal.js slide presentation** — `docs/slides.html` standalone presentation for live demos and team onboarding
- Product spec: `docs/specs/F-0095-github-pages-site.md`

### Notes
- Documentation milestone — no Go code changes, no version bump to CLI binary
- Enable GitHub Pages in repo settings: Settings > Pages > Source: Deploy from branch, main, /docs

---

## v0.13.1 (2026-03-17)

### Changed
- `/cuj-test` skill: added Phase 6 (Cleanup) — browser close, artifact management, test data cleanup, state reset, cleanup report
- SWE-QA agent template: cleanup step added to workflow
- Version bumped to 0.13.1

---

## v0.13.0 (2026-03-18)

### Added
- **Enhanced SWE-QA agent template** — Complete rewrite of `SWEQATemplate` with CUJ-oriented testing philosophy:
  - **Critical User Journey (CUJ) inventory management** — Maintain a structured `docs/CUJ.md` file with CUJ IDs (CUJ-001, CUJ-002, etc.), journey descriptions, step-by-step user actions, expected outcomes, priorities (P0/P1/P2), and last tested date/result
  - **Headless Chromium / Puppeteer instructions** — Detailed browser automation guidance: page navigation, element selection, form interaction, screenshot capture, waiting strategies (network idle, element visibility), console error monitoring, multi-page flows, viewport configuration
  - **Structured test reporting** — Per-CUJ pass/fail status, screenshot evidence, console errors, timing, failure reproduction steps
  - **Lighthouse integration** — Performance, accessibility, SEO, and best practices audits
- **`/cuj-list` skill** — New slash command for creating and updating the CUJ inventory (`docs/CUJ.md`). Reads existing CUJs to auto-increment IDs, collects journey details interactively, writes structured inventory, commits the change
  - New `SkillCUJListTemplate` in `internal/templates/skill_cuj_list.go`
- **`/cuj-test` skill** — New slash command for running CUJ tests via headless Chromium/Puppeteer. Loads CUJ inventory, lets user select targets (all, specific IDs, by priority), executes steps with screenshots, monitors console, generates pass/fail report, updates `docs/CUJ.md` with results
  - New `SkillCUJTestTemplate` in `internal/templates/skill_cuj_testing.go`
- Product spec: `docs/specs/F-0088-swe-qa-cuj-testing.md`

### Changed
- Wizard Step 7 (Skills) now shows 21 skills — added `/cuj-list` and `/cuj-test` under "Development Skills — Tier 2" (default N, since not all apps have frontends)
- `backfillSkillDefaults()` includes `"cuj-list"` and `"cuj-test"` with default `false` for backward compatibility with pre-v0.13.0 settings
- Generator `skillFiles` expanded to 21 predefined skills
- Version bumped to 0.13.0

### Added Files
- `internal/templates/skill_cuj_list.go` — `/cuj-list` skill template
- `internal/templates/skill_cuj_testing.go` — `/cuj-test` skill template

---

## v0.12.0 (2026-03-17)

### Added
- **Team sizing: lean / standard / full team configurations**
  - **Lean** — PM + 1 SWE + SWE-Test — simplified pipeline where PM handles backlog, progress log, and release notes directly (no TPM, no Reviewer)
  - **Standard** — PM + TPM + configurable SWEs + SWE-Test + Reviewer (current default, unchanged behavior)
  - **Full** — PM + TPM + up to 5 SWEs + SWE-Test + SWE-QA + Reviewer + Platform (all optional agents enabled by default)
- **Backlog enforcement (F-0080)** — Non-negotiable backlog tracking rule baked into CLAUDE.md template for all team sizes: every piece of work gets a `docs/BACKLOG.md` entry regardless of team configuration
- Product spec: `docs/specs/F-0079-team-sizing.md`

### Changed
- Wizard Step 5 (Agent Team) now starts with team size selection (`lean`/`standard`/`full`) before SWE count or optional agent prompts
- Lean wizard skips SWE count and optional agent prompts, hardcodes 1 SWE + SWE-Test
- Full wizard defaults SWE count to 5 and auto-enables all optional agents
- CLAUDE.md template generates conditional pipeline rules per team size (lean has simplified pipeline, standard unchanged, full includes all roles)
- Generator conditionally creates agent files based on team size
- `backfillDefaults` in settings.go defaults `TeamSize` to `"standard"` for backward compatibility with pre-v0.12.0 settings
- Version bumped to 0.12.0

---

## v0.11.0 (2026-03-17)

### Added
- **`/brainstorm` skill for PM agent** — Structured 5-phase product ideation sessions with PO covering vision check-in, feature ideation, competitive analysis, prioritization discussion, and action items
- New `SkillBrainstormTemplate` in `internal/templates/skill_brainstorm.go`
- At session end, PM offers to capture outcomes via `/roadmap` (backlog items) and `/spec` (product specs)
- Product spec: `docs/specs/F-0074-brainstorm-skill.md`

### Changed
- Wizard Step 7 (Skills) now shows 7 PM skills (added `/brainstorm` with default Y)
- `skillsSummary()` updated to include brainstorm in 19-skill total
- `backfillSkillDefaults()` includes `"brainstorm"` for backward compatibility with pre-v0.11.0 settings
- Generator `skillFiles` slice expanded to 19 predefined skills
- Version bumped to 0.11.0

---

## v0.10.0 (2026-03-17)

### Added
- **`-d <folder>` / `--dir <folder>` CLI flag** — specify target directory from the command line
- Auto-creates directory if it doesn't exist (including parent directories, like `mkdir -p`)
- Works with wizard: skips "Target directory" prompt in Step 1, pre-fills with `-d` value
- Works with `-r`: `appteam -r -d ./path` regenerates into a specific directory (overrides CWD)
- Product spec: `docs/specs/F-0070-target-dir-flag.md`

### Changed
- Flag parsing rewritten to loop-based iteration supporting multi-flag combinations in any order
- `wizard.Run()` accepts optional `targetDir` parameter to pre-set target directory
- `--help` output updated with `-d, --dir <folder>` option
- Version bumped to 0.10.0

---

## v0.9.0 (2026-03-17)

### Added
- **Custom agent definitions** — Users can now define their own agent roles beyond the built-in set (PM, TPM, SWE, Reviewer, SWE-Test, SWE-QA, Platform). Examples: Frontend Engineer, UX Designer, DBA, Data Analyst, DevOps Engineer, Security Engineer, Tech Writer
- Custom agents are defined during the wizard in Step 5 (Agent Team) with:
  - Name (kebab-case, used as filename)
  - Title (display name)
  - Description (role summary)
  - Instructions (multi-line bullet points for agent behavior)
- Each custom agent generates a `.claude/agents/<name>.md` file with role, description, instructions, project context, and git conventions
- Custom agents persist in `settings.json` and regenerate correctly with `-r`
- `CustomAgentConfig` struct with `Name`, `Title`, `Description`, `Instructions` fields added to config package
- `CustomAgents []CustomAgentConfig` field added to `ProjectConfig`
- `CustomAgentTemplate` constant and `CustomAgentTemplateData` struct in `internal/templates/custom_agent.go`
- Backward-compatible: pre-v0.9.0 `settings.json` files without `CustomAgents` field load without error (nil slice is zero value)
- 62 new tests (117 total across all packages, all pass): custom agent config round-trip, template rendering, wizard input parsing, generator file creation
- Product spec: `docs/specs/F-0064-custom-agents.md`

### Changed
- Wizard Step 5 (Agent Team) now includes a custom agents section after optional agents (Platform, Reviewer, SWE-Test, SWE-QA) and before model selection
- Generator loops through `cfg.CustomAgents` after built-in agents, using the same `fileSpec` pattern for consistent rendering
- Confirmation summary (Step 8) displays custom agents with name and title
- Version bumped to 0.9.0

### Added Files
- `internal/templates/custom_agent.go` — Custom agent template and data struct

---

## v0.8.0 (2026-03-17)

### Added
- **12 universal development skills** — New skill templates for common development workflows, generated as `.claude/skills/<name>.md`:
  - **Tier 1** (default Y — included by default):
    - `/debug` — Systematic debugging: reproduce, isolate root cause, fix, verify with test
    - `/test` — Write tests for a module/function: analyze code, generate table-driven tests, run and verify
    - `/review` — Code review checklist: security (OWASP top 10), performance, correctness, error handling, style
    - `/docs` — Generate/update documentation for a module, API, or feature
    - `/refactor` — Safe refactoring: run tests before, make changes, run tests after, verify no regressions
    - `/hotfix` — Emergency fix: create hotfix branch, fix, test, merge, tag patch release
  - **Tier 2** (default N — opt-in specialized workflows):
    - `/api-design` — Design REST/GraphQL endpoints: routes, request/response schemas, error codes, auth
    - `/schema` — Database schema design: tables, migrations, indexes, relationships
    - `/deploy` — Deployment checklist: pre-deploy checks, deploy, post-deploy verification, rollback plan
    - `/security` — Security audit: injection, XSS, auth issues, secrets, dependency vulnerabilities
    - `/adr` — Architecture Decision Record: context, decision, consequences
    - `/standup` — Standup summary: done (git log), next (backlog), blockers
- **Custom skills** — Define project-specific skills during the wizard by entering `name: description` pairs. Each custom skill generates a `.claude/skills/<name>.md` file and is persisted in `settings.json`
- `CustomSkillConfig` struct with `Name` and `Description` fields added to config package
- `CustomSkills []CustomSkillConfig` field added to `ProjectConfig`
- Backward-compatible skill defaults: `LoadSettings` backfills missing skill keys — PM + Tier 1 skills default true, Tier 2 default false (pre-v0.8.0 settings.json files upgrade seamlessly)
- 4 new tests (55 total across all packages, all pass): custom skill config round-trip, backward-compat skill backfill, custom skill generation, dev skill template rendering
- Product spec: `docs/specs/F-0057-universal-skills.md`

### Changed
- Wizard Step 7 (Skills) now displays skills in two groups with headers: "Project Management Skills:" (6 existing) and "Development Skills:" (12 new)
- `SelectedSkills` map grows from 6 keys to 18 (6 PM + 6 Tier 1 + 6 Tier 2)
- Generator `skillFiles` slice expanded from 6 to 18 predefined skills, plus dynamic custom skill generation
- Generator now outputs up to 30 skill files (18 predefined + custom)
- Version bumped to 0.8.0

### Added Files
- `internal/templates/skill_debug.go` — `/debug` skill template
- `internal/templates/skill_test_write.go` — `/test` skill template
- `internal/templates/skill_review.go` — `/review` skill template
- `internal/templates/skill_docs.go` — `/docs` skill template
- `internal/templates/skill_refactor.go` — `/refactor` skill template
- `internal/templates/skill_hotfix.go` — `/hotfix` skill template
- `internal/templates/skill_api_design.go` — `/api-design` skill template
- `internal/templates/skill_schema.go` — `/schema` skill template
- `internal/templates/skill_deploy.go` — `/deploy` skill template
- `internal/templates/skill_security.go` — `/security` skill template
- `internal/templates/skill_adr.go` — `/adr` skill template
- `internal/templates/skill_standup.go` — `/standup` skill template

---

## v0.7.0 (2026-03-17)

### Added
- **`/roadmap` skill** — New Claude Code slash command (`.claude/skills/roadmap.md`) that auto-determines the next sequential F-NNNN ID from `docs/BACKLOG.md` and appends new roadmap items to the Future Items table with feature name, priority, and description
- **Selectable skills in wizard** — New Step 7 (Skills) presents all 6 available skills as Y/n prompts, each defaulting to Yes:
  - `/spec` — Create product specs
  - `/release` — Generate release notes
  - `/pipeline` — Show agent pipeline diagram
  - `/status` — Milestone status summary
  - `/regenerate` — Regenerate all files
  - `/roadmap` — Add items to backlog roadmap
- `SelectedSkills map[string]bool` field in `ProjectConfig` for tracking skill selections
- Backward compatibility: nil/empty `SelectedSkills` map (pre-v0.7.0 `settings.json`) defaults to all skills enabled
- 6 new tests (51 total across all packages, all pass)
- Product spec: `docs/specs/F-0051-roadmap-selectable-skills.md`

### Changed
- Wizard expanded to 8 steps (new Skills step is Step 7, Confirm moved to Step 8)
- Generator conditionally generates only selected skill files (was unconditional)
- Configuration summary now displays selected skills list
- Generator now outputs up to 18 files (was 17)
- Version bumped to 0.7.0

### Fixed
- Added `coverage.out` to `.gitignore`

### Added Files
- `internal/templates/skill_roadmap.go` — `/roadmap` skill template

---

## v0.6.0 (2026-03-17)

### Added
- Comprehensive test suite — 45 tests across all 5 packages (1,510 lines of test code)
- `internal/config/settings_test.go` — 7 tests (85% coverage): settings path, save/load round-trip, directory creation, error handling, JSON format validation
- `internal/templates/templates_test.go` — 18 tests with 40+ subtests (100% coverage): all template functions (add, sub, seq, linkRange), Render, ClaudeMD with/without GCP, SWE variations, optional agents, conventions, smoke test of all 18 template constants
- `internal/generator/generator_test.go` — 5 tests (58% coverage): directory creation, file completeness, minimal config, settings persistence, content correctness
- `internal/wizard/style_test.go` — 6 tests: Styler color on/off, PadBold, Banner, StepHeader, Divider
- `internal/wizard/wizard_test.go` — 5 tests: ask/askBool/askInt helpers, full Run integration, cancel flow
- `main_test.go` — 4 tests: --help, --version, unknown flag, -r without settings (via exec.Command)
- Product spec: `docs/specs/F-0045-test-coverage.md`

---

## v0.5.1 (2026-03-17)

### Added
- `docs/guide/` — 8 user-facing documentation pages:
  - `getting-started.md` — install and full wizard walkthrough
  - `cli-reference.md` — flags, environment variables, exit codes
  - `generated-files.md` — every file appteam creates, organized by directory
  - `agent-roles.md` — PM, TPM, SWE, and optional agent roles explained
  - `pipeline.md` — mandatory development pipeline workflow
  - `skills.md` — the 5 generated Claude Code slash commands
  - `configuration.md` — settings.json schema, regeneration, manual editing

### Changed
- README.md rewritten with concise overview and links to full documentation

---

## v0.5.0 (2026-03-16)

### Added
- `.appteam/settings.json` — wizard saves full config as JSON for future reuse
- `--regenerate` / `-r` flag — regenerate all files from saved settings without wizard
- Saved config detection on startup with "Use saved config? (Y/n)" prompt
- Bootstrap skills generation — 5 Claude Code skills created in `.claude/skills/`:
  - `/spec` — create a new product spec (`docs/specs/F-NNNN-slug.md`)
  - `/release` — update release notes, commit, tag, push
  - `/pipeline` — spin up full agent team for a feature
  - `/status` — summarize current milestone from backlog
  - `/regenerate` — regenerate from saved settings
- Product specs: `docs/specs/F-0040-settings-json.md`, `docs/specs/F-0041-bootstrap-skills.md`

### Changed
- Generator now creates `.claude/skills/` directory and outputs 17 files (was 12)
- `--help` updated with `-r, --regenerate` option

### Fixed
- TargetDir portability: override with CWD when loading from settings.json

### Added Files
- `internal/config/settings.go` — SaveSettings, LoadSettings, SettingsExist, SettingsPath
- `internal/templates/skill_spec.go` — `/spec` skill template
- `internal/templates/skill_release.go` — `/release` skill template
- `internal/templates/skill_pipeline.go` — `/pipeline` skill template
- `internal/templates/skill_status.go` — `/status` skill template
- `internal/templates/skill_regenerate.go` — `/regenerate` skill template

---

## v0.4.2 (2026-03-16)

### Added
- PROGRESS.md, RELEASENOTES.md, and TAG nodes in MermaidJS pipeline diagram
- TPM → PROGRESS.md (session log) link in diagram
- PM → RELEASENOTES.md (version entry) and PM → TAG (tag release) links in diagram
- Mandatory pipeline steps for updating PROGRESS.md, RELEASENOTES.md, and creating git tags after every milestone

### Changed
- CLAUDE.md template pipeline now includes explicit steps for progress journaling, release notes, and tagging
- Pipeline Steps section uses numbered list with doc update and tagging steps at the end

---

## v0.4.1 (2026-03-16)

### Added
- `--help` / `-h` flag — prints usage text
- `--version` / `-v` flag — prints current version
- Unknown flag handling with error message and hint to `--help`
- Color-coded MermaidJS pipeline diagram: green for downward flow (request → execution), blue for upward flow (results → reporting), gray for side effects
- `add` and `linkRange` template functions for dynamic link index computation
- Legend section in generated PIPELINE.md

---

## v0.4.0 (2026-03-16)

### Added
- `docs/` directory structure for all project management context
- `docs/specs/` directory with `TEMPLATE.md` for PM product specs
- PM spec creation workflow — one spec file per feature/bug/enhancement (`F-NNNN-slug.md`)
- Spec column in backlog template with links to spec files
- Key References in CLAUDE.md template now includes `docs/specs/`, `docs/RELEASENOTES.md`, `docs/PIPELINE.md`
- Agent-to-file ownership matrix (PM owns specs, TPM owns backlog/progress)

### Changed
- Tracking files moved from project root to `docs/`: BACKLOG.md, PROGRESS.md, RELEASENOTES.md, PIPELINE.md
- Generator now writes to `docs/` and creates `docs/specs/` directory
- All agent templates (PM, TPM, SWE, SWE-Test, SWE-QA, Reviewer, Platform) updated with `docs/` paths
- PM template rewritten with spec creation responsibilities, workflow, and file naming conventions
- TPM template updated to reference specs and include spec links in backlog entries
- SWE, SWE-Test, SWE-QA, Reviewer templates updated to read specs for acceptance criteria
- CLAUDE.md template pipeline steps reference spec creation
- Root now only contains CLAUDE.md and README.md

### Added Files
- `internal/templates/spec_template.go` — Spec TEMPLATE.md template

---

## v0.3.0 (2026-03-16)

### Added
- MermaidJS pipeline diagram template (`PIPELINE.md`) showing the full agent workflow: PO → PM → TPM → SWE → Test/QA → Review → TPM → PM → PO
- Dynamic diagram rendering based on project config (SWE count, optional agents)
- Git repository initialization in wizard (new Step 2): detects existing `.git`, offers `git init`, `gh repo create`, or manual remote URL
- Model selection for agents: Opus 4.6 (default), Sonnet 4.6, Haiku 4.5
- BACKLOG.md now uses `F-0001` format IDs (sequential across all milestones, never reused) with Dependencies column

### Changed
- Wizard expanded to 7 steps (target directory moved to Step 1, git setup is Step 2, confirm is Step 7)
- Co-Authored-By lines parameterized with `{{.ModelName}}` across all agent templates (tpm, swe, swe-test, swe-qa, platform)
- CLAUDE.md template parameterized for model name and ID (no more hardcoded Opus 4.6)
- `setupGitRepo` now returns `error` — `git init` failure is fatal, `gh`/remote failures are warnings
- GitHub org default is now empty (no hardcoded "ahafin")

### Fixed
- Wizard git detection was checking CWD instead of target directory (target dir wasn't set yet)
- Dead code `_ = out` in generator git commands removed

### Added Files
- `internal/templates/pipeline.go` — MermaidJS pipeline diagram template

---

## v0.2.0 (2026-03-16)

### Added
- ANSI color support with automatic TTY detection
- `NO_COLOR` environment variable support (https://no-color.org/)
- Styled welcome banner with box-drawing characters
- Step progress indicators (`━━ Step N of 6 ━━ Title`)
- Green `▸` prompt markers with dim default values
- Dimmed `│` continuation lines for multi-line input
- Section dividers between wizard steps
- Box-drawn configuration summary with `✓ Yes` / `✗ No` boolean labels
- Styled generator output (green checkmarks, bold status messages)
- Generator now creates BACKLOG.md, PROGRESS.md, and RELEASENOTES.md starter files
- Comprehensive README with global/local install instructions and generated files documentation

### Changed
- `wizard.Run()` signature now accepts `color bool` parameter
- `generator.Generate()` signature now accepts `color bool` parameter
- `main.go` performs TTY detection and passes color flag through

### Fixed
- Banner box right-side `│` alignment (ANSI codes have zero display width)
- Summary column alignment with `%-14s` and ANSI-wrapped labels

### Added Files
- `internal/wizard/style.go` — Styler with color helpers and TTY detection
- `internal/templates/backlog.go` — BACKLOG.md template
- `internal/templates/progress.go` — PROGRESS.md template
- `internal/templates/releasenotes.go` — RELEASENOTES.md template

---

## v0.1.0 (2026-03-16)

### Added
- Interactive 6-step wizard for configuring agent teams
- Generates `CLAUDE.md` project instructions matching swole reference structure
- Generates agent definition files in `.claude/agents/`:
  - PM (Product Manager) — always generated
  - TPM (Technical Program Manager) — always generated
  - SWE-1 through SWE-5 — configurable count with custom titles and specialty bullets
  - SWE-Test (Test Engineer) — optional
  - SWE-QA (QA Engineer) — optional
  - Platform Engineer — optional
  - Reviewer (Code Review) — optional
- Optional GCP project configuration with conditional sections in CLAUDE.md
- Custom project conventions support
- Configurable target directory
- Confirmation summary before file generation
- Zero external dependencies (stdlib only)

### Tech Stack
- Go 1.25, `text/template`, `bufio`, `os`, `fmt`
