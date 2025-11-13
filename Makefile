.PHONY: build run test clean install tui example lint lint-fix fmt

BINARY_NAME=hrun
MAIN_PATH=cmd/hrun/main.go

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	./$(BINARY_NAME) $(ARGS)

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

install:
	go install $(MAIN_PATH)

tui: build
	./$(BINARY_NAME) tui examples/sample.http

example: build
	./$(BINARY_NAME) run examples/sample.http --request 1

lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run --timeout=5m; \
	fi

lint-fix:
	@echo "Running golangci-lint with auto-fix..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run --fix --timeout=5m; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run --fix --timeout=5m; \
	fi

fmt:
	go fmt ./...
	go mod tidy
