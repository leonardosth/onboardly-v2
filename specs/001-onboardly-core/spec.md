# Feature Specification: Onboardly Core Requirements

**Feature Branch**: `001-onboardly-core`

**Created**: 2026-06-04

**Status**: Draft

**Input**: User description: "implement onboardly core requirements covering auth, client management, projects, meetings, dashboard metrics, and notifications backlog"

## Clarifications

### Session 2026-06-04
- Q: Are the Notification Alerts (FR-020) and PDF Export (FR-021) requirements part of the initial MVP release, or should they be deferred to the backlog? → A: Defer both to the Backlog.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - User Authentication & Access Control (Priority: P1)

As a system user (Admin or Analista), I want to securely register and log in to the application so that I can access authorized resources matching my role.

**Why this priority**: Crucial for application security and separating administrative actions from analyst work.

**Independent Test**: Register a new user, log in, and verify that accessing unauthorized routes returns access denied (403 Forbidden).

**Acceptance Scenarios**:
1. **Given** an unregistered user, **When** they fill the registration form with email, role, and password, **Then** their password is saved in the database as a cryptographically secure bcrypt hash.
2. **Given** a registered user, **When** they submit valid email and password credentials, **Then** the system returns a secure JWT token for session tracking.
3. **Given** an authenticated user with the "Analista" role, **When** they attempt to access Admin-restricted routes (such as deleting a client), **Then** they receive a 403 Forbidden error response.

---

### User Story 2 - Client Management & Salesforce / Google Sheets Integration (Priority: P1)

As an Admin or Analista, I want to manage client records and automatically ingest new client data via webhooks triggered from Google Sheets (containing exported reports from Salesforce) so that the database remains in sync.

**Why this priority**: Required to initialize the client portfolio before setting up projects.

**Independent Test**: Send a valid JSON payload to the webhook endpoint and verify the client is correctly created with name and CNPJ validation.

**Acceptance Scenarios**:
1. **Given** an authenticated user, **When** they create a new client with empty details or an invalid CNPJ format, **Then** the system rejects the operation and displays validation errors.
2. **Given** Google Sheets (acting as a pipeline for Salesforce data), **When** an Apps Script or integration tool posts a valid JSON payload of clients/contracts to the webhook route, **Then** the client is automatically registered in Onboardly.

---

### User Story 3 - Projects & Meeting Management (Priority: P2)

As an Analista, I want to create deployment projects for clients and schedule meetings for those projects using standard dates so that the onboarding process runs smoothly.

**Why this priority**: Serves as the primary operational flow of the onboarding system.

**Independent Test**: Schedule a meeting for a project and confirm the date is saved in ISO 8601 format, and that scheduling a meeting for a non-existent project is rejected.

**Acceptance Scenarios**:
1. **Given** a valid client, **When** an Analista creates a project, **Then** the project starts in the Backlog status.
2. **Given** an active project, **When** an Analista schedules a meeting with a date-time in ISO 8601 format, **Then** the meeting is saved successfully.
3. **Given** a non-existent project ID, **When** an Analista schedules a meeting, **Then** the system blocks the action and maintains database referential integrity.

---

### User Story 4 - Dashboard & Analytics BI (Priority: P3)

As a manager, I want to view a real-time dashboard with consolidated metrics and activity feeds to track activation and no-show performance.

**Why this priority**: Provides business intelligence and status overview for administrative decision-making.

**Independent Test**: Load the dashboard and verify that the calculated KPIs match the operational database records exactly.

**Acceptance Scenarios**:
1. **Given** project records in the database, **When** the dashboard is rendered, **Then** the system calculates the project activation rate (target ~86%) and no-show rate (target <10%) based on real-time data.
2. **Given** activities recorded over the past month, **When** loading the dashboard, **Then** the recent activity feed displays sorted logs of client, project, and meeting updates.

## Edge Cases

- **Duplicate CNPJ Ingestion**: When the ERP webhook sends a payload with a CNPJ that already exists, the system updates the existing client details instead of creating a duplicate entry.
- **Meeting to Project Deletion**: If a project is deleted, any scheduled meetings referencing that project must be deleted cascadingly or blocked to prevent orphan records.
- **Invalid Date Input**: If a meeting request contains a non-ISO 8601 date format, the API rejects the request with a client error (400 Bad Request).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST allow the registration of new users with specific access roles (e.g., Admin, Analista).
- **FR-002**: The system MUST authenticate users via email and password, hashing the password using bcrypt.
- **FR-003**: The system MUST generate and validate JWT session tokens for secure request tracking.
- **FR-004**: The system MUST enforce Role-Based Access Control (RBAC), limiting specific routes/actions by role.
- **FR-005**: The system MUST require a Name and CNPJ when registering a new client manually.
- **FR-006**: The system MUST expose a webhook receiver endpoint to process client data in JSON format from Google Sheets (used to sync Salesforce report exports). Webhook requests MUST be authenticated using a Shared Secret Token passed in the request header (e.g., `X-Webhook-Token`).
- **FR-007**: The system MUST support standard CRUD operations (create, read, update, delete) for clients.
- **FR-008**: The system MUST allow the creation of deployment projects linked to a client.
- **FR-009**: The system MUST support updating project status through predefined stages (e.g., Backlog, Em andamento, Go-Live).
- **FR-010**: The system MUST filter the project dashboard by the assigned analyst or the project's current status.
- **FR-011**: The system MUST calculate and track the activation status (active/inactive) of each project. A project MUST be dynamically calculated as "active" if its current status is in a stage other than "Go-Live" or completed (e.g., stages like "Backlog" or "Em andamento" are considered active), and "inactive" once it reaches the "Go-Live" or completed stage.
- **FR-012**: The system MUST allow scheduling meetings linked to specific projects.
- **FR-013**: The system MUST register the Analista responsible for conducting each meeting.
- **FR-014**: The system MUST enforce database referential integrity, blocking meetings scheduled for non-existent projects.
- **FR-015**: The system MUST accept and store all datetime values in ISO 8601/RFC3339 format.
- **FR-016**: The system MUST display a dashboard showing consolidated, real-time KPI metrics.
- **FR-017**: The system MUST calculate the project activation rate and meeting no-show rate metrics.
- **FR-018**: The system MUST display a monthly implementation progress chart spanning the last 6 months.
- **FR-019**: The system MUST generate a feed of "Atividades Recentes" combining client, project, and meeting updates.
- **FR-020**: [DEFERRED TO BACKLOG] The system SHOULD emit notifications/visual alerts for projects close to their deadline.
- **FR-021**: [DEFERRED TO BACKLOG] The system SHOULD allow exporting productivity reports in PDF format.

### Key Entities *(include if feature involves data)*

- **User**: Represents a platform user containing email, password (hashed), and role (Admin or Analista).
- **Client**: Represents an onboarding client target containing Name, CNPJ, and active status.
- **Project**: Represents a deployment process linked to a Client, with a progress status and active state.
- **Meeting**: Represents a client touchpoint linked to a Project, with a datetime, analyst, and no-show flag.
- **Activity**: A system audit log entry tracking updates for the dashboard feed.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The login endpoint responds with a valid JWT in less than 500ms under normal load.
- **SC-002**: Webhook payloads are parsed, validated, and saved in the database in less than 1.5 seconds.
- **SC-003**: The dashboard metrics are calculated and rendered within 2 seconds of loading.
- **SC-004**: System successfully blocks 100% of meeting creation requests that reference invalid project IDs.

## Assumptions

- The application uses a relational database to ensure referential integrity constraints.
- The standard user timezone is represented correctly in UTC format on the server side.
- Webhook payloads contain standard client properties conforming to the system requirements.
- Notification alerts (FR-020) and PDF reports export (FR-021) are deferred to the project backlog and are out of scope for the initial MVP release.
