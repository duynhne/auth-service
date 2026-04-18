# auth-service

Authentication microservice for user login, registration, and JWT token management.

## Features

- User login with JWT tokens
- User registration
- Token validation
- Session management

## API Endpoints

All routes follow Variant A naming — single path for browser and in-cluster callers. See [homelab naming convention](https://github.com/duynhlab/homelab/blob/main/docs/api/api-naming-convention.md).

| Method | Path | Audience |
|--------|------|----------|
| `POST` | `/auth/v1/public/login` | public |
| `POST` | `/auth/v1/public/register` | public |
| `GET` | `/auth/v1/private/me` | private |

- Browser: `https://gateway.duynhne.me/auth/v1/…`
- Service-to-service (JWT validation): `http://auth.auth.svc.cluster.local:8080/auth/v1/private/me`

## Tech Stack

- Go + Gin framework
- PostgreSQL 17 (auth-db cluster, HA)
- PgBouncer connection pooling
- OpenTelemetry tracing

## Development

### Prerequisites

- Go 1.25+
- [golangci-lint](https://golangci-lint.run/welcome/install/) v2+

### Local Development

```bash
# Install dependencies
go mod tidy
go mod download

# Build
go build ./...

# Test
go test ./...

# Lint (must pass before PR merge)
golangci-lint run --timeout=10m

# Run locally (requires .env or env vars)
go run cmd/main.go
```

### Pre-push Checklist

```bash
go build ./... && go test ./... && golangci-lint run --timeout=10m
```

## License

MIT
