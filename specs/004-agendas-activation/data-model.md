# Data Model: Agendas and Activation Tracking

## Schema Changes

### Migration: `000002_activation.up.sql`

This migration modifies the existing `meetings` and `projects` tables, and adds indexes to support the new features.

---

### 1. Meetings Table — Modified

**New columns**:

| Column | Type | Default | Nullable | Description |
|--------|------|---------|----------|-------------|
| `status` | `VARCHAR(20)` | `'scheduled'` | NOT NULL | Meeting lifecycle: `scheduled`, `completed`, `cancelled` |
| `completed_at` | `TIMESTAMPTZ` | — | YES | Timestamp when the meeting was marked as completed |

**CHECK Constraint**: `status IN ('scheduled', 'completed', 'cancelled')`

**Backfill rule**: Existing meetings get `status = 'scheduled'` (the default). Rows where `no_show = TRUE` remain `scheduled` with their `no_show` flag — they represent meetings that were scheduled but nobody attended.

**Index**: `idx_meetings_analyst_status` on `(analyst_id, status)` — supports the "My Meetings" query filtered by analyst and status.

**Notes**: The `no_show` column is preserved for backward compatibility. Over time, `no_show` may be deprecated in favor of `status = 'cancelled'` but that's out of scope.

---

### 2. Projects Table — Modified

**New columns**:

| Column | Type | Default | Nullable | Description |
|--------|------|---------|----------|-------------|
| `activated_at` | `TIMESTAMPTZ` | — | YES | Timestamp when the client was marked as active during a meeting. NULL = not yet activated. |

**Index**: `idx_projects_activated_at` on `(activated_at)` — supports cohort and funnel queries.

---

## Entity Relationships (Updated)

```text
clients (1) ──< projects (N)
   │                │
   │                └──< meetings (N)
   │                        │
   created_at               analyst_id ──> users (1)
   (cohort key)
```

### Entity: Meeting (Updated)

| Field | Type | Constraints |
|-------|------|-------------|
| id | UUID | PK, auto-generated |
| project_id | UUID | FK → projects.id, NOT NULL |
| analyst_id | UUID | FK → users.id, NULLABLE |
| title | VARCHAR(255) | NOT NULL |
| scheduled_at | TIMESTAMPTZ | NOT NULL |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'scheduled', CHECK IN ('scheduled', 'completed', 'cancelled') |
| completed_at | TIMESTAMPTZ | NULLABLE |
| no_show | BOOLEAN | NOT NULL, DEFAULT FALSE |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**State Transitions**:
- `scheduled` → `completed` (via "Mark Active + Complete" or simple "Complete Meeting")
- `scheduled` → `cancelled` (future use)
- Completed/cancelled meetings cannot revert to scheduled.

### Entity: Project (Updated)

| Field | Type | Constraints |
|-------|------|-------------|
| id | UUID | PK, auto-generated |
| client_id | UUID | FK → clients.id, NOT NULL |
| name | VARCHAR(255) | NOT NULL |
| status | VARCHAR(50) | NOT NULL, DEFAULT 'Backlog', CHECK IN ('Backlog', 'Em andamento', 'Go-Live') |
| is_active | BOOLEAN | NOT NULL, DEFAULT TRUE |
| activated_at | TIMESTAMPTZ | NULLABLE |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Activation semantics**: `activated_at` is orthogonal to `status`. A project can be "activated" (client did the tasks) while still being "Em andamento". Finalization (status → Go-Live) is a separate manual action.

### Entity: Client (Unchanged)

| Field | Type | Constraints |
|-------|------|-------------|
| id | UUID | PK, auto-generated |
| name | VARCHAR(255) | NOT NULL |
| cnpj | VARCHAR(18) | UNIQUE, NOT NULL |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Cohort derivation**: `TO_CHAR(clients.created_at, 'YYYY-MM')` = purchase month (safra).

## Dashboard Aggregation Queries (Conceptual)

### Funnel Metrics

```sql
-- Registered: all projects
SELECT COUNT(*) FROM projects;

-- Participants: projects with at least 1 completed meeting
SELECT COUNT(DISTINCT m.project_id) FROM meetings m WHERE m.status = 'completed';

-- Active: projects where client was activated
SELECT COUNT(*) FROM projects WHERE activated_at IS NOT NULL;
```

### Cohort (Safra) Report

```sql
SELECT
  TO_CHAR(c.created_at, 'YYYY-MM') AS cohort,
  COUNT(p.id) AS total,
  COUNT(p.activated_at) AS activated
FROM projects p
JOIN clients c ON c.id = p.client_id
GROUP BY cohort
ORDER BY cohort DESC;
```

### 30-Day Activation Rate

```sql
SELECT
  COUNT(*) FILTER (WHERE p.activated_at <= c.created_at + INTERVAL '30 days') AS activated_in_30d,
  COUNT(*) FILTER (WHERE EXISTS (
    SELECT 1 FROM meetings m WHERE m.project_id = p.id AND m.status = 'completed'
  )) AS participants
FROM projects p
JOIN clients c ON c.id = p.client_id;
-- Rate = activated_in_30d / participants * 100
```

### Abandonment Rate

```sql
SELECT
  COUNT(*) FILTER (WHERE NOT EXISTS (
    SELECT 1 FROM meetings m WHERE m.project_id = p.id AND m.status = 'completed'
  )) AS abandoned,
  COUNT(*) AS total
FROM projects p;
-- Rate = abandoned / total * 100
```

### First-Meeting Activation Rate

```sql
WITH first_meetings AS (
  SELECT DISTINCT ON (m.project_id)
    m.project_id, m.id AS meeting_id, m.completed_at
  FROM meetings m
  WHERE m.status = 'completed'
  ORDER BY m.project_id, m.scheduled_at ASC
)
SELECT
  COUNT(*) FILTER (WHERE p.activated_at IS NOT NULL AND p.activated_at <= fm.completed_at + INTERVAL '1 minute') AS first_meeting_activated,
  COUNT(*) AS total_participants
FROM first_meetings fm
JOIN projects p ON p.id = fm.project_id;
-- Rate = first_meeting_activated / total_participants * 100
```
