# auth-service

> AI Agent context for understanding this repository

## 📋 Overview

Authentication microservice for the monitoring platform. Provides user login, registration, and session management via REST API.

## 🏗️ Architecture

```
auth-service/
├── cmd/
│   └── main.go              # Entry point, graceful shutdown
├── config/
│   └── config.go            # Environment-based configuration
├── db/migrations/
│   └── sql/                  # Flyway SQL migrations
├── internal/
│   ├── core/
│   │   ├── database.go      # PostgreSQL connection pool (pgx)
│   │   └── domain/
│   │       └── user.go      # Domain models
│   ├── logic/v1/
│   │   ├── service.go       # Business logic layer
│   │   └── errors.go        # Domain errors
│   └── web/v1/
│       └── handler.go       # HTTP handlers (Gin)
├── middleware/
│   ├── logging.go           # Request logging
│   ├── profiling.go         # Pyroscope integration
│   ├── prometheus.go        # Metrics middleware
│   └── tracing.go           # OpenTelemetry tracing
└── Dockerfile
```

## 🔌 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/login` | User login, returns JWT token |
| `POST` | `/api/v1/auth/register` | User registration |
| `GET` | `/api/v1/auth/me` | Get current user from token |
| `GET` | `/health` | Liveness probe (always 200) |
| `GET` | `/ready` | Readiness probe (503 during shutdown) |
| `GET` | `/metrics` | Prometheus metrics |

## 📐 3-Layer Architecture

| Layer | Location | Responsibility |
|-------|----------|----------------|
| **Web** | `internal/web/v1/handler.go` | HTTP handling, validation, DTO mapping, error translation |
| **Logic** | `internal/logic/v1/service.go` | Business rules, transaction orchestration |
| **Core** | `internal/core/` | Domain models, repository implementations, database |

**Constraints:**
- Web calls Logic only (not Core directly)
- Logic Layer: ❌ NO SQL queries, ❌ NO `database.GetDB()`, ❌ NO HTTP handling
- Core owns all database queries

## 🗄️ Database

| Component | Value |
|-----------|-------|
| **Cluster** | auth-db (Zalando Postgres Operator) |
| **PostgreSQL** | 17 |
| **HA** | 3 nodes (1 leader + 2 standbys) |
| **Pooler** | PgBouncer Sidecar (2 instances) |
| **Endpoint** | `auth-db-pooler.auth.svc.cluster.local:5432` |
| **Pool Mode** | Transaction |
| **Driver** | pgx/v5 (SimpleProtocol mode) |

**Dual Connection Pattern:**
- **Main container**: PgBouncer (`auth-db-pooler:5432`) - for transactions
- **Init container**: Direct (`auth-db:5432`) - for DDL migrations (no pooler)

## 🚀 Graceful Shutdown

**VictoriaMetrics Pattern:**
1. Signal received → `isShuttingDown.Store(true)`
2. `/ready` returns 503 → K8s stops routing traffic
3. Sleep `READINESS_DRAIN_DELAY` (5s) → propagation delay
4. Sequential cleanup: HTTP Server → Database → Tracer

**Config:**
- `SHUTDOWN_TIMEOUT`: 10s (default)
- `READINESS_DRAIN_DELAY`: 5s (default)
- `terminationGracePeriodSeconds`: 30

## 🔧 Tech Stack

| Component | Technology |
|-----------|------------|
| **Framework** | Gin v1.11 |
| **Database** | PostgreSQL 17 via pgx/v5 |
| **Logging** | Zerolog (from `github.com/duynhne/pkg`) |
| **Tracing** | OpenTelemetry with OTLP exporter |
| **Metrics** | Prometheus client |
| **Profiling** | Pyroscope |
| **Passwords** | bcrypt |

## 🛠️ Development

```bash
go mod download
go test -v ./...
go build -o auth-service ./cmd/main.go
```

## 🚀 CI/CD

Uses reusable GitHub Actions from [shared-workflows](https://github.com/duyhenryer/shared-workflows):
- `go-check.yml` - Tests and linting
- `sonarqube.yml` - SonarCloud analysis
- `docker-build.yml` - Build and push to GHCR
