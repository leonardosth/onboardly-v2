# Quickstart Validation Guide: Clients List Table Update & Unit Tests

**Feature Branch**: `005-clients-table-tests` | **Date**: 2026-06-15

## Prerequisites

- PostgreSQL running on `localhost:5433` with database `onboardlyv2`
- Go 1.22+ installed
- Node.js 18+ installed
- Backend `.env` file configured (see `backend/.env.example`)

## Scenario 1: Verify Aggregated Clients API

### Setup

1. Start the backend server:
   ```bash
   cd backend && go run cmd/server/main.go
   ```

2. Ensure seed data exists (at least 1 client with a project and meetings).

### Validation

```bash
# Login to get a JWT token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@onboardly.com","password":"admin123"}' | jq -r '.token')

# Call the updated clients endpoint
curl -s http://localhost:8080/api/clients \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### Expected Outcome

Response should be an array of objects containing:
- `id`, `name`, `cnpj` (original fields)
- `project_name`, `project_status`, `project_is_active` (latest project info)
- `responsible` (analyst email from latest meeting)
- `completed_agendas`, `total_agendas` (meeting counts)

Clients without projects should show `null` for project fields and `0` for agenda counts.

See [API contract](contracts/get-clients.md) for full field definitions.

## Scenario 2: Verify Frontend Table View

### Setup

1. Start both backend and frontend dev servers.
2. Login as Admin or Analista.

### Validation

1. Navigate to **Clientes** page.
2. Verify the page displays a **table** (not cards) with columns:
   - Nome do Cliente
   - CNPJ
   - Projeto
   - Responsável
   - Status
   - Situação (Ativo/Inativo)
   - Agendas Realizadas
3. Verify clients without projects show appropriate placeholder values.
4. Verify the "Ver Detalhes" link still works.
5. Verify the "Excluir" button is visible only for Admin role.

### Expected Outcome

The table should be styled consistently with the dark theme used across the app (matching `#0f172a` background and `#38bdf8` accent).

## Scenario 3: Verify Unit Tests

### Setup

```bash
cd backend
go test ./... -v
```

### Expected Outcome

- All tests pass with `PASS` status.
- Tests should include at minimum:
  - `client.Validate()` — valid and invalid CNPJ/name
  - `GetClientsWithDetails()` — mocked SQL returning expected aggregated data
  - `auth.HashPassword()` / `auth.CheckPasswordHash()` — password hashing round-trip
  - `auth.GenerateJWT()` — token generation and claims verification

See [data-model.md](data-model.md) for entity details referenced in tests.
