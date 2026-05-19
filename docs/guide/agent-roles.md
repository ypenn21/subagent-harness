# Agent Roles

appteam generates a team of specialized AI agents, each with a defined role and set of responsibilities. Agents are spawned as interactive teammates in separate tmux panes using Claude Code's `TeamCreate` and `Agent` tools.

## Core Agents (always generated)

### PM — Product Manager

**File:** `.claude/agents/pm.md`

The PM is the bridge between the Product Owner (the human) and the engineering team.

**Responsibilities:**
- Receives feedback, feature requests, and bug reports from the PO
- Creates product specs in `docs/specs/F-NNNN-slug.md` with detailed requirements and acceptance criteria
- Works with the TPM to create backlog items
- Updates `docs/RELEASENOTES.md` when milestones complete
- Reports completed work summaries back to the PO

**Key files owned:** `docs/specs/*.md`, `docs/RELEASENOTES.md`

### TPM — Technical Program Manager

**File:** `.claude/agents/tpm.md`

The TPM coordinates the engineering team and tracks overall progress.

**Responsibilities:**
- Assigns work items from the backlog to SWE agents
- Tracks blockers, dependencies, and parallel work streams
- Waits for all milestone items to complete before reporting to PM
- Updates `docs/BACKLOG.md` with status changes
- Updates `docs/PROGRESS.md` with session details at the end of every work session

**Key files owned:** `docs/BACKLOG.md`, `docs/PROGRESS.md`

## SWE Agents (1–5, configurable)

**Files:** `.claude/agents/swe-1.md` through `.claude/agents/swe-5.md`

Software Engineers implement features, fix bugs, and write code. Each SWE has a configurable title and specialty bullets that focus their expertise.

**Responsibilities:**
- Pick up assigned work items from the backlog
- Implement on feature branches
- Read specs from `docs/specs/` for acceptance criteria
- Hand off completed work to testing
- Update `docs/BACKLOG.md` to mark items as done

**Example specializations:**
- SWE-1: "Backend & API" — REST endpoints, database queries
- SWE-2: "Frontend & UI" — React components, CSS
- SWE-3: "Infrastructure" — Docker, CI/CD, deployment scripts

## Optional Agents

### SWE-Test — Test Engineer

**File:** `.claude/agents/swe-test.md`

Runs automated tests after SWE implementation. Ensures existing tests still pass and new tests are added for new functionality.

### SWE-QA — QA Engineer

**File:** `.claude/agents/swe-qa.md`

Performs end-to-end testing, browser QA with headless Chromium, Lighthouse audits, and visual verification. Validates that features work correctly from the user's perspective.

### Platform Engineer

**File:** `.claude/agents/platform.md`

Owns all infrastructure: GCP deployment, Docker, IAM/service accounts, monitoring, billing, CI/CD pipelines. Best suited for projects with cloud infrastructure requirements.

### Reviewer — Code Reviewer

**File:** `.claude/agents/reviewer.md`

Reviews code for quality, security, and performance. Checks that implementations match spec requirements and follow project conventions. Approves or requests changes before merging.

## How Agents Interact

Agents don't directly call each other. Instead, they communicate through the orchestrator (the main Claude Code context) using `SendMessage`. The typical flow is:

```
PO (human) → PM → TPM → SWE-1, SWE-2, ... → SWE-Test → Reviewer → TPM → PM → PO
```

See [Pipeline Workflow](pipeline.md) for the full mandatory workflow.

## Model Selection

All agents use the same AI model, configured during the wizard. The available models are:

| Model | ID | Best for |
|-------|-----|----------|
| Opus 4.6 | `claude-opus-4-6` | Complex reasoning, architecture decisions, thorough code review |
| Sonnet 4.6 | `claude-sonnet-4-6` | Balanced speed and quality for most tasks |
| Haiku 4.5 | `claude-haiku-4-5-20251001` | Fast, lightweight tasks where speed matters more than depth |

The model is set during the wizard and stored in `.appteam/settings.json`. It's embedded in `CLAUDE.md` so the orchestrator enforces it when spawning agents.
