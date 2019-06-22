.PHONY: migrate up build format run test setup-k8s helm-install helm-debug

migrate:
	@scripts/migrate.sh

up:
	@docker-compose -f build/docker-compose.yml up

run:
	@scripts/run-all.sh

build:
	@go build -ldflags="-s -w" -o bin/gateway services/feeds/cmd/main.go
	@go build -ldflags="-s -w" -o bin/gateway services/followers/cmd/main.go
	@go build -ldflags="-s -w" -o bin/gateway services/gateway/cmd/main.go
	@go build -ldflags="-s -w" -o bin/gateway services/tweets/cmd/main.go
	@go build -ldflags="-s -w" -o bin/gateway services/users/cmd/main.go

format:
	@scripts/gofmt.sh

test:
	@scripts/run-integration-tests.sh

setup-k8s:
	@scripts/setup-k8s.sh

helm-install:
	@helm install ./helm --tiller-namespace=twtr-dev --namespace=twtr-dev

helm-debug:
	@helm install --dry-run --debug ./helm --tiller-namespace=twtr-dev --namespace=twtr-dev