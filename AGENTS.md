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
