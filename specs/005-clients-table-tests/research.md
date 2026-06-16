# Research: Clients List Table Update & Unit Tests

**Feature Branch**: `005-clients-table-tests` | **Date**: 2026-06-15

## R1: Aggregated Client Table Query Strategy

**Decision**: Create a new backend service function `GetClientsWithDetails()` that uses a single SQL query with LEFT JOINs across `clients`, `projects`, `meetings`, and `users` tables, returning only the latest active project per client via a subquery.

**Rationale**: A single query with JOINs and subqueries is the most efficient approach for this kind of aggregation. The existing codebase already uses this pattern in the `meeting` package (see `GetMeetingsByAnalyst()` which JOINs projects and clients). Using a subquery to pick the latest active project (via `DISTINCT ON` or a CTE with `ROW_NUMBER()`) avoids N+1 problems and multiple round-trips.

**Alternatives considered**:
- Multiple separate queries (fetch clients, then loop to fetch projects, then meetings): Rejected due to N+1 query problem and poor performance at scale.
- Database VIEW: Unnecessarily heavyweight; a function-level query is consistent with existing patterns.

## R2: Unit Testing Framework for Go Backend

**Decision**: Use Go's built-in `testing` package combined with `github.com/stretchr/testify` for assertions and mocking.

**Rationale**: `testify` is the de facto standard library for Go testing assertions (analogous to Jest for JavaScript). It provides `assert`, `require`, and `mock` packages. The Go standard library's `testing` package provides the test runner. No separate test runner is needed — `go test ./...` runs all tests. This is the most natural and commonly adopted approach in the Go ecosystem.

**Alternatives considered**:
- `gocheck`: Less popular, less maintained.
- `ginkgo/gomega`: BDD-style — adds unnecessary complexity for this project's scope.
- Standard `testing` only (no assertion library): Verbose; `testify` provides cleaner assertions.

## R3: Test Strategy Without Database Dependency

**Decision**: Use interface-based dependency injection and testify mocks. The current code uses the global `db.DB` directly, so unit tests for service logic will require refactoring service functions to accept a database interface or use `sqlmock` (`github.com/DATA-DOG/go-sqlmock`) to mock `*sql.DB`.

**Rationale**: `go-sqlmock` allows testing SQL-based service functions without a real database. Combined with `testify/assert`, it provides a complete testing workflow. This is the standard pattern in Go backend codebases that use `database/sql`.

**Alternatives considered**:
- Full integration tests with a test database: Too complex for initial setup; can be added later.
- Refactoring to repository interfaces: Good long-term goal, but too invasive for this feature's scope. `sqlmock` works with the existing architecture.

## R4: Frontend Table Component

**Decision**: Build a custom HTML `<table>` inside the existing `ClientsList.vue`, replacing the current cards grid. No external table library is needed since the requirements are straightforward (display columns, no complex sorting/filtering).

**Rationale**: The existing codebase uses plain Vue 3 components with scoped CSS (no component library like Vuetify or PrimeVue). Adding a heavy table library would be inconsistent. A styled HTML table matches the existing design system.

**Alternatives considered**:
- AG Grid / Vuetify DataTable: Overkill; would introduce a heavy dependency.
- Keep cards with extra info: Doesn't match user's explicit request for table format.
