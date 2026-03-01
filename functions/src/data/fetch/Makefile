.PHONY: tests build benchmarks format lint clean

tests:
	go test -v -race ./...

build:
	go build ./...

benchmarks:
	go test -bench=. -benchmem ./...

format:
	gofmt -s -w .
	@if command -v golines >/dev/null 2>&1; then \
		golines -w .; \
	else \
		echo "golines not installed, skipping line wrapping"; \
	fi

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping lint"; \
	fi

clean:
	find . -type f -name "*.test" -delete
	find . -type f -name "coverage.out" -delete
	find . -type f -name "coverage.html" -delete
