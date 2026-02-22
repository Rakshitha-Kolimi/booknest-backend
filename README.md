# BookNest Platform

Backend API for BookNest (Go + Gin + PostgreSQL + GORM/pgx).

## Architecture

- `main.go`, `main_setup.go`: bootstrap + dependency wiring
- `internal/http/controller`: HTTP handlers and route registration
- `internal/service`: business logic
- `internal/repository`: DB access
- `internal/domain`: entities, enums, interfaces
- `internal/http/database/migrations`: SQL migrations
- `internal/middleware`: JWT auth, admin role guard, error/logging middleware

## Prerequisites

- Go 1.24+
- PostgreSQL 15+
- `migrate` CLI (for local DB migrations)

## Environment

Create/update `.env` in this folder with:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=booknest
DB_NAME=booknest
DB_PORT=5432
DB_URL=postgres://postgres:booknest@localhost:5432/booknest?sslmode=disable
JWT_SECRET=booknest_secret
SWAGGER_USER=booknest
SWAGGER_PASSWORD=<your-password>
```

Note: `JWT_AUTH_SECRET` is still supported for backward compatibility, but `JWT_SECRET` is the primary key.

## Run (Interview-Safe)

From this folder:

```bash
go mod download
go test ./...
go run .
```

- API: `http://localhost:8080`
- Health: `GET /health`
- Swagger: `http://localhost:8080/swagger/index.html`

## Migrations

Using Makefile targets:

```bash
make migrate-up
make migrate-down
make migrate-version
```

## Docker Option

```bash
docker compose up --build
```

Starts API + Postgres and runs migrations via `entrypoint.sh`.
