# Scholarship Management System Makefile

.PHONY: build run test clean swagger-init swagger-gen docs dev install migrate migrate-down help migrate-status migrate-init build-migrate

# Variables
BINARY_NAME=scholarship-api
MAIN_PATH=./cmd/server
DOCS_PATH=./docs

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf $(DOCS_PATH)/*.go
	@rm -rf $(DOCS_PATH)/*.json
	@rm -rf $(DOCS_PATH)/*.yaml

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Install swag CLI tool
swagger-init:
	@echo "Installing swag CLI tool..."
	@go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger-gen:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs --parseInternal --parseDependency

# Generate docs (alias for swagger-gen)
docs: swagger-gen

# Development mode 
dev:
	@./dev.sh

# Database migration up
migrate:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go -action=migrate

# Database migration down (if you have down migrations)
migrate-down:
	@echo "Rolling back database migration..."
	@echo "Rollback migrations not implemented yet"

# Install air for development
install-air:
	@echo "Installing air for hot reload..."
	@go install github.com/cosmtrek/air@v1.40.4

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t scholarship-api .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 scholarship-api

docker-compose-up:
	@echo "Starting services with docker-compose..."
	@docker-compose up -d

docker-compose-down:
	@echo "Stopping docker-compose services..."
	@docker-compose down

# Check migration status
migrate-status:
	@echo "Checking migration status..."
	go run cmd/migrate/main.go -action=status

# Initialize migration table
migrate-init:
	@echo "Initializing migration table..."
	go run cmd/migrate/main.go -action=init

# Build migration binary
build-migrate:
	@echo "Building migration binary..."
	go build -o bin/migrate cmd/migrate/main.go

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install dependencies"
	@echo "  swagger-init   - Install swag CLI tool"
	@echo "  swagger-gen    - Generate Swagger documentation"
	@echo "  docs           - Alias for swagger-gen"
	@echo "  dev            - Start development server with hot reload"
	@echo "  migrate        - Run database migrations"
	@echo "  migrate-down   - Rollback database migration"
	@echo "  install-air    - Install air for hot reload"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-compose-up   - Start services with docker-compose"
	@echo "  docker-compose-down - Stop docker-compose services"
	@echo "  migrate-status - Check migration status"
	@echo "  migrate-init   - Initialize migration table"
	@echo "  build-migrate  - Build migration binary"
	@echo "  help           - Show this help message"

# Default target
default: help