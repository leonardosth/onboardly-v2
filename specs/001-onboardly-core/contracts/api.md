# API Endpoint Contracts: Onboardly Core

This document defines the HTTP API endpoint schemas and contracts exposed by the backend Go application.

## Authentication & Authorization

All endpoints except `/api/auth/register`, `/api/auth/login`, and `/api/webhooks/*` require a valid JWT token passed in the `Authorization: Bearer <JWT>` header.

### 1. Register User
- **Route**: `POST /api/auth/register`
- **Request Body**:
  ```json
  {
    "email": "user@onboardly.com",
    "password": "SecurePassword123",
    "role": "Analista"
  }
  ```
- **Response** (`201 Created`):
  ```json
  {
    "id": "d3b07384-d113-4ec8-a5b6-72aa8e698889",
    "email": "user@onboardly.com",
    "role": "Analista"
  }
  ```

### 2. Login User
- **Route**: `POST /api/auth/login`
- **Request Body**:
  ```json
  {
    "email": "user@onboardly.com",
    "password": "SecurePassword123"
  }
  ```
- **Response** (`200 OK`):
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAb25ib2FyZGx5LmNvbSIsImNvbXBhbnkiOiJvbmJvYXJkbHkifQ...",
    "role": "Analista"
  }
  ```

---

## Client Management

### 3. Create Client
- **Route**: `POST /api/clients`
- **Request Body**:
  ```json
  {
    "name": "Acme Corporation",
    "cnpj": "12.345.678/0001-99"
  }
  ```
- **Response** (`201 Created`):
  ```json
  {
    "id": "e4b07384-d113-4ec8-a5b6-72aa8e698999",
    "name": "Acme Corporation",
    "cnpj": "12.345.678/0001-99",
    "created_at": "2026-06-04T15:00:00Z"
  }
  ```

### 4. List Clients
- **Route**: `GET /api/clients`
- **Response** (`200 OK`):
  ```json
  [
    {
      "id": "e4b07384-d113-4ec8-a5b6-72aa8e698999",
      "name": "Acme Corporation",
      "cnpj": "12.345.678/0001-99",
      "created_at": "2026-06-04T15:00:00Z"
    }
  ]
  ```

### 5. Webhook Sync (Google Sheets / Salesforce)
- **Route**: `POST /api/webhooks/clients`
- **Headers**:
  - `X-Webhook-Token`: `your-shared-secret-token`
- **Request Body**:
  ```json
  {
    "name": "Salesforce Client Inc",
    "cnpj": "98.765.432/0001-11"
  }
  ```
- **Response** (`200 OK`):
  ```json
  {
    "status": "success",
    "message": "Client synced successfully"
  }
  ```

---

## Projects & Meetings

### 6. Create Project
- **Route**: `POST /api/projects`
- **Request Body**:
  ```json
  {
    "client_id": "e4b07384-d113-4ec8-a5b6-72aa8e698999",
    "name": "Acme Core Deployment"
  }
  ```
- **Response** (`201 Created`):
  ```json
  {
    "id": "f5b07384-d113-4ec8-a5b6-72aa8e698777",
    "client_id": "e4b07384-d113-4ec8-a5b6-72aa8e698999",
    "name": "Acme Core Deployment",
    "status": "Backlog",
    "is_active": true
  }
  ```

### 7. Update Project Status
- **Route**: `PUT /api/projects/:id`
- **Request Body**:
  ```json
  {
    "status": "Go-Live"
  }
  ```
- **Response** (`200 OK`):
  ```json
  {
    "id": "f5b07384-d113-4ec8-a5b6-72aa8e698777",
    "status": "Go-Live",
    "is_active": false
  }
  ```

### 8. Schedule Meeting
- **Route**: `POST /api/meetings`
- **Request Body**:
  ```json
  {
    "project_id": "f5b07384-d113-4ec8-a5b6-72aa8e698777",
    "analyst_id": "d3b07384-d113-4ec8-a5b6-72aa8e698889",
    "title": "Onboarding Kickoff",
    "scheduled_at": "2026-06-10T14:30:00Z"
  }
  ```
- **Response** (`201 Created`):
  ```json
  {
    "id": "a1b07384-d113-4ec8-a5b6-72aa8e698666",
    "project_id": "f5b07384-d113-4ec8-a5b6-72aa8e698777",
    "analyst_id": "d3b07384-d113-4ec8-a5b6-72aa8e698889",
    "title": "Onboarding Kickoff",
    "scheduled_at": "2026-06-10T14:30:00Z",
    "no_show": false
  }
  ```

---

## Dashboard Analytics

### 9. Get Dashboard
- **Route**: `GET /api/dashboard`
- **Response** (`200 OK`):
  ```json
  {
    "metrics": {
      "activation_rate": 85.5,
      "no_show_rate": 9.2
    },
    "history": [
      { "month": "2025-12", "deployments": 5 },
      { "month": "2026-01", "deployments": 8 },
      { "month": "2026-02", "deployments": 12 },
      { "month": "2026-03", "deployments": 10 },
      { "month": "2026-04", "deployments": 15 },
      { "month": "2026-05", "deployments": 14 }
    ],
    "recent_activities": [
      {
        "id": "c1b07384-d113-4ec8-a5b6-72aa8e698444",
        "entity_type": "Project",
        "description": "Project 'Acme Core Deployment' updated status to 'Go-Live'",
        "created_at": "2026-06-04T12:00:00Z"
      }
    ]
  }
  ```
