# Tasks: Agendas and Activation Tracking

**Input**: Design documents from `specs/004-agendas-activation/`

**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅, quickstart.md ✅

**Tests**: Not explicitly requested — test tasks omitted.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Database migration and shared model updates that all user stories depend on.

- [x] T001 Create migration file `backend/migrations/000002_activation.up.sql` — add `status VARCHAR(20) NOT NULL DEFAULT 'scheduled'` and `completed_at TIMESTAMPTZ` columns to `meetings` table, add `activated_at TIMESTAMPTZ` column to `projects` table, add CHECK constraint on `meetings.status`, add indexes `idx_meetings_analyst_status(analyst_id, status)` and `idx_projects_activated_at(activated_at)`
- [x] T002 Apply migration to local database using psql or migration tool

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Update existing Go models and core service logic that multiple user stories depend on.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T003 [P] Update Meeting model in `backend/internal/meeting/model.go` — add `Status string` and `CompletedAt *time.Time` fields with JSON tags
- [x] T004 [P] Update Project model in `backend/internal/project/model.go` — add `ActivatedAt *time.Time` field with JSON tag
- [x] T005 Update Meeting service queries in `backend/internal/meeting/service.go` — update `GetMeetingsByProject` and `CreateMeeting` to include `status` and `completed_at` columns in SELECT and INSERT statements
- [x] T006 Update Project service queries in `backend/internal/project/service.go` — update `GetProjects`, `GetProjectByID`, `CreateProject`, and `UpdateProjectStatus` to include `activated_at` column in all SELECT/INSERT/UPDATE statements

**Checkpoint**: Foundation ready — models and existing queries updated with new columns. Existing functionality preserved.

---

## Phase 3: User Story 1 — Analyst Meeting View (Priority: P1) 🎯 MVP

**Goal**: Provide analysts with a dedicated view of their own completed meetings ("Minhas Agendas").

**Independent Test**: Log in as an analyst, navigate to the meetings page, and verify only that analyst's completed meetings are displayed.

### Implementation for User Story 1

- [x] T007 [US1] Implement `GetMeetingsByAnalyst(analystID, status string)` function in `backend/internal/meeting/service.go` — query meetings filtered by `analyst_id` and `status`, JOIN with `projects` and `clients` tables to include `project_name` and `client_name` in results
- [x] T008 [US1] Create `MeetingWithDetails` response struct in `backend/internal/meeting/model.go` — includes base Meeting fields plus `ProjectName string` and `ClientName string`
- [x] T009 [US1] Implement `ListMyMeetingsHandler` in `backend/internal/meeting/handler.go` — extract analyst ID from JWT context (`auth.UserEmailKey`), accept `status` query parameter (default `completed`), call `GetMeetingsByAnalyst`
- [x] T010 [US1] Register `GET /api/meetings/mine` route in `backend/internal/api/router.go` — inside the authenticated group, before existing `/api/meetings` routes
- [x] T011 [P] [US1] Create meeting API service in `frontend/src/services/meeting.js` — implement `getMyMeetings(status)` calling `GET /api/meetings/mine?status=`
- [x] T012 [US1] Update meetings Pinia store in `frontend/src/stores/meetings.js` — add `myMeetings` state, `fetchMyMeetings(status)` action using the new service
- [x] T013 [US1] Create "Minhas Agendas" page in `frontend/src/pages/MeetingsList.vue` — display table of completed meetings with columns: title, client name, project name, scheduled date, status. Include status filter (completed/scheduled/all)
- [x] T014 [US1] Register `/meetings` route in `frontend/src/router/index.js` — lazy-load `MeetingsList.vue` inside the MainLayout children
- [x] T015 [US1] Add "Agendas" nav item in `frontend/src/layouts/MainLayout.vue` — add `router-link` to `/meetings` with calendar icon 📅, positioned after "Projetos" link

**Checkpoint**: Analysts can log in and view their own completed meetings in a dedicated page. No other analyst's data is visible.

---

## Phase 4: User Story 2 — Client Activation Process (Priority: P1)

**Goal**: Allow analysts to mark a client as active during a meeting (atomically completing the meeting) and separately finalize the project.

**Independent Test**: Complete a meeting with activation, verify meeting status is `completed` and project has `activated_at` set. Then manually finalize the project and verify status is `Go-Live`.

### Implementation for User Story 2

- [x] T016 [US2] Implement `CompleteMeeting(meetingID string, activateClient bool)` function in `backend/internal/meeting/service.go` — set meeting `status = 'completed'` and `completed_at = NOW()`. If `activateClient` is true, also set `projects.activated_at = NOW()` on the associated project (if not already set). Validate meeting exists and is in `scheduled` status
- [x] T017 [US2] Implement `CompleteMeetingHandler` in `backend/internal/meeting/handler.go` — parse `{id}` path param and `activate_client` boolean from request body, call `CompleteMeeting`, return meeting + activation status
- [x] T018 [US2] Implement `FinalizeProject(projectID string)` function in `backend/internal/project/service.go` — set `status = 'Go-Live'`, `is_active = false`, `updated_at = NOW()`. Validate project exists and is not already `Go-Live`
- [x] T019 [US2] Implement `FinalizeProjectHandler` in `backend/internal/project/handler.go` — parse `{id}` path param, call `FinalizeProject`, return updated project
- [x] T020 [US2] Register new routes in `backend/internal/api/router.go` — `POST /api/meetings/{id}/complete` and `POST /api/projects/{id}/finalize` inside the authenticated group
- [x] T021 [P] [US2] Add `completeMeeting(meetingId, activateClient)` and `finalizeProject(projectId)` functions in `frontend/src/services/meeting.js` and `frontend/src/services/project.js` respectively
- [x] T022 [US2] Add "Marcar como Ativo" button and "Finalizar Projeto" button in `frontend/src/pages/ProjectDetail.vue` — "Marcar como Ativo" opens a modal to select meeting + activate; "Finalizar Projeto" shows confirmation dialog and calls finalize endpoint. Buttons conditionally shown based on project state

**Checkpoint**: Full activation workflow works end-to-end. A meeting can be completed with activation. The project can be finalized separately.

---

## Phase 5: User Story 3 — Dashboard Funnel Tracking (Priority: P2)

**Goal**: Display the implementation funnel (Inscritos > Participantes > Ativos) on the dashboard.

**Independent Test**: Verify the dashboard displays correct funnel counts matching the database state.

### Implementation for User Story 3

- [x] T023 [US3] Add funnel query functions in `backend/internal/dashboard/service.go` — implement `getFunnelData()` returning `FunnelData{Registered, Participants, Active int}` using the SQL queries from data-model.md (COUNT projects, COUNT DISTINCT project_id from completed meetings, COUNT projects with activated_at)
- [x] T024 [US3] Add `FunnelData` struct to dashboard service in `backend/internal/dashboard/service.go` — add to `DashboardData` struct as `Funnel FunnelData` field with `json:"funnel"`
- [x] T025 [US3] Update `GetDashboardData()` in `backend/internal/dashboard/service.go` — call `getFunnelData()` and include in response
- [x] T026 [US3] Add funnel visualization widget in `frontend/src/pages/Dashboard.vue` — display a horizontal funnel bar showing Inscritos → Participantes → Ativos with counts and conversion percentages between stages. Use the existing dark design system with gradient colors

**Checkpoint**: Dashboard shows the real-time funnel with accurate counts.

---

## Phase 6: User Story 4 — Cohort Analysis and KPIs (Priority: P2)

**Goal**: Display cohort (safra) activation data, abandonment rate, 30-day activation rate, and first-meeting activation rate.

**Independent Test**: Verify the dashboard KPI cards display correct percentages and cohort table shows per-month data.

### Implementation for User Story 4

- [x] T027 [P] [US4] Add KPI query functions in `backend/internal/dashboard/service.go` — implement `getAbandonmentRate()`, `getActivation30dRate()`, `getFirstMeetingActivationRate()` using SQL from data-model.md
- [x] T028 [P] [US4] Add cohort query function in `backend/internal/dashboard/service.go` — implement `getCohortData()` returning `[]CohortItem{Month, Total, Activated, ActivationRate}` using JOIN on clients.created_at
- [x] T029 [US4] Add new structs `CohortItem` and extend `Metrics` struct in `backend/internal/dashboard/service.go` — add `AbandonmentRate`, `FirstMeetingActivationRate`, `Activation30dRate` fields to Metrics; add `Cohorts []CohortItem` field to DashboardData
- [x] T030 [US4] Update `GetDashboardData()` in `backend/internal/dashboard/service.go` — call all new KPI and cohort functions, include in response
- [x] T031 [US4] Add KPI metric cards in `frontend/src/pages/Dashboard.vue` — add cards for: Abandonment Rate (target ≤ 20%, red/green threshold), 30-Day Activation Rate (target ≥ 80%), First-Meeting Activation Rate. Follow existing `metric-card` design pattern with progress bars and meta-status indicators
- [x] T032 [US4] Add cohort (Safra) table in `frontend/src/pages/Dashboard.vue` — display monthly cohort data as a styled table showing month, total projects, activated count, and activation percentage. Highlight rows meeting the 80% target

**Checkpoint**: Dashboard shows complete KPI suite and cohort analysis. All metrics match database state.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories.

- [x] T033 [P] Update seed script in `backend/cmd/seed/main.go` — add seed data for: meetings with various statuses (scheduled, completed), projects with varying `activated_at` values (some NULL, some within 30 days, some after), clients with `created_at` in different months for cohort testing
- [x] T034 [P] Ensure responsive layout for new dashboard widgets in `frontend/src/pages/Dashboard.vue` — verify funnel, KPI cards, and cohort table render correctly on smaller screens
- [x] T035 Run quickstart.md validation scenarios end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (migration applied) — BLOCKS all user stories
- **US1 (Phase 3)**: Depends on Phase 2 — backend query + frontend page
- **US2 (Phase 4)**: Depends on Phase 2 — activation endpoints + UI buttons
- **US3 (Phase 5)**: Depends on Phase 2 — dashboard funnel queries
- **US4 (Phase 6)**: Depends on Phase 2 — dashboard KPI + cohort queries
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2 — No dependency on other stories
- **US2 (P1)**: Can start after Phase 2 — No dependency on other stories
- **US3 (P2)**: Can start after Phase 2 — No dependency on US1/US2 (reads data directly)
- **US4 (P2)**: Can start after Phase 2 — No dependency on US1/US2/US3 (reads data directly)

### Within Each User Story

- Backend models/service before handlers
- Handlers before route registration
- Routes before frontend service
- Frontend service before store
- Store before page component
- Page before router registration

### Parallel Opportunities

- T003 and T004 can run in parallel (different model files)
- US1 and US2 can start in parallel after Phase 2 (different endpoints/pages)
- US3 and US4 can start in parallel after Phase 2 (different dashboard sections)
- T027 and T028 can run in parallel (different query functions, same file but different functions)
- All four user stories can theoretically run in parallel with different developers

---

## Parallel Example: Phase 2 (Foundational)

```bash
# Launch model updates in parallel (different files):
Task T003: "Update Meeting model in backend/internal/meeting/model.go"
Task T004: "Update Project model in backend/internal/project/model.go"
```

## Parallel Example: User Stories 1 & 2 (After Phase 2)

```bash
# Developer A: User Story 1 (Analyst Meeting View)
Task T007-T015: Meeting listing backend + frontend

# Developer B: User Story 2 (Client Activation)
Task T016-T022: Activation + finalization endpoints + UI
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Migration
2. Complete Phase 2: Model updates
3. Complete Phase 3: User Story 1 (Analyst Meeting View)
4. **STOP and VALIDATE**: Analyst can view their completed meetings
5. Deploy/demo if ready

### Incremental Delivery

1. Phase 1 + Phase 2 → Foundation ready
2. Add US1 → Analyst meeting view works → Deploy (MVP!)
3. Add US2 → Activation workflow works → Deploy
4. Add US3 → Dashboard funnel visible → Deploy
5. Add US4 → Full KPI suite → Deploy
6. Phase 7 → Polished and seeded → Final deploy

### Parallel Team Strategy

With 2 developers:

1. Team completes Phase 1 + Phase 2 together
2. Once Phase 2 is done:
   - Developer A: US1 (Meeting View) + US3 (Funnel)
   - Developer B: US2 (Activation) + US4 (KPIs)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
