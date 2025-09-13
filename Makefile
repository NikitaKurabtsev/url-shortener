DB_URL=sqlite3://storage/storage.db
SCHEMA_DIR=./schema

migrate-up:
	migrate -path $(SCHEMA_DIR) -database $(DB_URL) up

migrate-down:
	migrate -path $(SCHEMA_DIR) -database $(DB_URL) down -all

migrate-version:
	migrate -path $(SCHEMA_DIR) -database $(DB_URL) version

run:
	CONFIG_PATH=./config/local.yml go run cmd/url-shortener/main.go
