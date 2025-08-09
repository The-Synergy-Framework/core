# Core Library Makefile
# Provides common development tasks for the core library

# Variables
COVERAGE_DIR=./coverage
LINT_CONFIG=.golangci.yml

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Core Library Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development tasks
.PHONY: build
build: ## Build/compile the library to check for compilation errors
	@echo "Building core library..."
	go build ./...
	@echo "Build successful - no compilation errors"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(COVERAGE_DIR)
	@go clean -cache -testcache
	@echo "Clean complete"

# Testing tasks
.PHONY: test
test: ## Run all tests
	@echo "Running all tests..."
	go test -v ./...

.PHONY: test-all
test-all: test-race test-coverage test-benchmark ## Run all tests including race detection, coverage, and benchmarks

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -race -v ./...

.PHONY: test-short
test-short: ## Run only short tests
	@echo "Running short tests..."
	go test -short -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Code quality tasks
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (if installed)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running linter with auto-fix..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

.PHONY: check
check: fmt vet lint ## Run all code quality checks

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

.PHONY: deps-check
deps-check: ## Check for outdated dependencies
	@echo "Checking for outdated dependencies..."
	@echo "Note: go-outdated tool is no longer available"
	@echo "Use 'go list -u -m all' to check for updates"

# Security tasks
.PHONY: security
security: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install with: go install github.com/cosmos/gosec/v2/cmd/gosec@latest"; \
	fi

.PHONY: vulncheck
vulncheck: ## Check for vulnerabilities
	@echo "Checking for vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "govulncheck not found. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi

# Documentation tasks
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "Starting godoc server on http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

.PHONY: docs-check
docs-check: ## Check documentation coverage
	@echo "Checking documentation coverage..."
	@go list -f '{{.Dir}}' ./... | xargs -I {} sh -c 'echo "Checking {}"; godoc -analysis=type,pointer -build=false {}'

# Development workflow
.PHONY: dev-setup
dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f $(LINT_CONFIG) ]; then \
		echo "Creating golangci-lint config..."; \
		golangci-lint config > $(LINT_CONFIG) 2>/dev/null || echo "golangci-lint config created"; \
	fi
	@echo "Development environment ready"

.PHONY: pre-commit
pre-commit: fmt vet lint test ## Run pre-commit checks
	@echo "Pre-commit checks completed"

.PHONY: ci
ci: deps check test-all security ## Run CI pipeline
	@echo "CI pipeline completed"

# Utility tasks
.PHONY: version
version: ## Show version information
	@echo "Version: $(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')"
	@echo "Go version: $(shell go version)"
	@echo "Build time: $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')"

.PHONY: info
info: ## Show project information
	@echo "Project: Core Library"
	@echo "Module: $(shell go list -m)"
	@echo "Go version: $(shell go version | cut -d' ' -f3)"
	@echo "OS/Arch: $(shell go env GOOS)/$(shell go env GOARCH)"
	@echo "GOPATH: $(shell go env GOPATH)"
	@echo "GOROOT: $(shell go env GOROOT)"

.PHONY: modules
modules: ## List all modules
	@echo "Available modules:"
	@go list -m all

# Cleanup tasks
.PHONY: clean-all
clean-all: clean ## Clean everything including go mod cache
	@echo "Cleaning everything..."
	go clean -modcache
	@echo "All clean"

# Install development tools
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/godoc@latest
	go install github.com/cosmos/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "Development tools installed"

# Watch mode for development
.PHONY: watch
watch: ## Watch for changes and run tests (requires fswatch)
	@echo "Watching for changes..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . | xargs -n1 -I{} make test; \
	else \
		echo "fswatch not found. Install with: brew install fswatch (macOS) or apt-get install fswatch (Ubuntu)"; \
	fi

# Default target
.DEFAULT_GOAL := help 