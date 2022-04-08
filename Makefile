# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: tidy test

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: test
test:
	go test -race -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out
	@rm coverage.out