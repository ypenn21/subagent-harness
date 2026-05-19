package templates

const ReleaseNotesTemplate = `# Release Notes — {{.ProjectName}}

## v0.1 — Initial Release

{{.Description}}

### Features
- Project scaffolding generated with appteam
- Multi-agent team structure configured
- Development pipeline and workflow established

### Team
{{- range .SWEs}}
- SWE-{{.Number}}: {{.Title}}
{{- end}}
{{- if .IncludeSWETest}}
- SWE-Test: Automated testing
{{- end}}
{{- if .IncludeSWEQA}}
- SWE-QA: E2E testing & QA
{{- end}}
{{- if .IncludePlatform}}
- Platform Engineer: Infrastructure & deployment
{{- end}}
{{- if .IncludeReviewer}}
- Reviewer: Code review & quality
{{- end}}
`
