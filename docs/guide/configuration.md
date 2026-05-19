# Configuration

## Saved Settings

After running the wizard, appteam saves your configuration to `.appteam/settings.json`. This file captures everything you entered so you can regenerate files without re-answering the wizard.

### Settings Schema

```json
{
  "ProjectName": "myapp",
  "Description": "A REST API for managing widgets",
  "TechStack": "Go, PostgreSQL, Chi router",
  "InitGit": true,
  "CreateRepo": false,
  "RepoURL": "https://github.com/org/myapp.git",
  "GitHubOrg": "org",
  "OwnerName": "Jane Smith",
  "OwnerEmail": "jane@example.com",
  "OwnerGitHub": "jsmith",
  "GCP": {
    "Enabled": false,
    "ProjectID": "",
    "ProjectNumber": "",
    "Organization": "",
    "Region": ""
  },
  "SWEs": [
    {
      "Number": 1,
      "Title": "Backend & API",
      "Bullets": [
        "REST endpoint implementation",
        "Database queries and migrations"
      ]
    },
    {
      "Number": 2,
      "Title": "Frontend & UI",
      "Bullets": [
        "React components",
        "CSS and responsive layout"
      ]
    }
  ],
  "IncludePlatform": false,
  "IncludeReviewer": true,
  "IncludeSWEQA": false,
  "IncludeSWETest": true,
  "Conventions": [
    "All API responses use JSON:API format"
  ],
  "ModelName": "Opus 4.6",
  "ModelID": "claude-opus-4-6",
  "TargetDir": "/path/to/project"
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `ProjectName` | string | Project name used in headers and references |
| `Description` | string | One-line project description |
| `TechStack` | string | Languages, frameworks, and tools |
| `InitGit` | bool | Whether to run `git init` |
| `CreateRepo` | bool | Whether to create a GitHub repo via `gh` |
| `RepoURL` | string | Git remote URL |
| `GitHubOrg` | string | GitHub org or username for repo creation |
| `OwnerName` | string | Product Owner's name (used in git config and agent instructions) |
| `OwnerEmail` | string | Product Owner's email |
| `OwnerGitHub` | string | Product Owner's GitHub username |
| `GCP.Enabled` | bool | Whether to include GCP configuration sections |
| `GCP.ProjectID` | string | GCP project ID |
| `GCP.ProjectNumber` | string | GCP project number |
| `GCP.Organization` | string | GCP organization |
| `GCP.Region` | string | GCP region (e.g., `us-west1`) |
| `SWEs` | array | List of SWE agent configurations |
| `SWEs[].Number` | int | SWE number (1–5) |
| `SWEs[].Title` | string | Specialty title |
| `SWEs[].Bullets` | array | Specialty bullet points |
| `IncludePlatform` | bool | Generate Platform Engineer agent |
| `IncludeReviewer` | bool | Generate Code Reviewer agent |
| `IncludeSWEQA` | bool | Generate SWE-QA agent |
| `IncludeSWETest` | bool | Generate SWE-Test agent |
| `Conventions` | array | Project-specific conventions |
| `ModelName` | string | Display name for the AI model |
| `ModelID` | string | Claude model ID used when spawning agents |
| `TargetDir` | string | Directory where files were generated (overridden by CWD on load) |

## Regeneration

### From saved settings (non-interactive)

```bash
appteam -r
```

This reads `.appteam/settings.json` and regenerates all files. The `TargetDir` is always overridden with the current working directory, so settings files are portable across machines.

### On startup (interactive prompt)

When you run `appteam` without flags in a directory with existing settings, it offers to reuse them:

```
● Found existing .appteam/settings.json
  Project: myapp  |  SWEs: 2  |  Model: Opus 4.6

  ▸ Use saved config? (Y/n):
```

- Press Enter or `y` to regenerate from saved settings
- Press `n` to run the full wizard (your answers will overwrite the saved settings)

### Important: tracking files

Regeneration overwrites `docs/BACKLOG.md`, `docs/PROGRESS.md`, and `docs/RELEASENOTES.md` with fresh templates. If you have project-specific content in these files, restore them after regeneration:

```bash
git checkout -- docs/BACKLOG.md docs/PROGRESS.md docs/RELEASENOTES.md
```

The `/regenerate` skill handles this automatically.

## Editing Settings Manually

You can edit `.appteam/settings.json` directly to change your configuration without re-running the wizard. Common changes:

- **Add a new SWE** — append to the `SWEs` array with the next number
- **Change the model** — update `ModelName` and `ModelID`
- **Toggle optional agents** — set `IncludeReviewer`, `IncludeSWETest`, etc. to `true`/`false`
- **Update conventions** — add or remove strings from the `Conventions` array

After editing, run `appteam -r` to regenerate files with the new settings.

## `.gitignore` Considerations

The `.appteam/` directory should generally be committed to version control so team members can regenerate files. However, `TargetDir` contains an absolute path — it's overridden at load time, so different machines work correctly regardless of the stored path.
