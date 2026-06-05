# Implementation Plan: Environment Variable and Database Configuration via .env

**Branch**: `002-configure-env-vars` | **Date**: 2026-06-04 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/002-configure-env-vars/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

The goal of this feature is to secure local application configuration by extracting credentials and environment variables into a dedicated `.env` file, ensuring that database configurations (port `5433`, password `Pelopes818*`, etc.) are read at startup. We will implement this in the Go backend using the standard `github.com/joho/godotenv` library and verify that Git ignores the sensitive `.env` file via the existing root `.gitignore` configuration.

## Technical Context

**Language/Version**: Go 1.22+ (Backend)

**Primary Dependencies**: `github.com/joho/godotenv` v1.5.1

**Storage**: PostgreSQL (configured on port `5433` with user `postgres` and password `Pelopes818*`)

**Testing**: `go test` (Backend configuration tests)

**Target Platform**: Go Runtime / Linux Server (Backend)

**Project Type**: Backend service configuration component

**Performance Goals**: File loading and configuration binding on startup under 50ms.

**Constraints**: Sensitive credentials must not be committed to Git. Failures to parse config or connect to the database must log descriptive errors and exit gracefully.

**Scale/Scope**: Go Backend configuration management.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Gate 1: Access Control & Cryptography**: Verified. Centralizing secret keys like `JWT_SECRET` in an external environment configuration prevents hardcoding of security credentials.
- **Gate 2: Database Referential Integrity**: N/A for this configuration-only feature.
- **Gate 3: Universal Chronology**: N/A.
- **Gate 4: Integration Security**: Verified. External webhook token `WEBHOOK_TOKEN` will be read from the environment variables, avoiding plain-text exposure in code files.
- **Gate 5: Real-time Metrics**: N/A for this configuration-only feature.

*Post-Design Re-Check (Phase 1 complete)*: All gates verified. No violations introduced by design artifacts.

## Scope Clarifications

Per session clarification (2026-06-04) and **FR-007**: Functional improvements such as client registration fixes, user management/creation UI, and meeting/agenda management are explicitly **out of scope** for this feature and must be implemented in a separate branch.

## Project Structure

### Documentation (this feature)

```text
specs/002-configure-env-vars/
├── plan.md              # This file
├── research.md          # Decision log
├── data-model.md        # Environment schema and validation
├── quickstart.md        # Scenario verification guide
└── contracts/
    └── env-variables.md # Input & struct contracts
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── server/
│       └── main.go      # Add godotenv.Load() call
├── internal/
│   └── config/
│       └── config.go    # Unchanged LookupEnv calls (will pick loaded vars)
├── .env                 # Local env overrides (ignored by Git)
├── .env.example         # Template file for developers
├── go.mod               # Add godotenv dependency
└── go.sum
```

**Structure Decision**: Web application layout. The configuration modifications will target the Go backend codebase (`backend/` subfolder), adding `.env` and `.env.example` in the backend root directory.

## Complexity Tracking

No constitution violations detected.

