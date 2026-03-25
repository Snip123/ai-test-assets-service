SERVICE_NAME := ai-test-assets-service

.PHONY: run test test-bdd generate check-generate migrate-up migrate-down lint build docker-build

run:
	go run ./cmd/server

test:
	go test ./... -short -count=1 -race

test-bdd:
	go test ./internal/... -v -count=1 -run TestBDD

generate:
	oapi-codegen -config .oapi-codegen.yaml docs/api/openapi.yaml
	sqlc generate

check-generate:
	$(MAKE) generate
	git diff --exit-code internal/generated/

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down 1

lint:
	golangci-lint run ./...

build:
	go build -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate

docker-build:
	docker build -t $(SERVICE_NAME):local .
