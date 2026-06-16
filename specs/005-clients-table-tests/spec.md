# Feature Specification: Clients List Table Update & Unit Tests

**Feature Branch**: `005-clients-table-tests`

**Created**: 2026-06-15

**Status**: Draft

**Input**: User description: "preciso solicitar uma alteração na tela de clientes: preciso que a listagem de clientes traga informações de número do projeto, nome do cliente, responsável, situação, status, agendas realizadas, etc.. em formato de tabela. Preciso também configurar os testes unitários do sistema"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View comprehensive clients list (Priority: P1)

As an Analyst or Admin, I want to see the clients list in a comprehensive table format so that I can quickly assess the status, project number, responsible analyst, and completed agendas for each client.

**Why this priority**: The core value of the clients screen is to quickly understand the status of clients and their projects. A table view provides this information at a glance.

**Independent Test**: Can be fully tested by navigating to the Clients screen and verifying that the data is displayed in a table with columns for project number, client name, responsible, situation, status, and completed agendas.

**Acceptance Scenarios**:

1. **Given** I am logged in and have clients with associated projects and agendas, **When** I navigate to the Clients list, **Then** I should see a table containing "número do projeto", "nome do cliente", "responsável", "situação", "status", and "agendas realizadas".
2. **Given** a client has no active projects or agendas, **When** they are listed in the table, **Then** the table should display appropriate empty/default values for project-related columns.

---

### User Story 2 - Automated Unit Testing Configuration (Priority: P2)

As a Developer, I want to have a unit testing framework configured for the system so that I can write tests to ensure code quality and prevent regressions, especially for core business logic.

**Why this priority**: Unit tests are essential for long-term maintainability and reliability of the codebase, ensuring new changes do not break existing functionality.

**Independent Test**: Can be fully tested by running the test command in the terminal and verifying that the test runner executes and reports the results.

**Acceptance Scenarios**:

1. **Given** the repository, **When** I run the unit test command, **Then** the test framework should execute all configured unit tests and output a success/failure report.

### Edge Cases

- What happens when a client has multiple projects? The table will show only the latest active project to keep the interface clean. Historical or concurrent projects can be viewed in the client's detailed view.
- How does the system handle very long names or project numbers in the table view? They should be truncated with an ellipsis and full text shown on hover.
- What happens if the test suite fails during a CI/CD pipeline run? The build/deployment process will be halted.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST display the clients list in a table format.
- **FR-002**: The clients table MUST include columns for: Project Number, Client Name, Responsible Analyst, Situation, Status, and Completed Agendas (agendas realizadas).
- **FR-003**: The backend MUST provide an API endpoint or update the existing clients endpoint to return this aggregated data (client + project + agendas info).
- **FR-004**: The system MUST have a unit testing framework configured for the backend (Go). It should use standard tooling like the built-in `testing` package along with a library like `testify` for assertions, mirroring the developer experience of tools like Jest.

### Key Entities *(include if feature involves data)*

- **Client**: Represents the customer, containing the base data (name, etc.).
- **Project**: Represents the implementation project linked to the client, providing the project number, situation, status, and responsible analyst.
- **Meeting/Agenda**: Represents the interactions with the client, from which "completed agendas" count is derived.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The clients page successfully displays the table with all required columns populated with accurate data.
- **SC-002**: Unit test commands execute successfully in both local and CI environments.
- **SC-003**: At least one sample unit test is written and passes to prove the configuration works.

## Assumptions

- We assume the backend database already stores the required relationships (Client -> Project -> Meeting/Agenda) to aggregate this data.
- We assume standard table pagination or scrolling will be implemented to handle a large number of clients.
- We assume standard testing frameworks will be used (e.g., standard `testing` package for Go backend, Vitest/Jest for Vue frontend).
