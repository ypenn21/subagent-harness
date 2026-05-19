# F-0070: Target Directory CLI Flag (`-d`/`--dir`)

**Type:** Feature
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-17

## Problem

Users must always enter the target directory interactively via the wizard (Step 1) or rely on CWD for `-r`. There is no way to specify the target directory from the command line, which limits scripting, automation, and quick one-liner usage.

## Requirements

1. Add `-d <folder>` / `--dir <folder>` CLI flag to set the target directory
2. If the specified directory does not exist, create it (including parents, like `mkdir -p`)
3. When used with the wizard: skip the "Target directory" prompt in Step 1 and pre-fill `cfg.TargetDir` with the `-d` value
4. When used with `-r`: `appteam -r -d ./my-project` regenerates into the specified directory, overriding the CWD-based default and any `TargetDir` in `settings.json`
5. When used alone: `appteam -d ./my-project` runs the wizard with the target directory pre-set
6. Update `--help` output to include the `-d, --dir` option
7. Error and exit non-zero if `-d` is provided with an empty string value

## Current Behavior

- `main.go:31` — Target directory is always CWD (for `-r`) or wizard-prompted
- `wizard.go:31` — `cfg.TargetDir = ask(scanner, w, s, "Target directory", ".")` always prompts
- Flag parsing is a simple `os.Args` switch in `main.go:17-46` — no flag library

## Implementation Notes

- The existing flag parsing uses `os.Args` manually (no `flag` package). The implementation should extend this pattern or migrate to `flag` — SWE's discretion
- The wizard needs a way to know that `-d` was provided so it can skip the prompt. Options: add a `TargetDir` field to a new options struct passed to `Run()`, or pre-set `cfg.TargetDir` and pass a flag
- Directory creation should happen early (in `main.go`) before wizard or regenerate logic runs

## Acceptance Criteria

- [ ] `appteam -d ./new-project` creates the directory if missing and runs the wizard with target dir pre-set (not prompted)
- [ ] `appteam -r -d ./other-dir` regenerates into the specified directory (loads settings.json from `./other-dir`)
- [ ] `appteam --dir ./foo` long-form flag works identically to `-d`
- [ ] `appteam --help` output includes `-d, --dir <folder>` with a description
- [ ] `appteam -d ""` errors with a non-zero exit code
- [ ] Unit tests cover: flag parsing, directory creation, integration with `-r`, wizard prompt skipping
- [ ] Existing tests continue to pass

## Out of Scope

- Changing the wizard flow beyond skipping the target-dir prompt
- Validating directory permissions or disk space
- Supporting `~` expansion (shell handles this before the process receives it)

## Dependencies

- None

## Open Questions

- None
