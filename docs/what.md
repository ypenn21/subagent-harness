# What It Does

appteam generates a complete multi-agent team configuration for Claude Code projects. Here's what you get.

## Agent Roles

Every team has a set of specialized AI agents, each with a defined role:

| Agent | Role |
|-------|------|
| **PM** | Product Manager — specs, release notes, PO communication |
| **TPM** | Technical Program Manager — backlog, progress tracking, coordination |
| **SWE-1 ... SWE-5** | Software Engineers — implementation on feature branches |
| **SWE-Test** | Test Engineer — automated test verification |
| **SWE-QA** | QA Engineer — browser testing, CUJ validation, Lighthouse audits |
| **Reviewer** | Code Reviewer — quality, security, performance checks |
| **Platform** | Platform Engineer — infrastructure, deployment, CI/CD |
| **Custom** | Define your own roles — DBA, UX Designer, Tech Writer, anything |

## Team Sizing

Choose the team size that fits your project:

| Size | Agents | Best For |
|------|--------|----------|
| **Lean** | PM + 1 SWE + SWE-Test | Solo projects, small features, quick prototypes |
| **Standard** | PM + TPM + 2 SWEs + SWE-Test + Reviewer | Most projects (default) |
| **Full** | PM + TPM + 5 SWEs + SWE-Test + SWE-QA + Reviewer + Platform | Large projects, complex infrastructure |

## 21 Skills

Skills are Claude Code slash commands generated as `.claude/skills/<name>.md`:

### PM Skills (7)

| Skill | Description |
|-------|-------------|
| `/spec` | Create a product spec with requirements and acceptance criteria |
| `/release` | Update release notes, commit, tag, and push |
| `/roadmap` | Add items to the backlog roadmap |
| `/status` | Summarize current milestone status |
| `/pipeline` | Spin up the full agent team |
| `/regenerate` | Regenerate all files from saved settings |
| `/brainstorm` | Structured product ideation session with PO |

### Dev Tier 1 (6, included by default)

| Skill | Description |
|-------|-------------|
| `/debug` | Systematic debugging: reproduce, isolate, fix, verify |
| `/test` | Write tests for a module or function |
| `/review` | Code review checklist (OWASP, performance, correctness) |
| `/docs` | Generate or update documentation |
| `/refactor` | Safe refactoring with test verification |
| `/hotfix` | Emergency fix: branch, fix, test, merge, tag |

### Dev Tier 2 (8, opt-in)

| Skill | Description |
|-------|-------------|
| `/api-design` | Design REST/GraphQL endpoints |
| `/schema` | Database schema and migration design |
| `/deploy` | Deployment checklist with rollback plan |
| `/security` | Security audit (injection, XSS, auth, dependencies) |
| `/adr` | Architecture Decision Record |
| `/standup` | Daily standup summary from git log and backlog |
| `/cuj-list` | Create and manage Critical User Journey inventory |
| `/cuj-test` | Run CUJ tests via headless Chromium |

Plus **custom skills** — define your own during the wizard.

## Generated Files

appteam generates 33+ files across four directories:

| Directory | Contents |
|-----------|----------|
| Root | `CLAUDE.md` — project instructions with team workflow, pipeline rules, conventions |
| `.claude/agents/` | Agent role definitions (PM, TPM, SWEs, optional agents, custom agents) |
| `.claude/skills/` | Slash command skills (21 predefined + custom) |
| `docs/` | BACKLOG.md, PROGRESS.md, RELEASENOTES.md, PIPELINE.md, specs/TEMPLATE.md |
| `.appteam/` | settings.json — saved wizard config for regeneration |

## Mandatory Pipeline

Every piece of work follows a strict pipeline — no shortcuts:

```
PO --> PM --> TPM --> SWEs --> Test --> Review --> TPM --> PM --> PO
```

1. **PO** provides feedback to the **PM**
2. **PM** creates a product spec with requirements and acceptance criteria
3. **PM + TPM** create backlog items with priority and dependencies
4. **TPM** assigns work to **SWE** agents
5. **SWEs** implement on feature branches
6. **SWE-Test** verifies all tests pass
7. **Reviewer** checks quality, security, performance
8. **TPM** updates backlog and progress log
9. **PM** updates release notes and reports back to **PO**

## Backlog Enforcement

Every piece of work gets a `docs/BACKLOG.md` entry — no exceptions. This is baked into the generated `CLAUDE.md` as a non-negotiable rule. Items are tracked with:

- Sequential IDs (`F-0001`, `F-0002`, ...)
- Priority, status, owner, dependencies
- Links to product specs
- Version milestones
