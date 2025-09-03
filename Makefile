.PHONY: build run test clean install

BINARY_NAME=cassie
MAIN_PATH=cmd/cassie/main.go

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