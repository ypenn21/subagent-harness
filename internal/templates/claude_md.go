package templates

const ClaudeMDTemplate = `# {{.ProjectName}} — Project Instructions

## Project Overview

{{.Description}}

## Key References

- [README.md](README.md) — Project overview, tech stack
- [docs/BACKLOG.md](docs/BACKLOG.md) — Feature backlog with priorities, status, dependencies, and feedback{{if eq .TeamSize "lean"}} (owned by PM){{else}} (owned by TPM){{end}}
- [docs/PROGRESS.md](docs/PROGRESS.md) — Session-by-session development log (update every session)
- [docs/RELEASENOTES.md](docs/RELEASENOTES.md) — Version history (Keep a Changelog format, owned by PM{{if ne .TeamSize "lean"}}/TPM{{end}})
- [docs/PIPELINE.md](docs/PIPELINE.md) — MermaidJS agent workflow diagram
- [docs/specs/](docs/specs/) — Product requirement specs (one per feature, owned by PM)
- [docs/specs/TEMPLATE.md](docs/specs/TEMPLATE.md) — Spec template for PM
{{- if .GCP.Enabled}}

## GCP Project

- **Project ID:** ` + "`{{.GCP.ProjectID}}`" + `
- **Project Number:** ` + "`{{.GCP.ProjectNumber}}`" + `
- **Organization:** ` + "`{{.GCP.Organization}}`" + `
- **Region:** ` + "`{{.GCP.Region}}`" + ` (primary — matches existing Cloud Run services)

### Service Accounts

| Account | Email | Role | Purpose |
|---------|-------|------|---------|
| **owner-sa** | ` + "`owner-sa@{{.GCP.ProjectID}}.iam.gserviceaccount.com`" + ` | Owner | Full project admin — used by Platform Engineer for all GCP operations |

### Credentials

- **` + "`owner-sa-key.json`" + `** — Owner SA key. Used for all ` + "`gcloud`" + ` CLI interactions. **Never use ` + "`GOOGLE_APPLICATION_CREDENTIALS`" + `** — always reference the key file directly (e.g., ` + "`--key-file=owner-sa-key.json`" + ` or load explicitly in code).
- Protected by ` + "`.gitignore`" + ` pattern ` + "`*-sa-key.json`" + `. **Never commit these.**
{{- end}}

## Product Owner / CEO

- **Name:** {{.OwnerName}}
- **Role:** Product Owner (PO) and CEO — the human in the loop
- **GitHub:** {{.OwnerGitHub}}
- **Email:** {{.OwnerEmail}}
- **Git config:** Always use ` + "`git -c user.name=\"{{.OwnerName}}\" -c user.email=\"{{.OwnerEmail}}\"`" + ` for commits so authorship is consistent in history.

## Version Control

- **Commit frequently** — after each meaningful change (new feature, bug fix, refactor, config change). Small, focused commits over large monolithic ones.
- **Write verbose commit messages** — first line is a concise summary (imperative mood, under 72 chars), followed by a blank line and a detailed body explaining *what* changed and *why*. Include context that won't be obvious from the diff.
- **Never commit secrets** — ` + "`.gitignore`" + ` protects ` + "`*-sa-key.json`" + ` and ` + "`.env`" + `. Verify with ` + "`git status`" + ` before committing.
- **Review before pushing** — use ` + "`git diff --staged`" + ` to review staged changes before committing.
{{- if eq .Framework "scion"}}
- **Keep ` + "`main`" + ` stable** — each agent works on its own worktree branch and merges back to ` + "`main`" + ` when ready.
{{- else}}
- **Keep ` + "`main`" + ` stable** — use feature branches for non-trivial work, merge back to ` + "`main`" + ` when ready.
{{- end}}
- **Tag milestones** — use annotated git tags for significant releases or milestones.

## Progress Journaling

- **Always update docs/PROGRESS.md** at the end of every session with:
  - Date and session number
  - What was accomplished (with specifics — files changed, features added, bugs fixed)
  - Key decisions made and rationale
  - Next steps / open items
- Commit the docs/PROGRESS.md update as part of the session's final commit

## Team Workflow

This project uses a multi-agent team structure.
{{if eq .TeamSize "lean"}}
### Feedback → Backlog → Execution Pipeline

**This is the mandatory workflow for all user feedback and requests:**

1. **PO ({{.OwnerName}})** provides feedback, feature requests, or bug reports to the **PM**
2. **PM** translates PO feedback into a product spec in ` + "`docs/specs/F-NNNN-slug.md`" + ` with detailed requirements and acceptance criteria
3. **PM** creates/updates work items in docs/BACKLOG.md with priority, scope, and dependencies
4. **PM** assigns work to **SWE-1**
{{- if eq .Framework "scion"}}
5. **SWE-1** implements directly on their worktree branch, following the existing codebase conventions
{{- else}}
5. **SWE-1** implements on a feature branch, following the existing codebase conventions
{{- end}}
6. **SWE-1** hands off completed work to **SWE-Test** for verification
7. **SWE-1** updates docs/BACKLOG.md marking items as completed, tested, and verified
8. **PM** updates docs/PROGRESS.md and docs/RELEASENOTES.md
9. **PM** reports back to the **PO** with a summary of completed work
10. **Tag the release** with ` + "`git tag -a vX.Y.Z`" + ` after PO approval and push tags

**Every piece of feedback goes through this pipeline — no skipping steps.**

**Always use the pipeline for all bug fixes and new features — never ask the PO for confirmation on whether to use the pipeline. Just do it.**

### Agent-Only Execution Rule (Non-Negotiable)

**All project work must be performed by a designated Agent role.** No work is done directly — it is always delegated to the appropriate agent (PM, SWE-1, SWE-Test). If a task requires a role or specialization that does not exist in the current team roster, **stop and check with the PO ({{.OwnerName}})** before proceeding. The PO will decide whether to create a new agent role or reassign the work.

{{if eq .Framework "scion"}}### Agent Orchestration via Scion (Non-Negotiable)

**All agent work MUST use Scion for multi-agent orchestration.**

Agents are defined as templates in ` + "`.scion/templates/`" + ` and run in isolated containers with git worktrees. The correct workflow is:

1. **Launch all agents simultaneously** — Start every agent in the pipeline at once with ` + "`scion start <name> \"task\" --type <template-name>`" + `
2. **Agents signal via ` + "`sciontool status task_completed`" + `** — Agents inside containers cannot use ` + "`scion message`" + ` (Docker is not available). Instead, they signal completion via ` + "`sciontool status task_completed \"message\"`" + `
3. **Orchestrator relays messages** — When an agent signals completion, use ` + "`scion message <next-agent> \"message\"`" + ` from the host to hand off to the next agent in the pipeline
4. **Monitor progress** — Use ` + "`scion list`" + ` to see running agents and ` + "`scion attach <name>`" + ` for interactive sessions
5. **Stop all agents** — When PM signals that the milestone is complete, stop all agents with ` + "`scion stop <name>`" + `

**All agents run in containers with isolated git worktrees.** Each agent works on its own copy of the repo and merges results back when complete.
{{else}}### Interactive Agent Teams via Tmux (Non-Negotiable)

**All agent work MUST use interactive Agent Teams (TeamCreate), NOT subprocess agents.**

Agents must be spawned as interactive teammates in separate tmux panes so the PO can observe and interact with each agent in real time. The correct workflow is:

1. **Create a team** with ` + "`TeamCreate`" + ` (e.g., ` + "`team_name: \"feature-xyz\"`" + `)
2. **Create tasks** with ` + "`TaskCreate`" + ` — one per work item, with clear descriptions
3. **Spawn teammates** using the ` + "`Agent`" + ` tool with ` + "`team_name`" + ` parameter — this launches each agent in its own tmux pane
4. **Assign tasks** via ` + "`TaskUpdate`" + ` with ` + "`owner`" + ` set to the agent name
5. **Coordinate** via ` + "`SendMessage`" + ` — agents report progress and results back to the team lead
6. **Shutdown gracefully** — send ` + "`shutdown_request`" + ` to each agent when work is complete
7. **Clean up** with ` + "`TeamDelete`" + ` after all agents have shut down

**Never use background subprocess agents (Agent tool without ` + "`team_name`" + `).** The PO must always be able to see agent activity in tmux panes. Parallel work should be visible, not hidden.

**All agents MUST use the {{.ModelName}} model (` + "`{{.ModelID}}`" + `).** When spawning teammates with the ` + "`Agent`" + ` tool, always set ` + "`model: \"{{.ModelID}}\"`" + ` unless the PO explicitly specifies a different model. No agent should default to a lesser model.
{{end}}

### Mandatory Development Pipeline (Non-Negotiable)

**All PO feedback and feature requests MUST follow this pipeline — no shortcuts, no exceptions:**

1. **PO → PM**: PO provides feedback, feature requests, or bug reports to the PM Agent
2. **PM → Spec**: PM creates a product spec in ` + "`docs/specs/F-NNNN-slug.md`" + ` (copy from ` + "`docs/specs/TEMPLATE.md`" + `) with detailed requirements and acceptance criteria
3. **PM → Backlog**: PM creates/updates work items in docs/BACKLOG.md with priority, scope, and dependencies (linking to the spec)
{{- if eq .Framework "scion"}}
4. **PM → SWE-1**: PM assigns the work item to SWE-1. SWE-1 implements directly on their worktree branch
{{- else}}
4. **PM → SWE-1**: PM assigns the work item to SWE-1. SWE-1 implements on a feature branch
{{- end}}
5. **SWE-1 → Testing**: After implementation, SWE-1 hands off to SWE-Test (runs all automated tests — existing tests must pass, new tests added for new functionality)
6. **SWE-1 → Backlog Update**: Once coding, functionality, and testing are complete, SWE-1 updates docs/BACKLOG.md marking items as completed, tested, and verified
7. **SWE-1 → PM**: SWE-1 informs PM that the work item is done
8. **PM → Progress + Release Notes**: PM updates docs/PROGRESS.md with session details and docs/RELEASENOTES.md with the new version entry
9. **PM → PO**: PM creates a summary of all completed work and reports back to the PO
10. **Tag**: After PO approval, tag the release with ` + "`git tag -a vX.Y.Z -m \"description\"`" + ` and push tags
11. **Mandatory updates**: docs/BACKLOG.md, docs/PROGRESS.md, and docs/RELEASENOTES.md MUST be updated every milestone. Git tags MUST be created for every release
12. **No direct code changes**: The orchestrator (main Claude context) MUST NEVER write or edit application code directly. Only SWE agents write code. Only PM agents update backlog/progress/release docs

**Violating this pipeline is a process failure.** If time pressure tempts a shortcut, stop and confirm with the PO first.

### Roles
- **PO / CEO** ({{.OwnerName}}) — Product Owner, the human in the loop. Provides feedback, feature requests, and bug reports. Approves direction, tests the app
- **PM** — Receives PO feedback. Translates to specs. Manages backlog, progress log, and release notes. Coordinates with SWE-1 and SWE-Test
- **SWE-1** — {{(index .SWEs 0).Title}}
- **SWE-Test** — Test coverage and quality assurance. Runs all automated tests after SWE implementation. Ensures existing tests pass and new tests are added for new functionality
{{else if eq .TeamSize "full"}}
### Feedback → Backlog → Execution Pipeline

**This is the mandatory workflow for all user feedback and requests:**

1. **PO ({{.OwnerName}})** provides feedback, feature requests, or bug reports to the **PM**
2. **PM** translates PO feedback into a product spec in ` + "`docs/specs/F-NNNN-slug.md`" + ` with detailed requirements and acceptance criteria, then works with the **TPM** to create/update work items in docs/BACKLOG.md with priority, scope, and dependencies
3. **TPM** assigns individual work items to the appropriate **SWE agents** (scaling from 1–{{len .SWEs}} SWEs as needed) and coordinates parallel execution
{{- if eq .Framework "scion"}}
4. **SWE agents** implement directly on their worktree branches, following the existing codebase conventions
{{- else}}
4. **SWE agents** implement on feature branches, following the existing codebase conventions
{{- end}}
5. **SWE agents** hand off completed work to **SWE-Test** (automated tests) and **SWE-QA** (E2E testing) for end-to-end verification
6. Once coding, functionality, and testing are complete, SWEs update docs/BACKLOG.md marking items as completed, tested, and verified, then inform the **TPM**
7. **TPM** updates docs/PROGRESS.md with session details, waits for all work items in the milestone to be completed, then informs the **PM**
8. **PM** updates docs/RELEASENOTES.md with the new version, creates a summary of all completed work, and reports back to the **PO**
9. **Tag the release** with ` + "`git tag -a vX.Y.Z`" + ` after PO approval and push tags

**Every piece of feedback goes through this pipeline — no skipping steps.**

**Always use the pipeline for all bug fixes and new features — never ask the PO for confirmation on whether to use the pipeline. Just do it.**

### Agent-Only Execution Rule (Non-Negotiable)

**All project work must be performed by a designated Agent role.** No work is done directly — it is always delegated to the appropriate agent (PM, TPM, SWE-1 through SWE-{{len .SWEs}}, SWE-Test, SWE-QA, Platform Engineer, Reviewer). If a task requires a role or specialization that does not exist in the current team roster, **stop and check with the PO ({{.OwnerName}})** before proceeding. The PO will decide whether to create a new agent role or reassign the work.

{{if eq .Framework "scion"}}### Agent Orchestration via Scion (Non-Negotiable)

**All agent work MUST use Scion for multi-agent orchestration.**

Agents are defined as templates in ` + "`.scion/templates/`" + ` and run in isolated containers with git worktrees. The correct workflow is:

1. **Launch all agents simultaneously** — Start every agent in the pipeline at once with ` + "`scion start <name> \"task\" --type <template-name>`" + `
2. **Agents signal via ` + "`sciontool status task_completed`" + `** — Agents inside containers cannot use ` + "`scion message`" + ` (Docker is not available). Instead, they signal completion via ` + "`sciontool status task_completed \"message\"`" + `
3. **Orchestrator relays messages** — When an agent signals completion, use ` + "`scion message <next-agent> \"message\"`" + ` from the host to hand off to the next agent in the pipeline
4. **Monitor progress** — Use ` + "`scion list`" + ` to see running agents and ` + "`scion attach <name>`" + ` for interactive sessions
5. **Stop all agents** — When PM signals that the milestone is complete, stop all agents with ` + "`scion stop <name>`" + `

**All agents run in containers with isolated git worktrees.** Each agent works on its own copy of the repo and merges results back when complete.
{{else}}### Interactive Agent Teams via Tmux (Non-Negotiable)

**All agent work MUST use interactive Agent Teams (TeamCreate), NOT subprocess agents.**

Agents must be spawned as interactive teammates in separate tmux panes so the PO can observe and interact with each agent in real time. The correct workflow is:

1. **Create a team** with ` + "`TeamCreate`" + ` (e.g., ` + "`team_name: \"feature-xyz\"`" + `)
2. **Create tasks** with ` + "`TaskCreate`" + ` — one per work item, with clear descriptions
3. **Spawn teammates** using the ` + "`Agent`" + ` tool with ` + "`team_name`" + ` parameter — this launches each agent in its own tmux pane
4. **Assign tasks** via ` + "`TaskUpdate`" + ` with ` + "`owner`" + ` set to the agent name
5. **Coordinate** via ` + "`SendMessage`" + ` — agents report progress and results back to the team lead
6. **Shutdown gracefully** — send ` + "`shutdown_request`" + ` to each agent when work is complete
7. **Clean up** with ` + "`TeamDelete`" + ` after all agents have shut down

**Never use background subprocess agents (Agent tool without ` + "`team_name`" + `).** The PO must always be able to see agent activity in tmux panes. Parallel work should be visible, not hidden.

**All agents MUST use the {{.ModelName}} model (` + "`{{.ModelID}}`" + `).** When spawning teammates with the ` + "`Agent`" + ` tool, always set ` + "`model: \"{{.ModelID}}\"`" + ` unless the PO explicitly specifies a different model. No agent should default to a lesser model.
{{end}}

### Mandatory Development Pipeline (Non-Negotiable)

**All PO feedback and feature requests MUST follow this pipeline — no shortcuts, no exceptions:**

1. **PO → PM**: PO provides feedback, feature requests, or bug reports to the PM Agent
2. **PM → Spec**: PM creates a product spec in ` + "`docs/specs/F-NNNN-slug.md`" + ` (copy from ` + "`docs/specs/TEMPLATE.md`" + `) with detailed requirements and acceptance criteria
3. **PM → TPM**: PM works with TPM to create work items in docs/BACKLOG.md with priority, scope, and dependencies (linking to the spec)
{{- if eq .Framework "scion"}}
4. **TPM → SWE**: TPM assigns individual work items to SWE agents (1–{{len .SWEs}} SWEs, scaled based on workload). Each SWE picks up their assigned item and implements directly on their worktree branch
{{- else}}
4. **TPM → SWE**: TPM assigns individual work items to SWE agents (1–{{len .SWEs}} SWEs, scaled based on workload). Each SWE picks up their assigned item and implements on a feature branch
{{- end}}
5. **SWE → Testing**: After implementation, SWEs hand off to SWE-Test (runs all automated tests — existing tests must pass, new tests added for new functionality) and SWE-QA (E2E testing)
6. **SWE → Backlog Update**: Once coding, functionality, and testing are complete, SWEs update docs/BACKLOG.md marking items as completed, tested, and verified
7. **SWE → TPM**: SWEs inform TPM that their work items are done
8. **TPM → Progress**: TPM updates docs/PROGRESS.md with session details (what was done, decisions, next steps)
9. **TPM → PM**: TPM waits for all work items in the milestone to be completed, then informs PM
10. **PM → Release Notes**: PM updates docs/RELEASENOTES.md with the new version entry (Added, Changed, Fixed sections)
11. **PM → PO**: PM creates a summary of all completed work and reports back to the PO
12. **Tag**: After PO approval, tag the release with ` + "`git tag -a vX.Y.Z -m \"description\"`" + ` and push tags
13. **Mandatory updates**: docs/BACKLOG.md, docs/PROGRESS.md, and docs/RELEASENOTES.md MUST be updated every milestone. Git tags MUST be created for every release
14. **No direct code changes**: The orchestrator (main Claude context) MUST NEVER write or edit application code directly. Only SWE agents write code. Only PM/TPM agents update backlog/progress/release docs

**Violating this pipeline is a process failure.** If time pressure tempts a shortcut, stop and confirm with the PO first.

### Roles
- **PO / CEO** ({{.OwnerName}}) — Product Owner, the human in the loop. Provides feedback, feature requests, and bug reports. Approves direction, tests the app
- **PM** — Receives all PO feedback. Translates it into detailed product requirements with acceptance criteria. Works with TPM to create backlog items. Creates completion summaries and reports back to PO
- **TPM** — Coordinates between agents. Allocates individual work items to SWEs. Tracks blockers and dependencies. Waits for all milestone items to complete before reporting to PM. Maintains docs/PROGRESS.md and docs/BACKLOG.md
{{- range .SWEs}}
- **SWE-{{.Number}}** — {{.Title}}
{{- end}}
- **SWE-Test** — Test coverage and quality assurance. Runs all automated tests after SWE implementation. Ensures existing tests pass and new tests are added for new functionality
- **SWE-QA** — QA and browser testing. Headless Chromium screenshots via puppeteer-core, visual verification, Lighthouse audits, E2E smoke tests. Validates end-to-end functionality
- **Platform Engineer (PE)** — GCP expert (DevOps + SRE). Owns all infrastructure: Cloud Run deployment, Dockerfile, IAM/service accounts, monitoring, billing, free tier quota tracking, troubleshooting via GCP logs, reliability engineering
- **Reviewer** — Code review, quality/security/performance checks
{{else}}
### Feedback → Backlog → Execution Pipeline

**This is the mandatory workflow for all user feedback and requests:**

1. **PO ({{.OwnerName}})** provides feedback, feature requests, or bug reports to the **PM**
2. **PM** translates PO feedback into a product spec in ` + "`docs/specs/F-NNNN-slug.md`" + ` with detailed requirements and acceptance criteria, then works with the **TPM** to create/update work items in docs/BACKLOG.md with priority, scope, and dependencies
3. **TPM** assigns individual work items to the appropriate **SWE agents** (scaling from 1–{{len .SWEs}} SWEs as needed) and coordinates parallel execution
{{- if eq .Framework "scion"}}
4. **SWE agents** implement directly on their worktree branches, following the existing codebase conventions
{{- else}}
4. **SWE agents** implement on feature branches, following the existing codebase conventions
{{- end}}
5. **SWE agents** hand off completed work to {{- if .IncludeSWETest}} **SWE-Test** (automated tests){{end}}{{if and .IncludeSWETest .IncludeSWEQA}} and{{end}}{{if .IncludeSWEQA}} **SWE-QA** (E2E testing){{end}}{{if not (or .IncludeSWETest .IncludeSWEQA)}} testing{{end}} for end-to-end verification
6. Once coding, functionality, and testing are complete, SWEs update docs/BACKLOG.md marking items as completed, tested, and verified, then inform the **TPM**
7. **TPM** updates docs/PROGRESS.md with session details, waits for all work items in the milestone to be completed, then informs the **PM**
8. **PM** updates docs/RELEASENOTES.md with the new version, creates a summary of all completed work, and reports back to the **PO**
9. **Tag the release** with ` + "`git tag -a vX.Y.Z`" + ` after PO approval and push tags

**Every piece of feedback goes through this pipeline — no skipping steps.**

**Always use the pipeline for all bug fixes and new features — never ask the PO for confirmation on whether to use the pipeline. Just do it.**

### Agent-Only Execution Rule (Non-Negotiable)

**All project work must be performed by a designated Agent role.** No work is done directly — it is always delegated to the appropriate agent (PM, TPM, SWE-1 through SWE-{{len .SWEs}}{{if .IncludeSWETest}}, SWE-Test{{end}}{{if .IncludeSWEQA}}, SWE-QA{{end}}{{if .IncludePlatform}}, Platform Engineer{{end}}{{if .IncludeReviewer}}, Reviewer{{end}}). If a task requires a role or specialization that does not exist in the current team roster, **stop and check with the PO ({{.OwnerName}})** before proceeding. The PO will decide whether to create a new agent role or reassign the work.

{{if eq .Framework "scion"}}### Agent Orchestration via Scion (Non-Negotiable)

**All agent work MUST use Scion for multi-agent orchestration.**

Agents are defined as templates in ` + "`.scion/templates/`" + ` and run in isolated containers with git worktrees. The correct workflow is:

1. **Launch all agents simultaneously** — Start every agent in the pipeline at once with ` + "`scion start <name> \"task\" --type <template-name>`" + `
2. **Agents signal via ` + "`sciontool status task_completed`" + `** — Agents inside containers cannot use ` + "`scion message`" + ` (Docker is not available). Instead, they signal completion via ` + "`sciontool status task_completed \"message\"`" + `
3. **Orchestrator relays messages** — When an agent signals completion, use ` + "`scion message <next-agent> \"message\"`" + ` from the host to hand off to the next agent in the pipeline
4. **Monitor progress** — Use ` + "`scion list`" + ` to see running agents and ` + "`scion attach <name>`" + ` for interactive sessions
5. **Stop all agents** — When PM signals that the milestone is complete, stop all agents with ` + "`scion stop <name>`" + `

**All agents run in containers with isolated git worktrees.** Each agent works on its own copy of the repo and merges results back when complete.
{{else}}### Interactive Agent Teams via Tmux (Non-Negotiable)

**All agent work MUST use interactive Agent Teams (TeamCreate), NOT subprocess agents.**

Agents must be spawned as interactive teammates in separate tmux panes so the PO can observe and interact with each agent in real time. The correct workflow is:

1. **Create a team** with ` + "`TeamCreate`" + ` (e.g., ` + "`team_name: \"feature-xyz\"`" + `)
2. **Create tasks** with ` + "`TaskCreate`" + ` — one per work item, with clear descriptions
3. **Spawn teammates** using the ` + "`Agent`" + ` tool with ` + "`team_name`" + ` parameter — this launches each agent in its own tmux pane
4. **Assign tasks** via ` + "`TaskUpdate`" + ` with ` + "`owner`" + ` set to the agent name
5. **Coordinate** via ` + "`SendMessage`" + ` — agents report progress and results back to the team lead
6. **Shutdown gracefully** — send ` + "`shutdown_request`" + ` to each agent when work is complete
7. **Clean up** with ` + "`TeamDelete`" + ` after all agents have shut down

**Never use background subprocess agents (Agent tool without ` + "`team_name`" + `).** The PO must always be able to see agent activity in tmux panes. Parallel work should be visible, not hidden.

**All agents MUST use the {{.ModelName}} model (` + "`{{.ModelID}}`" + `).** When spawning teammates with the ` + "`Agent`" + ` tool, always set ` + "`model: \"{{.ModelID}}\"`" + ` unless the PO explicitly specifies a different model. No agent should default to a lesser model.
{{end}}

### Mandatory Development Pipeline (Non-Negotiable)

**All PO feedback and feature requests MUST follow this pipeline — no shortcuts, no exceptions:**

1. **PO → PM**: PO provides feedback, feature requests, or bug reports to the PM Agent
2. **PM → Spec**: PM creates a product spec in ` + "`docs/specs/F-NNNN-slug.md`" + ` (copy from ` + "`docs/specs/TEMPLATE.md`" + `) with detailed requirements and acceptance criteria
3. **PM → TPM**: PM works with TPM to create work items in docs/BACKLOG.md with priority, scope, and dependencies (linking to the spec)
{{- if eq .Framework "scion"}}
4. **TPM → SWE**: TPM assigns individual work items to SWE agents (1–{{len .SWEs}} SWEs, scaled based on workload). Each SWE picks up their assigned item and implements directly on their worktree branch
{{- else}}
4. **TPM → SWE**: TPM assigns individual work items to SWE agents (1–{{len .SWEs}} SWEs, scaled based on workload). Each SWE picks up their assigned item and implements on a feature branch
{{- end}}
5. **SWE → Testing**: After implementation, SWEs hand off to {{- if .IncludeSWETest}} SWE-Test (runs all automated tests — existing tests must pass, new tests added for new functionality){{end}}{{if and .IncludeSWETest .IncludeSWEQA}} and{{end}}{{if .IncludeSWEQA}} SWE-QA (E2E testing){{end}}{{if not (or .IncludeSWETest .IncludeSWEQA)}} testing{{end}}
6. **SWE → Backlog Update**: Once coding, functionality, and testing are complete, SWEs update docs/BACKLOG.md marking items as completed, tested, and verified
7. **SWE → TPM**: SWEs inform TPM that their work items are done
8. **TPM → Progress**: TPM updates docs/PROGRESS.md with session details (what was done, decisions, next steps)
9. **TPM → PM**: TPM waits for all work items in the milestone to be completed, then informs PM
10. **PM → Release Notes**: PM updates docs/RELEASENOTES.md with the new version entry (Added, Changed, Fixed sections)
11. **PM → PO**: PM creates a summary of all completed work and reports back to the PO
12. **Tag**: After PO approval, tag the release with ` + "`git tag -a vX.Y.Z -m \"description\"`" + ` and push tags
13. **Mandatory updates**: docs/BACKLOG.md, docs/PROGRESS.md, and docs/RELEASENOTES.md MUST be updated every milestone. Git tags MUST be created for every release
14. **No direct code changes**: The orchestrator (main Claude context) MUST NEVER write or edit application code directly. Only SWE agents write code. Only PM/TPM agents update backlog/progress/release docs

**Violating this pipeline is a process failure.** If time pressure tempts a shortcut, stop and confirm with the PO first.

### Roles
- **PO / CEO** ({{.OwnerName}}) — Product Owner, the human in the loop. Provides feedback, feature requests, and bug reports. Approves direction, tests the app
- **PM** — Receives all PO feedback. Translates it into detailed product requirements with acceptance criteria. Works with TPM to create backlog items. Creates completion summaries and reports back to PO
- **TPM** — Coordinates between agents. Allocates individual work items to SWEs. Tracks blockers and dependencies. Waits for all milestone items to complete before reporting to PM. Maintains docs/PROGRESS.md and docs/BACKLOG.md
{{- range .SWEs}}
- **SWE-{{.Number}}** — {{.Title}}
{{- end}}
{{- if .IncludeSWETest}}
- **SWE-Test** — Test coverage and quality assurance. Runs all automated tests after SWE implementation. Ensures existing tests pass and new tests are added for new functionality
{{- end}}
{{- if .IncludeSWEQA}}
- **SWE-QA** — QA and browser testing. Headless Chromium screenshots via puppeteer-core, visual verification, Lighthouse audits, E2E smoke tests. Validates end-to-end functionality
{{- end}}
{{- if .IncludePlatform}}
- **Platform Engineer (PE)** — GCP expert (DevOps + SRE). Owns all infrastructure: Cloud Run deployment, Dockerfile, IAM/service accounts, monitoring, billing, free tier quota tracking, troubleshooting via GCP logs, reliability engineering
{{- end}}
{{- if .IncludeReviewer}}
- **Reviewer** — Code review, quality/security/performance checks
{{- end}}
{{end}}
### Backlog Tracking (Non-Negotiable)

**Every piece of work gets a backlog entry in ` + "`docs/BACKLOG.md`" + ` — no exceptions.**
Regardless of team size, all features, bug fixes, and enhancements must be tracked
in the backlog before implementation begins and updated when completed.

### Other Conventions
{{- if eq .Framework "scion"}}
- **Branching:** Each Scion agent commits directly to its worktree branch (e.g., ` + "`pm`" + `, ` + "`swe-1`" + `, ` + "`swe-test`" + `). Do NOT create separate feature branches — worktrees already provide isolation
{{- else}}
- **Branching:** Feature branches (` + "`feature/<name>`" + `) off ` + "`main`" + `
{{- end}}
{{- if .IncludePlatform}}
- **Platform Engineer (PE) owns all GCP interactions:** Cloud Run deployment, Dockerfile, IAM/service accounts, logging, monitoring, billing, free tier quota tracking, reliability, troubleshooting
{{- end}}
- **Co-Author:** All commits include ` + "`Co-Authored-By: Claude {{.ModelName}} <noreply@anthropic.com>`" + `
{{- if .GCP.Enabled}}

## GCP Free Tier (Non-Negotiable)

- **This app must stay within the GCP free tier. Zero additional billing.**
- Single user app ({{.OwnerName}} only) — no need for high availability or scale
- Cloud Run config: **256Mi memory, 0.5 vCPU, maxScale=1, minInstances=0, request-based CPU**
- Region: **{{.GCP.Region}}** (matches existing services)
- Clean up old Artifact Registry images to stay within 0.5 GB free storage
{{- end}}
{{- if .Conventions}}

## Important Conventions
{{range .Conventions}}
- {{.}}
{{- end}}
{{- end}}
`
