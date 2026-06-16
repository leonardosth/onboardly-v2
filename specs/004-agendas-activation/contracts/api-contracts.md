# API Contracts: Agendas and Activation Tracking

## New Endpoints

### 1. `GET /api/meetings/mine` — List Analyst's Completed Meetings

Returns all meetings assigned to the authenticated analyst, filtered by status.

**Auth**: JWT (any role)

**Query Parameters**:

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `status` | string | `completed` | Filter by meeting status: `scheduled`, `completed`, `cancelled`, or empty for all |

**Response**: `200 OK`

```json
[
  {
    "id": "uuid",
    "project_id": "uuid",
    "analyst_id": "uuid",
    "title": "Reunião de Implantação - Cliente X",
    "scheduled_at": "2026-06-08T14:00:00Z",
    "status": "completed",
    "completed_at": "2026-06-08T15:30:00Z",
    "no_show": false,
    "created_at": "2026-06-01T10:00:00Z",
    "project_name": "Implantação Cliente X",
    "client_name": "Empresa X"
  }
]
```

**Notes**: The response includes `project_name` and `client_name` joined from related tables for display convenience.

---

### 2. `POST /api/meetings/{id}/complete` — Complete Meeting + Optionally Activate Client

Marks a meeting as completed. Optionally marks the associated project's client as active.

**Auth**: JWT (any role)

**Path Parameters**: `id` — Meeting UUID

**Request Body**:

```json
{
  "activate_client": true
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `activate_client` | boolean | No (default: false) | If true, also sets `projects.activated_at = NOW()` for the associated project |

**Response**: `200 OK`

```json
{
  "meeting": {
    "id": "uuid",
    "project_id": "uuid",
    "status": "completed",
    "completed_at": "2026-06-08T15:30:00Z"
  },
  "project_activated": true
}
```

**Error Responses**:
- `404 Not Found`: Meeting not found
- `400 Bad Request`: Meeting is already completed or cancelled
- `400 Bad Request`: Project already activated (when `activate_client: true` and `activated_at` is not null)

---

### 3. `POST /api/projects/{id}/finalize` — Finalize Implementation Project

Manually finalizes a project by setting status to `Go-Live`.

**Auth**: JWT (any role)

**Path Parameters**: `id` — Project UUID

**Request Body**: None required.

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "name": "Implantação Cliente X",
  "status": "Go-Live",
  "is_active": false,
  "activated_at": "2026-06-08T15:30:00Z",
  "created_at": "2026-06-01T10:00:00Z",
  "updated_at": "2026-06-08T16:00:00Z"
}
```

**Error Responses**:
- `404 Not Found`: Project not found
- `400 Bad Request`: Project is already finalized (status = Go-Live)

**Notes**: This reuses the existing `UpdateProjectStatus` logic but provides a dedicated endpoint for the "finalize" action per the clarified spec.

---

## Modified Endpoints

### 4. `GET /api/dashboard` — Dashboard Data (Extended)

The existing endpoint is extended to return additional metrics.

**Auth**: JWT (any role)

**Response**: `200 OK` (extended schema)

```json
{
  "metrics": {
    "activation_rate": 72.5,
    "no_show_rate": 8.3,
    "abandonment_rate": 15.2,
    "first_meeting_activation_rate": 45.0,
    "activation_30d_rate": 68.0
  },
  "funnel": {
    "registered": 100,
    "participants": 75,
    "active": 54
  },
  "cohorts": [
    {
      "month": "2026-06",
      "total": 20,
      "activated": 14,
      "activation_rate": 70.0
    },
    {
      "month": "2026-05",
      "total": 25,
      "activated": 20,
      "activation_rate": 80.0
    }
  ],
  "history": [
    { "month": "2026-01", "deployments": 5 }
  ],
  "recent_activities": []
}
```

**New fields**:
- `metrics.abandonment_rate`: Percentage of projects with no completed meetings
- `metrics.first_meeting_activation_rate`: Percentage of participants activated in first meeting
- `metrics.activation_30d_rate`: Percentage of participants activated within 30 calendar days
- `funnel`: Object with `registered`, `participants`, `active` counts
- `cohorts`: Array of monthly cohort data with activation rates

**Backward compatibility**: Existing `metrics.activation_rate`, `metrics.no_show_rate`, `history`, and `recent_activities` fields remain unchanged.
