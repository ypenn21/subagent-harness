# appteam — Backlog

**Maintained by:** TPM
**Last updated:** 2026-03-26 (v0.18.0 — completed)

---

## How to Read This Backlog

- **ID:** Unique feature identifier (`F-0001`, `F-0002`, etc.) — sequential across all milestones, never reused
- **Priority:** P0 (critical path), P1 (important), P2 (nice to have)
- **Status:** `TODO` | `IN PROGRESS` | `DONE` | `BLOCKED`
- **Owner:** Assigned team member
- **Dependencies:** Other feature IDs that must complete first

---

## v0.18.0 (Completed — 2026-03-26)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0122 | Add YAML frontmatter to all 8 agent templates (PM, TPM, SWE, SWE-Test, SWE-QA, Reviewer, Platform, Custom) | P1 | DONE | SWE | — | [spec](docs/specs/F-0122-agent-skill-frontmatter.md) | name, description, model fields; disallowedTools for PM/TPM/Reviewer |
| F-0123 | Add YAML frontmatter to all 21 skill templates | P1 | DONE | SWE | — | [spec](docs/specs/F-0122-agent-skill-frontmatter.md) | name, description, user-invocable; disable-model-invocation for workflow skills |
| F-0124 | Update tests for frontmatter in agent and skill template output | P1 | DONE | SWE | F-0122, F-0123 | [spec](docs/specs/F-0122-agent-skill-frontmatter.md) | Verify frontmatter markers, name/description fields present in rendered output |
| F-0125 | Version bump to 0.18.0 | P1 | DONE | SWE | F-0122, F-0123, F-0124 | [spec](docs/specs/F-0122-agent-skill-frontmatter.md) | Version constant in main.go |
| F-0126 | Test verification | P1 | DONE | SWE-Test | F-0125 | [spec](docs/specs/F-0122-agent-skill-frontmatter.md) | 97 tests, 0 failures |

---

## v0.17.1 (Completed — 2026-03-19)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0121 | Make all CLAUDE.md "feature branch" references conditional on framework | P0 | DONE | — | F-0117 | — | 7 pipeline/version-control sections now say "worktree branch" for Scion, "feature branch" for Claude Code |

---

## v0.17.0 (Completed — 2026-03-19)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0117 | Add Worktree Isolation section to all Scion agent templates with git merge instructions | P0 | DONE | SWE-2 | — | [spec](docs/specs/F-0117-scion-merge-instructions.md) | PM: no merge; TPM: merge pm; SWE: merge pm (lean) or tpm (std/full); SWE-Test: merge swe-N; SWE-QA: merge swe-test; Reviewer: merge swe-qa or swe-test; Platform: merge pm or tpm; Custom: generic |
| F-0118 | Tests for Worktree Isolation merge instructions | P0 | DONE | SWE-2 | F-0117 | [spec](docs/specs/F-0117-scion-merge-instructions.md) | Assert "Worktree Isolation" section in all templates; verify team-size-conditional merge commands |
| F-0119 | Version bump to 0.17.0 | P0 | DONE | SWE-2 | F-0117, F-0118 | [spec](docs/specs/F-0117-scion-merge-instructions.md) | Version constant in main.go |
| F-0120 | Test verification | P0 | DONE | SWE-Test | F-0119 | [spec](docs/specs/F-0117-scion-merge-instructions.md) | All tests pass, no regressions |

---

## v0.16.2 (In Progress)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0114 | Generate `.scion/grove-id` (UUID v4) when Framework == "scion" | P0 | DONE | SWE-2 | — | [spec](docs/specs/F-0114-scion-grove-id.md) | Use crypto/rand; also create `.scion/agents/` directory |
| F-0115 | Tests for grove-id generation | P0 | DONE | SWE-2 | F-0114 | [spec](docs/specs/F-0114-scion-grove-id.md) | UUID format validation, conditional generation tests, uniqueness |
| F-0116 | Test verification | P0 | TODO | SWE-Test | F-0115 | [spec](docs/specs/F-0114-scion-grove-id.md) | All tests pass, no regressions |

---

## v0.16.1 (Completed — 2026-03-19)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0109 | Fix Scion signaling: replace `scion message` with `sciontool status task_completed` in agent templates | P0 | DONE | SWE-2 | F-0103 | [spec](docs/specs/F-0109-scion-signaling-fix.md) | Docker CLI unavailable inside containers; agents use `sciontool status` instead; Docker unavailability warning added |
| F-0110 | Update `/pipeline` skill and CLAUDE.md for orchestrator relay model | P0 | DONE | SWE-2 | F-0109 | [spec](docs/specs/F-0109-scion-signaling-fix.md) | Agents signal completion, orchestrator relays from host; all 3 team sizes (lean/standard/full) updated |
| F-0111 | Fix `.gitignore` entry: `.scion/` → `.scion/agents/` | P0 | DONE | SWE-2 | — | [spec](docs/specs/F-0109-scion-signaling-fix.md) | Keeps `.scion/templates/` tracked in git |
| F-0112 | Version bump to 0.16.1 + update Scion test assertions | P0 | DONE | SWE-2 | F-0109, F-0110, F-0111 | [spec](docs/specs/F-0109-scion-signaling-fix.md) | Version constant, updated test assertions for signaling protocol content |
| F-0113 | Test verification | P0 | DONE | SWE-Test | F-0112 | [spec](docs/specs/F-0109-scion-signaling-fix.md) | All Scion-related tests pass with updated assertions; no regressions in non-Scion tests |

---

## v0.16.0 (Completed — 2026-03-18)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0103 | Update `/pipeline` skill Scion section for parallel agent launch | High | DONE | SWE-1 | — | [spec](docs/specs/F-0103-scion-parallel-launch.md) | Launch ALL agents simultaneously; orchestrator monitors only; `scion stop` all on completion |
| F-0104 | Add inter-agent messaging protocol to Scion agent templates | High | DONE | SWE-1 | F-0103 | [spec](docs/specs/F-0103-scion-parallel-launch.md) | Each role gets "Messaging Protocol" section: who to wait for, who to message, exact `scion message` commands |
| F-0105 | Update CLAUDE.md Scion orchestration sections for parallel launch | High | DONE | SWE-1 | F-0103 | [spec](docs/specs/F-0103-scion-parallel-launch.md) | All three team sizes (lean/standard/full): parallel launch, monitor-only orchestrator |
| F-0106 | Team size variations for messaging protocol | High | DONE | SWE-1 | F-0104 | [spec](docs/specs/F-0103-scion-parallel-launch.md) | Lean: PM→SWE-1 direct; Standard: full chain with TPM; Full: includes SWE-QA in chain |
| F-0107 | Version bump to 0.16.0 + tests | High | DONE | SWE-1 | F-0103, F-0104, F-0105, F-0106 | [spec](docs/specs/F-0103-scion-parallel-launch.md) | Version constant, 13 new test functions with 20+ subtests, no regressions |
| F-0108 | Test verification | High | DONE | SWE-Test | F-0107 | [spec](docs/specs/F-0103-scion-parallel-launch.md) | All tests pass; messaging protocol content verified in template output |

---

## v0.15.0 (Completed — 2026-03-18)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0096 | Add `Framework` field to ProjectConfig (claude-code / scion) | High | DONE | SWE-1 | — | [spec](docs/specs/F-0096-scion-framework.md) | `Framework string` + `DefaultHarness string` fields, backfill to "claude-code" |
| F-0097 | Add framework selection to wizard (new Step 2) | High | DONE | SWE-1 | F-0096 | [spec](docs/specs/F-0096-scion-framework.md) | Framework choice + harness prompt for Scion; renumber steps 2-8 to 3-9 |
| F-0098 | Create Scion template constants (YAML, system-prompt, agents.md) | High | DONE | SWE-2 | — | [spec](docs/specs/F-0096-scion-framework.md) | `scion_agent_yaml.go`, `scion_system_prompt.go`, `scion_agents_md.go` |
| F-0099 | Update generator for Scion output (.scion/templates/) | High | DONE | SWE-2 | F-0096, F-0098 | [spec](docs/specs/F-0096-scion-framework.md) | Conditional directory structure, skills embedding, team size interaction |
| F-0100 | Update CLAUDE.md template for Scion agent management | High | DONE | SWE-2 | F-0096 | [spec](docs/specs/F-0096-scion-framework.md) | Replace tmux/TeamCreate with scion start/list/attach/message/stop commands |
| F-0101 | Version bump to 0.15.0 + tests | High | DONE | SWE-1 | F-0096, F-0097, F-0098, F-0099, F-0100 | [spec](docs/specs/F-0096-scion-framework.md) | Version constant, tests for config/wizard/templates/generator/CLAUDE.md |
| F-0102 | Test verification | High | DONE | SWE-Test | F-0096, F-0097, F-0098, F-0099, F-0100, F-0101 | [spec](docs/specs/F-0096-scion-framework.md) | 148 tests across 5 packages — all pass; backward compat verified |

---

## v0.14.0 (Completed — 2026-03-18)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0095 | GitHub Pages site with Docsify | High | DONE | SWE-1 | — | [spec](docs/specs/F-0095-github-pages-site.md) | Docsify SPA, dark theme, Why/What/How/Roadmap pages, Reveal.js slide deck |

---

## v0.13.1 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0094 | Add cleanup phase to /cuj-test skill template | High | DONE | SWE-1 | — | [spec](docs/specs/F-0094-cuj-cleanup-step.md) | Phase 6: browser close, artifact cleanup, test data cleanup, state reset, cleanup report; SWE-QA workflow updated; version 0.13.1 |

---

## v0.13.0 (Completed — 2026-03-18)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0088 | Rewrite SWE-QA agent template with CUJ-oriented testing | High | DONE | SWE-1 | — | [spec](docs/specs/F-0088-swe-qa-cuj-testing.md) | Comprehensive Puppeteer/Chromium instructions, CUJ inventory management, test reporting, Lighthouse |
| F-0089 | Create `/cuj-list` skill template | High | DONE | SWE-1 | — | [spec](docs/specs/F-0088-swe-qa-cuj-testing.md) | New SkillCUJListTemplate for CUJ inventory creation/update |
| F-0090 | Create `/cuj-test` skill template | High | DONE | SWE-1 | — | [spec](docs/specs/F-0088-swe-qa-cuj-testing.md) | New SkillCUJTestTemplate for headless Chromium CUJ testing |
| F-0091 | Register cuj-list and cuj-test in wizard, generator, settings | High | DONE | SWE-1 | F-0089, F-0090 | [spec](docs/specs/F-0088-swe-qa-cuj-testing.md) | Tier 2 default N, backfillSkillDefaults, skillFiles, allSkills |
| F-0092 | Version bump to 0.13.0 + tests | High | DONE | SWE-1 | F-0088, F-0089, F-0090, F-0091 | [spec](docs/specs/F-0088-swe-qa-cuj-testing.md) | Version constant, template rendering tests, backfill tests |
| F-0093 | Test verification | High | DONE | SWE-Test | F-0088, F-0089, F-0090, F-0091, F-0092 | [spec](docs/specs/F-0088-swe-qa-cuj-testing.md) | 148 tests (80 top-level + subtests) across 5 packages — all pass; new CUJ template tests |

---

## v0.12.0 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0079 | Add `TeamSize` field to ProjectConfig + settings backfill | High | DONE | SWE-1 | — | [spec](docs/specs/F-0079-team-sizing.md) | TeamSize string (lean/standard/full), backfill to "standard" |
| F-0083 | Add team size selection to wizard Step 5 | High | DONE | SWE-1 | F-0079 | [spec](docs/specs/F-0079-team-sizing.md) | Team size prompt + conditional agent prompts for lean/standard/full |
| F-0084 | Update CLAUDE.md template with conditional pipeline rules | High | DONE | SWE-2 | F-0079 | [spec](docs/specs/F-0079-team-sizing.md) | Lean/standard/full pipeline variants, roles per team size |
| F-0080 | Add backlog enforcement section to CLAUDE.md template | High | DONE | SWE-2 | F-0084 | [spec](docs/specs/F-0079-team-sizing.md) | Non-negotiable backlog tracking rule in all team sizes |
| F-0085 | Update generator for conditional agent file generation | High | DONE | SWE-2 | F-0079, F-0084 | [spec](docs/specs/F-0079-team-sizing.md) | Agent files generated based on TeamSize |
| F-0086 | Version bump to 0.12.0 + tests | High | DONE | SWE-1 | F-0079, F-0083 | [spec](docs/specs/F-0079-team-sizing.md) | Version constant, new tests for all three team sizes |
| F-0087 | Test verification | High | DONE | SWE-Test | F-0079, F-0083, F-0084, F-0085, F-0086 | [spec](docs/specs/F-0079-team-sizing.md) | 143 tests (72 top-level + subtests) across 5 packages — all pass; new team sizing tests |

---

## v0.11.0 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0074 | Create `SkillBrainstormTemplate` in `internal/templates/skill_brainstorm.go` | High | DONE | SWE-1 | — | [spec](docs/specs/F-0074-brainstorm-skill.md) | 5-phase brainstorm flow template with project context variables |
| F-0075 | Register `/brainstorm` in wizard Step 7 + `skillsSummary()` | High | DONE | SWE-1 | F-0074 | [spec](docs/specs/F-0074-brainstorm-skill.md) | Under "Project Management Skills" with default Y; add to allSkills |
| F-0076 | Register brainstorm in generator `skillFiles` + `backfillSkillDefaults()` | High | DONE | SWE-1 | F-0074 | [spec](docs/specs/F-0074-brainstorm-skill.md) | Generator entry + backward-compat backfill default true |
| F-0077 | Version bump to 0.11.0 + tests | High | DONE | SWE-1 | F-0074, F-0075, F-0076 | [spec](docs/specs/F-0074-brainstorm-skill.md) | Version constant, new template rendering tests, verify all existing tests pass |
| F-0078 | Test verification | High | DONE | SWE-Test | F-0074, F-0075, F-0076, F-0077 | [spec](docs/specs/F-0074-brainstorm-skill.md) | 125 tests (61 top-level + 64 subtests) across 5 packages — all pass; new brainstorm template tests, backward compat |

---

## v0.10.0 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0070 | Add `-d`/`--dir` flag for target directory with auto-creation | High | DONE | SWE-1 | — | [spec](docs/specs/F-0070-target-dir-flag.md) | Parse `-d <folder>` / `--dir <folder>` CLI flag; create directory (incl. parents) if missing; update `--help` output; error on empty value |
| F-0071 | Update wizard to skip target dir prompt when `-d` is provided | High | DONE | SWE-1 | F-0070 | [spec](docs/specs/F-0070-target-dir-flag.md) | Pre-fill `cfg.TargetDir` with `-d` value; skip Step 1 target directory prompt; support `-r -d` combo |
| F-0072 | Version bump to 0.10.0 + tests | High | DONE | SWE-1 | F-0071 | [spec](docs/specs/F-0070-target-dir-flag.md) | Version constant → "0.10.0" in main.go; unit tests for flag parsing, dir creation, wizard skip, `-r -d` integration |
| F-0073 | Test verification | High | DONE | SWE-Test | F-0070, F-0071, F-0072 | [spec](docs/specs/F-0070-target-dir-flag.md) | 64 tests across 5 packages — all pass; 6 new `-d` flag tests; existing tests pass |

---

## v0.9.0 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0064 | Add `CustomAgentConfig` struct + `CustomAgents` field to `ProjectConfig` | High | DONE | SWE-1 | — | [spec](docs/specs/F-0064-custom-agents.md) | New `CustomAgentConfig{Name, Title, Description, Instructions}` struct; `CustomAgents []CustomAgentConfig` in ProjectConfig |
| F-0065 | Add custom agent input to wizard Step 5 | High | DONE | SWE-1 | F-0064 | [spec](docs/specs/F-0064-custom-agents.md) | After optional agents (Platform, Reviewer, SWE-Test, SWE-QA), before model selection; loop with name/title/description/instructions |
| F-0066 | Create `CustomAgentTemplate` | High | DONE | SWE-2 | — | [spec](docs/specs/F-0064-custom-agents.md) | `internal/templates/custom_agent.go`; `CustomAgentTemplateData` struct; role/description/instructions/project context/git conventions |
| F-0067 | Update generator for custom agent file generation | High | DONE | SWE-2 | F-0064, F-0066 | [spec](docs/specs/F-0064-custom-agents.md) | Loop through `cfg.CustomAgents`, generate `.claude/agents/<name>.md` using `CustomAgentTemplate` + `fileSpec` pattern |
| F-0068 | Version bump to 0.9.0 + test updates | High | DONE | SWE-1 | F-0065 | [spec](docs/specs/F-0064-custom-agents.md) | Version constant → "0.9.0" in main.go; update existing tests for custom agent support |
| F-0069 | Test verification and coverage | High | DONE | SWE-Test | F-0064, F-0065, F-0066, F-0067, F-0068 | [spec](docs/specs/F-0064-custom-agents.md) | 117 tests across 5 packages — all pass; new tests for custom agent config round-trip, template rendering, wizard input, generator output |

---

## v0.8.0 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0057 | Add `CustomSkillConfig` struct + `CustomSkills` field to `ProjectConfig` | High | DONE | SWE-1 | — | [spec](docs/specs/F-0057-universal-skills.md) | New `CustomSkillConfig{Name, Description}` struct; `CustomSkills []CustomSkillConfig` in ProjectConfig |
| F-0058 | Update wizard Step 7 with grouped skills display + custom skills input | High | DONE | SWE-1 | F-0057 | [spec](docs/specs/F-0057-universal-skills.md) | Two groups: "Project Management Skills:" / "Development Skills:" headers; custom skill `name: description` input loop |
| F-0059 | Create 6 Tier 1 dev skill templates (debug, test, review, docs, refactor, hotfix) | High | DONE | SWE-2 | — | [spec](docs/specs/F-0057-universal-skills.md) | 6 new template files in `internal/templates/`; Tier 1 defaults Y |
| F-0060 | Create 6 Tier 2 dev skill templates (api-design, schema, deploy, security, adr, standup) | High | DONE | SWE-2 | — | [spec](docs/specs/F-0057-universal-skills.md) | 6 new template files in `internal/templates/`; Tier 2 defaults N |
| F-0061 | Update generator for 18 predefined skills + custom skill generation | High | DONE | SWE-2 | F-0057, F-0059, F-0060 | [spec](docs/specs/F-0057-universal-skills.md) | Expanded `skillFiles` slice (18 entries); custom skill loop creating `.claude/skills/<name>.md` |
| F-0062 | Version bump to 0.8.0 + update tests | High | DONE | SWE-1 | F-0058 | [spec](docs/specs/F-0057-universal-skills.md) | Version constant → "0.8.0" in main.go; updated existing tests for new skill count |
| F-0063 | Test verification and coverage | High | DONE | SWE-Test | F-0057, F-0058, F-0059, F-0060, F-0061, F-0062 | [spec](docs/specs/F-0057-universal-skills.md) | 55 tests across 5 packages — all pass; new tests for custom skill config, backward-compat backfill, custom skill generation |

---

## v0.7.0 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0051 | Add `SelectedSkills` field to ProjectConfig | High | DONE | SWE-1 | — | [spec](docs/specs/F-0051-roadmap-selectable-skills.md) | `SelectedSkills map[string]bool` in config.go |
| F-0052 | Add skills selection step to wizard | High | DONE | SWE-1 | F-0051 | [spec](docs/specs/F-0051-roadmap-selectable-skills.md) | New Step 7 (Skills), bump totalSteps to 8, update Confirm to Step 8 |
| F-0053 | Create /roadmap skill template | High | DONE | SWE-2 | — | [spec](docs/specs/F-0051-roadmap-selectable-skills.md) | `SkillRoadmapTemplate` in `internal/templates/skill_roadmap.go` |
| F-0054 | Update generator for conditional skill generation | High | DONE | SWE-2 | F-0051, F-0053 | [spec](docs/specs/F-0051-roadmap-selectable-skills.md) | Replace unconditional skill append with `SelectedSkills` check; backward-compat nil-map fallback |
| F-0055 | Update tests for new functionality | High | DONE | SWE-Test | F-0051, F-0052, F-0053, F-0054 | [spec](docs/specs/F-0051-roadmap-selectable-skills.md) | 6 new tests: skill_roadmap template, skills wizard step, conditional generation, backward compat |
| F-0056 | Version bump to 0.7.0 | Low | DONE | SWE-1 | — | [spec](docs/specs/F-0051-roadmap-selectable-skills.md) | Version constant updated to "0.7.0" in main.go |

---

## v0.6.0 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0045 | Config package tests (settings.go, config.go) | High | DONE | SWE-1 | — | [spec](docs/specs/F-0045-test-coverage.md) | SettingsPath, SettingsExist, SaveSettings/LoadSettings round-trip, error cases |
| F-0046 | Template package tests (Render, template functions, all templates) | High | DONE | SWE-2 | — | [spec](docs/specs/F-0045-test-coverage.md) | add/sub/seq/linkRange, ClaudeMD ±GCP ±agents, SWE ±bullets, all 18+ templates |
| F-0047 | Generator package tests (file creation, directory structure) | High | DONE | SWE-2 | — | [spec](docs/specs/F-0045-test-coverage.md) | Directory creation, core/SWE/optional/skill files, settings persistence |
| F-0048 | Wizard package tests (Styler, ask helpers, integration) | High | DONE | SWE-1 | — | [spec](docs/specs/F-0045-test-coverage.md) | Styler color on/off, PadBold, Banner, StepHeader, Divider, ask/askBool/askInt |
| F-0049 | Main package tests (flag parsing, CLI behavior) | High | DONE | SWE-1 | — | [spec](docs/specs/F-0045-test-coverage.md) | --help, --version, unknown flag, -r without settings |
| F-0050 | Test verification and coverage report | High | DONE | SWE-Test | F-0045, F-0046, F-0047, F-0048, F-0049 | [spec](docs/specs/F-0045-test-coverage.md) | 45 tests, all pass. Coverage: config 85%, templates 100%, generator 58%, wizard 82% |

---

## v0.5.1 (Completed — 2026-03-17)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0043 | User-facing documentation | High | DONE | SWE-1 | F-0041 | — | 8 guide pages in `docs/guide/`: getting started, CLI ref, generated files, agent roles, pipeline, skills, configuration |
| F-0044 | README.md update | High | DONE | SWE-1 | F-0043 | — | Concise overview with links to full guide |

## v0.5.0 (Completed — 2026-03-16)

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0040 | `.appteam/settings.json` non-interactive mode | High | DONE | SWE-1 | — | [spec](docs/specs/F-0040-settings-json.md) | Save config as JSON, load on startup, `-r` flag |
| F-0041 | Bootstrap skills generation | High | DONE | SWE-2 | F-0040 | [spec](docs/specs/F-0041-bootstrap-skills.md) | 5 Claude Code skills: `/spec`, `/release`, `/pipeline`, `/status`, `/regenerate` |
| F-0042 | TargetDir portability fix | Medium | DONE | Reviewer | F-0040 | — | Override TargetDir with CWD when loading from settings.json |

## v0.4.2 (Completed — 2026-03-16)

| ID | Feature | Priority | Status | Owner | Dependencies | Notes |
|----|---------|----------|--------|-------|--------------|-------|
| F-0039 | Add PROGRESS, RELEASENOTES, TAG to pipeline | High | DONE | SWE-1 | F-0037 | TPM→PROGRESS, PM→RELEASENOTES, PM→TAG nodes + mandatory workflow steps |

## v0.4.1 (Completed — 2026-03-16)

| ID | Feature | Priority | Status | Owner | Dependencies | Notes |
|----|---------|----------|--------|-------|--------------|-------|
| F-0037 | Color pipeline diagram links | Medium | DONE | SWE-1 | F-0036 | Green downward, blue upward, gray side effects; `add`/`linkRange` funcs |
| F-0038 | `--help` / `--version` flags | Low | DONE | SWE-1 | — | `-h`/`--help`, `-v`/`--version`, unknown flag error handling |

## v0.4.0 (Completed — 2026-03-16)

| ID | Feature | Priority | Status | Owner | Dependencies | Notes |
|----|---------|----------|--------|-------|--------------|-------|
| F-0031 | Move tracking files to `docs/` directory | High | DONE | SWE-1 | — | BACKLOG, PROGRESS, RELEASENOTES, PIPELINE → `docs/` |
| F-0032 | PM spec creation system | High | DONE | SWE-1 | F-0031 | `docs/specs/` directory + TEMPLATE.md + PM workflow |
| F-0033 | Update all agent templates with `docs/` paths | High | DONE | SWE-1 | F-0031 | PM, TPM, SWE, SWE-Test, SWE-QA, Reviewer, Platform |
| F-0034 | Add Spec column to backlog template | High | DONE | SWE-1 | F-0032 | Links to `docs/specs/F-NNNN-slug.md` |
| F-0035 | Update CLAUDE.md template Key References | High | DONE | SWE-1 | F-0031, F-0032 | All `docs/` paths + specs directory |
| F-0036 | Update pipeline template for `docs/` paths | Medium | DONE | SWE-1 | F-0031 | BACKLOG.md node → `docs/BACKLOG.md` |

## v0.3.0 (Completed — 2026-03-16)

| ID | Feature | Priority | Status | Owner | Dependencies | Notes |
|----|---------|----------|--------|-------|--------------|-------|
| F-0022 | MermaidJS pipeline diagram template | High | DONE | SWE-2 | — | `internal/templates/pipeline.go`, generates PIPELINE.md |
| F-0023 | Git repo initialization in wizard | High | DONE | SWE-1 | — | New Step 2, git init, gh repo create, remote add |
| F-0024 | Model selection for agents | High | DONE | SWE-1 | — | Opus 4.6 / Sonnet 4.6 / Haiku 4.5, parameterized templates |
| F-0025 | Code review of v0.3.0 changes | High | DONE | Reviewer | F-0022, F-0023, F-0024 | Approved |

## v0.2.0 (Completed — 2026-03-16)

| ID | Feature | Priority | Status | Owner | Dependencies | Notes |
|----|---------|----------|--------|-------|--------------|-------|
| F-0008 | ANSI color support with TTY detection | High | DONE | SWE-1 | — | `style.go` — Styler, `IsTTY()`, `NO_COLOR` support |
| F-0009 | Styled welcome banner | High | DONE | SWE-1 | — | Box-drawing chars, bold cyan title, dim subtitle |
| F-0010 | Step progress indicators | High | DONE | SWE-1 | — | `━━ Step N of 6 ━━ Title` format |
| F-0011 | Styled prompts | High | DONE | SWE-1 | — | Green `▸`, dim defaults, `│` continuation |
| F-0012 | Section dividers | Medium | DONE | SWE-1 | — | Dimmed `─` horizontal rules between steps |
| F-0013 | Polished configuration summary | High | DONE | SWE-1 | — | Box-drawing border, `✓ Yes` / `✗ No`, `•` bullets |
| F-0014 | Styled generator output | Medium | DONE | SWE-1 | — | Green `✓`, bold headers, color flag passthrough |
| F-0015 | Fix banner box misalignment | High | DONE | SWE-1 | F-0009 | `boxLine()` helper pads visible text before ANSI |
| F-0016 | Fix summary `%-14s` ANSI alignment | High | DONE | SWE-1 | F-0013 | `PadBold()` pads raw text before wrapping |
| F-0017 | BACKLOG.md template | High | DONE | SWE-1 | — | `internal/templates/backlog.go` |
| F-0018 | PROGRESS.md template | High | DONE | SWE-1 | — | `internal/templates/progress.go` |
| F-0019 | RELEASENOTES.md template | High | DONE | SWE-1 | — | `internal/templates/releasenotes.go` |
| F-0020 | Update README.md | High | DONE | SWE-1 | — | Install, usage, generated files docs |
| F-0021 | Code review | High | DONE | Reviewer | F-0008–F-0020 | Approved — all checks pass |

## v0.1.0 (Completed — 2026-03-16)

| ID | Feature | Priority | Status | Owner | Dependencies | Notes |
|----|---------|----------|--------|-------|--------------|-------|
| F-0001 | Project scaffolding | High | DONE | — | — | go.mod, .gitignore, directory structure |
| F-0002 | Config data structures | High | DONE | — | — | ProjectConfig, SWEConfig, GCPConfig |
| F-0003 | Template system | High | DONE | — | F-0002 | 8 template constants + Render() helper |
| F-0004 | Interactive wizard | High | DONE | — | F-0002 | 6-step flow with ask/askBool/askInt |
| F-0005 | File generator | High | DONE | — | F-0003, F-0004 | CLAUDE.md + agent file generation |
| F-0006 | Dogfooding | Medium | DONE | — | F-0005 | Generated CLAUDE.md + agents for appteam itself |
| F-0007 | README.md | Medium | DONE | — | — | Install, usage, tech stack docs |

---

## Future Items

| ID | Feature | Priority | Status | Dependencies | Notes |
|----|---------|----------|--------|--------------|-------|
| F-0026 | CLI flags for non-interactive mode | Medium | SUPERSEDED | — | Replaced by F-0040 (`.appteam/settings.json` + `-r` flag) |
| F-0027 | `--dry-run` flag | Low | TODO | — | Preview output without writing files |
| F-0028 | Custom agent definitions | Low | DONE | — | Implemented as F-0064–F-0069 in v0.9.0 |
| F-0029 | Update/regenerate mode | Medium | SUPERSEDED | — | Replaced by F-0040 (`.appteam/settings.json` + `-r` flag) |
| F-0030 | `--version` flag | Low | DONE | — | Implemented as F-0038 |
| F-0079 | Team sizing — lean / standard / full team configurations | P1 | DONE | — | Implemented as F-0079, F-0083–F-0087 in v0.12.0 |
| F-0080 | Backlog enforcement in all team modes | P1 | DONE | F-0079 | Implemented as part of F-0080 in v0.12.0 |
| F-0081 | Agent self-assessment / retrospective skill | P1 | TODO | — | `/retrospective` skill — agents review their own performance after milestones, capture lessons learned, feed improvements back into agent `.md` files over time |
| F-0082 | Multi-project shared conventions | P2 | TODO | — | Global profile (`~/.appteam/global.json`) merged into every project's settings — shared git config, preferred patterns, common skills across projects |
