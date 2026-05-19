# Pipeline Workflow

Every piece of work in an appteam project follows a mandatory pipeline. No shortcuts, no exceptions.

## The Pipeline

```
PO → PM → TPM → SWEs → Test/QA → Review → TPM → PM → PO
```

### Step-by-step

1. **PO → PM** — The Product Owner (human) provides feedback, feature requests, or bug reports to the PM agent
2. **PM → Spec** — PM creates a product spec in `docs/specs/F-NNNN-slug.md` using the template. The spec includes detailed requirements and acceptance criteria
3. **PM → TPM** — PM works with TPM to create work items in `docs/BACKLOG.md` with priority, scope, dependencies, and a link to the spec
4. **TPM → SWEs** — TPM assigns individual work items to SWE agents (1–5 SWEs, scaled to workload). Each SWE picks up their assigned item and implements on a feature branch
5. **SWEs → Testing** — After implementation, SWEs hand off to SWE-Test (automated tests) and/or SWE-QA (E2E testing)
6. **SWEs → Backlog** — SWEs update `docs/BACKLOG.md` marking items as completed
7. **SWEs → TPM** — SWEs inform TPM that their work is done
8. **TPM → Progress** — TPM updates `docs/PROGRESS.md` with session details
9. **TPM → PM** — TPM waits for all items in the milestone to complete, then informs PM
10. **PM → Release Notes** — PM updates `docs/RELEASENOTES.md` with the new version entry
11. **PM → PO** — PM creates a summary and reports back to the PO
12. **Tag** — After PO approval, tag the release with `git tag -a vX.Y.Z` and push

### Mandatory updates

Every milestone must include updates to:
- `docs/BACKLOG.md` — status changes for all work items
- `docs/PROGRESS.md` — session log with what was done, decisions, next steps
- `docs/RELEASENOTES.md` — version entry with Added/Changed/Fixed sections
- Git tag — annotated tag for the release

## Product Specs

Before any work enters the backlog, the PM creates a spec file:

```
docs/specs/F-0042-dark-mode.md
```

Specs are copied from `docs/specs/TEMPLATE.md` and include:
- Problem statement
- Numbered requirements (specific, testable)
- Acceptance criteria (checkboxes)
- Out of scope
- Dependencies

The spec ID matches the backlog ID (`F-0042`), and the backlog links back to the spec.

## Backlog Structure

`docs/BACKLOG.md` organizes work by milestone:

```markdown
## v0.6.0

| ID | Feature | Priority | Status | Owner | Dependencies | Spec | Notes |
|----|---------|----------|--------|-------|--------------|------|-------|
| F-0043 | Dark mode | High | IN PROGRESS | SWE-1 | — | [spec](docs/specs/F-0043-dark-mode.md) | ... |
| F-0044 | Theme picker | Medium | TODO | — | F-0043 | [spec](docs/specs/F-0044-theme-picker.md) | ... |
```

IDs are sequential across all milestones (`F-0001`, `F-0002`, ...) and never reused.

## Agent Teams via tmux

All agent work uses interactive Agent Teams, not background subprocesses. This means:

1. **`TeamCreate`** — creates a named team (e.g., `feature-dark-mode`)
2. **`TaskCreate`** — creates work items with clear descriptions
3. **`Agent`** with `team_name` — spawns each agent in its own tmux pane so the PO can observe
4. **`TaskUpdate`** — assigns tasks to agents
5. **`SendMessage`** — agents report progress and results
6. **`shutdown_request`** — graceful shutdown when work completes
7. **`TeamDelete`** — cleanup after all agents shut down

The PO can always see agent activity in tmux panes. Parallel work is visible, not hidden.

## Pipeline Diagram

appteam generates a MermaidJS flowchart in `docs/PIPELINE.md` that visualizes the full pipeline. The diagram is dynamic — it renders based on your agent configuration (SWE count, optional agents).

Links are color-coded:
- **Green** — downward flow (request → execution)
- **Blue** — upward flow (results → reporting)
- **Gray** — side effects (backlog updates, progress logging)

## Why the pipeline is mandatory

The pipeline exists to ensure:
- **Every feature has a spec** with clear requirements before coding starts
- **Work is tracked** in the backlog with status, owner, and dependencies
- **Progress is logged** so anyone can understand what happened and why
- **Release notes are maintained** so the PO and users know what changed
- **Releases are tagged** so any version can be checked out and reviewed

Skipping steps leads to undocumented work, lost context, and untraceable changes. The pipeline prevents that.
