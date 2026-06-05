# Research: Environment Variable and Database Configuration via .env

## Decision 1: Configuration Loading Mechanism
* **Chosen Approach**: Use the popular and lightweight Go library `github.com/joho/godotenv` to read and load the `.env` file into the process environment variables at startup.
* **Rationale**: `godotenv` is the de-facto standard in the Go ecosystem. It robustly parses comments, quoted strings, and exports variables, which prevents custom parsing bugs.
* **Alternatives Considered**: 
  * Custom parser: Rejected due to maintenance overhead and edge cases in parsing quotes/comments.
  * Loading via shell execution (e.g. `source .env`): Rejected because it makes execution platform-dependent (Windows vs Linux).

## Decision 2: Location of .env and .env.example
* **Chosen Approach**: Place both `.env` and `.env.example` in the `backend/` directory.
* **Rationale**: The backend service runs from the `backend/` working directory (as defined in [quickstart.md](file:///C:/Users/leona/Documents/Projects/onboardly-v2/specs/001-onboardly-core/quickstart.md)). Loading `.env` from the execution working directory is standard and matches existing deployment patterns.
* **Alternatives Considered**: 
  * Repository root: Rejected because the backend is structured as a subfolder and runs independently.

## Decision 3: Database URL Mapping
* **Chosen Approach**: Configure the database credentials directly via the `DATABASE_URL` environment variable format already supported by [config.go](file:///C:/Users/leona/Documents/Projects/onboardly-v2/backend/internal/config/config.go).
  * Format: `postgresql://postgres:Pelopes818*@localhost:5433/onboardlyv2?sslmode=disable`
* **Rationale**: Matches the existing backend config structure. Modifying the config loader to parse individual host/port/user/password fields would increase complexity without added benefit.
* **Alternatives Considered**:
  * Individual variables (`DB_HOST`, `DB_PORT`, etc.): Rejected to keep [config.go](file:///C:/Users/leona/Documents/Projects/onboardly-v2/backend/internal/config/config.go) unchanged and simpler.

## Decision 4: Git Exclusion
* **Chosen Approach**: Confirm and rely on the root [.gitignore](file:///C:/Users/leona/Documents/Projects/onboardly-v2/.gitignore) which already has a rule for `.env*` at line 25.
* **Rationale**: The rule `.env*` correctly covers `.env`, `.env.local`, etc., ensuring that developers do not accidentally commit local secrets.
