# F-0040: Settings JSON — Non-Interactive Mode

**Type:** Feature
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-16

## Problem

Currently, every run of `appteam` requires completing the full interactive wizard. For projects that have already been configured, this is tedious and error-prone — users must re-enter identical values each time they want to regenerate files. There is no way to script or automate `appteam` in CI/CD pipelines.

A saved settings file would let users regenerate agent configurations instantly without re-answering prompts, and enable non-interactive usage in automated workflows.

## Requirements

1. **Save on wizard completion** — After the wizard finishes and files are generated, the system must serialize the full `ProjectConfig` struct as JSON and write it to `.appteam/settings.json` inside the target directory
2. **JSON structure** — The settings file must be a well-formatted (indented) JSON file that faithfully represents all `ProjectConfig` fields including nested structs (`GCPConfig`, `SWEConfig` slices, `Conventions` slice)
3. **Load on startup** — When `appteam` is run in a directory that contains `.appteam/settings.json`, the system must detect it and offer the user a choice: reuse the saved config (skip wizard) or start fresh (run wizard, overwrite settings)
4. **`--regenerate` / `-r` flag** — A new CLI flag that reads `.appteam/settings.json` from the current directory and regenerates all files without any interactive prompts. The flag must bypass the wizard entirely
5. **Error handling for `-r`** — If `-r` is used but no `.appteam/settings.json` exists, print a clear error message (`No settings.json found. Run appteam first to create one.`) and exit with code 1
6. **Roundtrip fidelity** — Loading a saved `settings.json` and regenerating must produce byte-identical output to the original wizard-driven generation (same config values = same files)
7. **Settings committed to repo** — The `.appteam/settings.json` file is intended to be committed to version control so team members can regenerate without re-entering config
8. **Update `--help` output** — Add `-r, --regenerate` to the help text with description `Regenerate from saved .appteam/settings.json`
9. **`.appteam/` directory** — Create the `.appteam/` directory in the target dir if it doesn't exist. This directory will also be used by future features (skills in F-0041)

## Acceptance Criteria

- [ ] Running the wizard saves `.appteam/settings.json` in the target directory after file generation
- [ ] The JSON file contains all `ProjectConfig` fields with correct values (including SWEs, GCP, Conventions)
- [ ] Running `appteam` in a directory with existing `settings.json` prompts the user to reuse or start fresh
- [ ] Choosing "reuse" skips the wizard and regenerates using saved config
- [ ] Choosing "start fresh" runs the wizard normally and overwrites `settings.json`
- [ ] `appteam -r` reads `settings.json` from `.appteam/settings.json` in the current directory and regenerates without prompts
- [ ] `appteam -r` in a directory without `settings.json` prints an error and exits with code 1
- [ ] `--help` output includes the `-r` / `--regenerate` flag
- [ ] Roundtrip test: wizard config saved and reloaded produces identical generated files
- [ ] JSON file is human-readable (indented with 2 spaces)

## Out of Scope

- Schema migration / versioning of settings.json (future)
- Partial config overrides via CLI flags (e.g., `--project-name=foo -r`)
- Merging settings with existing generated files (conflict resolution)

## Dependencies

- None (builds on existing wizard and generator)

## Open Questions

- None
