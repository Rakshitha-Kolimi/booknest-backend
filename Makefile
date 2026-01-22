-include .env
export

# Migration path
MIGRATIONS_DIR=internal/http/database/migrations
MIGRATE=migrate -path $(MIGRATIONS_DIR) -database $(DB_URL)


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
	migrate create -ext sql -dir $(MIGRATIONS_DIR)  $(name)
