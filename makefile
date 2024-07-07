.SILENT:
include .env

.PHONY: run
run:
	go run ./cmd/api

.PHONY: build
build:
	go build -ldflags='-s' -o=./bin/api ./cmd/api

.PHONY: up
up:
	migrate -path=./migrations -database="$(DB_URL)" up

.PHONY: down
down:
	migrate -path=./migrations -database="$(DB_URL)" down

.PHONY: force
force:
	migrate -path=./migrations -database="$(DB_URL)" force 1

.PHONY: audit
audit:
	go mod tidy
	go mod verify
	go vet ./...
	go test -race -vet=off ./...
