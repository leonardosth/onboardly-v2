# Research: Agendas and Activation Tracking

## Decision Log

### 1. Meeting Status Model

**Decision**: Add a `status` column to the `meetings` table with values `scheduled` | `completed` | `cancelled` instead of using only the existing `no_show` boolean.

**Rationale**: The existing `no_show` boolean is insufficient to represent the full meeting lifecycle. We need `completed` to track participantes (clients who attended at least one meeting). A status enum provides clearer semantics and extensibility.

**Alternatives considered**:
- Using `no_show = FALSE` as "completed": Ambiguous — a scheduled meeting that hasn't happened yet also has `no_show = FALSE`. Cannot distinguish "completed" from "scheduled".
- Adding a `completed_at` TIMESTAMPTZ nullable column: Works but makes queries more complex (`WHERE completed_at IS NOT NULL`). A status column is more readable and extensible.

### 2. Client Activation Storage

**Decision**: Add `activated_at TIMESTAMPTZ` to the `projects` table (nullable). A project is considered "activated" when `activated_at IS NOT NULL`.

**Rationale**: Activation is per project (implementation), not per client. Using a timestamp rather than a boolean enables time-based KPI calculations (e.g., "activated within 30 calendar days of `clients.created_at`"). The project already has `is_active` and a status flow; `activated_at` is orthogonal — it marks when the client was confirmed active during a meeting, not the project's workflow status.

**Alternatives considered**:
- Boolean `client_activated` on projects: Loses the time dimension needed for KPI calculation.
- Separate `activations` table: Over-engineering for a single-event recording. Can be refactored later if activation history becomes needed.

### 3. Funnel Definitions (Derived from Spec)

**Decision**: The funnel stages are computed dynamically via SQL, not stored:

| Stage | Definition | SQL approximation |
|-------|------------|-------------------|
| Inscritos (Registered) | Projects that exist | `COUNT(*) FROM projects` |
| Participantes (Participants) | Projects with at least 1 meeting with status `completed` | `COUNT(DISTINCT project_id) FROM meetings WHERE status = 'completed'` |
| Ativos (Active) | Projects where `activated_at IS NOT NULL` | `COUNT(*) FROM projects WHERE activated_at IS NOT NULL` |

**Rationale**: These are simple aggregate queries. No materialized views or caching needed at current scale (<10k projects).

### 4. Cohort (Safra) Derivation

**Decision**: The cohort is determined by `clients.created_at` truncated to month (`YYYY-MM`). The spec confirms that contract date = purchase date = `created_at`.

**Rationale**: The client's `created_at` timestamp already serves as the contract/purchase date. No new column is needed. Cohort grouping uses `TO_CHAR(clients.created_at, 'YYYY-MM')`.

### 5. "Activate + Complete" Combined Action

**Decision**: The API endpoint `POST /api/meetings/{id}/complete` accepts an optional `activate_client` boolean. When true, it atomically:
1. Sets the meeting status to `completed` and records `completed_at`
2. Sets `projects.activated_at = NOW()` (if not already set)

Project finalization remains a separate `PUT /api/projects/{id}` call with `status: "Go-Live"`.

**Rationale**: Matches the clarified spec (C: "Mark active + complete meeting is one step, but project finalization is separate"). A single endpoint reduces network round-trips and ensures atomicity.

### 6. "First Meeting" Determination

**Decision**: A client's "first meeting" for a project is the meeting with the earliest `scheduled_at` that has `status = 'completed'`. If the client was activated during that specific meeting, they count as "activated on first meeting".

**Rationale**: Using `scheduled_at` ordering (not `completed_at`) reflects the intended chronological meeting order. This can be computed with a window function `ROW_NUMBER() OVER (PARTITION BY project_id ORDER BY scheduled_at)`.

### 7. Abandonment Calculation

**Decision**: Abandoned = projects with zero meetings having `status = 'completed'`. The formula: `abandoned_count / total_projects * 100`. Target: ≤ 20%.

**Rationale**: Matches the spec definition: "clientes que se inscreveram na implantação mas não participaram de nenhuma reunião". A project without any completed meeting is considered abandoned.
