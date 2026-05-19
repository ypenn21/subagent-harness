# Scion Integration — Lessons Learned

Field notes from building Scion framework support into appteam and running
multiple end-to-end pipeline tests (PM → SWE → SWE-Test) with real agents
using Vertex AI auth. Written March 2026.

---

## Table of Contents

1. [Architecture Mental Model](#architecture-mental-model)
2. [What Works Well](#what-works-well)
3. [Bugs and Friction Points](#bugs-and-friction-points)
4. [Design Observations](#design-observations)
5. [Recommendations for the Scion Team](#recommendations-for-the-scion-team)
6. [Integration Guide for Tool Authors](#integration-guide-for-tool-authors)

---

## Architecture Mental Model

After extensive testing, here is how Scion actually works under the hood:

```
Host Machine
├── .scion/
│   ├── grove-id              ← UUID identifying this project grove
│   ├── templates/            ← Agent definitions (committed to git)
│   │   ├── pm/
│   │   │   ├── scion-agent.yaml
│   │   │   ├── system-prompt.md
│   │   │   └── agents.md
│   │   └── swe-1/
│   │       └── ...
│   └── agents/               ← Runtime state (gitignored)
│       ├── pm/
│       │   ├── prompt.md
│       │   ├── scion-agent.json
│       │   └── workspace/    ← Git worktree for this agent
│       └── swe-1/
│           └── ...
├── CLAUDE.md
└── (project files)

Docker Container (per agent)
├── /home/scion/              ← Agent's workspace (mounted worktree)
│   ├── CLAUDE.md
│   ├── .scion/templates/     ← Visible but read-only context
│   └── (project files from worktree)
├── /home/scion/agent.log     ← sciontool writes here
├── sciontool                 ← In-container tool (status, hooks)
└── claude                    ← The AI agent process
    └── ~/.config/gcloud/     ← Mounted read-only from host for Vertex AI
```

### Key Architectural Facts

- **Each agent gets its own git worktree.** Agent `pm` works on branch `pm`,
  agent `swe-1` works on branch `swe-1`. They cannot see each other's
  uncommitted work. To see another agent's commits, an agent must explicitly
  `git merge <branch-name>`.

- **Agents run inside Docker containers.** They have no access to Docker CLI,
  no access to the host filesystem (except mounted volumes), and no ability
  to run `scion` commands. The `scion` CLI is a host-only tool.

- **`sciontool` is the in-container counterpart.** It provides `sciontool status`,
  `sciontool hook`, etc. It does NOT provide messaging capabilities.

- **The orchestrator is the human (or script) on the host.** The orchestrator
  monitors agent status and relays messages between agents using
  `scion message <name> "text"` from the host side.

- **Harness configs** live at `~/.scion/harness-configs/<name>/` and define the
  container image, user, and startup files (`.bashrc`, `.claude.json`,
  `.claude/settings.json`). The `claude` harness is the default.

---

## What Works Well

### Vertex AI Authentication
Auth just works. Scion automatically mounts `~/.config/gcloud` read-only into
containers. Auth resolves as `vertex-ai` with zero configuration. Agents use
`claude-opus-4-6@default` without any extra setup. This was initially suspected
as the source of errors but turned out to be perfectly solid.

### Agent Isolation
Git worktree isolation is genuinely useful. Each agent works on its own branch
without conflicts. The merge-based handoff (`git merge pm` to get PM's spec)
is clean and audit-friendly — you can see exactly what each agent contributed.

### Template System
The template structure (`scion-agent.yaml` + `system-prompt.md` + `agents.md`)
is well-designed. YAML config is minimal and correct. The separation of
agent instructions from system prompts is clean.

### Container Lifecycle
`scion start`, `scion stop`, `scion list` all work reliably. Container startup
is fast (~5-10 seconds to agent activity). `scion stop` cleanly shuts down
containers.

### sciontool Status Signaling
`sciontool status task_completed "message"` works perfectly inside containers.
The `[COMPLETED]` log entry is easy to parse and monitor. This is the correct
primitive for agent-to-orchestrator communication.

---

## Bugs and Friction Points

### 1. `scion message` Doesn't Work Inside Containers (Critical, Fixed in appteam)

**Impact:** Any agent template that tells agents to use `scion message` to
communicate with other agents will fail 100% of the time.

**Root cause:** `scion message` uses `docker exec` under the hood to inject
text into the target agent's tmux session. Docker CLI is not available inside
containers.

**Error message:**
```
Error: docker ps failed: exec: "docker": executable file not found in $PATH
```

**Our fix:** Switched all agent templates to use `sciontool status task_completed`
for signaling and documented the orchestrator relay pattern (orchestrator on host
uses `scion message` to relay between agents).

**Recommendation for Scion:** Either:
- (a) Document prominently that `scion message` is host-only and agents cannot
  message each other directly, OR
- (b) Add a `sciontool message <target> "text"` command that works inside
  containers (perhaps via a sidecar API or socket)

Option (b) would be transformative — it would enable truly autonomous multi-agent
pipelines without human orchestration.

### 2. `scion init` Conflicts with Pre-existing `.scion/` Directory (Fixed in appteam)

**Impact:** If a tool (like appteam) creates `.scion/templates/` before
`scion init` runs, `scion init` refuses with "already inside a scion project"
but the grove isn't actually initialized (no `grove-id`).

**Root cause:** `scion init` checks for `.scion/` directory existence but
doesn't check for `grove-id`. Finding the directory, it assumes initialization
is complete.

**Our fix:** appteam now generates `.scion/grove-id` (UUID v4) and
`.scion/agents/` directory alongside templates.

**Recommendation for Scion:** Make `scion init` more resilient:
- If `.scion/` exists but `grove-id` is missing, offer to complete initialization
  rather than refusing outright
- Or at minimum, `scion init --force` should work in this scenario (currently
  it also refuses)

### 3. `.gitignore` Warning is Incorrect (Cosmetic)

**Impact:** Every `scion start` prints:
```
Warning: '.scion/' is not in .gitignore. Run 'scion init' to fix this.
```
even when `.scion/agents/` IS in `.gitignore`.

**Root cause:** The warning check looks for `.scion/` literally. But `.scion/`
in `.gitignore` would also ignore `.scion/templates/`, which need to be committed.
The correct entry is `.scion/agents/` (runtime state only), which is what
`scion init` itself creates.

**Recommendation for Scion:** Fix the gitignore check to accept `.scion/agents/`
as sufficient. The check should verify that runtime state is ignored, not that
the entire `.scion/` tree is ignored.

### 4. `scion logs` Path Mismatch (Minor)

**Impact:** `scion logs <agent>` fails with:
```
Error: log file not found: /tmp/project/.scion/agents/pm/home/agent.log
```

**Root cause:** `scion logs` looks for the log file on the host filesystem at
`.scion/agents/<name>/home/agent.log`, but the log is written inside the
container at `/home/scion/agent.log`. The host path mapping doesn't match.

**Workaround:** Use `docker exec <agent> tail /home/scion/agent.log` instead.

**Recommendation for Scion:** Fix the log path resolution, or mount the log
file back to the host so `scion logs` can find it.

### 5. `scion message` Requires `--interrupt` for Idle Agents (Friction)

**Impact:** If an agent is idle (at a prompt, waiting for input), a plain
`scion message <name> "text"` is silently ignored. The message never appears
in the agent's session.

**Root cause:** The message is injected into tmux but the agent process isn't
actively reading. The `--interrupt` flag forces injection into the tmux pane.

**Recommendation for Scion:** Make `--interrupt` the default behavior, or at
least auto-detect idle agents and interrupt them. A silently dropped message
is a confusing failure mode — the orchestrator thinks the message was delivered
but the agent never receives it.

### 6. `prompt.md` Task Conflicts on Re-runs (Minor)

**Impact:** If a previous agent run left behind `.scion/agents/<name>/prompt.md`,
the next `scion start` fails with "task conflict: both prompt.md and start
options provide a task."

**Workaround:** `rm -rf .scion/agents/` between runs.

**Recommendation for Scion:** `scion start` should clean up stale agent state
from previous runs, or `scion stop` should remove `prompt.md` on shutdown.

### 7. Agents Create Feature Branches Instead of Committing to Worktree Branch (Critical, Fixed in appteam)

**Impact:** When agents are told to "implement on a feature branch" (a common
convention in CLAUDE.md / system prompts), they create `feature/<name>` branches
inside their worktree. This breaks the `git merge <agent-name>` handoff model —
downstream agents merge the agent's worktree branch (e.g., `swe-1`) but the
actual code is on `feature/scaffold` (unreachable from `swe-1`).

**Root cause:** Template/instructions saying "feature branches" — agents follow
instructions literally. In a Scion worktree model, each agent already has an
isolated branch named after itself. Creating sub-branches defeats the purpose.

**Our fix:** Made ALL branching references in CLAUDE.md, agent templates, and
pipeline descriptions framework-conditional. Scion projects now consistently say
"commit directly to your worktree branch — do NOT create feature branches."

**Key insight for tool authors:** If you generate project scaffolds that include
CLAUDE.md or similar AI agent instructions, you must audit EVERY reference to
"feature branches" — not just the branching convention section, but also pipeline
step descriptions, workflow sections, and agent role definitions. A single
contradictory reference can override all your other instructions.

---

## Design Observations

### The Orchestrator Relay Pattern

Scion's current architecture requires a human (or script) orchestrator on the
host to relay messages between agents. This is because:

1. Agents can signal completion (`sciontool status task_completed`)
2. But they can't send messages to other agents (no `scion message` inside containers)
3. So the orchestrator monitors for completion, then relays using `scion message` from host

This creates a **hub-and-spoke** model:

```
        ┌─────────────┐
        │ Orchestrator │  (host)
        │   (human)    │
        └──────┬───────┘
         ╱     │     ╲
        ╱      │      ╲
   ┌────┐  ┌─────┐  ┌──────────┐
   │ PM │  │SWE-1│  │ SWE-Test │  (containers)
   └────┘  └─────┘  └──────────┘
```

Each agent signals up (via sciontool), orchestrator relays across (via scion message).
Agents never talk directly to each other.

**Implications:**
- The orchestrator must be attentive and relay promptly, or the pipeline stalls
- The orchestrator must know the pipeline order (PM → SWE → SWE-Test)
- This is fine for human-supervised workflows but limits full automation
- A `sciontool message` command (see recommendation #1) would unlock
  true agent-to-agent automation

### Git Worktree Isolation is a Feature AND a Friction Point

**Feature:** Each agent has a clean workspace. No merge conflicts during parallel
work. Clear audit trail of who changed what.

**Friction:** Agents can't see each other's work without explicit `git merge`.
This must be documented in every agent template that depends on another agent's
output (e.g., SWE needs PM's spec, SWE-Test needs SWE's code).

The merge instruction is the single most important thing in an agent's template.
If you forget to tell SWE-1 to `git merge pm`, it will try to implement from
the CLAUDE.md alone without the spec — and it might succeed, but the result
won't match the spec.

### Template Discovery is Implicit

Scion finds templates by scanning `.scion/templates/` directories. The mapping
from template name to agent role is purely by directory name convention. There's
no manifest or registry. This is simple and works well, but means:

- Template names must match what `scion start --type <name>` expects
- There's no validation that a template is well-formed until you try to start it
- Typos in `--type` silently fall back to the `default` template

---

## Recommendations for the Scion Team

### Priority 1: In-Container Messaging (High Impact)

Add `sciontool message <target> "text"` that works inside containers. This would:
- Eliminate the need for human orchestrators in simple pipelines
- Enable fully autonomous multi-agent workflows
- Make Scion competitive with framework-native agent teams (Claude Code Agent Teams)

Implementation options:
- Unix socket or HTTP sidecar in each container that forwards to host's `scion message`
- A shared message bus (Redis, file-based, or Unix socket) mounted into containers
- A `sciontool` subcommand that writes to a well-known path, with a host-side
  watcher that picks up and delivers messages

### Priority 2: Fix `scion init` Resilience (Medium Impact)

When `.scion/` exists but `grove-id` is missing, complete the initialization
instead of refusing. This unblocks any tool that generates templates before
`scion init` runs. The `--force` flag should also handle this case.

### Priority 3: Fix `.gitignore` Warning (Low Impact, High Annoyance)

Accept `.scion/agents/` as a valid gitignore entry. The current check that
looks for `.scion/` literally is incorrect — ignoring all of `.scion/` would
also ignore templates, which should be committed.

### Priority 4: Default `--interrupt` for `scion message` (Medium Impact)

Silent message drops when an agent is idle is a confusing failure mode. Either:
- Make `--interrupt` the default
- Auto-detect idle state and interrupt automatically
- Return an error/warning if the message couldn't be delivered

### Priority 5: Fix `scion logs` Path Resolution (Low Impact)

Map the in-container log path (`/home/scion/agent.log`) to the host path
correctly so `scion logs <agent>` works without `docker exec`.

### Priority 6: Clean Up Stale Agent State (Low Impact)

`scion stop` should remove `prompt.md` (or `scion start` should handle existing
`prompt.md` gracefully) to prevent "task conflict" errors on re-runs.

### Priority 7: Template Validation (Nice to Have)

Add `scion template validate <name>` to check that a template directory has
the required files (`scion-agent.yaml`, etc.) and the YAML is well-formed.
Would catch typos and missing files before agent launch.

---

## Integration Guide for Tool Authors

If you're building a tool that generates Scion project configs (like appteam),
here's what you need to create:

### Required Files

```
.scion/
├── grove-id                  ← UUID v4 string + newline
├── agents/                   ← Empty directory (runtime state)
└── templates/
    └── <agent-name>/
        ├── scion-agent.yaml  ← Agent metadata
        ├── system-prompt.md  ← System prompt for the agent
        └── agents.md         ← Agent instructions (CLAUDE.md equivalent)
```

### scion-agent.yaml Format

```yaml
schema_version: "1"
description: "Agent role description"
agent_instructions: agents.md
system_prompt: system-prompt.md
default_harness_config: claude
```

### .gitignore Entry

```
.scion/agents/
```

Do NOT use `.scion/` — that would ignore templates too.

### Grove-ID Generation

Generate a UUID v4 and write it (with trailing newline) to `.scion/grove-id`.
Use `crypto/rand` or equivalent — don't use timestamp-based UUIDs.

### Agent Template Best Practices

1. **Always include merge instructions.** Every agent that depends on another
   agent's output must be told to `git merge <branch>` first.

2. **Use `sciontool status task_completed` for signaling.** Never use
   `scion message` in agent templates — it doesn't work inside containers.

3. **Warn about Docker unavailability.** Agents may try to use Docker commands.
   Explicitly state that Docker CLI is not available inside containers.

4. **Document the orchestrator relay model.** Make it clear that the orchestrator
   (human on host) relays messages between agents. Agents cannot communicate
   directly with each other.

5. **Include `--interrupt` in relay documentation.** Any instructions for the
   orchestrator should mention the `--interrupt` flag for idle agents.

---

*Last updated: 2026-03-19 — Based on appteam v0.16.2 and Scion (current as of March 2026)*
