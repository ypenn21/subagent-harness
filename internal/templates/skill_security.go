package templates

const SkillSecurityTemplate = `---
name: security
description: Perform a security audit against OWASP top 10 and project-specific risks
user-invocable: true
---

# /security — Security audit

## Trigger

User invokes ` + "`/security`" + ` with a target file, module, or the entire project to audit.

## Instructions

1. **Scan for injection vulnerabilities:**
   - SQL injection — parameterized queries vs string concatenation
   - Command injection — ` + "`exec.Command`" + ` with unsanitized input
   - XSS — user input rendered in HTML without escaping
   - Path traversal — file paths constructed from user input
2. **Check authentication and authorization:**
   - Authentication is required on protected endpoints
   - Authorization checks enforce least privilege
   - Session management is secure (expiry, rotation, invalidation)
   - Password hashing uses bcrypt/scrypt/argon2 (not MD5/SHA)
3. **Check for secrets exposure:**
   - No hardcoded API keys, passwords, or tokens in source code
   - ` + "`.env`" + ` and credential files are in ` + "`.gitignore`" + `
   - Secrets are loaded from environment variables or a secrets manager
   - Git history does not contain accidentally committed secrets
4. **Review dependency vulnerabilities:**
   - Check for known CVEs in dependencies
   - Verify dependencies are from trusted sources
   - Review ` + "`go.sum`" + ` for unexpected changes
5. **Check HTTP security headers:**
   - CORS policy is restrictive (not ` + "`*`" + `)
   - CSP (Content-Security-Policy) is configured
   - HSTS, X-Frame-Options, X-Content-Type-Options are set
   - Cookies use Secure, HttpOnly, SameSite attributes
6. **Output findings** — Provide a structured report:
   - **Critical** — Actively exploitable vulnerabilities
   - **High** — Vulnerabilities requiring immediate remediation
   - **Medium** — Issues to address in the next release
   - **Low** — Informational or hardening recommendations
   - Each finding includes: description, location, and remediation steps

## Project Context

- **Project:** {{.ProjectName}}
- **Owner:** {{.OwnerName}} ({{.OwnerEmail}})
`
