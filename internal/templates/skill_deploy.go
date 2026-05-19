package templates

const SkillDeployTemplate = `---
name: deploy
description: Run a deployment checklist for staging or production
user-invocable: true
---

# /deploy — Deployment checklist

## Trigger

User invokes ` + "`/deploy`" + ` with a target environment (e.g., staging, production) or deployment context.

## Instructions

1. **Pre-deploy checks:**
   - All tests pass: ` + "`go test ./...`" + `
   - No uncommitted changes: ` + "`git status`" + ` is clean
   - On the correct branch (main for production, develop for staging)
   - Version/tag is set correctly
   - Dependencies are up to date
2. **Review what is being deployed:**
   - List changes since last deployment: ` + "`git log --oneline <last-tag>..HEAD`" + `
   - Check for breaking changes, migrations, or config updates
   - Verify environment variables and secrets are configured
3. **Run the deployment:**
   - Execute the deploy command for the target environment
   - Monitor deployment progress and logs
   - Watch for errors or warnings during rollout
4. **Post-deploy verification:**
   - Health check endpoints return 200
   - Smoke tests pass (critical user flows work)
   - Logs show no unexpected errors
   - Metrics and monitoring dashboards look normal
5. **Document the rollback plan:**
   - Command to revert to the previous version
   - Database migration rollback steps (if applicable)
   - Feature flags to disable (if applicable)
   - Who to contact if rollback is needed
6. **Report back** — Confirm: what was deployed, to which environment, verification results, and the rollback plan

## Project Context

- **Project:** {{.ProjectName}}
- **Owner:** {{.OwnerName}} ({{.OwnerEmail}})
`
