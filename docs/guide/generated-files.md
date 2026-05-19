# Generated Files

appteam generates up to 17 files across 4 directories. The exact count depends on how many SWE agents and optional agents you configure.

## Project Root

| File | Description |
|------|-------------|
| `CLAUDE.md` | Project instructions — the main file Claude Code reads for context. Contains project overview, team workflow, pipeline rules, role definitions, version control conventions, and progress journaling instructions |

## `docs/`

| File | Description |
|------|-------------|
| `BACKLOG.md` | Feature backlog with milestone tables, priority, status, owner, dependencies, and spec links. Owned by the TPM |
| `PROGRESS.md` | Session-by-session development log. Updated at the end of every work session by the TPM |
| `RELEASENOTES.md` | Version history in Keep a Changelog format (Added, Changed, Fixed). Owned by the PM |
| `PIPELINE.md` | MermaidJS flowchart showing the full agent workflow. Renders dynamically based on your agent configuration |
| `specs/TEMPLATE.md` | Product spec template. The PM copies this for each new feature (`F-NNNN-slug.md`) |

## `.claude/agents/`

Agent definition files tell Claude Code how each agent should behave when spawned as a teammate.

### Always generated

| File | Description |
|------|-------------|
| `pm.md` | **Product Manager** — translates PO feedback into specs, manages requirements, updates release notes |
| `tpm.md` | **Technical Program Manager** — coordinates agents, assigns work, tracks progress, maintains backlog |

### Per-SWE (1–5 agents)

| File | Description |
|------|-------------|
| `swe-1.md` | Software Engineer 1 with your configured title and specialty |
| `swe-2.md` | Software Engineer 2 (if configured) |
| `swe-3.md` | Software Engineer 3 (if configured) |
| `swe-4.md` | Software Engineer 4 (if configured) |
| `swe-5.md` | Software Engineer 5 (if configured) |

### Optional agents

| File | Included when | Description |
|------|---------------|-------------|
| `swe-test.md` | SWE-Test enabled | **Test Engineer** — runs automated tests, ensures coverage, verifies acceptance criteria |
| `swe-qa.md` | SWE-QA enabled | **QA Engineer** — E2E testing, browser QA, Lighthouse audits, visual verification |
| `platform.md` | Platform enabled | **Platform Engineer** — infrastructure, deployment, GCP, Docker, CI/CD |
| `reviewer.md` | Reviewer enabled | **Code Reviewer** — code review, quality checks, security review |

## `.claude/skills/`

Skill files define slash commands that Claude Code can execute. All 5 skills are always generated.

| File | Slash command | Description |
|------|--------------|-------------|
| `spec.md` | `/spec` | Create a new product spec in `docs/specs/` |
| `release.md` | `/release` | Update release notes, commit, tag, and push |
| `pipeline.md` | `/pipeline` | Spin up the full agent team for a feature |
| `status.md` | `/status` | Summarize current milestone status from the backlog |
| `regenerate.md` | `/regenerate` | Regenerate all files from saved settings |

See [Skills](skills.md) for detailed usage of each slash command.

## `.appteam/`

| File | Description |
|------|-------------|
| `settings.json` | Saved wizard configuration. Used by `-r` flag and auto-detected on startup |

See [Configuration](configuration.md) for the settings schema.

## File Counts

| Configuration | Files generated |
|---------------|----------------|
| 1 SWE, no optional agents | 12 |
| 2 SWEs + Reviewer + SWE-Test | 15 |
| 3 SWEs + all optional agents | 18 |
| 5 SWEs + all optional agents | 20 |

The count is: 1 (CLAUDE.md) + 5 (docs) + 2 (PM/TPM) + N (SWEs) + optional agents + 5 (skills) + 1 (settings.json).
