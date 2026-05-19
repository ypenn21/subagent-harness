# F-0045: Comprehensive Test Coverage

**Type:** Enhancement
**Priority:** P0 (critical)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-16

## Problem

The appteam codebase has grown to 5 packages (config, wizard, templates, generator, main) with 19 template files and zero automated tests. The Reviewer has flagged this gap in sessions 4, 5, and 6. Without tests, regressions go undetected, and refactoring carries high risk. Test coverage is the #1 priority for v0.6.0.

## Requirements

### 1. Config Package Tests (`internal/config/`)

1.1. **Struct initialization** — Verify that `ProjectConfig`, `SWEConfig`, and `GCPConfig` zero values are correct and fields are accessible after assignment.

1.2. **SettingsPath** — `SettingsPath(dir)` must return `<dir>/.appteam/settings.json` for any input directory path.

1.3. **SettingsExist** — Must return `true` when `.appteam/settings.json` exists in the given directory, `false` otherwise.

1.4. **SaveSettings / LoadSettings round-trip** — Save a fully populated `ProjectConfig` to a temp directory, then load it back. All fields must match the original, including nested `GCPConfig` and `[]SWEConfig` slices.

1.5. **SaveSettings creates directory** — Calling `SaveSettings` on a directory that doesn't yet contain `.appteam/` must create the directory automatically.

1.6. **LoadSettings error cases** — Must return an error when the file doesn't exist. Must return an error when the file contains invalid JSON.

1.7. **JSON format** — Saved file must be valid JSON with indentation (2-space indent as produced by `json.MarshalIndent`) and a trailing newline.

### 2. Template Package Tests (`internal/templates/`)

2.1. **Template function `add`** — Verify `add(2, 3)` returns `5`.

2.2. **Template function `sub`** — Verify `sub(5, 3)` returns `2`.

2.3. **Template function `seq`** — Verify `seq(1, 5)` returns `[]int{1, 2, 3, 4, 5}`. Verify `seq(3, 3)` returns `[]int{3}`.

2.4. **Template function `linkRange`** — Verify `linkRange(0, 3)` returns `"0,1,2"`. Verify `linkRange(5, 0)` returns `""`.

2.5. **Render function** — Pass a simple template string with data and verify the output matches expected text. Verify that an invalid template returns an error. Verify that a template referencing a missing field returns an error.

2.6. **ClaudeMDTemplate** — Render with a minimal `ProjectConfig` (GCP disabled) and verify:
  - Output contains `# <ProjectName> — Project Instructions`
  - Output does NOT contain `## GCP Project` section
  - Output contains owner name, email, GitHub username
  - Output contains SWE role listing for each SWE in the config
  - Output contains model name and model ID

2.7. **ClaudeMDTemplate with GCP enabled** — Render with `GCP.Enabled = true` and verify:
  - Output contains `## GCP Project` section
  - Output contains GCP Project ID, Project Number, Organization, Region
  - Output contains `## GCP Free Tier` section

2.8. **ClaudeMDTemplate with optional agents** — Render with all combinations of `IncludeSWETest`, `IncludeSWEQA`, `IncludePlatform`, `IncludeReviewer` and verify the corresponding role descriptions appear/don't appear in output.

2.9. **ClaudeMDTemplate with conventions** — Render with non-empty `Conventions` slice and verify `## Important Conventions` section appears. Render with empty `Conventions` and verify the section is absent.

2.10. **SWETemplate** — Render with `SWETemplateData` containing bullets and verify:
  - Output contains `# SWE-<N> Agent — <Title>`
  - Output contains each bullet
  - Output contains `Your specialty is <Title>`

2.11. **SWETemplate without bullets** — Render with empty `Bullets` and verify:
  - Output contains `You are additional engineering capacity`
  - Output contains default bullets (`General full-stack development`, etc.)

2.12. **All remaining templates render without error** — For each of the following template constants, render with a fully populated `ProjectConfig` and verify no error is returned: `PMTemplate`, `TPMTemplate`, `SWETestTemplate`, `SWEQATemplate`, `PlatformTemplate`, `ReviewerTemplate`, `BacklogTemplate`, `ProgressTemplate`, `ReleaseNotesTemplate`, `PipelineTemplate`, `SpecTemplateFile`, `SkillSpecTemplate`, `SkillReleaseTemplate`, `SkillPipelineTemplate`, `SkillStatusTemplate`, `SkillRegenerateTemplate`.

2.13. **Template output contains expected markers** — For key templates (PM, TPM, SWE-Test, Reviewer, Platform), verify the output contains the agent's role name and project name.

### 3. Generator Package Tests (`internal/generator/`)

3.1. **Directory creation** — Call `Generate` with a temp directory as `TargetDir` and verify:
  - `.claude/agents/` directory is created
  - `.claude/skills/` directory is created
  - `docs/specs/` directory is created

3.2. **Core files generated** — Verify these files are created after generation:
  - `CLAUDE.md`
  - `docs/BACKLOG.md`
  - `docs/PROGRESS.md`
  - `docs/RELEASENOTES.md`
  - `docs/PIPELINE.md`
  - `docs/specs/TEMPLATE.md`
  - `.claude/agents/pm.md`
  - `.claude/agents/tpm.md`

3.3. **SWE agent files** — With `N` SWEs configured, verify `N` files are created: `.claude/agents/swe-1.md` through `.claude/agents/swe-N.md`.

3.4. **Optional agent files** — Verify `.claude/agents/swe-test.md` is created only when `IncludeSWETest = true`. Same for `swe-qa.md` (`IncludeSWEQA`), `platform.md` (`IncludePlatform`), `reviewer.md` (`IncludeReviewer`).

3.5. **Skill files** — Verify these 5 skill files are always created:
  - `.claude/skills/spec.md`
  - `.claude/skills/release.md`
  - `.claude/skills/pipeline.md`
  - `.claude/skills/status.md`
  - `.claude/skills/regenerate.md`

3.6. **Settings persistence** — After generation, `.appteam/settings.json` must exist in the target directory.

3.7. **File content correctness** — At minimum, verify `CLAUDE.md` starts with `# <ProjectName>` and each SWE file starts with `# SWE-<N> Agent`.

3.8. **Note on git/GitHub operations** — Tests for `setupGitRepo` are out of scope for this milestone because they require `git` and `gh` CLI binaries and real filesystem state. The `Generate` function should be tested with `InitGit = false` and `CreateRepo = false` to avoid triggering git operations.

### 4. Wizard Package Tests (`internal/wizard/`)

4.1. **Styler with color disabled** — Create `NewStyler(false)`. Verify `Bold("x")` returns `"x"` (no ANSI codes), same for `Dim`, `Green`, `BoldCyan`, `BoldGreen`, `BoldWhite`.

4.2. **Styler with color enabled** — Create `NewStyler(true)`. Verify `Bold("x")` returns `"\033[1mx\033[0m"`. Verify `Green("x")` returns `"\033[32mx\033[0m"`. Verify `BoldCyan("x")` returns `"\033[1m\033[36mx\033[0m"`.

4.3. **PadBold** — With color disabled, verify `PadBold("Hi:", 10)` returns `"Hi:       "` (padded to 10 chars). With color enabled, verify the result contains ANSI codes wrapping the padded text.

4.4. **Banner** — Verify `Banner()` returns a string containing `"appteam"` and `"Claude Code Agent Team Generator"`. Verify it contains box-drawing characters (`"┌"`, `"└"`).

4.5. **StepHeader** — Verify `StepHeader(2, 7, "Git Repository")` contains `"Step 2 of 7"` and `"Git Repository"`.

4.6. **Divider** — Verify `Divider()` returns a string containing `"─"` characters.

4.7. **ask helper** — Provide a `bufio.Scanner` backed by a `strings.Reader`. Verify:
  - User input `"myproject"` returns `"myproject"`
  - Empty input with default `"fallback"` returns `"fallback"`
  - Input with whitespace is trimmed

4.8. **askBool helper** — Verify:
  - `"y"` returns `true`, `"n"` returns `false`
  - `"yes"` returns `true`, `"no"` returns `false`
  - Empty input with `defaultVal=true` returns `true`
  - Empty input with `defaultVal=false` returns `false`

4.9. **askInt helper** — Verify:
  - `"3"` returns `3`
  - Empty input with `defaultVal=5` returns `5`
  - Non-numeric input `"abc"` returns `defaultVal`

4.10. **Note on `ask`, `askBool`, `askInt` visibility** — These are unexported functions. Tests must be in the `wizard` package (file `internal/wizard/wizard_test.go`) to access them directly. Alternatively, they can be tested indirectly through the `Run` function.

4.11. **Wizard `Run` integration test** — Provide a `strings.Reader` with pre-scripted input (one answer per line, matching the wizard's prompt order) and verify the returned `ProjectConfig` has the correct values. Test with both "proceed" and "cancel" scenarios.

### 5. Main Package Tests (`main_test.go`)

5.1. **`--help` flag** — Build and execute the binary with `--help`. Verify exit code 0 and output contains `"Usage: appteam"`.

5.2. **`-h` flag** — Same as `--help`.

5.3. **`--version` flag** — Execute with `--version`. Verify exit code 0 and output contains `"appteam v"`.

5.4. **`-v` flag** — Same as `--version`.

5.5. **Unknown flag** — Execute with `--bogus`. Verify exit code 1 and stderr contains `"Unknown option"`.

5.6. **`-r` / `--regenerate` without settings** — Execute with `-r` in a temp directory with no `.appteam/settings.json`. Verify exit code 1 and stderr contains `"No settings.json found"`.

5.7. **Note on testing approach** — Main package tests should use `exec.Command` to run the compiled binary (via `go build` + temp binary, or `go run`) to test actual CLI behavior including exit codes.

## Acceptance Criteria

- [ ] All tests pass: `go test ./...` exits 0 with no failures
- [ ] Config package: minimum 6 test functions covering SettingsPath, SettingsExist, SaveSettings, LoadSettings, round-trip, and error cases
- [ ] Template package: minimum 10 test functions covering all template functions (add, sub, seq, linkRange), Render, ClaudeMD (with/without GCP, optional agents, conventions), SWE (with/without bullets), and error-free rendering of all 18+ template constants
- [ ] Generator package: minimum 3 test functions covering directory creation, file generation completeness, and optional agent conditionals
- [ ] Wizard package: minimum 6 test functions covering Styler methods (color on/off), PadBold, Banner, StepHeader, Divider, and ask/askBool/askInt helpers
- [ ] Main package: minimum 4 test functions covering --help, --version, unknown flag, and -r without settings
- [ ] No test relies on external network, git operations, or GitHub CLI
- [ ] Tests use `t.TempDir()` for filesystem operations (automatic cleanup)
- [ ] Test files follow Go convention: `*_test.go` in the same package as the code under test

## Out of Scope

- `setupGitRepo` testing (requires git/gh binaries and real filesystem git state)
- Integration tests that run the full wizard interactively against a terminal
- Benchmark tests
- Fuzz testing
- Coverage percentage enforcement (target is reasonable coverage, not a hard percentage gate)
- CI/CD pipeline setup (future milestone)

## Dependencies

- None (this is a standalone test effort against the existing codebase)

## Open Questions

- None. The codebase is stable and well-understood. Testing can proceed immediately.
