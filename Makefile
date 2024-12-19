# Makefile for CloudToggle Development Environment

.PHONY: up down run build clean

# Environment Variables
DB_CONTAINER_NAME=cloudtoggle_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cloudtoggle
DB_IMAGE=postgres:latest

WAIT_CMD := $(if $(findstring Windows, $(OS)), timeout /T 10 /NOBREAK, sleep 10)

# Start development environment
up:
	@echo "Starting development environment..."
	docker run --name $(DB_CONTAINER_NAME) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		-d $(DB_IMAGE)
	@echo "Waiting for PostgreSQL to be ready..."
	$(WAIT_CMD)
	@echo "Copying SQL migration files to container..."
	@docker cp migrations/001_create_tables.sql $(DB_CONTAINER_NAME):/tmp/001_create_tables.sql
	@echo "Running database migrations..."
	@docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -f /tmp/001_create_tables.sql
	@echo "Development environment is ready."

# Stop and remove development environment
down:
	@echo "Stopping and removing development containers..."
	docker stop $(DB_CONTAINER_NAME) || true
	docker rm $(DB_CONTAINER_NAME) || true
	@echo "Cleaned up development environment."

# Run the application
run:
	@echo "Running CloudToggle application..."
	go run ./cmd/main.go

# Build the application
build:
	@echo "Building CloudToggle binary..."
	go build -o cloudtoggle ./cmd/main.go
	@echo "Build complete."

# Clean up generated files and artifacts
clean:
	@echo "Cleaning up build artifacts..."
	@rm -f cloudtoggle
	@echo "Clean up complete."
