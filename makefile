# Makefile
.PHONY: dev build clean run test help swagger migrate-up migrate-down migrate-status migrate-create migrate-reset db-setup fetch-data

# Default target
.DEFAULT_GOAL := help

# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

dev: ## Start development server with hot reload
	@echo "Starting development server with hot reload..."
	air

build: ## Build the application
	@echo "Building application..."
	go build -o ./bin/main ./cmd/api/

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf ./bin/

run: ## Run the application without hot reload
	@echo "Running application..."
	go run ./cmd/api/

test: ## Run tests
	@echo "Running tests..."
	go test ./...

install: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go
	@echo "Fixing compatibility issues..."
	@grep -v "LeftDelim\|RightDelim" docs/docs.go > docs/docs.go.tmp && mv docs/docs.go.tmp docs/docs.go
	@echo "✅ Swagger docs generated at docs/"

migrate-up: ## Run database migrations
	@echo "Running database migrations..."
	@mkdir -p cmd/migrate/migrations
	goose -dir cmd/migrate/migrations postgres "$(DB_URL)" up
	@echo "✅ Migrations completed"

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	goose -dir cmd/migrate/migrations postgres "$(DB_URL)" down
	@echo "✅ Migration rolled back"

migrate-status: ## Check migration status
	@echo "Checking migration status..."
	goose -dir cmd/migrate/migrations postgres "$(DB_URL)" status

migrate-create: ## Create new migration (usage: make migrate-create NAME=create_users_table)
	@echo "Creating new migration: $(NAME)"
	@mkdir -p cmd/migrate/migrations
	goose -dir cmd/migrate/migrations create $(NAME) sql
	@echo "✅ Migration file created in cmd/migrate/migrations/"

migrate-reset: ## Reset database (careful!)
	@echo "⚠️  Resetting database - this will drop all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		goose -dir cmd/migrate/migrations postgres "$(DB_URL)" reset; \
		echo "✅ Database reset completed"; \
	else \
		echo "❌ Database reset cancelled"; \
	fi

db-setup: ## Create database if it doesn't exist
	@echo "Setting up database..."
	@createdb $(DB_NAME) 2>/dev/null || echo "Database $(DB_NAME) already exists"
	@echo "✅ Database setup completed"


fetch-data: ## Fetch all hotel data from Cupid API
	@echo "Fetching hotel data from Cupid API..."
	go run ./cmd/fetch/main.go
	@echo "✅ Data fetching completed"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t cupid-api .

help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort