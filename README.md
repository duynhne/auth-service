# auth-service

Authentication microservice for user login, registration, and JWT token management.

## Features

- User login with JWT tokens
- User registration
- Token validation
- Session management

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/login` | User login |
| `POST` | `/api/v1/auth/register` | User registration |
| `GET` | `/api/v1/auth/me` | Get current user |

## Tech Stack

- Go + Gin framework
- PostgreSQL 17 (auth-db cluster, HA)
- PgBouncer connection pooling
- OpenTelemetry tracing

## Development

```bash
go mod download
go test ./...
go run cmd/main.go
```

## License

MIT
