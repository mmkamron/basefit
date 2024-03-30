.SILENT:
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

run:
	go run cmd/main.go

up:
	@migrate -database "$(DB_URL)" -path migrations up

down:
	@migrate -database "$(DB_URL)" -path migrations down
