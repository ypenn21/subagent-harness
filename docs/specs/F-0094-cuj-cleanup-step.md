# F-0094: Add Cleanup Phase to /cuj-test Skill Template

**Type:** Enhancement
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-17

## Problem

The /cuj-test skill and SWE-QA agent template don't include a cleanup step after test execution. Test artifacts (screenshots, browser instances, temporary test data, console logs) accumulate across runs and need to be cleaned up. Without an explicit cleanup phase, agents may leave browser processes running, fill up disk with screenshots, and leave test data polluting the application state.

## Requirements

1. Add a "Phase 6: Cleanup" to the `/cuj-test` skill template in `internal/templates/skill_cuj_testing.go` that:
   - Closes browser instances (`await browser.close()`)
   - Optionally removes screenshots directory if user wants clean runs
   - Cleans up any temporary test data created during tests (test users, test records, etc.)
   - Resets application state if tests modified it
   - Reports what was cleaned up
2. Update the SWE-QA agent template in `internal/templates/swe_qa.go` to mention cleanup as part of the testing workflow
3. Bump version to 0.13.1 in `main.go`
4. Verify all existing tests still pass

## Acceptance Criteria

- [ ] `/cuj-test` template (`SkillCUJTestTemplate`) contains a Phase 6: Cleanup section with browser close, artifact cleanup, state reset, and cleanup report
- [ ] SWE-QA template (`SWEQATemplate`) mentions cleanup in its testing workflow
- [ ] Version constant in `main.go` is "0.13.1"
- [ ] All 148 existing tests pass

## Out of Scope

- Adding new test cases for cleanup (covered by existing template rendering tests)
- Changing the /cuj-list skill template
- Modifying wizard or generator logic

## Dependencies

- F-0090 (/cuj-test skill template — already done in v0.13.0)

## Open Questions

- None
