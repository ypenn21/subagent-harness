# /schema — Database schema design

## Trigger

User invokes `/schema` with an entity, feature, or data model to design.

## Instructions

1. **Define entities and relationships:**
   - List all tables/collections needed
   - Identify primary keys (prefer UUIDs or auto-increment IDs)
   - Map relationships: one-to-one, one-to-many, many-to-many
   - Draw out the entity-relationship model
2. **Design the schema** for each table:
   - Column name, data type, and constraints
   - NOT NULL, UNIQUE, DEFAULT values
   - Foreign key references with ON DELETE/ON UPDATE behavior
   - Check constraints where appropriate
3. **Create a migration file:**
   - Use sequential or timestamped naming (e.g., `001_create_users.sql`)
   - Include both UP (create) and DOWN (rollback) migrations
   - Wrap in transactions where supported
4. **Add indexes** for common query patterns:
   - Foreign key columns (for JOIN performance)
   - Columns used in WHERE clauses and ORDER BY
   - Composite indexes for multi-column queries
   - Unique indexes for business-rule uniqueness
5. **Document constraints and decisions:**
   - Why each data type was chosen
   - Nullability rationale
   - Cascading behavior (CASCADE, SET NULL, RESTRICT)
   - Any denormalization trade-offs
6. **Report back** — Summarize the schema design: tables created, relationships, indexes, and any open questions

## Project Context

- **Project:** appteam
- **Owner:** Ameer Abbas (ameer00@gmail.com)
