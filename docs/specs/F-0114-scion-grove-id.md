# F-0114: Generate .scion/grove-id for Scion Projects

**Type:** Bug
**Priority:** P0 (critical)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-19

## Problem

When appteam generates a Scion framework project, it creates `.scion/templates/` but does NOT create `.scion/grove-id`. Without grove-id:
- `scion init` refuses to run ("already inside a scion project" because `.scion/` exists)
- `scion list`, `scion start`, etc. can't recognize the grove
- Users must manually run `uuidgen > .scion/grove-id` as a workaround

The generator also does not create `.scion/agents/` directory, which is needed at runtime.

## Requirements

1. When `Framework == "scion"`, the generator must create `.scion/grove-id` containing a UUID (v4 format)
2. The generator must create the `.scion/agents/` directory (empty, with a `.gitkeep` if needed)
3. UUID generation must use Go stdlib only (no external dependencies) -- use `crypto/rand` to generate a v4 UUID
4. The grove-id file must contain only the UUID string (no trailing newline required but acceptable)

## Acceptance Criteria

- [ ] `.scion/grove-id` is created when Framework == "scion"
- [ ] The file contains a valid v4 UUID (8-4-4-4-12 hex format)
- [ ] `.scion/agents/` directory is created
- [ ] `.scion/grove-id` is NOT created when Framework != "scion"
- [ ] UUID is generated using `crypto/rand` (not `math/rand`)
- [ ] Existing tests still pass
- [ ] New tests verify grove-id generation and UUID format

## Out of Scope

- Modifying scion template content
- Changes to the wizard flow
- Any external dependency additions

## Dependencies

- None (builds on existing Scion support from F-0096)

## Open Questions

- None
