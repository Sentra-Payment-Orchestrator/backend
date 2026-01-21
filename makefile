# Create database url variable
-include .env
export

.PHONY: migrate

# Migrate the database up to the latest version
migrate-up:
	migrate -path ./migrations -database $(DATABASE_URL) up

# Migrate the database down by one version
migrate-down:
	migrate -path ./migrations -database $(DATABASE_URL) down 1

migrate-drop:
	migrate -path ./migrations -database $(DATABASE_URL) drop

# Force the database to a specific version
# take version as an argument (e.g., make migrate-force version=2)
migrate-force:
	migrate -path ./migrations -database $(DATABASE_URL) force $(version)

# Show the current version of the database
migrate-version:
	migrate -path ./migrations -database $(DATABASE_URL) version

# Create a new migration file
# take name as an argument (e.g., make create-migration name=create_users_table)
create-migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

run-dev:
	go run *.go