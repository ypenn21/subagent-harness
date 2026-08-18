# appteam

CLI tool that generates Claude Code agent team configurations via an interactive wizard. Zero external dependencies — stdlib only.

## What It Does

`appteam` scaffolds a complete multi-agent team setup for Claude Code projects. Run the wizard, answer a few questions, and it generates all the files you need: project instructions (`CLAUDE.md`), [...]

## Install

### Prerequisite: Go

appteam is written in Go and requires a Go toolchain to build or install. Minimum supported version: Go 1.20. We recommend using the latest stable Go release when possible.

Install Go using your platform's package manager or the official binaries:

macOS (Homebrew):

```bash
brew install go
```

Ubuntu/Debian (apt):

```bash
sudo apt update
sudo apt install -y golang-go
```

Fedora/CentOS (dnf/yum):

```bash
sudo dnf install golang
# or
sudo yum install golang
```

Windows (winget):

```powershell
winget install GoLang.Go
```

Or download and install the official binaries from the Go website: https://golang.org/dl

After installing, verify the toolchain and version:

```bash
go version
# example output: go version go1.20.7 darwin/amd64
```


```bash
# Global install
go install github.com/ahafin/appteam@latest

# Or build locally
git clone https://github.com/ahafin/appteam.git
cd appteam
go build -o appteam .
```

## Quick Start

```bash
cd /path/to/your/project
appteam
```

The interactive wizard walks through 7 steps: project basics, git setup, product owner, GCP config (optional), agent team composition, conventions, and confirmation. After confirmation, all files [...]

```bash
# Regenerate without the wizard
appteam -r

# Show help
appteam -h
```

## Generated Files

| Directory | Files | Description |
|-----------|-------|-------------|
| Root | `CLAUDE.md` | Project instructions with team workflow and pipeline rules |
| `docs/` | `BACKLOG.md`, `PROGRESS.md`, `RELEASENOTES.md`, `PIPELINE.md`, `specs/TEMPLATE.md` | Project management context |
| `.claude/agents/` | `pm.md`, `tpm.md`, `swe-N.md`, + optional agents | Agent role definitions |
| `.claude/skills/` | `spec.md`, `release.md`, `pipeline.md`, `status.md`, `regenerate.md` | Claude Code slash commands |
| `.appteam/` | `settings.json` | Saved wizard config for regeneration |

## Documentation

Full documentation is in [`docs/guide/`](docs/guide/):

- [Getting Started](docs/guide/getting-started.md) — install and first run
- [CLI Reference](docs/guide/cli-reference.md) — flags and environment variables
- [Generated Files](docs/guide/generated-files.md) — every file appteam creates
- [Agent Roles](docs/guide/agent-roles.md) — what each agent does
- [Pipeline Workflow](docs/guide/pipeline.md) — the mandatory development pipeline
- [Skills](docs/guide/skills.md) — the 5 generated slash commands
- [Configuration](docs/guide/configuration.md) — saved settings, regeneration

## Tech Stack

- **Go** (stdlib only — zero external dependencies)
- `text/template` for markdown rendering
- `bufio` for interactive prompts
- `encoding/json` for settings persistence
- ANSI escape codes for terminal styling (auto-disabled when not a TTY, respects `NO_COLOR`)
