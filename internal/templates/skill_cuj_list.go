package templates

const SkillCUJListTemplate = `---
name: cuj-list
description: Create or update the Critical User Journey inventory
user-invocable: true
disable-model-invocation: true
---

# /cuj-list — Create or update Critical User Journey inventory

## Trigger

User invokes ` + "`/cuj-list`" + ` to define, review, or update the project's CUJ inventory.

## Instructions

You are helping the user build and maintain a structured inventory of Critical User Journeys (CUJs) — the end-to-end flows that define whether the product works for real users. This inventory lives in ` + "`docs/CUJ.md`" + ` and is the source of truth for all QA testing.

### Phase 1: Read Existing Inventory

1. Check if ` + "`docs/CUJ.md`" + ` exists
2. If it exists, read it and identify:
   - All existing CUJ IDs (` + "`CUJ-001`" + `, ` + "`CUJ-002`" + `, etc.)
   - The highest CUJ-NNN ID (to auto-increment new entries)
   - Current coverage — which user flows are already documented
3. If it doesn't exist, start fresh with CUJ-001

### Phase 2: Collect User Journeys

Ask the user about the application's key user flows. Guide the conversation with questions like:

- What are the most critical things a user does in the app? (These become P0 CUJs)
- What are the main workflows — signup, login, core feature usage, settings, etc.?
- Are there flows that involve multiple pages or steps?
- Which flows, if broken, would mean the product is unusable?
- Are there admin or power-user flows that should be tested?

For each journey, collect:

1. **Name** — short descriptive title (e.g., "User Registration Flow")
2. **Description** — one sentence explaining the journey
3. **Priority** — P0 (critical path), P1 (important), P2 (nice to have)
4. **Steps** — numbered list of specific user actions:
   - Navigate to a page
   - Fill in a form field
   - Click a button
   - Wait for a response
   - Verify a result
5. **Expected outcomes** — what should happen at key steps (success messages, page transitions, data changes)

### Phase 3: Write the Inventory

Write or update ` + "`docs/CUJ.md`" + ` with the following structure:

` + "```" + `markdown
# Critical User Journeys — {{.ProjectName}}

> This file is the source of truth for all CUJ-based QA testing.
> Use ` + "`/cuj-test`" + ` to execute these journeys against the running application.

## Summary

| ID | Journey | Priority | Last Tested | Result |
|----|---------|----------|-------------|--------|
| CUJ-001 | User Registration | P0 | — | — |
| CUJ-002 | Login & Dashboard | P0 | — | — |
| CUJ-003 | Feature X Workflow | P1 | — | — |

---

## CUJ-001: User Registration Flow

- **Priority:** P0 (critical path)
- **Description:** New user signs up and reaches the dashboard
- **Last Tested:** —
- **Result:** —

### Steps

1. Navigate to /signup
2. Fill in name, email, password
3. Click "Create Account"
4. Verify confirmation message appears
5. Navigate to /login and log in
6. Verify dashboard loads

### Expected Outcomes

- Step 3: Form submits, confirmation message shown
- Step 5: Login succeeds, redirects to dashboard
- Step 6: Dashboard displays welcome message
` + "```" + `

### Phase 4: Commit

1. Stage ` + "`docs/CUJ.md`" + `
2. Commit with the message:
   ` + "```" + `
   Add/update CUJ inventory: N journeys (N new, N updated)
   ` + "```" + `
3. Report back with a summary of what was added or changed

## Tips

- Start with P0 journeys — the flows that absolutely must work
- Keep steps specific and actionable (a QA agent will execute them literally)
- Include expected outcomes for steps where something user-visible should happen
- For existing inventories, ask the user if any journeys need updating or removing
- The summary table at the top of ` + "`docs/CUJ.md`" + ` should always be in sync with the detailed entries below

## Project Context

- **Project:** {{.ProjectName}}
- **Owner:** {{.OwnerName}} ({{.OwnerEmail}})
- **CUJ File:** ` + "`docs/CUJ.md`" + `
`
