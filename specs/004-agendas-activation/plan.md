# Implementation Plan: Agendas and Activation Tracking

**Branch**: `004-agendas-activation` | **Date**: 2026-06-08 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/004-agendas-activation/spec.md`

## Summary

This feature adds analyst agenda views, a client activation workflow, and dashboard KPI/funnel analytics. The backend extends the existing Go/Chi/PostgreSQL stack with new columns, a migration, new API endpoints for meeting completion + client activation, and new dashboard aggregation queries. The frontend adds a "Minhas Agendas" page (Vue 3 + Pinia) and expands the Dashboard with funnel, cohort, and KPI widgets.

## Technical Context

**Language/Version**: Go 1.22+, JavaScript (ES2020+)

**Primary Dependencies**: go-chi/chi v5, Pinia, Vue 3, axios

**Storage**: PostgreSQL (existing database with migrations)

**Testing**: Manual validation via quickstart scenarios

**Target Platform**: Web (browser) + Linux/Windows server

**Project Type**: Web application (backend API + frontend SPA)

**Performance Goals**: Dashboard queries under 500ms for typical data volumes (<10k projects)

**Constraints**: Must follow existing migration pattern (sequential `000NNN_*.up.sql`). All datetime fields in ISO 8601/RFC3339.

**Scale/Scope**: Small-to-medium deployment (dozens of analysts, thousands of clients)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Security & RBAC | ✅ PASS | New endpoints protected by existing JWT auth middleware. Analyst-scoped queries use `auth.UserEmailKey` context. |
| II. Referential Integrity | ✅ PASS | New `activated_at` column on `projects` ties to existing FK chain. Meeting `status` column uses CHECK constraint. |
| III. ISO 8601 Dates | ✅ PASS | All new timestamp columns use TIMESTAMPTZ with RFC3339 serialization in Go. |
| IV. Robust External Integrations | ✅ PASS (N/A) | No new external integrations. |
| V. Real-time Consolidated Metrics | ✅ PASS | Funnel, cohort, and KPI metrics derived via live database queries (no arbitrary estimates). |

## Project Structure

### Documentation (this feature)

```text
specs/004-agendas-activation/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
backend/
├── migrations/
│   └── 000002_activation.up.sql       # New migration
├── internal/
│   ├── meeting/
│   │   ├── model.go                   # Add Status, CompletedAt fields
│   │   ├── service.go                 # Add CompleteMeeting, ListByAnalyst
│   │   └── handler.go                 # Add complete + list-my-meetings endpoints
│   ├── project/
│   │   ├── model.go                   # Add ActivatedAt field
│   │   ├── service.go                 # Add ActivateClient, FinalizeProject
│   │   └── handler.go                 # Add activate + finalize endpoints
│   ├── dashboard/
│   │   ├── service.go                 # Add funnel, cohort, KPI queries
│   │   └── handler.go                 # (no handler changes needed)
│   └── api/
│       └── router.go                  # Register new routes
└── tests/                             # (future)

frontend/
├── src/
│   ├── pages/
│   │   ├── MeetingsList.vue           # NEW: "Minhas Agendas" page
│   │   └── Dashboard.vue              # Extend with funnel + cohort + KPI widgets
│   ├── stores/
│   │   └── meetings.js                # Add fetchMyMeetings, completeMeeting actions
│   ├── services/
│   │   └── meeting.js                 # NEW: meeting API service
│   ├── router/
│   │   └── index.js                   # Add /meetings route
│   └── layouts/
│       └── MainLayout.vue             # Add "Agendas" nav item
```

**Structure Decision**: Follows existing web application pattern with `backend/` and `frontend/` separation. New files follow the established Go package conventions (model.go / service.go / handler.go) and Vue conventions (pages + stores + services).
