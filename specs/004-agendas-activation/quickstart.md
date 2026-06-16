# Quickstart: Agendas and Activation Tracking

## Prerequisites

- Backend running: `cd backend && go run cmd/server/main.go`
- Frontend running: `cd frontend && npm run dev`
- PostgreSQL accessible with migration `000002_activation.up.sql` applied
- At least one user registered (Analista role)

## Scenario 1: Apply Migration

```bash
# From backend/ directory
# Migration file: migrations/000002_activation.up.sql
# Apply using psql or your migration tool
psql -U <user> -d <database> -f migrations/000002_activation.up.sql
```

**Expected**: Tables `meetings` and `projects` updated with new columns. No data loss.

## Scenario 2: Complete a Meeting and Activate Client

### Step 1: Create a client and project (if not exists)

```bash
# Create client
curl -X POST http://localhost:8080/api/clients \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Empresa Teste", "cnpj": "12.345.678/0001-99"}'

# Create project for client
curl -X POST http://localhost:8080/api/projects \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"client_id": "<client_id>", "name": "Implantação Teste"}'
```

### Step 2: Schedule a meeting

```bash
curl -X POST http://localhost:8080/api/meetings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"project_id": "<project_id>", "title": "Reunião de Implantação", "scheduled_at": "2026-06-10T14:00:00Z"}'
```

### Step 3: Complete the meeting and activate the client

```bash
curl -X POST http://localhost:8080/api/meetings/<meeting_id>/complete \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"activate_client": true}'
```

**Expected**: Meeting status changes to `completed`, `completed_at` is set, and `projects.activated_at` is set.

### Step 4: Finalize the project (separate action)

```bash
curl -X POST http://localhost:8080/api/projects/<project_id>/finalize \
  -H "Authorization: Bearer <token>"
```

**Expected**: Project status changes to `Go-Live`, `is_active` set to `false`.

## Scenario 3: View "My Meetings" (Analyst Agenda)

```bash
curl http://localhost:8080/api/meetings/mine?status=completed \
  -H "Authorization: Bearer <token>"
```

**Expected**: Returns a list of completed meetings assigned to the authenticated analyst, with `project_name` and `client_name` included.

## Scenario 4: Dashboard Funnel and KPIs

```bash
curl http://localhost:8080/api/dashboard \
  -H "Authorization: Bearer <token>"
```

**Expected**: JSON response includes:
- `funnel.registered`, `funnel.participants`, `funnel.active` counts
- `metrics.abandonment_rate` (target: ≤ 20%)
- `metrics.activation_30d_rate` (target: ≥ 80%)
- `metrics.first_meeting_activation_rate`
- `cohorts[]` array with per-month activation data

## Scenario 5: Frontend Validation

1. **Login** as an Analyst
2. **Navigate to "Agendas"** in the sidebar → Verify the "Minhas Agendas" page displays completed meetings
3. **Navigate to Dashboard** → Verify the funnel widget shows Inscritos > Participantes > Ativos
4. **Verify cohort table** shows monthly data with activation rates
5. **Verify KPI cards** display abandonment rate, 30-day activation rate, and first-meeting activation rate

## Data Setup for Full KPI Testing

To fully validate all KPIs, seed data with:
- Multiple clients with `created_at` in different months (for cohort testing)
- Projects with varying activation states (some activated, some not)
- Meetings with some completed and some not (for abandonment testing)
- At least one project activated within 30 days and one after 30 days
- At least one project activated on first meeting

See `backend/cmd/seed/main.go` for the existing seed script — extend it with the new fields.
