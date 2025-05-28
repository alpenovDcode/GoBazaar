# GoBazaar - Microservices-Based E-commerce Platform

[![CI](https://github.com/yourusername/GoBazaar/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/GoBazaar/actions)
[![Coverage](https://img.shields.io/codecov/c/github/yourusername/GoBazaar)](https://codecov.io/gh/yourusername/GoBazaar)
[![Docker Pulls](https://img.shields.io/docker/pulls/yourusername/gobazaar)](https://hub.docker.com/r/yourusername/gobazaar)

## Overview
GoBazaar is a modern, scalable e-commerce platform built using a microservices architecture in Go.

## Features
- **Authentication & Authorization** (JWT)
- **Product Catalog** (REST & gRPC)
- **Shopping Cart** stored in Redis
- **Order Processing** with event-driven workflow
- **Payment Integration** via Stripe (sandbox)
- **Monitoring & Tracing** (Prometheus, Grafana, OpenTelemetry)
- **Containerization & Orchestration** (Docker, Kubernetes, Helm)

## Microservices
- **Auth Service**: User registration, login, JWT issuance
- **Product Service**: CRUD operations for products, caching
- **Cart Service**: User cart management in Redis
- **Order Service**: Order creation and storage
- **Payment Service**: Payment intent and webhook handling

## Tech Stack
| Component           | Technology / Library           |
|---------------------|--------------------------------|
| Language            | Go 1.21+                       |
| HTTP / gRPC         | Gin, grpc-go                   |
| SQL Database        | PostgreSQL                     |
| Cache / NoSQL       | Redis                          |
| Message Broker      | NATS (or RabbitMQ)             |
| Authorization       | JWT (RS256)                    |
| Containerization    | Docker, Docker Compose         |
| Orchestration       | Kubernetes, Helm               |
| Monitoring          | Prometheus, Grafana, OpenTelemetry |
| CI/CD               | GitHub Actions                 |

## Requirements
- Go 1.21 or higher
- Docker & Docker Compose
- Kubernetes (for production)
- Make
- Git

## Quick Start
### Clone Repository
```bash
git clone https://github.com/yourusername/GoBazaar.git
cd GoBazaar
```

### Install Dependencies
```bash
go mod download
```

### Run Locally with Docker Compose
```bash
docker-compose up -d
```
All services will be available at `http://localhost` on ports defined in `docker-compose.yml`.

```

## Project Structure
```
.
├── api/              # API definitions (OpenAPI specs, proto files)
├── cmd/              # Entry points for each service
├── internal/         # Application logic
│   ├── auth/
│   ├── product/
│   ├── cart/
│   ├── order/
│   ├── payment/
│   └── common/
├── deployments/      # Docker Compose, Helm charts
├── docs/             # Documentation and diagrams
├── Makefile
└── README.md
```

## Day 1 Milestones
- [ ] Initialize repository and create service directories
- [ ] Add base `Makefile` and `Dockerfile` for each service
- [ ] Configure CI pipeline with linters and tests
- [ ] Publish initial commit with project skeleton

## Contributing
1. Fork this repository
2. Create a branch `feature/<feature-name>`
3. Implement changes and add tests
4. Open a Pull Request for review

## Code of Conduct
This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). Please review the code of conduct.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
