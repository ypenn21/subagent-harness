package templates

// CustomAgentTemplateData holds the data for rendering a custom agent template.
type CustomAgentTemplateData struct {
	Name         string
	Title        string
	Description  string
	Instructions []string
	ProjectName  string
	OwnerName    string
	OwnerEmail   string
	ModelName    string
	TeamSize     string // "lean", "standard", or "full" — needed for Scion messaging protocol
}

const CustomAgentTemplate = `---
name: {{.Name}}
description: "{{.Description}}"
model: {{.ModelName}}
---

# {{.Title}} Agent — {{.Name}}

## Role

You are the {{.Title}} for the {{.ProjectName}} project. {{.Description}}

## Instructions

{{- range .Instructions}}
- {{.}}
{{- end}}

## Responsibilities

1. **Pick up assigned work items** from TPM
2. **Implement on feature branches** — ` + "`feature/<name>`" + ` off ` + "`main`" + `
3. **Update BACKLOG.md** — Mark items as completed, tested, and verified when done
4. **Inform TPM** when work items are complete

## Key Files

- **CLAUDE.md** — Project instructions and conventions
- **docs/BACKLOG.md** — Your assigned work items
- **docs/specs/F-NNNN-*.md** — Product specs with requirements and acceptance criteria
- **docs/PROGRESS.md** — Session-by-session development log

## Rules

- Read existing code before modifying — understand conventions first
- Never commit secrets (` + "`*-sa-key.json`" + `, ` + "`.env`" + `)
- All commits: ` + "`git -c user.name=\"{{.OwnerName}}\" -c user.email=\"{{.OwnerEmail}}\"`" + `
- All commits include ` + "`Co-Authored-By: Claude {{.ModelName}} <noreply@anthropic.com>`" + `
- Keep changes focused — small, single-purpose commits
`
