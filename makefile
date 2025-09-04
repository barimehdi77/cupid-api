# Makefile
.PHONY: dev build clean run test help

# Default target
.DEFAULT_GOAL := help

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

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t cupid-api .

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'