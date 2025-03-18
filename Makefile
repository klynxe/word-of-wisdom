LINT_VERSION = v1.64.8

up-server:
	docker compose up --build -d wow-server
.PHONY: up-server

up-client:
	docker compose up --build wow-client
.PHONY: up-client

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(LINT_VERSION) run --allow-parallel-runners --timeout 10m
.PHONY: lint

test:
	go test -v ./...
.PHONY: test