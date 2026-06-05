# Quickstart Validation Guide: Onboardly Core

This document outlines the steps to verify the end-to-end setup and functionality of the Onboardly core application.

## Prerequisites

- **Go**: Version 1.22+ installed.
- **Node.js**: Version 18+ and `npm` installed.
- **PostgreSQL**: Running locally with a database named `onboardlyv2` and schema created per the [Data Model](file:///C:/Users/leona/Documents/Projects/onboardly-v2/specs/001-onboardly-core/data-model.md).

---

## Service Startup Instructions

### 1. Start Go Backend API
```bash
# Navigate to backend directory
cd backend

# Install Go module dependencies
go mod tidy

# Start the server (binds to http://localhost:8080)
go run cmd/server/main.go
```

### 2. Start Vue.js 3 Frontend
```bash
# Navigate to frontend directory
cd frontend

# Install package dependencies
npm install

# Start the Vite dev server (runs on http://localhost:5173)
npm run dev
```

---

## End-to-End Validation Scenarios

Follow these scenarios sequentially to validate all requirements. You can use any API client (e.g., `curl`) to run the commands.

### Scenario 1: User Registration and Login (FR-001, FR-002, FR-003, FR-004)
1. **Register Admin**:
   ```bash
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@onboardly.com","password":"AdminPassword123","role":"Admin"}'
   ```
   *Expected Outcome*: Status `201 Created` with the registered user profile returned (password omitted).

2. **Login & Extract JWT**:
   ```bash
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@onboardly.com","password":"AdminPassword123"}'
   ```
   *Expected Outcome*: Status `200 OK` returning a JSON payload containing the JWT token. Save this token for subsequent calls.

---

### Scenario 2: Webhook Sync from Google Sheets (FR-006)
1. **Send Webhook Sync**:
   ```bash
   curl -X POST http://localhost:8080/api/webhooks/clients \
     -H "Content-Type: application/json" \
     -H "X-Webhook-Token: your-shared-secret-token" \
     -d '{"name":"Google Sheets Client","cnpj":"12.345.678/0001-99"}'
   ```
   *Expected Outcome*: Status `200 OK` returning `{"status":"success","message":"Client synced successfully"}`.

2. **Verify Client Created**:
   Run a `GET /api/clients` request using the JWT authorization token and confirm "Google Sheets Client" is listed.

---

### Scenario 3: Projects & Meeting Lifecycle (FR-008, FR-011, FR-012, FR-014, FR-015)
1. **Create Project**:
   Create a project linked to the client ID returned in Scenario 2.
   ```bash
   curl -X POST http://localhost:8080/api/projects \
     -H "Authorization: Bearer <JWT_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"client_id":"<CLIENT_UUID>","name":"Acme Integration"}'
   ```
   *Expected Outcome*: Status `201 Created` returning the project with state `Backlog` and `is_active: true`.

2. **Schedule Meeting (Standard Date)**:
   ```bash
   curl -X POST http://localhost:8080/api/meetings \
     -H "Authorization: Bearer <JWT_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"project_id":"<PROJECT_UUID>","analyst_id":"<USER_UUID>","title":"Kickoff Meeting","scheduled_at":"2026-06-10T14:30:00Z"}'
   ```
   *Expected Outcome*: Status `201 Created`.

3. **Verify Referential Integrity Constraint**:
   Attempt to schedule a meeting with a non-existent project UUID:
   ```bash
   curl -X POST http://localhost:8080/api/meetings \
     -H "Authorization: Bearer <JWT_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"project_id":"00000000-0000-0000-0000-000000000000","analyst_id":"<USER_UUID>","title":"Orphan Meeting","scheduled_at":"2026-06-10T14:30:00Z"}'
   ```
   *Expected Outcome*: Status `400 Bad Request` or `422 Unprocessable Entity` indicating database constraint failure.

---

### Scenario 4: Dashboard & BI Metrics (FR-016, FR-017, FR-019)
1. **Fetch Dashboard Statistics**:
   ```bash
   curl -X GET http://localhost:8080/api/dashboard \
     -H "Authorization: Bearer <JWT_TOKEN>"
   ```
   *Expected Outcome*: Returns JSON containing current activation rates, no-show rates, 6-month historical graphs, and the recent activity feed.
