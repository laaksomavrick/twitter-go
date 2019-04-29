.PHONY: migrate up build format

migrate:
	./scripts/migrate.sh

up:
	docker-compose -f build/docker-compose.yml up

build:
	go build -ldflags="-s -w" -o bin/gateway services/gateway/cmd/main.go

format:
	./scripts/gofmt.sh