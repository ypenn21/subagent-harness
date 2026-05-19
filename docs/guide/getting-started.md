# Getting Started

## Install

### Global install (recommended)

```bash
go install github.com/ahafin/appteam@latest
```

Make sure `$GOPATH/bin` (or `$GOBIN`) is in your `PATH`.

### Local build

```bash
git clone https://github.com/ahafin/appteam.git
cd appteam
go build -o appteam .
```

## Run the wizard

```bash
cd /path/to/your/project
appteam
```

The wizard walks through 7 steps:

### Step 1 — Project Basics

```
▸ Project name: myapp
▸ Short description: A REST API for managing widgets
▸ Tech stack: Go, PostgreSQL, Chi router
▸ Target directory [.]: .
```

The target directory is where all files will be generated. Use `.` for the current directory.

### Step 2 — Git Repository

appteam detects whether a `.git` directory already exists:

- **Existing repo** — asks for an optional remote URL
- **No repo** — offers to run `git init` and optionally create a GitHub repo via `gh repo create`

### Step 3 — Product Owner

```
▸ Name [Ameer Abbas]: Jane Smith
▸ Email [ameer00@gmail.com]: jane@example.com
▸ GitHub username [ahafin]: jsmith
```

The product owner is the human in the loop — all agent work flows through this person.

### Step 4 — GCP Configuration

Optional. If your project uses Google Cloud Platform, answer `y` and provide:
- GCP Project ID
- GCP Project Number
- Organization
- Region

If your project doesn't use GCP, answer `n` (the default) and skip this step entirely.

### Step 5 — Agent Team

```
▸ Number of SWE agents (1-5) [3]: 2

SWE-1
▸ Title/specialty: Backend & API
▸ Specialty bullets (blank line to finish)
  │ REST endpoint implementation
  │ Database queries and migrations
  │

SWE-2
▸ Title/specialty: Frontend & UI
▸ Specialty bullets (blank line to finish)
  │ React components
  │ CSS and responsive layout
  │

▸ Include Platform Engineer? (y/N): n
▸ Include Code Reviewer? (Y/n): y
▸ Include SWE-Test? (Y/n): y
▸ Include SWE-QA? (Y/n): n
```

Choose the number of SWE agents (1–5) and give each one a title and specialty bullets. Then choose which optional agents to include.

You'll also pick the default AI model for all agents:

```
Default model for agents:
  1. Opus 4.6 (claude-opus-4-6)
  2. Sonnet 4.6 (claude-sonnet-4-6)
  3. Haiku 4.5 (claude-haiku-4-5-20251001)
▸ Select model [1]: 1
```

### Step 6 — Conventions

Add project-specific conventions as free-text lines. These appear in the generated `CLAUDE.md` under "Important Conventions":

```
▸ App-specific conventions (blank line to finish)
  │ All API responses use JSON:API format
  │ Database migrations must be reversible
  │
```

### Step 7 — Confirm

The wizard displays a summary of your configuration. Review it and confirm to generate files.

## What happens next

After confirmation, appteam generates all files in the target directory (see [Generated Files](generated-files.md)) and saves your configuration to `.appteam/settings.json` for future use.

The next time you run `appteam` in the same directory, it detects the saved config and asks:

```
● Found existing .appteam/settings.json
  Project: myapp  |  SWEs: 2  |  Model: Opus 4.6

  ▸ Use saved config? (Y/n):
```

Press Enter to regenerate from saved settings, or `n` to run the wizard again.

You can also regenerate non-interactively:

```bash
appteam -r
```

See [Configuration](configuration.md) for more on saved settings and regeneration.
