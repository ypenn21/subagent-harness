package templates

const PipelineTemplate = `# Development Pipeline — {{.ProjectName}}

> How features, bugs, and enhancements flow through the agent team.

{{- /* Compute link counts for styling */ -}}
{{- $n := len .SWEs -}}
{{- $downCount := add 2 $n -}}
{{- if .IncludeSWETest}}{{$downCount = add $downCount $n}}{{end -}}
{{- if .IncludeSWEQA}}{{$downCount = add $downCount $n}}{{end -}}
{{- if .IncludeReviewer}}{{$downCount = add $downCount $n}}{{end -}}
{{- /* upCount: feedback links + SWE→TPM + TPM→PM + PM→PO */ -}}
{{- $upCount := add $n 2 -}}
{{- if .IncludeSWETest}}{{$upCount = add $upCount $n}}{{end -}}
{{- if .IncludeSWEQA}}{{$upCount = add $upCount $n}}{{end -}}
{{- if .IncludeReviewer}}{{$upCount = add $upCount $n}}{{end -}}
{{- /* sideCount: PM→BACKLOG + SWE→BACKLOG(n) + TPM→PROGRESS + PM→RELEASENOTES + TPM/PM→TAG + platform links */ -}}
{{- $sideCount := add 4 $n -}}
{{- if .IncludePlatform}}{{$sideCount = add $sideCount $n}}{{end}}

` + "```mermaid" + `
flowchart TD
    PO(["PO (Product Owner)"])
    PM["PM"]
    TPM["TPM"]
{{- range .SWEs}}
    SWE{{.Number}}["SWE-{{.Number}}: {{.Title}}"]
{{- end}}
{{- if .IncludeSWETest}}
    SWETEST["SWE-Test"]
{{- end}}
{{- if .IncludeSWEQA}}
    SWEQA["SWE-QA"]
{{- end}}
{{- if .IncludePlatform}}
    PLAT["Platform Engineer"]
{{- end}}
{{- if .IncludeReviewer}}
    REV["Reviewer"]
{{- end}}
    BACKLOG[("docs/BACKLOG.md")]
    PROGRESS[("docs/PROGRESS.md")]
    RELEASENOTES[("docs/RELEASENOTES.md")]
    TAG{{"{{"}}"git tag vX.Y.Z"{{"}}"}}

    %% ── Downward flow (request → execution) ──
    PO -->|"feedback / bugs / features"| PM
    PM -->|"spec + requirements"| TPM
{{- range .SWEs}}
    TPM -->|"assign work"| SWE{{.Number}}
{{- end}}
{{- if .IncludeSWETest}}
{{- range .SWEs}}
    SWE{{.Number}} -->|"hand off"| SWETEST
{{- end}}
{{- end}}
{{- if .IncludeSWEQA}}
{{- range .SWEs}}
    SWE{{.Number}} -->|"hand off"| SWEQA
{{- end}}
{{- end}}
{{- if .IncludeReviewer}}
{{- range .SWEs}}
    SWE{{.Number}} -->|"code review"| REV
{{- end}}
{{- end}}

    %% ── Upward flow (results → reporting) ──
{{- if .IncludeSWETest}}
{{- range .SWEs}}
    SWETEST -.->|"test results"| SWE{{.Number}}
{{- end}}
{{- end}}
{{- if .IncludeSWEQA}}
{{- range .SWEs}}
    SWEQA -.->|"QA results"| SWE{{.Number}}
{{- end}}
{{- end}}
{{- if .IncludeReviewer}}
{{- range .SWEs}}
    REV -.->|"review feedback"| SWE{{.Number}}
{{- end}}
{{- end}}
{{- range .SWEs}}
    SWE{{.Number}} -->|"work complete"| TPM
{{- end}}
    TPM -->|"milestone complete"| PM
    PM -->|"summary report"| PO

    %% ── Side effects (docs & tagging) ──
    PM -.->|"create / update items"| BACKLOG
{{- range .SWEs}}
    SWE{{.Number}} -.->|"update"| BACKLOG
{{- end}}
    TPM -.->|"session log"| PROGRESS
    PM -.->|"version entry"| RELEASENOTES
    PM -.->|"tag release"| TAG
{{- if .IncludePlatform}}
{{- range .SWEs}}
    PLAT -.->|"infra support"| SWE{{.Number}}
{{- end}}
{{- end}}

    %% ── Link colors ──
    linkStyle {{linkRange 0 $downCount}} stroke:#2ea043,stroke-width:2px
    linkStyle {{linkRange $downCount $upCount}} stroke:#58a6ff,stroke-width:2px
` + "```" + `

## Legend

- **Green arrows** (` + "`→`" + `) — Downward flow: request → execution
- **Blue arrows** (` + "`⇢`" + `) — Upward flow: results → reporting
- **Gray dashed** — Side effects (docs updates, tagging, infra support)

## Pipeline Steps

1. **PO** provides feedback, bug reports, or feature requests to the **PM**
2. **PM** creates a product spec (` + "`docs/specs/F-NNNN-slug.md`" + `) and translates feedback into detailed requirements
3. **PM** works with **TPM** to create and prioritize items in docs/BACKLOG.md
4. **TPM** assigns individual work items to **SWE** agents
5. **SWE** agents implement on feature branches
{{- if .IncludeSWETest}}
6. **SWE-Test** runs automated tests to verify implementation
{{- end}}
{{- if .IncludeSWEQA}}
7. **SWE-QA** performs end-to-end quality assurance
{{- end}}
{{- if .IncludeReviewer}}
8. **Reviewer** conducts code review for quality, security, and performance
{{- end}}
{{- if .IncludePlatform}}
9. **Platform Engineer** provides infrastructure and deployment support
{{- end}}
10. **SWE** agents update docs/BACKLOG.md and inform **TPM** when work is complete
11. **TPM** updates docs/PROGRESS.md with session details (what was done, decisions, next steps)
12. **TPM** waits for all milestone items to complete, then reports to **PM**
13. **PM** updates docs/RELEASENOTES.md with the new version entry (Added, Changed, Fixed)
14. **PM** creates a summary of completed work and reports back to the **PO**
15. **Tag release** — after PO approval, create annotated git tag (` + "`git tag -a vX.Y.Z`" + `) and push
`
