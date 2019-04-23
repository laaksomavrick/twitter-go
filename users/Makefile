.PHONY: build clean deploy

build:
	go mod vendor
	env GOOS=linux go build -mod=vendor -ldflags="-s -w" -o bin/create_handler create_handler/create_handler.go
	env GOOS=linux go build -mod=vendor -ldflags="-s -w" -o bin/get_handler get_handler/get_handler.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose
