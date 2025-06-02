.PHONY: help build test lint clean run-auth run-product run-cart run-order run-payment run-gateway docker-build docker-up docker-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build parameters
BINARY_DIR=bin
DOCKER_REGISTRY=gobazaar

# Services
SERVICES=auth product cart order payment gateway

# Default target
help: ## Show this help message
	@echo 'GoBazaar - E-commerce Microservices Platform'
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build all services
build: clean ## Build all microservices
	@echo "🔨 Building all services..."
	@mkdir -p $(BINARY_DIR)
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		$(GOBUILD) -o $(BINARY_DIR)/$$service ./cmd/$$service; \
	done
	@echo "✅ Build completed!"

# Test all packages
test: ## Run all tests with coverage
	@echo "🧪 Running tests..."
	@$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests completed! Coverage report: coverage.html"

# Quick test without coverage
test-quick: ## Run tests without coverage report
	@echo "🧪 Running quick tests..."
	@$(GOTEST) -short ./...
	@echo "✅ Quick tests completed!"

# Lint code
lint: ## Run all linters
	@echo "🔍 Running linters..."
	@$(GOFMT) ./...
	@$(GOVET) ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
	@echo "✅ Linting completed!"

# Format code
fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	@$(GOFMT) ./...
	@echo "✅ Code formatted!"

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BINARY_DIR)/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean completed!"

# Run individual services
run-auth: build ## Run Auth Service
	@echo "🚀 Starting Auth Service on :8080..."
	@./$(BINARY_DIR)/auth

run-product: build ## Run Product Service
	@echo "🚀 Starting Product Service on :8081..."
	@./$(BINARY_DIR)/product

run-cart: build ## Run Cart Service
	@echo "🚀 Starting Cart Service on :8082..."
	@./$(BINARY_DIR)/cart

run-order: build ## Run Order Service
	@echo "🚀 Starting Order Service on :8083..."
	@./$(BINARY_DIR)/order

run-payment: build ## Run Payment Service
	@echo "🚀 Starting Payment Service on :8084..."
	@./$(BINARY_DIR)/payment

run-gateway: build ## Run API Gateway
	@echo "🚀 Starting API Gateway on :8000..."
	@./$(BINARY_DIR)/gateway

# Docker commands
docker-build: ## Build Docker images for all services
	@echo "🐳 Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building $(DOCKER_REGISTRY)/$$service..."; \
		docker build -t $(DOCKER_REGISTRY)/$$service -f deployments/docker/$$service.Dockerfile .; \
	done
	@echo "✅ Docker images built!"

docker-up: ## Start all services with Docker Compose
	@echo "🐳 Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "✅ Services started! Gateway available at http://localhost:8000"

docker-down: ## Stop all services
	@echo "🐳 Stopping services..."
	@docker-compose down
	@echo "✅ Services stopped!"

# Development dependencies
deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "✅ Dependencies updated!"

# Install development tools
tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed!"

# Generate code (protobuf, etc.)
generate: ## Generate code from protobuf and other sources
	@echo "⚙️  Generating code..."
	@$(GOCMD) generate ./...
	@echo "✅ Code generation completed!"

# Full check (test + lint)
check: test lint ## Run all checks (test, lint)
	@echo "✅ All checks completed!"

# Show project status
status: ## Show project status
	@echo "📊 Project Status:"
	@echo "=================="
	@echo "Go version: $$(go version)"
	@echo "Project: $$(head -1 go.mod | cut -d' ' -f2)"
	@echo "Services: $(SERVICES)" 