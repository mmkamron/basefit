.SILENT:
include .env

run:
	go run ./cmd/api

up:
	migrate -path=./migrations -database="$(DB_URL)" up

down:
	migrate -path=./migrations -database="$(DB_URL)" down

force:
	migrate -path=./migrations -database="$(DB_URL)" force 1
