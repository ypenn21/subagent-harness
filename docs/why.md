# Why appteam

## The Problem

AI coding agents are powerful, but managing a team of them is chaotic without structure. Every project needs agent definitions, workflow rules, skill commands, and coordination protocols. Setting this up by hand is slow, repetitive, and inconsistent across projects.

## The Pain Points

**No standard team structure.** Each project invents its own agent roles from scratch. One project has a "reviewer" agent, another calls it "code-checker," a third skips review entirely. There's no shared vocabulary or consistent role definition.

**Pipeline rules get lost or ignored.** Without explicit workflow rules baked into the project, agents take shortcuts. Features ship without specs. Backlogs go stale. Release notes are an afterthought. Progress goes untracked.

**Skills need to be recreated for every project.** Slash commands like `/spec`, `/release`, `/debug`, and `/test` are useful everywhere, but there's no way to carry them across projects. Each new project starts from zero.

**Agent roles are ad-hoc.** Without clear role definitions, agents overlap, miss responsibilities, or work at cross-purposes. Who owns the backlog? Who writes specs? Who runs tests? Without explicit assignments, these questions get answered inconsistently.

**Backlog tracking is optional when it should be mandatory.** Work happens but isn't tracked. Items get completed without updating the backlog. New team members (human or AI) have no way to see what's been done, what's in progress, or what's next.

## The Solution

appteam is an opinionated generator that creates a complete, consistent team structure in seconds. Run the wizard once, and you get:

- **Defined agent roles** with clear responsibilities and boundaries
- **A mandatory pipeline** that every piece of work must follow — no shortcuts
- **21 ready-to-use skills** for common workflows (specs, releases, debugging, testing, security audits)
- **Backlog enforcement** baked into the project instructions — every piece of work gets tracked
- **Consistent structure** across all your projects — same roles, same pipeline, same conventions

The result is a team that works the same way every time, regardless of the project. New agents know their role. Work is tracked. Pipeline rules are enforced. And you can regenerate everything from saved settings whenever you need to update.
