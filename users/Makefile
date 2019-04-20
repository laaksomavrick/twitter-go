.PHONY: build clean deploy

build:
	env GOOS=linux go build -mod=vendor -ldflags="-s -w" -o bin/create src/create.go

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose
