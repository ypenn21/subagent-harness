# /standup — Generate standup summary

## Trigger

User invokes `/standup` to generate a daily standup summary.

## Instructions

1. **Gather "Done" items** — Read recent git log (last 24 hours):
   ```
   git log --oneline --since="24 hours ago"
   ```
   Summarize each commit as a bullet point. Group related commits together
2. **Gather "Next" items** — Read the backlog for current milestone work:
   - Check `docs/BACKLOG.md` for items with status `IN PROGRESS` or `TODO` in the current milestone
   - List the top 2-3 items to focus on next
3. **Identify blockers:**
   - Check for failing tests: `go test ./...`
   - Check for unresolved merge conflicts
   - Note any dependencies on other team members or external services
   - Flag any items that are stuck or need input
4. **Output in standup format:**

   ```
   ## Standup — YYYY-MM-DD

   ### Done
   - Completed item 1
   - Completed item 2

   ### Next
   - Next item 1
   - Next item 2

   ### Blockers
   - Blocker 1 (or "None")
   ```

5. **Report back** — Display the formatted standup summary

## Project Context

- **Project:** appteam
- **Owner:** Ameer Abbas (ameer00@gmail.com)
- **Backlog:** `docs/BACKLOG.md`
