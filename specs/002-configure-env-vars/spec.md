# Feature Specification: Environment Variable and Database Configuration via .env

**Feature Branch**: `002-configure-env-vars`

**Created**: 2026-06-04

**Status**: Draft

**Input**: User description: "vamos colocar as variaveis de ambiente em um arquivo dedicado .env, confira se o gitignore está considerando esse arquivo. nas variáveis de ambiente preciso definir os dados do banco: endereço localhost, porta 5433, usuario postgres, senha Pelopes818*"

## Clarifications

### Session 2026-06-04
- Q: O escopo da branch atual (configuração de variáveis de ambiente via .env) deve ser expandido para incluir a correção/implementação do fluxo completo de administração (como cadastro de clientes, gerenciamento de usuários e agendas), ou essas melhorias funcionais serão tratadas em uma branch/tarefa separada? → A: Tratar essas melhorias e correções como uma feature/branch separada, mantendo o escopo atual estritamente focado no isolamento de variáveis de ambiente.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Secure and Centralized Environment Configuration (Priority: P1)

As a developer, I want to define local application configurations (such as database credentials and port bindings) in a central `.env` file, so that sensitive connection parameters are kept separate from the source code and can be adjusted for different deployments.

**Why this priority**: This is the core requirement of the feature. Centralized configuration prevents hardcoding of sensitive credentials and enables deployment flexibility.

**Independent Test**: Verify that when a `.env` file is present, the application overrides its fallback configuration and successfully connects to the database at the specified address (`localhost:5433`), using the user `postgres` and password `Pelopes818*`.

**Acceptance Scenarios**:

1. **Given** no local `.env` file, **When** the application starts up, **Then** it falls back to the safe, predefined default configuration values.
2. **Given** a `.env` file with custom database credentials (host `localhost`, port `5433`, user `postgres`, password `Pelopes818*`), **When** the application starts up, **Then** it overrides the defaults and successfully connects to the database at the specified port and host with the defined credentials.

---

### User Story 2 - Git Exclusion Validation (Priority: P2)

As a project security administrator, I want to ensure that `.env` files containing actual secrets are ignored by Git, so that credentials are never accidentally pushed to a shared repository.

**Why this priority**: High priority to prevent security breaches and credential leaks in public or private repositories.

**Independent Test**: Verify that checking in files ignores the `.env` file.

**Acceptance Scenarios**:

1. **Given** a `.env` file created in the project root or backend directory, **When** running a git status or staging command, **Then** the `.env` file is excluded from git tracking.

---

### User Story 3 - Environment Template (Priority: P2)

As a new developer joining the project, I want a `.env.example` template file, so that I can immediately see which environment variables are required and easily create my own `.env` file.

**Why this priority**: Eases developer onboarding and documents the schema of required variables without disclosing real credentials.

**Independent Test**: Verify `.env.example` exists in the repository, does not contain actual production/sensitive secrets, but provides placeholder keys and/or safe developer defaults.

**Acceptance Scenarios**:

1. **Given** a `.env.example` file in the codebase, **When** a developer copies it to `.env` and fills in the actual values, **Then** the application works without further adjustments.

---

### Edge Cases

- **Malformed/Invalid Config File**: What happens if the `.env` file is present but contains invalid or malformed variable definitions?
  - The application should report a clear syntax/parsing error on startup or fall back gracefully if syntax allows, preventing silent/unexpected failures.
- **Database Unreachable**: What happens if the database is not running on port 5433 or the credentials are wrong when the application tries to connect?
  - The application should log a clear connection failure message on startup, referencing the configured host/port, and exit gracefully with a non-zero status code rather than panicking.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST support reading configuration variables from a `.env` file located in the root or appropriate service directory at startup.
- **FR-002**: The `.env` file MUST support configuring database connection details including host, port, user, password, and database name (either via individual variables or a combined `DATABASE_URL`).
- **FR-003**: The codebase MUST contain a `.env.example` template file in the relevant directory containing the keys of all supported environment variables.
- **FR-004**: The project's Git configuration (e.g. `.gitignore`) MUST ignore any `.env` files containing actual values while tracking `.env.example`.
- **FR-005**: If the `.env` file is missing, the system MUST fall back to safe default configuration parameters.
- **FR-006**: The system MUST log an explicit error and exit gracefully if it fails to connect to the database with the configuration retrieved from the `.env` file.
- **FR-007**: Functional improvements like client registration fixes, user management/creation, and meeting/agenda management are out of scope for this feature and MUST be implemented in a separate branch.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can configure a custom local environment in under 1 minute by copying a template file.
- **SC-002**: Zero connection configuration parameters or credentials leak into the Git history.
- **SC-003**: The application successfully starts up and connects to a PostgreSQL database configured on port `5433` with user `postgres` and password `Pelopes818*` using environment variables.

## Assumptions

- The backend application is the main consumer of the `.env` configuration file.
- The `.gitignore` file already has a rule to exclude `.env*` files, which covers `.env`, `.env.local`, and other environment-specific configurations.
- The default database name is `onboardlyv2` as defined in the existing code.
