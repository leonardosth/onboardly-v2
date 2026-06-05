# Contract: Environment Configuration Interface

This document specifies the interface contract between the environment (hosting environment, `.env` file) and the Go backend application.

## 1. Input Specification: The `.env` File
The application reads configuration from a `.env` file at startup. The format of the file must follow standard key-value pairs separated by `=` sign:

```env
PORT=8080
DATABASE_URL=postgresql://postgres:Pelopes818*@localhost:5433/onboardlyv2?sslmode=disable
JWT_SECRET=super-secret-onboardly-key
WEBHOOK_TOKEN=your-shared-secret-token
```

### Constraints:
* **Whitespace**: Leading or trailing whitespace in values must be trimmed.
* **Quotes**: Values wrapped in single (`'`) or double (`"`) quotes should be parsed correctly, stripping the outermost quotes.
* **Comments**: Lines starting with `#` are treated as comments and ignored.

## 2. API Context Configuration Interface
The loaded values are bound to the internal configuration struct used across the application:

```go
type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	WebhookToken string
}
```

This contract ensures that all modules (database connection pool, router, authentication middleware, and webhooks) obtain their parameters through this structured configuration object rather than invoking `os.Getenv` dynamically during execution.
