run:
	go run ./cmd/api

dev:
	air

build:
	go build -o bin/api ./cmd/server

test:
	go test ./...

ent:
	go generate ./ent

lint:
	golangci-lint run

fmt:
	go fmt ./...
