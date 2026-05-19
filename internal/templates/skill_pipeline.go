package templates

const SkillPipelineTemplate = `---
name: pipeline
description: Visualize the agent workflow pipeline as a MermaidJS diagram
user-invocable: true
---

# /pipeline — Spin up the full agent team

## Trigger

User invokes ` + "`/pipeline`" + ` with a feature description or milestone name.

## Instructions
{{if eq .Framework "scion"}}
1. **Launch ALL agents simultaneously:**
` + "```" + `
scion start pm "Create spec for <feature>. Signal completion via sciontool when done." --type pm
{{- if ne .TeamSize "lean"}}
scion start tpm "Wait for orchestrator to relay PM spec. Assign work to SWEs." --type tpm
{{- end}}
{{- range .SWEs}}
scion start swe-{{.Number}} "Wait for orchestrator to relay assignment. Implement, then signal completion." --type swe-{{.Number}}
{{- end}}
scion start swe-test "Wait for orchestrator to relay SWE completion. Run tests and signal results." --type swe-test
{{- if and (ne .TeamSize "lean") .IncludeSWEQA}}
scion start swe-qa "Wait for orchestrator to relay test results. Run QA and signal results." --type swe-qa
{{- end}}
{{- if and (ne .TeamSize "lean") .IncludeReviewer}}
scion start reviewer "Wait for orchestrator to relay test/QA results. Review code and signal results." --type reviewer
{{- end}}
` + "```" + `
2. **Orchestrator relays messages** — Agents signal completion via ` + "`sciontool status task_completed`" + ` inside their containers. You (the orchestrator on the host) relay messages between agents using ` + "`scion message <name> \"message\"`" + `
3. **Pipeline order — you enforce it by relaying:**
   - {{if eq .TeamSize "lean"}}PM → (you relay to SWE-1) → (you relay to SWE-Test) → (you relay to PM){{else}}PM → (you relay to TPM) → (you relay to SWEs) → (you relay to SWE-Test){{if .IncludeSWEQA}} + SWE-QA{{end}}{{if .IncludeReviewer}} → (you relay to Reviewer){{end}} → (you relay to PM){{end}}
   - docs/BACKLOG.md, docs/PROGRESS.md, and docs/RELEASENOTES.md MUST be updated
4. **Monitor progress** — ` + "`scion list`" + ` to see running agents and their status
5. **Relay messages** — When an agent signals completion, use ` + "`scion message <next-agent> \"message\"`" + ` to hand off to the next agent
6. **Attach to agents** — ` + "`scion attach <name>`" + ` for interactive sessions
7. **Completion** — When PM signals "Milestone complete", stop all agents with ` + "`scion stop <name>`" + ` for each agent
{{else}}
1. **Create a team** — Use ` + "`TeamCreate`" + ` with a descriptive team name based on the feature (e.g., ` + "`feature-dark-mode`" + `)
2. **Create tasks** — Use ` + "`TaskCreate`" + ` to create work items for the pipeline:
   - **Task 1:** PM — Define requirements and create spec in ` + "`docs/specs/`" + `
{{- if ne .TeamSize "lean"}}
   - **Task 2:** TPM — Create backlog items and coordinate SWE assignments
{{- end}}
   - **Task 3+:** SWE — Implement the feature (one task per work item)
   - **Task N:** SWE-Test — Run tests and verify acceptance criteria
{{- if .IncludeReviewer}}
   - **Task N+1:** Reviewer — Code review
{{- end}}
3. **Spawn agents** — Use the ` + "`Agent`" + ` tool with ` + "`team_name`" + ` parameter to launch each agent in its own tmux pane. Always use model ` + "`{{.ModelID}}`" + ` for all agents:
   - PM agent (reads ` + "`.claude/agents/pm.md`" + `)
{{- if ne .TeamSize "lean"}}
   - TPM agent (reads ` + "`.claude/agents/tpm.md`" + `)
{{- end}}
{{- range .SWEs}}
   - SWE-{{.Number}} agent (reads ` + "`.claude/agents/swe-{{.Number}}.md`" + `)
{{- end}}
   - SWE-Test agent (if included)
4. **Assign tasks** — Use ` + "`TaskUpdate`" + ` with ` + "`owner`" + ` set to each agent's name
5. **Follow the mandatory pipeline:**
   - PO feedback → PM creates spec → {{if ne .TeamSize "lean"}}TPM creates backlog → {{end}}SWEs implement → SWE-Test verifies{{if .IncludeReviewer}} → Reviewer approves{{end}}
   - docs/BACKLOG.md, docs/PROGRESS.md, and docs/RELEASENOTES.md MUST be updated
6. **Monitor and coordinate** — Use ` + "`SendMessage`" + ` to communicate with agents and track progress
7. **Shutdown gracefully** — Send ` + "`shutdown_request`" + ` to each agent when all work is complete
8. **Clean up** — Use ` + "`TeamDelete`" + ` after all agents have shut down
{{end}}
## Project Context

- **Project:** {{.ProjectName}}
- **Owner:** {{.OwnerName}} ({{.OwnerEmail}})
{{- if eq .Framework "scion"}}
- **Framework:** Scion
- **Harness:** ` + "`{{.DefaultHarness}}`" + `
- **Agent templates:** ` + "`.scion/templates/`" + `
{{- else}}
- **Framework:** Claude Code Agent Teams
- **Model:** ` + "`{{.ModelID}}`" + `
- **Agent definitions:** ` + "`.claude/agents/`" + `
{{- end}}
- **Pipeline diagram:** ` + "`docs/PIPELINE.md`" + `
`
