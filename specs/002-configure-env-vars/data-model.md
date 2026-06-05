# Configuration Data Model: Environment Variables

Since this feature introduces no new database tables or relationships, this document defines the schema of the environment variables used to configure the Onboardly backend service.

## Configuration Schema

| Environment Variable | Description | Default Fallback Value | Example Value (Local Development) |
|----------------------|-------------|-------------------------|-----------------------------------|
| `PORT` | The port the Go backend HTTP server binds to. | `8080` | `8080` |
| `DATABASE_URL` | The PostgreSQL connection string containing credentials, host, port, database, and SSL options. | `postgresql://postgres:postgres@localhost:5432/onboardlyv2?sslmode=disable` | `postgresql://postgres:Pelopes818*@localhost:5433/onboardlyv2?sslmode=disable` |
| `JWT_SECRET` | Secret key used to sign and verify JWT authentication tokens. | `super-secret-onboardly-key` | `super-secret-onboardly-key` |
| `WEBHOOK_TOKEN` | Shared secret token required in the `X-Webhook-Token` header for Google Sheets sync. | `your-shared-secret-token` | `your-shared-secret-token` |

## Validation Rules
* **`PORT`**: Must be a valid numeric port between 1 and 65535.
* **`DATABASE_URL`**: Must be a valid PostgreSQL connection URI scheme (`postgresql://` or `postgres://`).
* **`JWT_SECRET`**: Must not be empty. In production environments, it should be a cryptographically strong random string.
* **`WEBHOOK_TOKEN`**: Must not be empty.
