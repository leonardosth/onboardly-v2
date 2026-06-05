# Research: Onboardly Core Architecture

This document outlines the architectural decisions and technology stack selections for the Onboardly core platform.

## Decision: Technology Stack & Core Libraries

We select a decoupled web application structure with a Go backend and a Vue.js 3 frontend.

### Go Backend (API Service)
- **Framework/Router**: `go-chi/chi` v5. Preferred for its lightweight nature, standard library context compatibility, and zero-allocation routing.
- **Database Driver**: `jackc/pgx` v5 (via database/sql compatibility). Offers performance optimizations and native PostgreSQL type support.
- **Password Hashing**: `golang.org/x/crypto/bcrypt`.
- **Session Tokens**: `golang-jwt/jwt` v5 for generating and validating cryptographically secure JWTs.

### Vue.js 3 Frontend (Single Page Application)
- **Build Tool**: Vite (fast developer setup and optimized production bundling).
- **State Management**: Pinia (lightweight, modular, type-safe).
- **Routing**: Vue Router.
- **Styling**: Vanilla CSS with modern flex/grid layouts.

### Database
- **PostgreSQL**: Standard relational database to enforce critical constraints (CNPJ uniqueness, foreign keys on meetings, etc.).

---

## Rationale

1. **High Concurrency and Low Latency**: Go compiles to a native binary, allowing us to easily meet the <500ms login latency (SC-001) and <1.5s webhook processing (SC-002) success criteria.
2. **Referential Integrity**: PostgreSQL native foreign key constraints guarantee that meetings cannot reference non-existent projects (FR-014), fulfilling strict database referential integrity at the engine level.
3. **Responsive and Active UI**: Vue 3's reactive Composition API enables simple state management (using Pinia) to feed the live Dashboard metrics in real-time.

---

## Alternatives Considered

### Node.js (Express) vs Go
- **Evaluation**: Node.js is widely used, but Go provides superior CPU performance and true multi-threading out of the box, which ensures stable processing under concurrent webhook loads.
- **Decision**: Go selected to support the low-latency target.

### React vs Vue.js 3
- **Evaluation**: React is highly common but has more boilerplates. Vue 3 with Vite offers a faster startup cycle and simpler reactivity models.
- **Decision**: Vue 3 selected for ease of development.

### MongoDB vs PostgreSQL
- **Evaluation**: MongoDB is flexible for unstructured data, but Onboardly's core domain models (Clients, Projects, Meetings) have highly relational structures. Without strict foreign keys, preventing orphan records becomes difficult.
- **Decision**: PostgreSQL selected to enforce referential integrity.
