.PHONY: run test lint

all: run

run:
	@go run main.go

test:
	@go test ./...

lint:
	@golangci-lint run ./...

format:
	@golangci-lint fmt ./...
