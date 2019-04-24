.PHONY: migrate up build

migrate:
	./scripts/migrate.sh

up:
	docker-compose -f build/docker-compose.yml up

build:
	go build -ldflags="-s -w" -o bin/hello twitter-go/hello/main.go
