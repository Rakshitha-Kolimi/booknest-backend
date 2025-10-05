# Database connection URL
DB_URL=postgres://postgres:postgres@localhost:5432/booknest?sslmode=disable

# Migration path
MIGRATE=migrate -path internal/http/database/migrations -database $(DB_URL)

# Commands
migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-force:
	$(MIGRATE) force VERSION

migrate-new:
	migrate create -ext sql -dir internal/http/database/migrations -seq $(name)
