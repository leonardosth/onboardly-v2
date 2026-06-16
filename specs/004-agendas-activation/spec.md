# Feature Specification: Agendas and Activation Tracking

**Feature Branch**: `004-agendas-activation`

**Created**: 2026-06-08

**Status**: Draft

**Input**: User description: "preciso implementar a visão de agendas, para que o analista possa acompanhar as próprias agendas realizadas. Preciso também implementar o conceito de ativo, onde o cliente executa determinadas tarefas durante a reunião de implantação, caso o cliente seja ativo na primeira reunião, o analista pode marcar que o cliente está ativo, marcar a reunião como realizada e finalizar o projeto. No dashboard, o analista precisa acompanhar o funil, inscritos > participantes > ativos. Preciso também acompanhar a quantidade de clientes ativos de cada safra - a safra é determinada pelo mes de compra. A meta é ativar 80% dos clientes participantes em até 30 dias após a contratação. Temos também meta de 20% de abandono (clientes que se inscreveram na implantação mas não participaram de nenhuma reunião). Preciso acompanhar também a taxa de clientes ativos na primeira reunião."

## Clarifications

### Session 2026-06-08

- Q: Action Workflow - When the analyst marks the client as active during the first meeting, should the system automatically mark the meeting as completed and finalize the project in a single action, or are these three separate manual steps? → A: Mark active + complete meeting is one step, but project finalization is a separate manual step.
- Q: KPI Calculation - How should the 30-day activation limit be calculated for the KPI? → A: Exactly 30 calendar days from the contract creation date.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Analyst Meeting View (Priority: P1)

As an Analyst, I want to view my completed meetings (agendas realizadas) so that I can track my daily activities and project progression.

**Why this priority**: It is the foundation for tracking interactions and allows analysts to manage their workload.

**Independent Test**: Can be tested by having an analyst log in and verify if only their own past meetings appear in the view.

**Acceptance Scenarios**:

1. **Given** the analyst is logged in, **When** they navigate to the meetings view, **Then** they see a list of their completed meetings.
2. **Given** multiple analysts exist, **When** analyst A views their agendas, **Then** they do not see analyst B's agendas.

---

### User Story 2 - Client Activation Process (Priority: P1)

As an Analyst, I want to mark a client as "Active" during or after an implementation meeting if they have executed the necessary tasks. Project finalization is a separate manual step.

**Why this priority**: Defines the core business value of activating clients and progressing them through the implementation funnel.

**Independent Test**: Can be tested by marking a client as active during a meeting, verifying the meeting is marked as done, and subsequently verifying the project can be manually finalized.

**Acceptance Scenarios**:

1. **Given** an ongoing implementation meeting, **When** the analyst marks the client as active, **Then** the client is tagged as active and the meeting is marked as completed. The project remains ongoing until manually finalized.
2. **Given** the client is active in the first meeting, **When** the system records this, **Then** it counts towards the "Active on first meeting" metric.

---

### User Story 3 - Dashboard Funnel Tracking (Priority: P2)

As a Manager or Analyst, I want to view an implementation funnel (Registered > Participants > Active) on the dashboard so that I can monitor the conversion rates at each stage.

**Why this priority**: Provides essential visibility into the overall health of the implementation process and business goals.

**Independent Test**: Can be tested by creating dummy clients at each stage and verifying the funnel accurately reflects the numbers.

**Acceptance Scenarios**:

1. **Given** there are clients in various stages, **When** viewing the dashboard, **Then** the funnel displays the correct count for Registered, Participants, and Active clients.

---

### User Story 4 - Cohort Analysis and KPIs (Priority: P2)

As a Manager, I want to track active clients by cohort (purchase month), the 30-day activation rate, abandonment rate, and first-meeting activation rate.

**Why this priority**: Required for strategic decision-making and ensuring the company meets its KPIs (80% activation, 20% abandonment).

**Independent Test**: Can be tested by generating reports for a specific month and verifying the calculated percentages align with the expected data.

**Acceptance Scenarios**:

1. **Given** clients who purchased in a specific month, **When** viewing the cohort data, **Then** the number of active clients for that month is displayed.
2. **Given** the KPI dashboard, **When** analyzing activation, **Then** it shows the percentage of participating clients activated within 30 days.

### Edge Cases

- What happens when a meeting is marked as done but the project is not finalized?
- How does the system handle a client who changes their purchase date (safra)?
- What happens if an analyst is reassigned from a project after meetings have been completed?
- How is the activation limit calculated if the contract date is missing?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a view for analysts to see their own completed meetings.
- **FR-002**: System MUST allow analysts to mark a client as "Active".
- **FR-003**: System MUST provide a separate manual action to finalize the implementation project after a client is marked active.
- **FR-004**: System MUST calculate and display the implementation funnel: Registered (Inscritos) > Participants (Participantes) > Active (Ativos).
- **FR-005**: System MUST group and display active clients by cohort (Safra), based on their purchase month.
- **FR-006**: System MUST track the time between contract signing and activation (measured as exactly 30 calendar days) to evaluate the activation goal.
- **FR-007**: System MUST identify "Abandoned" clients (registered for implementation but with 0 meetings participated).
- **FR-008**: System MUST calculate the percentage of clients activated during their first meeting.

### Key Entities

- **Meeting (Agenda)**: Has a status (e.g., scheduled, completed), an assigned analyst, and links to a Project/Client.
- **Client**: Has a purchase date (determining the cohort) and an activation status (Active/Inactive).
- **Project**: Represents the implementation phase, has a status (e.g., ongoing, finalized).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The funnel correctly displays Registered > Participants > Active metrics in real-time.
- **SC-002**: System accurately reports the percentage of participating clients activated within exactly 30 calendar days of their contract creation date (Target: 80%).
- **SC-003**: System accurately identifies and reports the abandonment rate (Target: 20%).
- **SC-004**: System accurately reports the percentage of clients activated on their first meeting.
- **SC-005**: Analysts can view 100% of their completed meetings without seeing other analysts' meetings unless authorized.

## Assumptions

- "Registered" (Inscritos) refers to clients who have a created implementation project.
- "Participants" (Participantes) refers to clients who have had at least one meeting marked as completed.
- The contract date (data de contratação) is the same as the purchase date (data de compra) used to define the cohort.
- The concept of "tasks executed during the meeting" by the client does not require a granular checklist system initially; a single "Mark as Active" toggle is sufficient to represent this state.
