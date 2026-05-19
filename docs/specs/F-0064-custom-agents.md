# F-0064: Custom Agent Definitions

**Type:** Feature
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-17

## Problem

The current agent system provides a fixed set of built-in roles: PM, TPM, SWE-1 through SWE-5, Platform, Reviewer, SWE-Test, and SWE-QA. While these cover many workflows, teams often need specialized agents that don't map to the defaults — for example, a Frontend Engineer, UX Designer, DBA, Data Analyst, DevOps Engineer, Security Engineer, or Tech Writer. Currently, users must manually create these agent files after generation, and they are not persisted across regenerations.

This feature lets users define custom agent roles during the wizard, which are then generated as `.claude/agents/<name>.md` files and persisted in `settings.json` so they survive regeneration with `-r`.

## Requirements

### Config

1. Add a new `CustomAgentConfig` struct to `internal/config/config.go`:
   - `Name` string — agent identifier in kebab-case, used as the filename (`<name>.md`)
   - `Title` string — display title (e.g., "Frontend Engineer")
   - `Description` string — what this agent does / its area of responsibility
   - `Instructions` []string — specific bullet-point instructions for this agent
2. Add `CustomAgents []CustomAgentConfig` field to `ProjectConfig`
3. Custom agents must be persisted in `settings.json` via the existing `config.SaveSettings` / `config.LoadSettings` flow
4. Backward compatibility: pre-v0.9.0 `settings.json` files with no `CustomAgents` field must work (nil/empty slice is the zero value)

### Wizard UX

5. In Step 5 (Agent Team), after the existing optional agent prompts (Platform, Reviewer, SWE-Test, SWE-QA) and before the model selection, add a custom agents section:
   ```
   Custom Agents (blank line to skip):
     Agent name (kebab-case): frontend-eng
     Title: Frontend Engineer
     Description: React/TypeScript specialist for UI components
     Instructions (blank line to finish):
       | Build responsive UI components using React and TypeScript
       | Follow design system patterns and accessibility standards
       | Write Storybook stories for all new components
       |
     Add another custom agent? (Y/n):
   ```
6. The loop continues until the user enters a blank name or declines to add another
7. Validation: agent name must be non-empty and should be kebab-case (lowercase alphanumeric + hyphens)

### Template

8. Create a new `CustomAgentTemplate` constant in `internal/templates/custom_agent.go`
9. The template must follow the same structure as existing agent templates (e.g., `SWETemplate`):
   - Role title and description
   - Project context (name, owner, tech stack)
   - Specific instructions (rendered from user-provided bullets)
   - Git conventions (commit format, co-authored-by line)
   - Responsibilities section
10. Create a `CustomAgentTemplateData` struct with fields:
    - `Name` string
    - `Title` string
    - `Description` string
    - `Instructions` []string
    - `ProjectName` string
    - `OwnerName` string
    - `OwnerEmail` string
    - `ModelName` string

### Generator

11. In `internal/generator/generator.go`, after generating built-in agents (SWE-Test, SWE-QA, Platform, Reviewer), loop through `cfg.CustomAgents` and generate `.claude/agents/<name>.md` for each
12. Use the `CustomAgentTemplate` with a `CustomAgentTemplateData` struct populated from both `ProjectConfig` and `CustomAgentConfig` fields
13. Follow the same `fileSpec` pattern used for SWE agents (append to the `files` slice, rendered by the shared loop)

### Summary Display

14. The confirmation summary (Step 8) must list custom agents alongside built-in agents, showing name and title for each

## Acceptance Criteria

- [ ] `CustomAgentConfig` struct exists in `internal/config/config.go` with Name, Title, Description, Instructions fields
- [ ] `ProjectConfig` has `CustomAgents []CustomAgentConfig` field
- [ ] Wizard prompts for custom agents in Step 5, after optional agents and before model selection
- [ ] Wizard loop accepts name, title, description, and multi-line instructions for each custom agent
- [ ] Blank name or declining "Add another?" exits the loop
- [ ] `CustomAgentTemplate` constant exists in `internal/templates/custom_agent.go`
- [ ] Template renders role, description, instructions, project context, and git conventions
- [ ] Generator creates `.claude/agents/<name>.md` for each custom agent
- [ ] Custom agents persist in `settings.json` and regenerate correctly with `-r`
- [ ] Existing built-in agents are unaffected
- [ ] Confirmation summary shows custom agents with name and title
- [ ] Tests cover: config round-trip (save/load with custom agents), template rendering, wizard input parsing, generator output
- [ ] Backward compatibility: `settings.json` without `CustomAgents` field loads without error (nil slice)

## Out of Scope

- Agent-to-agent dependencies or communication rules (future work)
- Custom agent file ownership declarations (future work)
- Editing/removing individual custom agents during regeneration (future work)
- Custom agent template customization beyond the provided fields (future work)

## Dependencies

- None — builds on existing config, wizard, template, and generator infrastructure

## Open Questions

- None — design mirrors the established patterns from SWE agents and custom skills
