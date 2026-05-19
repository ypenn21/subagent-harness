package templates

const BacklogTemplate = `# Project Backlog — {{.ProjectName}}

**Maintained by:** TPM
**Last updated:** (initial generation)

---

## How to Read This Backlog

- **ID:** Unique feature identifier (` + "`F-0001`" + `, ` + "`F-0002`" + `, etc.) — sequential across all milestones, never reused
- **Priority:** P0 (critical path), P1 (important), P2 (nice to have)
- **Status:** ` + "`backlog`" + ` | ` + "`in-progress`" + ` | ` + "`in-review`" + ` | ` + "`done`" + ` | ` + "`blocked`" + `
- **Owner:** Assigned team member
- **Branch:** Git feature branch
- **Dependencies:** Other feature IDs that must complete first
- **Feedback:** Review notes, blockers, decisions — updated as work progresses

---

## Current Milestone

| ID | Feature | Spec | Priority | Status | Owner | Branch | Dependencies | Feedback |
|----|---------|------|----------|--------|-------|--------|--------------|----------|
| | | | | | | | | |

---

## Future Items

| ID | Feature | Spec | Priority | Status | Owner | Branch | Dependencies | Feedback |
|----|---------|------|----------|--------|-------|--------|--------------|----------|
| | | | | | | | | |

---

## Team Roster

| Role | Agent | Specialty |
|------|-------|-----------|
| PM | PM | Product requirements & PO communication |
| TPM | TPM | Backlog, coordination & progress tracking |
{{- range .SWEs}}
| SWE-{{.Number}} | SWE-{{.Number}} | {{.Title}} |
{{- end}}
{{- if .IncludeSWETest}}
| SWE-Test | SWE-Test | Automated testing & coverage |
{{- end}}
{{- if .IncludeSWEQA}}
| SWE-QA | SWE-QA | E2E testing & QA |
{{- end}}
{{- if .IncludePlatform}}
| Platform | Platform Engineer | Infrastructure & deployment |
{{- end}}
{{- if .IncludeReviewer}}
| Reviewer | Reviewer | Code review & quality |
{{- end}}
`
