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
│   │       └── user.go      # Domain models (User, LoginRequest, etc.)
│   ├── logic/v1/
│   │   ├── service.go       # Business logic layer
│   │   └── errors.go        # Domain errors
│   └── web/v1/
│       └── handler.go       # HTTP handlers (Gin)
├── middleware/
│   ├── logging.go           # Request logging middleware
│   ├── profiling.go         # Pyroscope integration
│   ├── prometheus.go        # Metrics middleware
│   ├── resource.go          # Resource limits
│   └── tracing.go           # OpenTelemetry tracing
├── go.mod
├── go.sum
└── Dockerfile
```

## 🔌 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/login` | User login, returns JWT token |
| `POST` | `/api/v1/auth/register` | User registration |
| `GET` | `/api/v1/auth/me` | Get current user from token |
| `GET` | `/health` | Health check |
| `GET` | `/ready` | Readiness probe (fails during shutdown) |
| `GET` | `/metrics` | Prometheus metrics |

## 🔧 Tech Stack

| Component | Technology |
|-----------|------------|
| **Framework** | Gin v1.11 |
| **Database** | PostgreSQL via pgx/v5 |
| **Logging** | Zerolog (from `github.com/duynhne/pkg`) |
| **Tracing** | OpenTelemetry with OTLP exporter |
| **Metrics** | Prometheus client |
| **Profiling** | Pyroscope |
| **Passwords** | bcrypt |

## 📦 Dependencies

- `github.com/duynhne/pkg` - Shared logger package
- `github.com/gin-gonic/gin` - HTTP framework
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `go.opentelemetry.io/otel` - Distributed tracing
- `github.com/prometheus/client_golang` - Metrics

## 🛠️ Development

### Prerequisites

- Go 1.25+
- PostgreSQL 15+
- Docker (optional)

### Local Build & Run

```bash
# Download dependencies
go mod download

# Run tests
go test -v ./...

# Build binary
go build -o auth-service ./cmd/main.go

# Run (requires PostgreSQL)
export DATABASE_URL="postgres://user:pass@localhost:5432/auth_db?sslmode=disable"
./auth-service
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVICE_NAME` | `auth-service` | Service identifier |
| `SERVICE_PORT` | `8080` | HTTP port |
| `LOG_LEVEL` | `info` | Logging level |
| `DATABASE_URL` | - | PostgreSQL connection string |
| `TRACING_ENABLED` | `false` | Enable OpenTelemetry |
| `TRACING_ENDPOINT` | - | OTLP endpoint |
| `PROFILING_ENABLED` | `false` | Enable Pyroscope |

### Docker Build

```bash
docker build -t auth-service -f Dockerfile .
docker run -p 8080:8080 -e DATABASE_URL="..." auth-service
```

## 🚀 CI/CD

Uses reusable GitHub Actions from [shared-workflows](https://github.com/duyhenryer/shared-workflows):

- **go-check.yml** - Tests and linting
- **sonarqube.yml** - SonarCloud analysis
- **docker-build.yml** - Build and push to GHCR

## 📐 Code Patterns

- **Layered architecture**: `handler` → `service` → `database`
- **Context-based tracing**: OpenTelemetry spans propagate through layers
- **Graceful shutdown**: Readiness probe fails first, then drain delay
- **Domain errors**: Custom error types (ErrUserNotFound, ErrInvalidCredentials)

## 🔗 Related Services

- Uses **pkg** for shared logging
- Authenticates users for all other services
