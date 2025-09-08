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

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-coverage-func: ## Run tests with function coverage
	@echo "Running tests with function coverage..."
	go test -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func=coverage.out

test-unit: ## Run only unit tests
	@echo "Running unit tests..."
	go test -short ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -run Integration ./...

test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	go test -bench=. ./...

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -race ./...

test-clean: ## Clean test artifacts
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html
	@echo "✅ Test artifacts cleaned"

sync-start: ## Start sync service with 12h interval
	@echo "Starting sync service with 12h interval..."
	@curl -X POST "http://localhost:8080/api/v1/admin/sync/start?interval=12h" || echo "❌ Failed to start sync service"

sync-start-24h: ## Start sync service with 24h interval
	@echo "Starting sync service with 24h interval..."
	@curl -X POST "http://localhost:8080/api/v1/admin/sync/start?interval=24h" || echo "❌ Failed to start sync service"

sync-stop: ## Stop sync service
	@echo "Stopping sync service..."
	@curl -X POST "http://localhost:8080/api/v1/admin/sync/stop" || echo "❌ Failed to stop sync service"

sync-now: ## Trigger immediate sync
	@echo "Triggering immediate sync..."
	@curl -X POST "http://localhost:8080/api/v1/admin/sync" || echo "❌ Failed to trigger sync"

sync-status: ## Check sync status
	@echo "Checking sync status..."
	@curl -s "http://localhost:8080/api/v1/admin/sync/status" | jq '.' || echo "❌ Failed to get sync status"

sync-health: ## Check sync health
	@echo "Checking sync health..."
	@curl -s "http://localhost:8080/api/v1/admin/sync/health" | jq '.' || echo "❌ Failed to get sync health"

sync-logs: ## Get sync logs
	@echo "Getting sync logs..."
	@curl -s "http://localhost:8080/api/v1/admin/sync/logs" | jq '.' || echo "❌ Failed to get sync logs"

sync-settings: ## Get sync settings
	@echo "Getting sync settings..."
	@curl -s "http://localhost:8080/api/v1/admin/sync/settings" | jq '.' || echo "❌ Failed to get sync settings"

test-integration-cupid: ## Run Cupid API integration tests
	@echo "Running Cupid API integration tests..."
	@if [ ! -f integration.env ]; then \
		echo "❌ integration.env not found. Copy integration.env.example to integration.env and set your API credentials."; \
		exit 1; \
	fi
	@echo "Loading integration test environment..."
	@export $$(grep -v '^#' integration.env | grep -v '^$$' | xargs) && go test -v -tags=integration ./internal/cupid/... -run "TestCupid.*Integration"

test-integration-cupid-connectivity: ## Test Cupid API connectivity only
	@echo "Testing Cupid API connectivity..."
	@if [ ! -f integration.env ]; then \
		echo "❌ integration.env not found. Copy integration.env.example to integration.env and set your API credentials."; \
		exit 1; \
	fi
	@echo "Loading integration test environment..."
	@export $$(grep -v '^#' integration.env | grep -v '^$$' | xargs) && go test -v -tags=integration ./internal/cupid/... -run "TestCupidAPIConnectivity"

test-integration-cupid-validation: ## Test Cupid API data validation
	@echo "Testing Cupid API data validation..."
	@if [ ! -f integration.env ]; then \
		echo "❌ integration.env not found. Copy integration.env.example to integration.env and set your API credentials."; \
		exit 1; \
	fi
	@echo "Loading integration test environment..."
	@export $$(grep -v '^#' integration.env | grep -v '^$$' | xargs) && go test -v -tags=integration ./internal/cupid/... -run "TestCupidDataValidation"

test-integration-cupid-performance: ## Test Cupid API performance
	@echo "Testing Cupid API performance..."
	@if [ ! -f integration.env ]; then \
		echo "❌ integration.env not found. Copy integration.env.example to integration.env and set your API credentials."; \
		exit 1; \
	fi
	@echo "Loading integration test environment..."
	@export $$(grep -v '^#' integration.env | grep -v '^$$' | xargs) && go test -v -tags=integration ./internal/cupid/... -run "TestCupidPerformance"

benchmark-cupid: ## Benchmark Cupid API performance
	@echo "Benchmarking Cupid API performance..."
	@if [ ! -f integration.env ]; then \
		echo "❌ integration.env not found. Copy integration.env.example to integration.env and set your API credentials."; \
		exit 1; \
	fi
	@echo "Loading integration test environment..."
	@export $$(grep -v '^#' integration.env | grep -v '^$$' | xargs) && go test -v -tags=integration -bench=BenchmarkCupidAPI ./internal/cupid/...

test-integration-all: ## Run all integration tests
	@echo "Running all integration tests..."
	@make test-integration-cupid

test-integration-setup: ## Setup integration test environment
	@echo "Setting up integration test environment..."
	@if [ ! -f integration.env ]; then \
		echo "📋 Creating integration.env from template..."; \
		cp integration.env.example integration.env; \
		echo "✅ integration.env created. Please edit it and set your API credentials."; \
		echo "🔑 Required: CUPID_API_KEY"; \
		echo "📝 Optional: CUPID_API_BASE_URL, CUPID_API_VERSION"; \
	else \
		echo "✅ integration.env already exists."; \
	fi

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
	@if [ -f "./fetch" ]; then \
		./fetch; \
	else \
		go run ./cmd/fetch/main.go; \
	fi
	@echo "✅ Data fetching completed"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t cupid-api .

help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort