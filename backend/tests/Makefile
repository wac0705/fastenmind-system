# Test Makefile for FastenMind Backend

.PHONY: test test-unit test-integration test-coverage test-setup test-cleanup help

# Test database configuration
TEST_DB_HOST ?= localhost
TEST_DB_PORT ?= 5432
TEST_DB_USER ?= postgres
TEST_DB_PASSWORD ?= password
TEST_DB_NAME ?= fastenmind_test

# Go test flags
TEST_FLAGS ?= -v -race
COVERAGE_FLAGS ?= -coverprofile=coverage.out -covermode=atomic

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test-setup: ## Setup test database
	@echo "Setting up test database..."
	@docker run -d --name fastenmind-test-db \
		-e POSTGRES_USER=$(TEST_DB_USER) \
		-e POSTGRES_PASSWORD=$(TEST_DB_PASSWORD) \
		-e POSTGRES_DB=$(TEST_DB_NAME) \
		-p $(TEST_DB_PORT):5432 \
		postgres:13-alpine || true
	@echo "Waiting for database to be ready..."
	@sleep 5

test-cleanup: ## Cleanup test database
	@echo "Cleaning up test database..."
	@docker stop fastenmind-test-db || true
	@docker rm fastenmind-test-db || true

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	@go test $(TEST_FLAGS) ./unit/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@TEST_DB_HOST=$(TEST_DB_HOST) \
	 TEST_DB_PORT=$(TEST_DB_PORT) \
	 TEST_DB_USER=$(TEST_DB_USER) \
	 TEST_DB_PASSWORD=$(TEST_DB_PASSWORD) \
	 TEST_DB_NAME=$(TEST_DB_NAME) \
	 go test $(TEST_FLAGS) ./integration/...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test $(TEST_FLAGS) $(COVERAGE_FLAGS) ./unit/... ./integration/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-report: test-coverage ## Generate and open coverage report
	@echo "Opening coverage report..."
	@open coverage.html || xdg-open coverage.html || start coverage.html

test-watch: ## Run tests in watch mode (requires entr)
	@echo "Running tests in watch mode..."
	@find . -name "*.go" | entr -c make test

test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	@go test -bench=. -benchmem ./...

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@go test -race ./...

test-short: ## Run short tests only
	@echo "Running short tests..."
	@go test -short ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v ./...

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

ci: mod-tidy fmt vet lint test ## Run CI pipeline (format, vet, lint, test)

clean: ## Clean test artifacts
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@go clean -testcache

# Test specific services
test-auth: ## Test auth service only
	@go test $(TEST_FLAGS) ./unit/service/auth_service_test.go ./integration/auth_integration_test.go

test-quote: ## Test quote service only
	@go test $(TEST_FLAGS) ./unit/service/quote_service_test.go

test-handler: ## Test handlers only
	@go test $(TEST_FLAGS) ./unit/handler/...

# Development helpers
install-deps: ## Install test dependencies
	@echo "Installing test dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go mod download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy