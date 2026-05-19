# appteam

**Generate opinionated AI agent teams for Claude Code in seconds.**

One command. Eight questions. A complete multi-agent team with roles, skills, workflows, and tracking — ready to go.

appteam scaffolds everything your Claude Code project needs: agent definitions, slash command skills, a mandatory development pipeline, backlog tracking, and project documentation. No configuration files to copy. No boilerplate to maintain. Just run the wizard and start building.

---

## At a Glance

| | |
|---|---|
| **21 skills** | PM workflows, debugging, testing, security audits, CUJ testing, and more |
| **8 agent roles** | PM, TPM, SWE (x5), SWE-Test, SWE-QA, Reviewer, Platform — plus custom agents |
| **33+ generated files** | CLAUDE.md, agents, skills, docs, pipeline diagram, settings |
| **3 team sizes** | Lean (3 agents), Standard (6+), Full (10+) — pick what fits |
| **Zero dependencies** | Pure Go stdlib. No npm, no pip, no Docker |

---

## What Gets Generated

```
your-project/
  CLAUDE.md                      # Project instructions + pipeline rules
  .claude/
    agents/
      pm.md                      # Product Manager
      tpm.md                     # Technical Program Manager
      swe-1.md ... swe-5.md      # Software Engineers
      swe-test.md                # Test Engineer
      reviewer.md                # Code Reviewer
      ...                        # + custom agents
    skills/
      spec.md                    # /spec — create product specs
      release.md                 # /release — cut a release
      pipeline.md                # /pipeline — spin up agent team
      debug.md                   # /debug — systematic debugging
      test.md                    # /test — write tests
      ...                        # 21 skills total + custom
  docs/
    BACKLOG.md                   # Feature backlog
    PROGRESS.md                  # Session-by-session log
    RELEASENOTES.md              # Version history
    PIPELINE.md                  # MermaidJS workflow diagram
    specs/TEMPLATE.md            # Product spec template
  .appteam/
    settings.json                # Saved config for regeneration
```

---

**[See how it works &#8594;](how.md)** &nbsp; | &nbsp; **[View slide deck &#8594;](slides.html ':ignore')**
