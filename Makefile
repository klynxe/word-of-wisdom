up-server:
	docker compose up --build -d wow-server
.PHONY: up-server

up-client:
	docker compose up --build wow-client
.PHONY: up-client

lint:
	golangci-lint run ./...
.PHONY: lint

test:
	go test -v ./...
.PHONY: test