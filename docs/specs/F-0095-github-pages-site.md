# F-0095: GitHub Pages Site with Docsify

**Type:** Feature
**Priority:** P1 (important)
**Status:** Approved
**Requested by:** PO
**Date:** 2026-03-18

## Problem

The PO needs a public webpage to present appteam to teammates -- covering why, what, how, and roadmap. Must be free, easily updatable, and shareable via link. Currently there is no public-facing website for the project.

## Requirements

1. Create `docs/index.html` -- Docsify SPA with clean, modern dark theme
2. Create `docs/.nojekyll` -- prevent Jekyll processing on GitHub Pages
3. Create `docs/_sidebar.md` -- navigation for presentation pages
4. Create `docs/home.md` -- landing page / hero with appteam overview and value prop
5. Create `docs/why.md` -- why appteam exists (problem it solves, pain points of managing AI agent teams manually)
6. Create `docs/what.md` -- what appteam does (features, generated files, team structure, skills system)
7. Create `docs/how.md` -- how it works (install, wizard walkthrough, regeneration, team sizing)
8. Create `docs/project-roadmap.md` -- roadmap (completed milestones, upcoming features)
9. Configure GitHub Pages to serve from `docs/` on main branch
10. Clean, simple, modern design -- dark theme, good typography, minimal

## Acceptance Criteria

- [ ] Site renders locally with `npx docsify-cli serve docs`
- [ ] All 4 presentation pages (why, what, how, roadmap) load and render correctly
- [ ] Landing page (home.md) displays hero with overview and value prop
- [ ] Sidebar navigation works across all pages
- [ ] Mobile responsive layout
- [ ] Existing docs/ files (BACKLOG.md, PROGRESS.md, specs/, guide/) are not affected
- [ ] Dark theme applied consistently

## Out of Scope

- Custom domain
- Analytics / tracking
- CI/CD pipeline for docs
- Search functionality beyond Docsify defaults

## Dependencies

- None

## Open Questions

- None
