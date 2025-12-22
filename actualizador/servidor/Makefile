.PHONY: run build lint format

run:
	@go run .

build:
	@go build .

lint:
	@golangci-lint run

format:
	@golangci-lint fmt
	@pg_format -i queries/**/*.sql
