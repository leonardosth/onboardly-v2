# Onboardly - Client Onboarding Deployment Tracker

Onboardly is a modern customer onboarding and deployment management platform designed to track client integrations, schedule meetings, and display BI KPIs in real-time.

## Technology Stack

- **Backend**: Go (Golang) 1.22+ utilizing `go-chi` and `pgx` driver.
- **Frontend**: Vue.js 3 Single Page Application (SPA) using Vite, Pinia, Vue Router, and Axios.
- **Database**: PostgreSQL for strict referential integrity.

---

## Backend Setup & Installation

1. **Configure Environment Variables**:
   Create a `.env` file inside the `backend/` directory by copying the template:
   ```bash
   cp backend/.env.example backend/.env
   ```
   Adjust the values in `.env` as needed. If no `.env` file is present, the backend falls back to standard defaults:
   - `PORT`: Server port (default: `8080`)
   - `DATABASE_URL`: Connection string (default: `postgresql://postgres:postgres@localhost:5432/onboardlyv2?sslmode=disable`)
   - `JWT_SECRET`: Signing key for JWT (default: `super-secret-onboardly-key`)
   - `WEBHOOK_TOKEN`: Shared secret webhook token (default: `your-shared-secret-token`)

2. **Initialize Database Schema**:
   Run the schema script located in `backend/migrations/000001_init.up.sql` against your PostgreSQL server to create the tables and indexes.

3. **Start backend API server**:
   ```bash
   cd backend
   go mod tidy
   go run cmd/server/main.go
   ```

---

## Frontend Setup & Installation

1. **Install Node Dependencies**:
   ```bash
   cd frontend
   npm install
   ```

2. **Run Dev server**:
   ```bash
   npm run dev
   ```
   The application will boot at `http://localhost:5173`.

---

## API & Webhooks Ingestion

- **Auth endpoints**: `POST /api/auth/register`, `POST /api/auth/login`.
- **Google Sheets Webhook Sync**: `POST /api/webhooks/clients` authenticated via `X-Webhook-Token: your-shared-secret-token`.
- **Dashboard API**: `GET /api/dashboard` (retrieves real-time KPIs and audit activity feed).
