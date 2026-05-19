package templates

const ProgressTemplate = `# Development Progress Log — {{.ProjectName}}

## Session 1

### Goals
- Initial project setup and configuration

### Completed
- Generated project scaffolding with appteam
  - CLAUDE.md with team workflow, conventions, and pipeline rules
  - Agent definitions for PM, TPM{{range .SWEs}}, SWE-{{.Number}}{{end}}{{if .IncludeSWETest}}, SWE-Test{{end}}{{if .IncludeSWEQA}}, SWE-QA{{end}}{{if .IncludePlatform}}, Platform Engineer{{end}}{{if .IncludeReviewer}}, Reviewer{{end}}
  - BACKLOG.md, PROGRESS.md, RELEASENOTES.md

### Next Steps
- Define initial feature backlog in BACKLOG.md
- Begin implementation of first milestone
`
