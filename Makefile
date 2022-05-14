TOOL_NAME   = waitforit
BIN         = ${TOOL_NAME}

.PHONY: build run test

all: build

build:
	@go build

## Start application
run:
	@sh -c "go build && ./$(BIN)"

## Perform all tests
test:
	@go test -v -cover ./...
