# How It Works

## Install

```bash
go install github.com/ahafin/appteam@latest
```

Or build from source:

```bash
git clone https://github.com/ahafin/appteam.git
cd appteam
go build -o appteam .
```

## Run the Wizard

```bash
cd /path/to/your/project
appteam
```

The interactive wizard walks through 8 steps:

### Step 1 — Project Basics

Project name, description, tech stack, and target directory.

```
▸ Project name: myapp
▸ Short description: A REST API for managing widgets
▸ Tech stack: Go, PostgreSQL, Chi router
▸ Target directory [.]: .
```

### Step 2 — Git Repository

Detects existing `.git` directory. Offers `git init` and optional GitHub repo creation via `gh`.

### Step 3 — Product Owner

Name, email, and GitHub username for the human in the loop.

### Step 4 — GCP Configuration

Optional. Provide GCP project ID, number, org, and region — or skip entirely.

### Step 5 — Agent Team

Choose team size (lean/standard/full), number of SWEs, specialties for each, and optional agents.

```
▸ Team size (lean/standard/full) [standard]: standard
▸ Number of SWE agents (1-5) [3]: 2

SWE-1
▸ Title/specialty: Backend & API

SWE-2
▸ Title/specialty: Frontend & UI

▸ Include Code Reviewer? (Y/n): y
▸ Include SWE-Test? (Y/n): y
```

### Step 6 — Conventions

Free-text project conventions that get embedded in `CLAUDE.md`.

### Step 7 — Skills

Select which skills to generate. PM and Tier 1 default to Yes, Tier 2 defaults to No.

### Step 8 — Confirm

Review the summary and confirm to generate all files.

## CLI Flags

```
appteam              Run the interactive wizard
appteam -r           Regenerate from saved settings (no wizard)
appteam -d ./path    Generate into a specific directory
appteam -r -d ./path Regenerate into a specific directory
appteam -v           Print version
appteam -h           Print help
```

## Regenerate

After the first run, your configuration is saved to `.appteam/settings.json`. Regenerate anytime without re-running the wizard:

```bash
appteam -r
```

This is useful after updating appteam itself — regenerate to pick up new templates, skills, or agent improvements.

## Team Sizing

| Size | What You Get |
|------|-------------|
| **Lean** | PM + 1 SWE + SWE-Test. PM handles backlog, progress, and release notes directly. Best for solo work or small projects |
| **Standard** | PM + TPM + configurable SWEs + SWE-Test + Reviewer. Full pipeline with coordination. The default |
| **Full** | PM + TPM + up to 5 SWEs + SWE-Test + SWE-QA + Reviewer + Platform. Everything enabled for large, complex projects |

The team size determines which agent files are generated and which pipeline rules appear in `CLAUDE.md`.
