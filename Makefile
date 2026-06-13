APP_BINARY := bin/server
ATLAS_MIGRATIONS_DIR := file://ent/migrations
ATLAS_SCHEMA_URL := ent://ent/schema
OBSERVABILITY_COMPOSE_FILE := deployments/observability/docker-compose.yml
APP_STACK_COMPOSE_FILE := deployments/app/docker-compose.yml

.PHONY: run dev build test ent ent-gen openapi lint fmt check-template migrate-diff migrate-status migrate-apply observability-up observability-down observability-logs app-stack-up app-stack-down app-stack-logs seed-admin

run:
	go run ./cmd/server

dev:
	air

build:
	mkdir -p bin
	go build -o $(APP_BINARY) ./cmd/server

test:
	go test ./...

ent: ent-gen

ent-gen:
	go generate ./ent

openapi:
	go run github.com/swaggo/swag/cmd/swag init --parseInternal -g main.go -d cmd/server,internal/modules/auth,internal/modules/health,internal/modules/users,internal/response -o internal/platform/docs

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

check-template:
	./scripts/check-template.sh

migrate-diff:
	@test -n "$(name)" || (echo "usage: make migrate-diff name=create_users" && exit 1)
	@test -n "$(ATLAS_DIFF_DATABASE_URL)" || (echo "ATLAS_DIFF_DATABASE_URL is required for migrate-diff" && exit 1)
	atlas migrate diff $(name) --dir "$(ATLAS_MIGRATIONS_DIR)" --to "$(ATLAS_SCHEMA_URL)" --dev-url "$(ATLAS_DIFF_DATABASE_URL)"

migrate-status:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required for migrate-status" && exit 1)
	atlas migrate status --dir "$(ATLAS_MIGRATIONS_DIR)" --url "$(DATABASE_URL)"

migrate-apply:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is required for migrate-apply" && exit 1)
	atlas migrate apply --dir "$(ATLAS_MIGRATIONS_DIR)" --url "$(DATABASE_URL)"

observability-up:
	docker compose -f $(OBSERVABILITY_COMPOSE_FILE) up -d

observability-down:
	docker compose -f $(OBSERVABILITY_COMPOSE_FILE) down

observability-logs:
	docker compose -f $(OBSERVABILITY_COMPOSE_FILE) logs -f

app-stack-up:
	docker compose -f $(APP_STACK_COMPOSE_FILE) up -d

app-stack-down:
	docker compose -f $(APP_STACK_COMPOSE_FILE) down

app-stack-logs:
	docker compose -f $(APP_STACK_COMPOSE_FILE) logs -f

seed-admin:
	go run ./cmd/seed
