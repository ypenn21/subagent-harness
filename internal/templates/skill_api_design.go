package templates

const SkillAPIDesignTemplate = `---
name: api-design
description: "Design an API endpoint with request/response schemas and validation"
user-invocable: true
---

# /api-design — Design API endpoints

## Trigger

User invokes ` + "`/api-design`" + ` with a resource name or feature that needs an API.

## Instructions

1. **Define the resource** — Identify the core entity (e.g., "users", "orders", "projects"). Determine its attributes, relationships, and lifecycle
2. **Design routes** following REST conventions:
   - ` + "`GET /resources`" + ` — List (with pagination, filtering, sorting)
   - ` + "`POST /resources`" + ` — Create
   - ` + "`GET /resources/:id`" + ` — Read
   - ` + "`PUT /resources/:id`" + ` — Update (full replace)
   - ` + "`PATCH /resources/:id`" + ` — Partial update
   - ` + "`DELETE /resources/:id`" + ` — Delete
   - Nested resources for relationships (e.g., ` + "`/users/:id/orders`" + `)
3. **Specify request/response schemas** in JSON:
   - Request body with required and optional fields
   - Response body with all fields and types
   - Use consistent naming (camelCase or snake_case — pick one)
4. **Define error codes:**
   - 400 Bad Request — validation errors (include field-level details)
   - 401 Unauthorized — missing or invalid authentication
   - 403 Forbidden — insufficient permissions
   - 404 Not Found — resource does not exist
   - 409 Conflict — duplicate or state conflict
   - 422 Unprocessable Entity — semantically invalid
   - 500 Internal Server Error — unexpected failures
5. **Document authentication and authorization:**
   - Which endpoints require authentication
   - What roles or permissions are needed for each endpoint
   - Token format and header (e.g., ` + "`Authorization: Bearer <token>`" + `)
6. **Output** — Produce an API spec document or code scaffolding with route definitions, handler stubs, and request/response types

## Project Context

- **Project:** {{.ProjectName}}
- **Owner:** {{.OwnerName}} ({{.OwnerEmail}})
`
