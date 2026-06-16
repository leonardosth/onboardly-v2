# Data Model: Clients List Table Update & Unit Tests

**Feature Branch**: `005-clients-table-tests` | **Date**: 2026-06-15

## Existing Entities (No Schema Changes Required)

The database schema already contains all the relationships needed. No migration is required.

### clients
| Field       | Type          | Constraints                     |
|-------------|---------------|---------------------------------|
| id          | UUID          | PK, default gen_random_uuid()   |
| name        | VARCHAR(255)  | NOT NULL                        |
| cnpj        | VARCHAR(18)   | UNIQUE, NOT NULL                |
| created_at  | TIMESTAMPTZ   | NOT NULL, default NOW()         |
| updated_at  | TIMESTAMPTZ   | NOT NULL, default NOW()         |

### projects
| Field        | Type          | Constraints                                |
|--------------|---------------|--------------------------------------------|
| id           | UUID          | PK, default gen_random_uuid()              |
| client_id    | UUID          | FK → clients(id) ON DELETE CASCADE         |
| name         | VARCHAR(255)  | NOT NULL                                   |
| status       | VARCHAR(50)   | CHECK (Backlog, Em andamento, Go-Live)     |
| is_active    | BOOLEAN       | NOT NULL, default TRUE                     |
| activated_at | TIMESTAMPTZ   | nullable                                   |
| created_at   | TIMESTAMPTZ   | NOT NULL, default NOW()                    |
| updated_at   | TIMESTAMPTZ   | NOT NULL, default NOW()                    |

### meetings
| Field        | Type          | Constraints                                |
|--------------|---------------|--------------------------------------------|
| id           | UUID          | PK, default gen_random_uuid()              |
| project_id   | UUID          | FK → projects(id) ON DELETE CASCADE        |
| analyst_id   | UUID          | FK → users(id) ON DELETE SET NULL          |
| title        | VARCHAR(255)  | NOT NULL                                   |
| scheduled_at | TIMESTAMPTZ   | NOT NULL                                   |
| status       | VARCHAR(20)   | CHECK (scheduled, completed, cancelled)    |
| completed_at | TIMESTAMPTZ   | nullable                                   |
| no_show      | BOOLEAN       | default FALSE                              |
| created_at   | TIMESTAMPTZ   | NOT NULL, default NOW()                    |

### users
| Field         | Type          | Constraints                     |
|---------------|---------------|---------------------------------|
| id            | UUID          | PK, default gen_random_uuid()   |
| email         | VARCHAR(255)  | UNIQUE, NOT NULL                |
| password_hash | VARCHAR(255)  | NOT NULL                        |
| role          | VARCHAR(50)   | CHECK (Admin, Analista)         |
| created_at    | TIMESTAMPTZ   | NOT NULL, default NOW()         |

## New Response Model (Backend Only — No DB Change)

### ClientWithDetails (Go struct)

A new struct for the aggregated API response. This is NOT a new database table.

| Field              | Type    | Source                                    |
|--------------------|---------|-------------------------------------------|
| id                 | UUID    | clients.id                                |
| name               | string  | clients.name                              |
| cnpj               | string  | clients.cnpj                              |
| project_number     | string  | projects.id (latest active project)       |
| project_name       | string  | projects.name                             |
| responsible_email  | string  | users.email (analyst from latest meeting) |
| situation          | string  | projects.status                           |
| is_active          | boolean | projects.is_active                        |
| completed_agendas  | integer | COUNT of meetings with status='completed' |
| total_agendas      | integer | COUNT of all meetings for the project     |
| created_at         | time    | clients.created_at                        |

## Relationships

```text
clients 1 ──────< N projects
projects 1 ──────< N meetings
meetings N >────── 1 users (analyst)
```

## Aggregation Query Strategy

The `GetClientsWithDetails()` function will use a query that:
1. Selects from `clients`
2. LEFT JOINs the latest project per client (using a subquery with `DISTINCT ON (client_id) ORDER BY created_at DESC`)
3. LEFT JOINs `meetings` on the project to count completed agendas
4. LEFT JOINs `users` to get the responsible analyst email from the most recent meeting

This returns one row per client, even if the client has no projects.
