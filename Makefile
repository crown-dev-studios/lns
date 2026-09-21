BINARY_NAME=lns
VERSION?=dev
COMMIT?=$(shell git rev-parse --short=12 HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell git show -s --format=%cI HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)
GORELEASER?=goreleaser

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Directories
BUILD_DIR=build
CMD_DIR=cmd/lns

.PHONY: all build build-fast clean test coverage deps lint fmt vet tidy snapshot release-check help

all: build

## Build

build: ## Build for current platform
	$(GOBUILD) -trimpath -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

build-fast: ## Build without optimizations (faster compile)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

## Development

run: ## Run the application
	$(GOCMD) run ./$(CMD_DIR) $(ARGS)

dev: ## Run with live reload (requires air)
	air

test: ## Run tests
	$(GOTEST) -v ./...

coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

fmt: ## Format code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

## Dependencies

deps: ## Download dependencies
	$(GOMOD) download

tidy: ## Tidy dependencies
	$(GOMOD) tidy

## Release

release-check: ## Validate the GoReleaser configuration
	$(GORELEASER) check
	bash scripts/test-homebrew-formula.sh

snapshot: ## Build local release artifacts without publishing
	$(GORELEASER) release --snapshot --clean

## Cleanup

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## Help

help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
