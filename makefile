.SILENT:
POSTGRES_USER := artix
POSTGRES_PASSWORD := 
POSTGRES_HOST := localhost
POSTGRES_PORT := 5432
POSTGRES_DATABASE := basefit

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

run:
	go run ./cmd/api

up:
	@migrate -path=./migrations -database="$(DB_URL)" up

down:
	@migrate -path=./migrations -database="$(DB_URL)" down

force:
	@migrate -path=./migrations -database="$(DB_URL)" force 1
