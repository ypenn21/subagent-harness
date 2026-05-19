# F-0109: Scion Inter-Agent Messaging Fix — Orchestrator Relay Model

**Type:** Bug
**Priority:** P0 (critical)
**Status:** Done
**Requested by:** PO
**Date:** 2026-03-19

## Problem

E2E testing of the Scion parallel agent workflow (v0.16.0, F-0103) revealed a critical bug: the `scion message` command cannot execute inside Docker containers because the Docker CLI is unavailable within Scion containers. Agents would fail at runtime with:

```
docker ps failed: exec: docker: executable file not found in $PATH
```

This made the entire inter-agent messaging protocol non-functional in any Scion deployment, since all agents run inside containers. The `scion message` approach assumed Docker CLI access from within containers, which is incorrect — containers do not have access to the Docker socket or CLI by default.

## Root Cause

The v0.16.0 messaging protocol instructed agents to use `scion message <target> "<text>"` for inter-agent communication. Under the hood, `scion message` shells out to `docker exec` on the target container, which requires the Docker CLI on the calling side. Since agents run inside containers without Docker installed, every `scion message` call fails.

## Solution

Replace the direct agent-to-agent messaging model with an **orchestrator relay model**:

- Agents signal task completion using `sciontool status task_completed` (an in-container tool that doesn't require Docker)
- The orchestrator (running on the host with Docker access) monitors agent status and relays information between agents
- Agents no longer attempt to communicate directly with each other

## Requirements

1. The Scion agent templates (`scion_agents_md.go`) must replace all "Messaging Protocol" sections with a "Signaling Protocol" that uses `sciontool status task_completed` instead of `scion message`
2. Each agent template must include a warning that Docker CLI is unavailable inside Scion containers
3. The `/pipeline` skill (`skill_pipeline.go`) must be updated to describe the orchestrator relay model for Scion workflows — agents signal completion, orchestrator relays between agents from the host
4. The CLAUDE.md template (`claude_md.go`) must update all three Scion orchestration sections (lean, standard, full team sizes) to describe the orchestrator relay workflow
5. The generator (`generator.go`) must fix the `.gitignore` entry from `.scion/` to `.scion/agents/` so that `.scion/templates/` remains tracked in git
6. All existing Scion messaging tests in `templates_test.go` must be updated to assert the new signaling protocol content
7. Version must be bumped from v0.16.0 to v0.16.1

## Acceptance Criteria

- [x] Agent templates contain `sciontool status task_completed` instead of `scion message`
- [x] Agent templates include Docker unavailability warning
- [x] `/pipeline` skill describes orchestrator relay model for Scion
- [x] CLAUDE.md lean/standard/full Scion sections all use orchestrator relay workflow
- [x] `.gitignore` entry is `.scion/agents/` (not `.scion/`)
- [x] All Scion-related tests pass with updated assertions
- [x] Version string is `0.16.1`
- [x] No regressions in non-Scion tests

## Files Changed

- `internal/templates/scion_agents_md.go` — Signaling protocol sections
- `internal/templates/skill_pipeline.go` — Orchestrator relay model
- `internal/templates/claude_md.go` — All 3 Scion orchestration sections
- `internal/generator/generator.go` — `.gitignore` fix
- `internal/templates/templates_test.go` — Updated test assertions
- `main.go` — Version bump

## Out of Scope

- Adding actual bidirectional messaging support to Scion containers (would require Docker-in-Docker or socket mounting)
- Changing the Scion container runtime architecture
- Adding new agent roles or templates

## Dependencies

- F-0103 (Scion parallel launch — the feature this bug was found in)

## Open Questions

- None — fix has been implemented and verified
