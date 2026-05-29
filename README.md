# Go Backend Boilerplate

[![Go Version](https://img.shields.io/github/go-mod/go-version/jonace-mpelule/go-backend-boilerplate)](https://golang.org)
[![Build Status](https://github.com/jonace-mpelule/go-backend-boilerplate/actions/workflows/ci.yml/badge.svg)](https://github.com/jonace-mpelule/go-backend-boilerplate/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonace-mpelule/go-backend-boilerplate)](https://goreportcard.com/report/github.com/jonace-mpelule/go-backend-boilerplate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A production-grade, modular monolith boilerplate for building scalable and maintainable backends in Go. This template focuses on developer productivity, observability, and clean architecture.

## 🚀 Features

- **Framework**: [Chi](https://github.com/go-chi/chi) for high-performance routing.
- **ORM**: [Ent](https://entgo.io/) for type-safe database modeling and queries.
- **Database**: PostgreSQL support out of the box.
- **Caching**: Redis integration with a simple cache-aside pattern.
- **Auth**: JWT-based authentication and flexible permission-based authorization.
- **Observability**: Built-in support for [Sentry](https://sentry.io/) (errors) and [PostHog](https://posthog.com/) (analytics).
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator) for request payload validation.
- **Development**: Hot reloading with [Air](https://github.com/cosmtrek/air).
- **CI/CD**: Pre-configured GitHub Actions for linting and testing.

## � Core Dependencies

This project leverages several industry-standard tools to ensure performance and reliability:

- **[Redis](https://redis.io/)**: Used as a high-performance caching layer. It implements a cache-aside pattern to reduce database load and improve response times for frequently accessed data.
- **[Sentry](https://sentry.io/)**: Integrated for real-time error tracking and monitoring. It captures unhandled exceptions and system errors, providing detailed stack traces and context for faster debugging.
- **[PostHog](https://posthog.com/)**: Used for product analytics. It allows tracking of key business events and user interactions without compromising sensitive PII, helping you understand how the application is used.
- **[Ent](https://entgo.io/)**: An entity framework for Go that provides a type-safe API for modeling and querying data, making database interactions more robust and maintainable.

## � Project Structure

```text
.
├── cmd/server          # Application entry point
├── ent/                # Ent schema and generated code
├── internal/
│   ├── app/            # Application composition root
│   ├── modules/        # Domain-specific logic (Modular Monolith)
│   ├── platform/       # Infrastructure & external services (DB, Cache, etc.)
│   └── utils/          # Shared helpers and error types
├── scripts/            # Development helper scripts
└── Makefile            # Common development commands
```

## 🛠️ Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) 1.26+
- [Docker](https://www.docker.com/get-started) & Docker Compose
- [Air](https://github.com/cosmtrek/air) (optional, for hot reloading)

### Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/jonace-mpelule/go-backend-boilerplate.git
   cd go-backend-boilerplate
   ```

2. **Rename the module** (Optional):
   To rename the project to your own module path (e.g., `github.com/youruser/yourproject`):
   ```bash
   chmod +x scripts/module_name.sh
   ./scripts/module_name.sh github.com/youruser/yourproject
   ```
   This script updates the `go.mod` file and all internal import paths throughout the project.

3. **Configure environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

### Running the Application

- **Development (with hot reload)**:
  ```bash
  make dev
  ```

- **Production Build**:
  ```bash
  make build
  ./bin/api
  ```

## 🧪 Commands

| Command | Description |
| :--- | :--- |
| `make dev` | Start development server with Air |
| `make test` | Run all tests |
| `make ent` | Generate Ent code from schema |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code with gofmt |

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
