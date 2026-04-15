# Build and Installation Settings
PLUGIN_NAME := docker-deps
DEST_DIR := $(HOME)/.docker/cli-plugins
BINARY_NAME := $(PLUGIN_NAME)

# Go configuration
GO := go
GOLINT := golangci-lint

.PHONY: all build install clean test lint help coverage plugin-build plugin-push

all: build test lint ## Run build, test, and lint

help: ## Show help for this Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the plugin binary for the current architecture
	$(GO) build -o $(BINARY_NAME)

install: build ## Build and install the plugin to the Docker CLI plugins directory
	mkdir -p $(DEST_DIR)
	cp $(BINARY_NAME) $(DEST_DIR)/

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	docker plugin rm $(PLUGIN_NAME) 2>/dev/null || true

test: ## Run unit tests
	$(GO) test -v ./...

coverage: ## Run tests and show coverage report
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

lint: ## Run linter checks
	$(GOLINT) run

dist: ## Build binaries for multiple platforms
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-windows-amd64.exe

plugin-build: ## Build the Docker managed plugin
	$(GO) build -o plugin/rootfs/docker-deps
	docker plugin create $(PLUGIN_NAME) plugin

plugin-run: plugin-build ## Build and run the plugin locally
	docker plugin enable $(PLUGIN_NAME)

plugin-push: plugin-build ## Build and push the plugin to Docker Hub
	docker plugin push $(PLUGIN_NAME)
