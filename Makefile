.PHONY: migrate up build format run test setup-k8s helm-install helm-debug helm-purge helm-upgrade docker-build

migrate:
	@scripts/migrate.sh

up:
	@docker-compose -f docker-compose.yml up

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
	@helm install ./helm --name=twtr-dev --tiller-namespace=twtr-dev --namespace=twtr-dev

helm-upgrade:
	@helm upgrade  twtr-dev ./helm --tiller-namespace=twtr-dev --namespace=twtr-dev

helm-debug:
	@helm install --dry-run --debug ./helm --name=twtr-dev --tiller-namespace=twtr-dev --namespace=twtr-dev

helm-purge:
	@helm ls --all --short --tiller-namespace=twtr-dev | xargs -L1 helm delete --purge --tiller-namespace=twtr-dev

docker-build:
	@scripts/build-docker-images.sh