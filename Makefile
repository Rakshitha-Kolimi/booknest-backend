-include .env
export

# Migration path
MIGRATE=migrate -path internal/http/database/migrations -database $(DB_URL)
MIGRATIONS_DIR=internal/http/database/migrations

# Commands
migrate-up:
	echo "Running migrations..."
	$(MIGRATE) up

migrate-down:
	echo "Remove migrations..."
	$(MIGRATE) down

migrate-force:
	echo "Forcing version..."
	$(MIGRATE) force $(VERSION)

migrate-version:
	echo "Fetching version..."
	$(MIGRATE) version

migrate-new:
	migrate create -ext sql -dir internal/http/database/migrations $(name)
