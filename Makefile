## Makefile for Tarmac Example Project

build:
	## Create build directory
	mkdir -p functions/build
	## Build Init Function
	docker run --rm -v `pwd`:/build -w /build/functions/build/init tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/init.wasm -target wasi /build/functions/src/init/main.go
	## Build CSV Fetch Function
	docker run --rm -v `pwd`:/build -w /build/functions/build/data/fetch tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/fetch.wasm -target wasi /build/functions/src/data/fetch/main.go
	## Build CSV Load Function
	docker run --rm -v `pwd`:/build -w /build/functions/build/data/load tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/load.wasm -target wasi /build/functions/src/data/load/main.go
	## Build HTTP Request Handler Function
	docker run --rm -v `pwd`:/build -w /build/functions/build/handler tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/handler.wasm -target wasi /build/functions/src/handler/main.go

.PHONY: tests
tests:
	## Run tests
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

docker-compose:
	docker compose up -d mysql
	sleep 15
	docker compose up example

docker-compose-background:
	docker compose up -d mysql
	sleep 15
	docker compose up -d example

run: build docker-compose
run-nobuild: docker-compose
run-background: build docker-compose-background

clean:
	rm -rf functions/build
	docker compose down --remove-orphans
