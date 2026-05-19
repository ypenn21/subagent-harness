package templates

const SkillADRTemplate = `---
name: adr
description: Create an Architecture Decision Record for a significant technical decision
user-invocable: true
---

# /adr — Architecture Decision Record

## Trigger

User invokes ` + "`/adr`" + ` with a decision title or architectural question.

## Instructions

1. **Determine the next ADR number:**
   - Check ` + "`docs/adr/`" + ` for existing ADR files
   - Auto-increment: if the highest is ` + "`0003-*.md`" + `, the next is ` + "`0004`" + `
   - If the directory does not exist, create it and start at ` + "`0001`" + `
2. **Create the ADR file** at ` + "`docs/adr/NNNN-kebab-case-title.md`" + `
3. **Fill in the ADR template:**

   ` + "```" + `markdown
   # NNNN. Title of Decision

   **Date:** YYYY-MM-DD
   **Status:** Proposed | Accepted | Deprecated | Superseded by [NNNN]

   ## Context

   What is the issue or question that motivates this decision?
   What constraints or forces are at play?

   ## Decision

   What is the change being proposed or adopted?
   State the decision clearly and concisely.

   ## Consequences

   ### Positive
   - Benefits of this decision

   ### Negative
   - Trade-offs or downsides of this decision

   ### Risks
   - What could go wrong and how to mitigate it
   ` + "```" + `

4. **Link related ADRs** — If this decision supersedes or depends on a previous ADR, add cross-references
5. **Report back** — Confirm the ADR number, file path, and a brief summary of the decision recorded

## Project Context

- **Project:** {{.ProjectName}}
- **Owner:** {{.OwnerName}} ({{.OwnerEmail}})
`
