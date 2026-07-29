.PHONY: help build test lint coverage clean install run

# Цвет для вывода help
BLUE := \033[36m
RESET := \033[0m

BINARY := depscheck
MAIN := ./cmd/depscheck

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE)%-15s$(RESET) %s\n", $$1, $$2}'

build: ## Build binary to ./bin/depscheck
	go build -o bin/$(BINARY) $(MAIN)

test: ## Run tests with race detector
	go test -race ./...

coverage: ## Generate HTML coverage report
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in browser"

lint: ## Run golangci-lint
	golangci-lint run ./...

tidy: ## Tidy go.mod
	go mod tidy

clean: ## Remove build artifacts
	rm -rf bin/ coverage.out coverage.html dist/

release: ## Build a local snapshot release (no publishing)
	goreleaser release --snapshot --clean

.DEFAULT_GOAL := help