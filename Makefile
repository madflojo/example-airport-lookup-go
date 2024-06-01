build:
	## Build Init Function
	mkdir -p functions/build/data
	docker run --rm -v `pwd`:/build -w /build/functions/build/data/init tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/data/init.wasm -target wasi /build/functions/src/data/init/main.go
	## Build CSV Fetch Function
	mkdir -p functions/build/data
	docker run --rm -v `pwd`:/build -w /build/functions/build/data/fetch tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/data/fetch.wasm -target wasi /build/functions/src/data/fetch/main.go
	## Build CSV Load Function
	mkdir -p functions/build/data
	docker run --rm -v `pwd`:/build -w /build/functions/build/data/load tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/data/load.wasm -target wasi /build/functions/src/data/load/main.go
	## Build HTTP Request Handler Function
	mkdir -p functions/build/handlers
	docker run --rm -v `pwd`:/build -w /build/functions/build/handlers/lookup/ tinygo/tinygo:0.25.0 tinygo build -o /build/functions/build/handlers/lookup.wasm -target wasi /build/functions/src/handlers/lookup/main.go

.PHONY: tests
tests:
	## Run tests
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	## Run tests for the lookup function
	$(MAKE) -C functions/src/handlers/lookup tests

docker-compose:
	docker compose up -d mysql
	sleep 15
	docker compose up data-manager lookup

docker-compose-background:

run: build docker-compose
run-nobuild: docker-compose
run-background: build docker-compose-background
run-stress: 
	docker compose up -d mysql
	sleep 15
	docker compose up data-manager lookup k6-stress --exit-code-from k6-stress
run-soak:
	docker compose up -d mysql
	sleep 15
	docker compose up data-manager lookup k6-soak --exit-code-from k6-soak

clean:
	rm -rf functions/build
	docker compose down --remove-orphans
