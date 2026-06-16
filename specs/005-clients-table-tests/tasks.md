# Tasks: Clients List Table Update & Unit Tests

**Input**: Design documents from `specs/005-clients-table-tests/`

**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/get-clients.md, quickstart.md

**Tests**: Included — the spec explicitly requests unit test configuration (User Story 2).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add testing dependencies to the Go project

- [x] T001 Add `github.com/stretchr/testify` dependency in `backend/go.mod`
- [x] T002 Add `github.com/DATA-DOG/go-sqlmock` dependency in `backend/go.mod`
- [x] T003 Run `go mod tidy` to resolve dependency graph in `backend/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Create the backend aggregated model and query that both user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Add `ClientWithDetails` struct to `backend/internal/client/model.go` with fields: ID, Name, CNPJ, ProjectName (*string), ProjectStatus (*string), ProjectIsActive (*bool), Responsible (*string), CompletedAgendas (int), TotalAgendas (int), CreatedAt (time.Time)
- [x] T005 Implement `GetClientsWithDetails()` function in `backend/internal/client/service.go` using LEFT JOIN LATERAL query (see plan.md SQL for details)
- [x] T006 Update `ListClientsHandler` in `backend/internal/client/handler.go` to call `GetClientsWithDetails()` instead of `GetAllClients()`

**Checkpoint**: Backend API now returns aggregated client data. Existing frontend will still render (it just won't show the new fields yet).

---

## Phase 3: User Story 1 - View Comprehensive Clients List Table (Priority: P1) 🎯 MVP

**Goal**: Replace the card-based clients grid with a comprehensive data table displaying project number, client name, responsible analyst, situation, status, and completed agendas.

**Independent Test**: Navigate to the Clients page and verify data table renders with all required columns populated.

### Implementation for User Story 1

- [x] T007 [US1] Replace `.clients-grid` card layout with `<table>` element in `frontend/src/pages/ClientsList.vue` template section with columns: Nome do Cliente, CNPJ, Projeto, Responsável, Status, Situação, Agendas Realizadas
- [x] T008 [US1] Add status badge rendering logic in `frontend/src/pages/ClientsList.vue` (Backlog → gray badge, Em andamento → blue badge, Go-Live → green badge)
- [x] T009 [US1] Add situação indicator in `frontend/src/pages/ClientsList.vue` (Ativo → green dot, Inativo → red dot, null → dash)
- [x] T010 [US1] Display agendas as "completed/total" format (e.g. "3/5") in `frontend/src/pages/ClientsList.vue`
- [x] T011 [US1] Handle null/empty project fields gracefully in `frontend/src/pages/ClientsList.vue` (show "—" for clients without projects)
- [x] T012 [US1] Style table with dark theme CSS matching existing design system in `frontend/src/pages/ClientsList.vue` scoped styles (background `#0f172a`, borders `rgba(255,255,255,0.05)`, accent `#38bdf8`)
- [x] T013 [US1] Ensure "Ver Detalhes" link and "Excluir" button (Admin only) remain functional in table row actions column in `frontend/src/pages/ClientsList.vue`
- [x] T014 [US1] Add responsive table wrapper with horizontal scroll for narrow viewports in `frontend/src/pages/ClientsList.vue` scoped styles

**Checkpoint**: Clients page fully functional with table view. Test by logging in and navigating to /clients.

---

## Phase 4: User Story 2 - Automated Unit Testing Configuration (Priority: P2)

**Goal**: Configure the Go backend unit testing framework with testify and go-sqlmock, and create sample tests to prove the configuration works.

**Independent Test**: Run `cd backend && go test ./... -v` and verify all tests pass.

### Tests for User Story 2

- [x] T015 [P] [US2] Create `backend/internal/client/model_test.go` with tests for `Client.Validate()`: valid CNPJ passes, empty name fails, invalid CNPJ format fails
- [x] T016 [P] [US2] Create `backend/internal/client/service_test.go` with mock test for `GetClientsWithDetails()` using go-sqlmock: mock SQL rows, verify correct struct mapping, verify empty result handling
- [x] T017 [P] [US2] Create `backend/internal/auth/service_test.go` with tests for: `HashPassword` produces valid hash, `CheckPasswordHash` returns true for correct password and false for wrong password, `GenerateJWT` produces token with correct claims (email, role, expiration)
- [x] T018 [P] [US2] Create `backend/internal/project/service_test.go` with test for `UpdateProjectStatus` validation: valid statuses accepted, invalid status rejected with error message

**Checkpoint**: All tests pass with `go test ./... -v`. Framework is configured and ready for future test additions.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [x] T019 Restart backend server (`go run cmd/server/main.go` in `backend/`) and verify no compilation errors
- [x] T020 Run `cd backend && go test ./... -v` and confirm all tests pass
- [x] T021 Run quickstart.md validation scenarios (Scenario 1: API, Scenario 2: UI, Scenario 3: Tests) per `specs/005-clients-table-tests/quickstart.md`
- [x] T022 Verify table text truncation with ellipsis for long client/project names in `frontend/src/pages/ClientsList.vue`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 completion — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Phase 2 (needs the updated API response)
- **User Story 2 (Phase 4)**: Depends on Phase 1 (needs testify/sqlmock deps) and Phase 2 (needs `GetClientsWithDetails` to test)
- **Polish (Phase 5)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends on Foundational (Phase 2) — frontend table consumes the new API shape
- **User Story 2 (P2)**: Depends on Foundational (Phase 2) — tests need the new service function to exist. Can run in parallel with US1 since they touch different files (backend test files vs. frontend Vue files).

### Within Each User Story

- US1: Template → Status badges → Situação → Agendas → Null handling → Styling → Actions → Responsive
- US2: All test files are independent ([P]) and can be written in parallel

### Parallel Opportunities

- T001 and T002 can run in parallel (both modify go.mod but are sequential adds)
- T015, T016, T017, T018 can ALL run in parallel (different test files)
- US1 (frontend) and US2 (backend tests) can run in parallel after Phase 2

---

## Parallel Example: User Story 2 (Tests)

```bash
# Launch all test file creation in parallel (different files, no dependencies):
Task T015: "Create model_test.go in backend/internal/client/"
Task T016: "Create service_test.go in backend/internal/client/"
Task T017: "Create service_test.go in backend/internal/auth/"
Task T018: "Create service_test.go in backend/internal/project/"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (add test deps)
2. Complete Phase 2: Foundational (aggregated query + handler update)
3. Complete Phase 3: User Story 1 (table UI)
4. **STOP and VALIDATE**: Test the clients table in the browser
5. Deploy/demo if ready — users can already see the improved table

### Incremental Delivery

1. Complete Setup + Foundational → Backend API updated
2. Add User Story 1 → Table UI → Test in browser → Deploy/Demo (MVP!)
3. Add User Story 2 → Unit tests → Run `go test ./...` → Deploy/Demo
4. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (frontend table)
   - Developer B: User Story 2 (backend tests)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Total: 22 tasks (3 setup + 3 foundational + 8 US1 + 4 US2 + 4 polish)
