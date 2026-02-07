# auth-service

> AI Agent context for understanding this repository

## 📋 Overview

Authentication microservice for the monitoring platform. Provides user login, registration, and session management via REST API.

## 🏗️ Architecture

```
auth-service/
├── cmd/main.go              # Entry point, graceful shutdown
├── config/config.go         # Environment-based configuration
├── db/migrations/sql/       # Flyway SQL migrations
├── internal/
│   ├── core/
│   │   ├── database.go      # PostgreSQL connection pool (pgx)
│   │   └── domain/user.go   # Domain models
│   ├── logic/v1/
│   │   ├── service.go       # Business logic layer
│   │   └── errors.go        # Domain errors
│   └── web/v1/handler.go    # HTTP handlers (Gin)
├── middleware/
└── Dockerfile
```

## 🔌 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/login` | User login, returns JWT token |
| `POST` | `/api/v1/auth/register` | User registration |
| `GET` | `/api/v1/auth/me` | Get current user from token |

## 📐 3-Layer Architecture

| Layer | Location | Responsibility |
|-------|----------|----------------|
| **Web** | `internal/web/v1/handler.go` | HTTP handling, validation, error translation |
| **Logic** | `internal/logic/v1/service.go` | Business rules (❌ NO SQL) |
| **Core** | `internal/core/` | Domain models, repositories, database |

## 🗄️ Database

| Component | Value |
|-----------|-------|
| **Cluster** | auth-db (Zalando Postgres Operator) |
| **PostgreSQL** | 17 |
| **HA** | 3 nodes (1 leader + 2 standbys) |
| **Pooler** | PgBouncer Sidecar (2 instances) |
| **Endpoint** | `auth-db-pooler.auth.svc.cluster.local:5432` |

**Dual Connection Pattern:**
- **Main container**: PgBouncer (`auth-db-pooler:5432`)
- **Init container**: Direct (`auth-db:5432`) - for DDL migrations

## 🚀 Production Patterns

### Graceful Shutdown
VictoriaMetrics pattern: `/ready` → 503 → drain delay (5s) → sequential cleanup (HTTP → DB → Tracer)

## 🔧 Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | Gin |
| Database | PostgreSQL 17 via pgx/v5 |
| Logging | Zerolog |
| Tracing | OpenTelemetry |
| Passwords | bcrypt |

## Code Quality

**MANDATORY**: All code changes MUST pass lint before committing.

- Linter: `golangci-lint` v2+ with `.golangci.yml` config (60+ linters enabled)
- Zero tolerance: PRs with lint errors will NOT be merged
- CI enforces: `go-check` job runs lint on every PR

### Commands (run in order)

```bash
go mod tidy              # Clean dependencies
go build ./...           # Verify compilation
go test ./...            # Run tests
golangci-lint run --timeout=10m  # Lint (MUST pass)
```

### Pre-commit One-liner

```bash
go build ./... && go test ./... && golangci-lint run --timeout=10m
```

### Common Lint Fixes

- `perfsprint`: Use `errors.New()` instead of `fmt.Errorf()` when no format verbs
- `nosprintfhostport`: Use `net.JoinHostPort()` instead of `fmt.Sprintf("%s:%s", host, port)`
- `errcheck`: Always check error returns (or explicitly `_ = fn()`)
- `goconst`: Extract repeated string literals to constants
- `gocognit`: Extract helper functions to reduce complexity
- `noctx`: Use `http.NewRequestWithContext()` instead of `http.NewRequest()`

## 3-Layer Coding Rules

**CRITICAL**: Strict layer boundaries. Violations will be rejected in code review.

### Layer Boundaries

| Layer | Location | ALLOWED | FORBIDDEN |
|-------|----------|---------|-----------|
| **Web** | `internal/web/v1/` | HTTP handling, JSON binding, DTO mapping, call Logic, aggregation | SQL queries, direct DB access, business rules |
| **Logic** | `internal/logic/v1/` | Business rules, call repository interfaces, domain errors | SQL queries, `database.GetPool()`, HTTP handling, `*gin.Context` |
| **Core** | `internal/core/` | Domain models, repository implementations, SQL queries, DB connection | HTTP handling, business orchestration |

### Dependency Direction

```
Web -> Logic -> Core (one-way only, never reverse)
```

- Web imports Logic and Core/domain
- Logic imports Core/domain and Core/repository interfaces
- Core imports nothing from Web or Logic

### DO

- Put HTTP handlers, request validation, error-to-status mapping in `web/`
- Put business rules, orchestration, transaction logic in `logic/`
- Put SQL queries in `core/repository/` implementations
- Use repository interfaces (defined in `core/domain/`) for data access in Logic layer
- Use dependency injection (constructor parameters) for all service dependencies

### DO NOT

- Write SQL or call `database.GetPool()` in Logic layer
- Import `gin` or handle HTTP in Logic layer
- Put business rules in Web layer (Web only translates and delegates)
- Call Logic functions directly from another service (use HTTP aggregation in Web layer)
- Skip the Logic layer (Web must not call Core/repository directly)
