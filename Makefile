build:
	## Build Init Function
	mkdir -p functions/build/data
	docker run --rm -v `pwd`:/build -w /build/functions/src/data/init tinygo/tinygo:0.38.0 tinygo build -o /build/functions/build/data/init.wasm -scheduler=none --no-debug -target=wasip1 -buildmode=c-shared main.go
	## Build CSV Fetch Function
	mkdir -p functions/build/data
	docker run --rm -v `pwd`:/build -w /build/functions/src/data/fetch tinygo/tinygo:0.38.0 tinygo build -o /build/functions/build/data/fetch.wasm -scheduler=none --no-debug -target=wasip1 -buildmode=c-shared main.go
	## Build CSV Load Function
	mkdir -p functions/build/data
	docker run --rm -v `pwd`:/build -w /build/functions/src/data/load tinygo/tinygo:0.38.0 tinygo build -o /build/functions/build/data/load.wasm -scheduler=none --no-debug -target=wasip1 -buildmode=c-shared main.go
	## Build US Seed Function
	mkdir -p functions/build/data
	docker run --rm -v `pwd`:/build -w /build/functions/src/data/seed tinygo/tinygo:0.38.0 tinygo build -o /build/functions/build/data/seed.wasm -scheduler=none --no-debug -target=wasip1 -buildmode=c-shared main.go
	## Build HTTP Request Handler Function
	mkdir -p functions/build/handlers
	docker run --rm -v `pwd`:/build -w /build/functions/src/handlers/lookup tinygo/tinygo:0.38.0 tinygo build -o /build/functions/build/handlers/lookup.wasm -scheduler=none --no-debug -target=wasip1 -buildmode=c-shared main.go

.PHONY: tests tidy
tests:
	## Run tests
	mkdir -p coverage
	go test -v -race -covermode=atomic -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	## Run tests for the lookup function
	$(MAKE) -C functions/src/handlers/lookup tests

tidy:
	## Run go mod tidy for all function modules
	@set -e; \
	for mod in $$(find functions/src -name go.mod | sort); do \
		dir=$$(dirname $$mod); \
		echo "==> $$dir"; \
		(cd $$dir && go mod tidy); \
	done

docker-compose:
	docker compose up -d mysql redis
	sleep 15
	docker compose up data-manager lookup

docker-compose-background:
	docker compose up -d mysql redis
	sleep 15
	docker compose up -d data-manager lookup

loadtest-setup:
	docker compose -f load-compose.yml up -d mysql
	sleep 15
	docker compose -f load-compose.yml up -d data-manager lookup

run: build docker-compose
run-nobuild: docker-compose
run-background: build docker-compose-background
run-stress: build loadtest-setup
	sleep 600
	k6 run --config tests/k6/stress.json tests/k6/script.js
run-soak: build loadtest-setup
	sleep 600
	k6 run --config tests/k6/soak.json tests/k6/script.js
run-steady: build docker-compose-background
	sleep 600
	k6 run --config tests/k6/steady.json tests/k6/script.js


clean:
	rm -rf functions/build
	docker compose down --remove-orphans
