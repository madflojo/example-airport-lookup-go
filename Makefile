## Makefile for Tarmac Example Project

build:
	## Create build directory
	mkdir -p functions/build
	## Run TinyGo build via Docker because its easier
	docker run --rm -v `pwd`:/build -w /build/functions/build/init tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/init.wasm -target wasi /build/functions/src/init/main.go
	docker run --rm -v `pwd`:/build -w /build/functions/build/data/fetch tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/fetch.wasm -target wasi /build/functions/src/data/fetch/main.go

tests:
	## Run tests
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

docker-compose:
	docker compose up

run: build docker-compose
