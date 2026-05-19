# Roadmap

## Completed Milestones

| Version | Date | Highlights |
|---------|------|------------|
| **v0.1.0** | 2026-03-16 | Interactive 6-step wizard, CLAUDE.md + agent generation, GCP support |
| **v0.2.0** | 2026-03-16 | ANSI color styling, TTY detection, backlog/progress/release note templates |
| **v0.3.0** | 2026-03-16 | MermaidJS pipeline diagram, git repo setup, model selection |
| **v0.4.0** | 2026-03-16 | `docs/` directory structure, PM spec system, pipeline enhancements |
| **v0.4.1** | 2026-03-16 | `--help`/`--version` flags, color-coded pipeline links |
| **v0.4.2** | 2026-03-16 | Progress, release notes, and tag nodes in pipeline diagram |
| **v0.5.0** | 2026-03-16 | Settings persistence (`.appteam/settings.json`), `-r` flag, 5 bootstrap skills |
| **v0.5.1** | 2026-03-17 | User-facing documentation (`docs/guide/`) |
| **v0.6.0** | 2026-03-17 | Comprehensive test suite (45 tests across 5 packages) |
| **v0.7.0** | 2026-03-17 | `/roadmap` skill, selectable skills in wizard |
| **v0.8.0** | 2026-03-17 | 12 universal dev skills, custom skills, grouped wizard UX |
| **v0.9.0** | 2026-03-17 | Custom agent definitions |
| **v0.10.0** | 2026-03-17 | `-d`/`--dir` target directory flag |
| **v0.11.0** | 2026-03-17 | `/brainstorm` skill for PM product ideation |
| **v0.12.0** | 2026-03-17 | Team sizing (lean/standard/full), backlog enforcement |
| **v0.13.0** | 2026-03-18 | Enhanced SWE-QA agent, `/cuj-list` + `/cuj-test` skills |
| **v0.13.1** | 2026-03-17 | CUJ test cleanup phase |

## Upcoming

| ID | Feature | Description |
|----|---------|-------------|
| **F-0081** | Agent retrospective skill | `/retrospective` — agents review their own performance after milestones, capture lessons learned, feed improvements back into agent definitions |
| **F-0082** | Multi-project shared conventions | Global profile (`~/.appteam/global.json`) merged into every project's settings — shared git config, preferred patterns, common skills across projects |

## Vision

appteam is heading toward **self-improving agent teams**. The long-term direction:

- **Retrospectives** — agents learn from each milestone and update their own definitions
- **Shared conventions** — common patterns and preferences carry across all your projects
- **Ecosystem growth** — community-contributed skills, agent templates, and team presets
- **Smarter defaults** — the wizard learns from your past projects to suggest better starting configurations

The goal is simple: make AI agent teams as easy to spin up as `npm init` — but with the structure and discipline that large projects demand.
