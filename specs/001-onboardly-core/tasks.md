# Tasks: Onboardly Core Requirements

**Input**: Design documents from `/specs/001-onboardly-core/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Standard unit and integration testing tasks are not explicitly requested by the specifications and are excluded from the core user stories to focus on implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `backend/src/` (or `backend/`), `frontend/src/`
- Paths shown below reflect this layout.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure.

- [x] T001 Initialize Go backend module structure in backend/go.mod
- [x] T002 Initialize Vue.js frontend project structure in frontend/package.json
- [x] T003 [P] Configure environment config structure in backend/internal/config/config.go
- [x] T004 [P] Configure formatting and lint rules in backend/.golangci.yml and frontend/eslint.config.js

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core database connections, migrations, and shared routing infrastructure that must be complete before any user story can begin.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T005 Initialize database connection pool in backend/internal/db/db.go using pgx
- [x] T006 Create database migration SQL scripts in backend/migrations/000001_init.up.sql
- [x] T007 Setup core HTTP router and middleware in backend/internal/api/router.go
- [x] T008 Implement default API error response payload models in backend/internal/api/errors.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel.

---

## Phase 3: User Story 1 - User Authentication & Access Control (Priority: P1) 🎯 MVP

**Goal**: Secure user signup, login, JWT session generation, and role-based access control (RBAC) middleware.

**Independent Test**: Register a new user, log in, and verify that accessing unauthorized routes returns 401/403.

### Implementation for User Story 1
- [x] T009 [P] [US1] Create User DB model and mapping logic in backend/internal/auth/model.go
- [x] T010 [US1] Implement user password hashing and token generation service in backend/internal/auth/service.go
- [x] T011 [US1] Implement user registration and login endpoints in backend/internal/auth/handler.go
- [x] T012 [US1] Create role-based access control (RBAC) middleware in backend/internal/auth/middleware.go
- [x] T013 [US1] Implement Vue frontend auth store in frontend/src/stores/auth.js
- [x] T014 [US1] Create login and registration views in frontend/src/pages/Login.vue and frontend/src/pages/Register.vue
- [x] T015 [US1] Implement Vue router navigation guards in frontend/src/router/index.js

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently.

---

## Phase 4: User Story 2 - Client Management & Salesforce / Google Sheets Integration (Priority: P1)

**Goal**: CRUD operations for clients, and webhook receiver to ingest client data from Google Sheets/Salesforce with shared secret token authorization.

**Independent Test**: Send client data via webhook to `/api/webhooks/clients`, verify it creates a client, and check client validation.

### Implementation for User Story 2
- [x] T016 [P] [US2] Create Client DB model and CNPJ validation rules in backend/internal/client/model.go
- [x] T017 [US2] Implement Client CRUD database query services in backend/internal/client/service.go
- [x] T018 [US2] Implement Client HTTP handlers in backend/internal/client/handler.go
- [x] T019 [US2] Implement external Google Sheets webhook sync receiver handler in backend/internal/client/webhook_handler.go
- [x] T020 [US2] Create Vue client API service and store in frontend/src/services/client.js and frontend/src/stores/clients.js
- [x] T021 [US2] Create Client list and detail components in frontend/src/pages/ClientsList.vue and frontend/src/pages/ClientDetail.vue

**Checkpoint**: At this point, User Stories 1 and 2 should both work independently.

---

## Phase 5: User Story 3 - Projects & Meeting Management (Priority: P2)

**Goal**: Create deployment projects linked to clients, update status, calculate activation, schedule meetings, and enforce referential integrity.

**Independent Test**: Create project, verify default "Backlog" status and `is_active` true. Schedule a meeting, verify referential integrity blocks invalid project IDs, and date uses ISO 8601.

### Implementation for User Story 3
- [x] T022 [P] [US3] Create Project and Meeting DB models in backend/internal/project/model.go and backend/internal/meeting/model.go
- [x] T023 [US3] Implement Project lifecycle service and status transition updates in backend/internal/project/service.go
- [x] T024 [US3] Implement Meeting scheduler logic and project checking in backend/internal/meeting/service.go
- [x] T025 [US3] Implement Project and Meeting HTTP endpoint handlers in backend/internal/project/handler.go and backend/internal/meeting/handler.go
- [x] T026 [US3] Create Vue projects and meetings Pinia stores in frontend/src/stores/projects.js and frontend/src/stores/meetings.js
- [x] T027 [US3] Create Project detail and status progression view in frontend/src/pages/ProjectDetail.vue
- [x] T028 [US3] Implement Meeting scheduler popup modal component in frontend/src/components/MeetingScheduler.vue

**Checkpoint**: At this point, User Stories 1, 2, and 3 should work.

---

## Phase 6: User Story 4 - Dashboard & Analytics BI (Priority: P3)

**Goal**: Display KPI metrics (activation and no-show rates), monthly historical progress, and a consolidated activity feed.

**Independent Test**: Load the dashboard, verify KPI rates match database count, and activity feed displays recent client/project/meeting log entries.

### Implementation for User Story 4
- [x] T029 [P] [US4] Create Activity log DB schema and logger service in backend/internal/activity/service.go
- [x] T030 [US4] Implement Dashboard KPI calculation and history analytics queries in backend/internal/dashboard/service.go
- [x] T031 [US4] Implement Dashboard API endpoint in backend/internal/dashboard/handler.go
- [x] T032 [US4] Create Vue Dashboard visualization page in frontend/src/pages/Dashboard.vue

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Code cleanup, configuration adjustments, and documentation.

- [x] T033 [P] Document environment configuration setup and deployment instructions in README.md
- [x] T034 Run response latency audit and log outcomes in backend/tests/performance_audit.txt
- [x] T035 Execute codebase-wide code formatting and lint verification across backend/ and frontend/

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories.
- **User Stories (Phase 3+)**: All depend on Foundational phase completion.
  - User stories proceed sequentially in priority order (P1 → P2 → P3).
- **Polish (Final Phase)**: Depends on all desired user stories being complete.

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2).
- **User Story 2 (P1)**: Can start after Foundational (Phase 2).
- **User Story 3 (P2)**: Can start after Foundational (Phase 2).
- **User Story 4 (P3)**: Can start after Foundational (Phase 2).

### Parallel Opportunities

- Setup tasks `T003`, `T004` can run in parallel.
- Foundational schema initialization and connection pool configuration can run in parallel.
- Once Foundational is done, Developers can work on `US1` and `US2` in parallel.
- Models within stories (e.g. `T009` and `T016`, or `T022` and `T029`) can be created in parallel.

---

## Parallel Example: User Story 1

```bash
# Implement auth store and routing guards in parallel:
Task: "Implement Vue frontend auth store in frontend/src/stores/auth.js"
Task: "Implement Vue router navigation guards in frontend/src/router/index.js"
```

---

## Implementation Strategy

### MVP First (User Story 1 & 2 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL)
3. Complete Phase 3: User Story 1 (Auth/RBAC)
4. Complete Phase 4: User Story 2 (Client CRUD & Webhook Sync)
5. **STOP and VALIDATE**: Test authentication and sync flow end-to-end (MVP target).
6. Proceed to subsequent stories (Projects & Dashboard).
