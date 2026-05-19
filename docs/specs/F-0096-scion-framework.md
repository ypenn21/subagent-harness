# F-0096: Scion as Alternative Team Framework

**Type:** Feature
**Priority:** P1 (important)
**Status:** Draft
**Requested by:** PO
**Date:** 2026-03-18

## Problem

appteam currently generates only Claude Code Agent Teams (files in `.claude/agents/` and `.claude/skills/`). Google Cloud Platform's **Scion** (https://github.com/GoogleCloudPlatform/scion) is a mature multi-agent orchestration tool that supports multiple harnesses (Gemini, Claude, OpenCode, Codex), container-based agent isolation, and git worktree separation. Teams that prefer Scion's container-based workflow, or want to use non-Claude harnesses, cannot use appteam today.

Adding Scion as an alternative framework option lets POs choose the orchestration tool that fits their infrastructure and preferences while reusing appteam's interactive wizard, team sizing, and skill content.

### What is Scion?

- Install: `go install github.com/GoogleCloudPlatform/scion/cmd/scion@latest`
- `scion init` creates a `.scion/` directory in the project
- Agents are defined as **templates** in `.scion/templates/<role>/`
- Each template directory contains:
  - `scion-agent.yaml` — config (name, description, agent_instructions, system_prompt, default_harness_config, env vars, resources)
  - `agents.md` — operational instructions for the agent
  - `system-prompt.md` — core persona and role definition
  - `home/` — optional config files
- Scion supports multiple **harnesses**: gemini, claude, opencode, codex
- Agents run in containers (Docker, Podman, K8s) with isolated git worktrees
- `scion start <name> "task" --type claude` launches an agent
- `scion list`, `scion attach`, `scion message`, `scion stop` for management
- Settings in `.scion/settings.json`

## Requirements

### Config

1. Add a `Framework string` field to `ProjectConfig` in `internal/config/config.go` with valid values: `"claude-code"`, `"scion"`
2. Update `backfillDefaults` logic in `internal/config/settings.go` to default `Framework` to `"claude-code"` when the field is empty or missing, ensuring backward compatibility with all pre-v0.15.0 `settings.json` files
3. The `Framework` field must be persisted in `.appteam/settings.json` and used by `-r` regeneration

### Wizard

4. Add a framework selection step to the wizard as **Step 2** (after Project Basics, before Git Repository), renumbering subsequent steps from 8 to 9 total:
   ```
   ━━ Step 2 of 9 ━━ Team Framework

   Which team framework?
     1) Claude Code Agent Teams  — .claude/agents/ + .claude/skills/ (default)
     2) Scion                    — .scion/templates/ with container-based agents

   Framework [1]:
   ```
5. When framework is `"scion"`, add a harness selection prompt:
   ```
   Default harness (claude/gemini/opencode/codex) [claude]:
   ```
   Store as `DefaultHarness string` in `ProjectConfig` (default `"claude"`)

6. The rest of the wizard flow remains the same — team sizing, agent roles, skills, conventions all apply regardless of framework. The wizard collects the same information; only the output format changes.

### Scion Templates

7. Create new Go template constants for Scion output files in `internal/templates/`:

   **`scion_agent_yaml.go`** — `ScionAgentYAMLTemplate`:
   ```yaml
   name: {{.Name}}
   description: {{.Description}}
   agent_instructions: agents.md
   system_prompt: system-prompt.md
   default_harness_config:
     type: {{.Harness}}
   ```

   **`scion_system_prompt.go`** — `ScionSystemPromptTemplate`:
   - Contains the role persona (same content as the current `.claude/agents/<role>.md` persona section)

   **`scion_agents_md.go`** — `ScionAgentsMDTemplate`:
   - Contains the operational instructions (same content as the current `.claude/agents/<role>.md` instructions section)

8. Each existing agent template (PM, TPM, SWE, SWE-Test, SWE-QA, Reviewer, Platform, custom agents) must have its content split into two conceptual parts for Scion output:
   - **System prompt** — the role definition, persona, and core responsibilities
   - **Agent instructions** — the operational procedures, workflows, and rules

   For Claude Code output, these remain combined in a single `.md` file as today. For Scion output, they are rendered into separate `system-prompt.md` and `agents.md` files.

### Generator

9. When `Framework == "scion"`, the generator must create the following directory structure instead of `.claude/agents/` and `.claude/skills/`:
   ```
   .scion/
     templates/
       pm/
         scion-agent.yaml
         agents.md
         system-prompt.md
       tpm/
         scion-agent.yaml
         agents.md
         system-prompt.md
       swe-1/
         scion-agent.yaml
         agents.md
         system-prompt.md
       swe-test/
         scion-agent.yaml
         agents.md
         system-prompt.md
       reviewer/
         scion-agent.yaml
         agents.md
         system-prompt.md
       ... (per team size and optional agents)
   ```

10. Agent template generation must respect `TeamSize` (lean/standard/full) the same way it does for Claude Code output — same agent inclusion rules apply.

11. Skills handling for Scion: Skills content is embedded into the agent's `agents.md` file as referenced procedures, since Scion agents don't have a separate `.claude/skills/` directory. Each agent's `agents.md` includes the skills relevant to their role.

### CLAUDE.md

12. CLAUDE.md is still generated when `Framework == "scion"` because Scion with the claude harness reads CLAUDE.md. However, the Interactive Agent Teams section must be adapted:
    - Replace tmux-based TeamCreate/Agent tool instructions with Scion commands:
      - `scion start <name> "task" --type <harness>` to launch agents
      - `scion list` to see running agents
      - `scion attach <name>` to interact with an agent
      - `scion message <name> "message"` to send messages
      - `scion stop <name>` to stop agents
    - Remove references to tmux panes, `TeamCreate`, `TeamDelete`, and subprocess agents
    - Keep the same pipeline structure (PO -> PM -> TPM -> SWE -> SWE-Test -> etc.)

13. When `Framework == "claude-code"`, CLAUDE.md output is unchanged from current behavior.

### Backward Compatibility

14. Default `Framework` to `"claude-code"` if not set in `settings.json` — all existing configs produce identical output
15. `DefaultHarness` defaults to `"claude"` if not set (only relevant when `Framework == "scion"`)

### Version

16. Version bump to `0.15.0` in `main.go`

### Tests

17. New tests must cover:
    - Config: `Framework` and `DefaultHarness` field serialization/deserialization, backfill defaults
    - Wizard: framework selection prompt parsing, harness selection for Scion, step renumbering
    - Templates: Scion YAML rendering, system-prompt rendering, agents.md rendering for each role
    - Generator: Scion directory structure creation (`.scion/templates/<role>/`), correct file contents, team size interaction
    - CLAUDE.md: Scion-specific agent management section (scion commands vs tmux)
    - Backward compat: pre-v0.15.0 settings load with `Framework` defaulting to `"claude-code"`, output unchanged

## Implementation Files

| File | Change |
|------|--------|
| `internal/config/config.go` | Add `Framework string` and `DefaultHarness string` fields to `ProjectConfig` |
| `internal/config/settings.go` | Backfill `Framework` to `"claude-code"` and `DefaultHarness` to `"claude"` |
| `internal/wizard/wizard.go` | New Step 2 (Framework), harness prompt, renumber steps 2-8 to 3-9 |
| `internal/templates/scion_agent_yaml.go` | New `ScionAgentYAMLTemplate` |
| `internal/templates/scion_system_prompt.go` | New `ScionSystemPromptTemplate` for each role |
| `internal/templates/scion_agents_md.go` | New `ScionAgentsMDTemplate` for each role |
| `internal/templates/claude_md.go` | Conditional agent management section (tmux vs scion commands) |
| `internal/generator/generator.go` | Conditional output: `.claude/agents/` vs `.scion/templates/<role>/`; skills embedding for Scion |
| `main.go` | Version bump to `0.15.0` |
| `internal/config/config_test.go` | Tests for Framework/DefaultHarness fields and backfill |
| `internal/templates/scion_*_test.go` | Tests for Scion template rendering |
| `internal/templates/claude_md_test.go` | Tests for CLAUDE.md Scion variant |
| `internal/generator/generator_test.go` | Tests for Scion directory structure and file generation |
| `internal/wizard/wizard_test.go` | Tests for framework selection wizard flow |

## Acceptance Criteria

- [ ] `ProjectConfig` has `Framework string` field with values `"claude-code"` and `"scion"`
- [ ] `ProjectConfig` has `DefaultHarness string` field with values `"claude"`, `"gemini"`, `"opencode"`, `"codex"`
- [ ] `backfillDefaults` sets `Framework` to `"claude-code"` and `DefaultHarness` to `"claude"` when empty/missing
- [ ] Pre-v0.15.0 `settings.json` files load correctly with defaults, producing identical output
- [ ] Wizard includes framework selection as Step 2 with descriptions of each option
- [ ] Wizard shows harness selection when Scion is chosen
- [ ] Total wizard steps increased from 8 to 9
- [ ] When `Framework == "scion"`, generator creates `.scion/templates/<role>/` directories
- [ ] Each Scion role directory contains `scion-agent.yaml`, `agents.md`, and `system-prompt.md`
- [ ] `scion-agent.yaml` contains correct name, description, and harness config
- [ ] Scion agent template generation respects `TeamSize` (lean/standard/full)
- [ ] Skills content is embedded into Scion agents' `agents.md` files
- [ ] CLAUDE.md is generated for both frameworks
- [ ] CLAUDE.md Scion variant references `scion start/list/attach/message/stop` instead of tmux/TeamCreate
- [ ] When `Framework == "claude-code"`, all output is identical to current behavior (no regression)
- [ ] Version reads `0.15.0`
- [ ] New tests cover config, wizard, template, generator, and CLAUDE.md changes
- [ ] All existing tests continue to pass

## Out of Scope

- Running `scion init` automatically (user runs this themselves before or after appteam)
- Generating Scion `settings.json` (distinct from appteam's `.appteam/settings.json`)
- Container/Docker configuration for Scion agents
- Scion-specific resource mounting or env var configuration beyond harness type
- Supporting both frameworks simultaneously in a single project (it's one or the other)
- Converting existing Claude Code output to Scion format (requires regeneration with `-r`)

## Dependencies

- F-0079 (v0.12.0 — Team Sizing) — team size logic must apply to both frameworks
- F-0064 (v0.9.0 — Custom Agents) — custom agents must generate Scion templates when framework is scion
- F-0057 (v0.8.0 — Skills) — skills must be embedded into Scion agent instructions

## Open Questions

- Should appteam generate a `.scion/settings.json` with default settings, or leave that to `scion init`? (Recommend: leave to `scion init` — out of scope for v0.15.0)
- Should the `home/` directory in each Scion template contain any default config files? (Recommend: skip for v0.15.0, add in a future iteration if needed)
- How should PM-only skills (spec, release, roadmap, etc.) be distributed across Scion agent templates? (Recommend: embed PM skills into PM's `agents.md`, dev skills into SWE `agents.md`, testing skills into SWE-Test/SWE-QA `agents.md`)
