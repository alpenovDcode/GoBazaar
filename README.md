# GoBazaar - Microservices-Based E-commerce Platform

[![CI](https://github.com/alpewa/GoBazaar/actions/workflows/ci.yml/badge.svg)](https://github.com/alpewa/GoBazaar/actions)
[![Coverage](https://img.shields.io/codecov/c/github/alpewa/GoBazaar)](https://codecov.io/gh/alpewa/GoBazaar)
[![Go Report Card](https://goreportcard.com/badge/github.com/alpewa/GoBazaar)](https://goreportcard.com/report/github.com/alpewa/GoBazaar)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview
GoBazaar is a modern, scalable e-commerce platform built using a microservices architecture in Go. This project demonstrates best practices for building distributed systems with proper CI/CD, monitoring, and deployment strategies.

## Features
- **Authentication & Authorization** (JWT)
- **Product Catalog** (REST & gRPC)
- **Shopping Cart** stored in Redis
- **Order Processing** with event-driven workflow
- **Payment Integration** via Stripe (sandbox)
- **API Gateway** for request routing and middleware
- **Monitoring & Tracing** (Prometheus, Grafana, OpenTelemetry)
- **Containerization & Orchestration** (Docker, Kubernetes, Helm)

## Microservices
- **Auth Service**: User registration, login, JWT issuance (Port: 8080)
- **Product Service**: CRUD operations for products, caching (Port: 8081)
- **Cart Service**: User cart management in Redis (Port: 8082)
- **Order Service**: Order creation and storage (Port: 8083)
- **Payment Service**: Payment intent and webhook handling (Port: 8084)
- **API Gateway**: Request routing and middleware (Port: 8000)

## Tech Stack
| Component           | Technology / Library           | Version |
|---------------------|--------------------------------|---------|
| Language            | Go                             | 1.23+   |
| HTTP Framework      | Gin                            | 1.10.1  |
| gRPC                | grpc-go                        | 1.68.2  |
| SQL Database        | PostgreSQL                     | 16      |
| Cache / NoSQL       | Redis                          | 7       |
| Message Broker      | NATS                           | 2.10    |
| Authorization       | JWT (golang-jwt/jwt/v5)        | 5.2.1   |
| Containerization    | Docker, Docker Compose         | Latest  |
| Orchestration       | Kubernetes, Helm               | Latest  |
| Monitoring          | Prometheus, Grafana, OpenTelemetry | Latest |
| CI/CD               | GitHub Actions                 | Latest  |

## Requirements
- Go 1.23 or higher
- Docker & Docker Compose
- Make
- Git

## Quick Start

### Automated Setup
```bash
git clone https://github.com/alpewa/GoBazaar.git
cd GoBazaar
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### Manual Setup
```bash
# Clone repository
git clone https://github.com/alpewa/GoBazaar.git
cd GoBazaar

# Install dependencies
make deps

# Install development tools
make tools

# Build all services
make build

# Run tests
make test-quick

# Start all services with Docker
make docker-up
```

All services will be available at `http://localhost` on ports defined in `docker-compose.yml`.

## Development Commands

### Essential Commands
```bash
make help          # Show all available commands
make build         # Build all microservices
make test          # Run tests with coverage
make test-quick    # Run tests without coverage
make lint          # Run linters
make fmt           # Format code
make clean         # Clean build artifacts
```

### Service Management
```bash
make run-auth      # Run Auth Service
make run-product   # Run Product Service
make run-cart      # Run Cart Service
make run-order     # Run Order Service
make run-payment   # Run Payment Service
make run-gateway   # Run API Gateway
```

### Docker Commands
```bash
make docker-build  # Build Docker images
make docker-up     # Start all services
make docker-down   # Stop all services
make status        # Show project status
```

## Project Structure
```
.
├── api/                    # API definitions (OpenAPI specs, proto files)
├── cmd/                    # Entry points for each service
│   ├── auth/              # Auth Service entry point
│   ├── product/           # Product Service entry point
│   ├── cart/              # Cart Service entry point
│   ├── order/             # Order Service entry point
│   ├── payment/           # Payment Service entry point
│   └── gateway/           # API Gateway entry point
├── internal/               # Application logic
│   ├── auth/              # Auth Service implementation
│   ├── product/           # Product Service implementation
│   ├── cart/              # Cart Service implementation
│   ├── order/             # Order Service implementation
│   ├── payment/           # Payment Service implementation
│   ├── gateway/           # API Gateway implementation
│   └── common/            # Shared code and models
├── pkg/                    # Public packages
│   ├── database/          # Database utilities
│   ├── logger/            # Logging utilities
│   ├── messaging/         # Message broker utilities
│   ├── jwt/               # JWT utilities
│   └── cache/             # Cache utilities
├── deployments/           # Deployment configurations
│   ├── docker/            # Dockerfiles for each service
│   ├── kubernetes/        # Kubernetes manifests
│   └── helm/              # Helm charts
├── scripts/               # Development scripts
├── tests/                 # Integration and E2E tests
├── docs/                  # Documentation and diagrams
├── .dockerignore          # Docker ignore file
├── .golangci.yml          # Linter configuration
├── docker-compose.yml     # Local development setup
├── Makefile              # Development automation
└── README.md
```

## Development Progress

### ✅ Day 1 - Project Initialization (Completed)
- [x] Initialize repository and create service directories
- [x] Create basic project structure with all microservices
- [x] Add base `main.go` files for each service
- [x] Configure `go.mod` with dependencies
- [x] Create initial documentation and `.gitignore`

### ✅ Day 2 - Development Tooling & Docker Setup (Completed)
- [x] Enhanced Makefile with 20+ development commands
- [x] Docker configuration for all services with multi-stage builds
- [x] Development tooling setup (golangci-lint, formatting, testing)
- [x] Created `.dockerignore` for optimized Docker builds
- [x] Added comprehensive linting configuration (`.golangci.yml`)
- [x] Implemented basic models and tests
- [x] Created setup script for development environment
- [x] Docker Compose configuration with health checks

### 🚧 Day 3 - CI/CD Pipeline (In Progress)
- [ ] GitHub Actions workflow for linters and tests
- [ ] Automated testing pipeline
- [ ] Code coverage reporting
- [ ] Security scanning integration

### 📋 Upcoming Days
- **Day 4-6**: Auth Service implementation
- **Day 7-10**: Product Service implementation  
- **Day 11-15**: Cart Service implementation
- **Day 16-20**: Order Service implementation
- **Day 21-25**: Payment Service implementation
- **Day 26-27**: API Gateway implementation
- **Day 28**: Local environment setup
- **Day 29**: Kubernetes deployment
- **Day 30**: Monitoring and final polish

## Environment Variables

Create a `.env` file in the root directory:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gobazaar

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# NATS
NATS_URL=nats://localhost:4222

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Stripe
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# Service URLs
AUTH_SERVICE_URL=localhost:8080
PRODUCT_SERVICE_URL=localhost:8081
CART_SERVICE_URL=localhost:8082
ORDER_SERVICE_URL=localhost:8083
PAYMENT_SERVICE_URL=localhost:8084
```

## API Documentation

Once the services are running, API documentation will be available at:

- **API Gateway**: http://localhost:8000/docs
- **Auth Service**: http://localhost:8080/docs
- **Product Service**: http://localhost:8081/docs
- **Cart Service**: http://localhost:8082/docs
- **Order Service**: http://localhost:8083/docs
- **Payment Service**: http://localhost:8084/docs

## Testing

```bash
# Run all tests
make test

# Run tests without coverage
make test-quick

# Run specific service tests
make test-auth
make test-product
make test-cart
make test-order
make test-payment

# Run benchmarks
make benchmark
```

## Contributing

1. Fork this repository
2. Create a branch `feature/<feature-name>`
3. Make your changes and add tests
4. Run `make check` to ensure code quality
5. Commit with conventional commit messages
6. Open a Pull Request for review

### Development Workflow

1. **Setup**: Run `./scripts/setup.sh` for initial setup
2. **Development**: Use `make help` to see available commands
3. **Testing**: Run `make test` before committing
4. **Linting**: Code is automatically formatted and linted
5. **Docker**: Test locally with `make docker-up`

## Monitoring

In production, the following monitoring tools will be available:

- **Prometheus**: Metrics collection
- **Grafana**: Metrics visualization
- **Jaeger**: Distributed tracing
- **ELK Stack**: Centralized logging

## Security

- JWT-based authentication with RS256 algorithm
- Input validation and sanitization
- Rate limiting on API endpoints
- Security headers middleware
- Regular dependency updates

## Performance

- Redis caching for frequently accessed data
- Connection pooling for databases
- Graceful shutdown handling
- Health checks for all services
- Horizontal scaling ready

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with modern Go practices and patterns
- Inspired by industry-standard microservices architectures
- Community-driven development approach
