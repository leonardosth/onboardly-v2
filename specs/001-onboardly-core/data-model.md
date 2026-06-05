# Data Model: Onboardly Core

This document outlines the database schema, entity structures, and constraints for the Onboardly platform.

## Relational Schema (PostgreSQL)

### 1. `users` Table
Represents platform operators.
- `id`: `UUID` (Primary Key, Defaults to `gen_random_uuid()`)
- `email`: `VARCHAR(255)` (Unique, Not Null, Indexed)
- `password_hash`: `VARCHAR(255)` (Not Null)
- `role`: `VARCHAR(50)` (Not Null, Constraint: `Admin` or `Analista`)
- `created_at`: `TIMESTAMPTZ` (Not Null, Defaults to `NOW()`)

### 2. `clients` Table
Represents customer organizations.
- `id`: `UUID` (Primary Key, Defaults to `gen_random_uuid()`)
- `name`: `VARCHAR(255)` (Not Null)
- `cnpj`: `VARCHAR(18)` (Unique, Not Null, Indexed)
- `created_at`: `TIMESTAMPTZ` (Not Null, Defaults to `NOW()`)
- `updated_at`: `TIMESTAMPTZ` (Not Null, Defaults to `NOW()`)

### 3. `projects` Table
Represents onboarding/deployment projects linked to a client.
- `id`: `UUID` (Primary Key, Defaults to `gen_random_uuid()`)
- `client_id`: `UUID` (Foreign Key referencing `clients(id)` ON DELETE CASCADE, Not Null)
- `name`: `VARCHAR(255)` (Not Null)
- `status`: `VARCHAR(50)` (Not Null, Default: `Backlog`. Predefined: `Backlog`, `Em andamento`, `Go-Live`)
- `is_active`: `BOOLEAN` (Not Null, Default: `true`. Derived from status: `true` if `status != 'Go-Live'`, otherwise `false`)
- `created_at`: `TIMESTAMPTZ` (Not Null, Defaults to `NOW()`)
- `updated_at`: `TIMESTAMPTZ` (Not Null, Defaults to `NOW()`)

### 4. `meetings` Table
Represents scheduling interactions associated with a project.
- `id`: `UUID` (Primary Key, Defaults to `gen_random_uuid()`)
- `project_id`: `UUID` (Foreign Key referencing `projects(id)` ON DELETE CASCADE, Not Null)
- `analyst_id`: `UUID` (Foreign Key referencing `users(id)` ON DELETE SET NULL, Not Null)
- `title`: `VARCHAR(255)` (Not Null)
- `scheduled_at`: `TIMESTAMPTZ` (Not Null)
- `no_show`: `BOOLEAN` (Not Null, Default: `false`)
- `created_at`: `TIMESTAMPTZ` (Not Null, Defaults to `NOW()`)

### 5. `activities` Table
Log of recent updates to populate the dashboard activity feed.
- `id`: `UUID` (Primary Key, Defaults to `gen_random_uuid()`)
- `entity_type`: `VARCHAR(50)` (Not Null, e.g., `Client`, `Project`, `Meeting`)
- `entity_id`: `UUID` (Not Null)
- `description`: `TEXT` (Not Null)
- `created_at`: `TIMESTAMPTZ` (Not Null, Defaults to `NOW()`)

---

## Validation Rules & State Transitions

### Validation Constraints
- **CNPJ**: Must match standard Brazilian format (14 digits, formatted as `XX.XXX.XXX/XXXX-XX`).
- **Email**: Must contain a valid email format.
- **Scheduled Time**: Must be parsed and stored in UTC conforming to ISO 8601/RFC3339.

### Project Status Lifecycle
```
[Backlog] --------> [Em andamento] --------> [Go-Live]
(is_active = true)  (is_active = true)      (is_active = false)
```
- Changing a project status to `Go-Live` automatically sets `is_active` to `false` and triggers recalculation of Dashboard metrics.
