# F-0117: Scion Worktree Isolation — Git Merge Instructions

**Type:** Enhancement
**Priority:** P0 (critical)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-19

## Problem

Scion agents run in isolated git worktrees — each agent works on its own branch and cannot see other agents' commits. When the pipeline advances (e.g., PM creates a spec, TPM assigns work, SWEs implement), downstream agents cannot read upstream agents' output because the files only exist on the upstream agent's branch.

Currently, the orchestrator must manually include "run git merge pm" or "run git merge swe-1" in every relay message. This is error-prone, repetitive, and easy to forget — leading to agents failing because they can't find specs, backlog entries, or implementation code that was committed on another branch.

The fix is to embed merge instructions directly into each Scion agent template so agents automatically know which branch to merge before starting their work.

## Solution

Add a "Worktree Isolation" section to each Scion agent template in `internal/templates/scion_agents_md.go`. This section tells the agent:

1. You are in an isolated git worktree on your own branch
2. Other agents' commits are on their own branches and invisible to you until you merge
3. Before reading another agent's output, run the appropriate `git merge <source-branch>` command

## Requirements

1. Each Scion agent template in `scion_agents_md.go` must include a new "Worktree Isolation" section placed **before** the "Signaling Protocol" section
2. The section must explain that the agent is in an isolated git worktree on its own branch
3. The section must include the specific `git merge` commands the agent needs to run, varying by role:

   - **PM**: No merge needed — PM starts the pipeline and creates the initial spec. Include a note: "You start the pipeline. No merge is needed before your work."
   - **TPM**: `git merge pm` — to read PM's spec in `docs/specs/` and create backlog entries
   - **SWE-1 through SWE-N**: Conditional on team size:
     - Lean teams (`TeamSize == "lean"`): `git merge pm` — PM assigns directly in lean teams
     - Standard/Full teams: `git merge tpm` — TPM assigns work in standard/full teams
   - **SWE-Test**: `git merge swe-1` as primary, plus a note to merge additional SWE branches if multiple SWEs contributed (e.g., `git merge swe-2` etc.) — to get implementation code for testing
   - **SWE-QA**: `git merge swe-test` — to get test results and verified code
   - **Reviewer**: Conditional on SWE-QA presence:
     - If SWE-QA exists (`IncludeSWEQA`): `git merge swe-qa`
     - Otherwise: `git merge swe-test`
   - **Platform**: Conditional on team size:
     - Lean teams: `git merge pm`
     - Standard/Full teams: `git merge tpm`
   - **Custom agents**: Generic instruction: "Merge the branch of whichever agent produces the input you need. Check the Signaling Protocol section for who you wait for, and merge that agent's branch."

4. The merge instructions must include a brief explanation of *why* the merge is needed (e.g., "to read the spec and backlog assignments")
5. All existing tests in `templates_test.go` must be updated to assert the presence of the "Worktree Isolation" section in rendered Scion agent templates
6. New test cases must verify the conditional merge commands vary correctly by team size
7. Version must be bumped to v0.17.0

## Template Section Format

Each agent's "Worktree Isolation" section should follow this structure:

```markdown
## Worktree Isolation

You are running in an isolated git worktree on your own branch. Other agents' work is on their branches and is **not visible** to you until you explicitly merge it.

**Before starting your work, run:**
` ``
git merge <source-branch>
` ``

This pulls in <description of what you're getting> from <source agent>.
```

## Acceptance Criteria

- [ ] `ScionPMAgentsMD` includes Worktree Isolation section with "no merge needed" note
- [ ] `ScionTPMAgentsMD` includes Worktree Isolation section with `git merge pm`
- [ ] `ScionSWEAgentsMD` includes Worktree Isolation section with conditional merge (pm for lean, tpm for standard/full)
- [ ] `ScionSWETestAgentsMD` includes Worktree Isolation section with `git merge swe-1` and note about additional SWE branches
- [ ] `ScionSWEQAAgentsMD` includes Worktree Isolation section with `git merge swe-test`
- [ ] `ScionReviewerAgentsMD` includes Worktree Isolation section with conditional merge (swe-qa if present, else swe-test)
- [ ] `ScionPlatformAgentsMD` includes Worktree Isolation section with conditional merge (pm for lean, tpm for standard/full)
- [ ] `ScionCustomAgentAgentsMD` includes Worktree Isolation section with generic merge guidance
- [ ] All existing Scion template tests pass with updated assertions
- [ ] New tests verify team-size-conditional merge commands render correctly
- [ ] Version string is `0.17.0`
- [ ] No regressions in non-Scion tests

## Files to Change

- `internal/templates/scion_agents_md.go` — Add Worktree Isolation section to all 8 agent templates
- `internal/templates/templates_test.go` — Update existing + add new test assertions
- `main.go` — Version bump to 0.17.0

## Out of Scope

- Automatic merge logic (agents deciding when to merge dynamically)
- Git conflict resolution instructions (assume clean merges for now — agents work on separate files)
- Changes to the orchestrator relay model or signaling protocol
- Changes to non-Scion (Claude Code Agent Teams) templates

## Dependencies

- F-0109 (Scion signaling fix — established the current Signaling Protocol section format)
- F-0103 (Scion parallel launch — established worktree-based agent isolation)

## Open Questions

- None — the merge relationships are deterministic based on the pipeline flow and team size
