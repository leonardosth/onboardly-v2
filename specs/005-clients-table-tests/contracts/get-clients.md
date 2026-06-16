# API Contract: GET /api/clients (Updated)

**Feature Branch**: `005-clients-table-tests` | **Date**: 2026-06-15

## Endpoint: List Clients with Details

**Method**: `GET`
**Path**: `/api/clients`
**Authentication**: Bearer JWT token required
**Roles**: Admin, Analista

### Description

Returns an aggregated list of all clients with their latest active project information, responsible analyst, and completed agendas count. This replaces the current response that only returns basic client fields.

### Request

No request body. No query parameters.

### Response: 200 OK

```json
[
  {
    "id": "uuid",
    "name": "Empresa XYZ Ltda",
    "cnpj": "12.345.678/0001-90",
    "project_name": "Implantação ERP",
    "project_status": "Em andamento",
    "project_is_active": true,
    "responsible": "analista@empresa.com",
    "completed_agendas": 3,
    "total_agendas": 5,
    "created_at": "2026-01-15T10:30:00Z"
  },
  {
    "id": "uuid",
    "name": "Startup ABC",
    "cnpj": "98.765.432/0001-10",
    "project_name": null,
    "project_status": null,
    "project_is_active": null,
    "responsible": null,
    "completed_agendas": 0,
    "total_agendas": 0,
    "created_at": "2026-02-20T14:00:00Z"
  }
]
```

### Response Fields

| Field              | Type           | Description                                          |
|--------------------|----------------|------------------------------------------------------|
| id                 | string (UUID)  | Client's unique identifier                           |
| name               | string         | Client's name                                        |
| cnpj               | string         | Client's CNPJ                                        |
| project_name       | string \| null | Name of the latest project (null if no projects)     |
| project_status     | string \| null | Status of the latest project (Backlog, Em andamento, Go-Live) |
| project_is_active  | bool \| null   | Whether the latest project is active                 |
| responsible        | string \| null | Email of the analyst on the latest meeting           |
| completed_agendas  | integer        | Count of completed meetings for the latest project   |
| total_agendas      | integer        | Count of all meetings for the latest project         |
| created_at         | string (ISO 8601) | Client creation timestamp                         |

### Error Responses

| Status | Body                                | Condition                   |
|--------|-------------------------------------|-----------------------------|
| 401    | `{"error": "authorization header missing"}` | No token provided    |
| 401    | `{"error": "invalid token"}`        | Invalid/expired JWT         |
| 500    | `{"error": "..."}`                  | Database error              |

### Backward Compatibility Note

This is a **breaking change** to the existing `GET /api/clients` endpoint. The response now includes additional fields beyond the original `id`, `name`, `cnpj`, `created_at`, `updated_at`. The frontend `ClientsList.vue` will be updated in the same feature to consume the new shape.
