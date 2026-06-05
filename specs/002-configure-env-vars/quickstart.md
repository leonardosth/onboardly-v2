# Quickstart Validation Guide: Environment & Database Config

This document outlines the steps to verify the end-to-end setup and functionality of the `.env` configuration implementation.

## Prerequisites

- **Go**: Version 1.22+ installed.
- **Git**: Installed and initialized in the repository.
- **PostgreSQL**: Running locally on port `5433` with a database named `onboardlyv2`, username `postgres`, and password `Pelopes818*`.

---

## Validation Scenarios

### Scenario 1: Setup `.env` from Template
1. Navigate to the `backend` directory:
   ```bash
   cd backend
   ```
2. Copy the `.env.example` template to `.env`:
   ```bash
   copy .env.example .env
   ```
3. Edit the `.env` file to contain the requested credentials:
   ```env
   PORT=8080
   DATABASE_URL=postgresql://postgres:Pelopes818*@localhost:5433/onboardlyv2?sslmode=disable
   JWT_SECRET=super-secret-onboardly-key
   WEBHOOK_TOKEN=your-shared-secret-token
   ```

*Expected Outcome*: The `.env` file is successfully created with correct parameters.

---

### Scenario 2: Verification of Startup & Database Connection
1. Run the Go backend server:
   ```bash
   go run cmd/server/main.go
   ```
2. Verify the log output.

*Expected Outcome*:
The logs must show that the server connected successfully to PostgreSQL on port `5433` and did not crash or panic:
```text
Connecting to PostgreSQL at postgresql://postgres:Pelopes818*@localhost:5433/onboardlyv2?sslmode=disable...
Database connection pool established successfully.
Starting HTTP server on port 8080...
```

---

### Scenario 3: Git Exclusion Verification
1. From the project root, run:
   ```bash
   git status
   ```

*Expected Outcome*:
The newly created `backend/.env` file MUST NOT appear under "Untracked files" or "Changes not staged for commit", verifying that it is correctly ignored by Git.
The template file `backend/.env.example` should be tracked (if newly added).
