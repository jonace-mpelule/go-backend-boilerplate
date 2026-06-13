# Go Backend Boilerplate

A production-oriented Go backend template for modular monoliths built with Chi, Ent, Postgres, Redis, zap, Sentry, and PostHog.

## What This Template Gives You

- Modular monolith structure under `internal/modules/<module>`
- Production-aware HTTP server timeouts and graceful shutdown
- Centralized config loading and validation
- Structured JSON responses via `internal/response`
- Typed application errors via `internal/errors`
- JWT auth middleware with permission guards
- Optional Redis cache through a single cache interface
- DB access through Ent with explicit startup health checks
- Liveness/readiness health endpoints
- Baseline CI for tests, linting, and template consistency
- Prometheus metrics endpoint for scrape-based monitoring
- Optional direct Loki log shipping for external Grafana/Loki stacks
- Resend as the default transactional email provider

## Project Layout

```text
.
├── cmd/server                # Application entrypoint
├── ent/                      # Ent schema and generated client code
├── internal/
│   ├── app/                  # Composition root and HTTP server bootstrapping
│   ├── config/               # Environment-backed application config
│   ├── errors/               # Typed reusable application errors
│   ├── middlewares/          # HTTP middleware stack
│   ├── modules/              # Independent feature modules
│   ├── permissions/          # Permission constants
│   ├── platform/             # External service adapters
│   ├── response/             # Standard API response helpers
│   └── types/                # Shared container types
├── scripts/                  # Helper scripts
└── .github/workflows/ci.yml  # Baseline CI
```

## Getting Started

### Prerequisites

- Go 1.26.2+
- PostgreSQL
- Redis (optional)
- `atlas` CLI for versioned migrations
- `golangci-lint` for local linting

### Setup

1. Clone the template.
2. Copy the example environment file:

```bash
cp .env.example .env
```

3. Update the placeholder module path if this is a new project:

```bash
./scripts/module_name.sh github.com/your-org/your-service
```

4. Provision local infrastructure. Example with Docker:

```bash
docker run --name template-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=template -p 5432:5432 -d postgres:17
docker run --name template-redis -p 6379:6379 -d redis:7
```

5. Generate Ent code if you change schema definitions:

```bash
make ent
```

### Run

```bash
make run
```

### Build

```bash
make build
./bin/server
```

## Migrations

This template defaults to versioned migrations instead of runtime schema creation.

- Generate Ent client code:

```bash
make ent
```

- Create a migration:

```bash
make migrate-diff name=create_users
```

- Check migration status:

```bash
make migrate-status
```

- Apply migrations:

```bash
make migrate-apply
```

`make migrate-diff` expects `ATLAS_DIFF_DATABASE_URL` to point to a separate disposable diff database.
`make migrate-status` and `make migrate-apply` expect `DATABASE_URL` to point to the target application database.

## Health Endpoints

- `GET /health/live` for liveness
- `GET /health/ready` for readiness
- `GET /health` as a readiness alias

## Observability

This template supports a portable observability layout:

- the app exposes Prometheus metrics on `/metrics`
- Prometheus scrapes the app over HTTP
- the app can optionally push structured logs directly to Loki
- Grafana is expected to live elsewhere and connect to Prometheus and Loki as external data sources

### Metrics

- Metrics are enabled by default.
- The route is configurable through `METRICS_PATH`.
- The metric namespace is configurable through `METRICS_NAMESPACE`.
- Prometheus should scrape the API directly; the app does not push metrics.

### Loki

- Loki push is disabled by default.
- Enable it with `LOKI_ENABLED=true`.
- Set `LOKI_URL` to your Loki base URL.
- Optional auth and tenancy are supported through:
  - `LOKI_TENANT_ID`
  - `LOKI_BASIC_AUTH_USERNAME`
  - `LOKI_BASIC_AUTH_PASSWORD`
- Additional static labels can be supplied via `LOKI_LABELS` using `key=value,key=value`.

### Local Prometheus/Loki Validation

This repo includes a lightweight local validation stack for Prometheus and Loki only.

```bash
make observability-up
make observability-logs
make observability-down
```

By default the included Prometheus config scrapes `host.docker.internal:8080`, so the Go app should be running on the host on port `8080` unless you update the scrape target.

## Commands

| Command | Description |
| :--- | :--- |
| `make run` | Run the server |
| `make dev` | Run with Air |
| `make build` | Build the binary into `bin/server` |
| `make test` | Run the test suite |
| `make ent` | Regenerate Ent client code |
| `make migrate-diff name=...` | Create a new Atlas migration using `ATLAS_DIFF_DATABASE_URL` |
| `make migrate-status` | Show Atlas migration status for `DATABASE_URL` |
| `make migrate-apply` | Apply Atlas migrations to `DATABASE_URL` |
| `make observability-up` | Start local Prometheus and Loki validation services |
| `make observability-down` | Stop local Prometheus and Loki validation services |
| `make observability-logs` | Follow local Prometheus and Loki logs |
| `make lint` | Run golangci-lint |
| `make fmt` | Run gofmt on the repo |
| `make check-template` | Check template consistency |

## Notes

- Redis is optional. If `REDIS_URL` is empty, the app uses a noop cache.
- Set a real `JWT_SECRET` before production. The default development secret is rejected in `APP_ENV=production`.
- The example `users` module is intentionally small, but it now follows the handler/service/repository/DTO split expected for new modules.
