.PHONY: all lint format

all:
	@go run .

lint:
	@golangci-lint run

format:
	@golangci-lint fmt
	@pg_format -i queries/*.sql
