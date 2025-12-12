.PHONY: help run build test clean swagger install db-create db-drop db-reset

# Default target
help:
	@echo "Available commands:"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo "  make install       - Install dependencies"
	@echo "  make db-create     - Create database"
	@echo "  make db-drop       - Drop database"
	@echo "  make db-reset      - Reset database (drop and create)"

# Run the application
run:
	@echo "Starting application..."
	go run cmd/api/main.go

# Build the application
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf tmp/
	go clean

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	swag init -g cmd/api/main.go -o docs

# Install/Update dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	go mod verify

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Database: Create
db-create:
	@echo "Creating database..."
	psql -U postgres -c "CREATE DATABASE taskmanagement;"

# Database: Drop
db-drop:
	@echo "Dropping database..."
	psql -U postgres -c "DROP DATABASE IF EXISTS taskmanagement;"

# Database: Reset (drop and create)
db-reset: db-drop db-create
	@echo "Database reset complete!"

# Development: Run with hot reload (requires air)
dev:
	@echo "Starting development server with hot reload..."
	air

# Install development tools
dev-tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t taskmanagement-api:latest .

docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down

docker-logs:
	@echo "Showing logs..."
	docker-compose logs -f api

docker-restart:
	@echo "Restarting services..."
	docker-compose restart

docker-clean:
	@echo "Cleaning up Docker resources..."
	docker-compose down -v
	docker system prune -f

docker-dev-up:
	@echo "Starting development PostgreSQL..."
	docker-compose -f docker-compose.dev.yml up -d

docker-dev-down:
	@echo "Stopping development PostgreSQL..."
	docker-compose -f docker-compose.dev.yml down

docker-shell:
	@echo "Opening shell in API container..."
	docker exec -it taskapi-app sh

docker-db-shell:
	@echo "Opening PostgreSQL shell..."
	docker exec -it taskapi-postgres psql -U postgres -d taskmanagement
