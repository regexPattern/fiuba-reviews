.PHONY: run test lint

all: run

run:
	@go run main.go

test:
	@go test -v ./...

lint:
	@golangci-lint run .

format:
	@golangci-lint fmt .
