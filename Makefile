.PHONY: migrate up build format run

migrate:
	./scripts/migrate.sh

up:
	docker-compose -f build/docker-compose.yml up

run:
	./scripts/run-all.sh

build:
	go build -ldflags="-s -w" -o bin/gateway services/gateway/cmd/main.go

format:
	./scripts/gofmt.sh

integration-test:
	./scripts/run-integration-tests.sh