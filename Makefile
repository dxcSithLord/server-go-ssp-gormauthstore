.PHONY: all test lint security build clean deps fmt tools help

# Go parameters
GO := go
GOFLAGS := -v
GOLANGCI_LINT := golangci-lint
GOSEC := gosec
GOVULNCHECK := govulncheck

# Project parameters
PACKAGE := github.com/sqrldev/server-go-ssp-gormauthstore
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Default target
all: deps lint test build

## help: Display this help message
help:
	@echo "Available targets:"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## /  /' | column -t -s ':'
	@echo ""

## deps: Download and tidy Go module dependencies
deps:
	@echo "==> Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	$(GO) mod verify

## lint: Run linters (requires golangci-lint)
lint:
	@echo "==> Running linters..."
	$(GOLANGCI_LINT) run --timeout=5m

## vet: Run go vet
vet:
	@echo "==> Running go vet..."
	$(GO) vet ./...

## fmt: Format Go code
fmt:
	@echo "==> Formatting code..."
	gofmt -s -w .
	goimports -w .

## fmt-check: Check if code is formatted
fmt-check:
	@echo "==> Checking code formatting..."
	@gofmt_output=$$(gofmt -l .); \
	if [ -n "$$gofmt_output" ]; then \
		echo "Files not formatted:"; \
		echo "$$gofmt_output"; \
		exit 1; \
	fi
	@echo "All files are properly formatted"

## test: Run tests with race detection
test:
	@echo "==> Running tests..."
	$(GO) test $(GOFLAGS) -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GO) tool cover -func=$(COVERAGE_FILE)

## test-short: Run tests without database integration
test-short:
	@echo "==> Running short tests..."
	$(GO) test $(GOFLAGS) -short ./...

## test-coverage: Generate HTML coverage report
test-coverage: test
	@echo "==> Generating coverage report..."
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

## test-coverage-check: Check if coverage meets threshold (70%)
test-coverage-check: test
	@echo "==> Checking coverage threshold..."
	@coverage=$$($(GO) tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$coverage%"; \
	if [ $$(echo "$$coverage < 70" | bc -l) -eq 1 ]; then \
		echo "ERROR: Coverage $$coverage% is below 70% threshold"; \
		exit 1; \
	else \
		echo "SUCCESS: Coverage meets threshold"; \
	fi

## bench: Run benchmarks
bench:
	@echo "==> Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

## security: Run security checks
security:
	@echo "==> Running security checks..."
	@echo "--- gosec ---"
	$(GOSEC) -fmt=text ./... || true
	@echo ""
	@echo "--- govulncheck ---"
	$(GOVULNCHECK) ./...

## security-ci: Run security checks for CI (strict)
security-ci:
	@echo "==> Running strict security checks..."
	$(GOSEC) -fmt=sarif -out=gosec-results.sarif ./...
	$(GOVULNCHECK) ./...

## build: Build the package
build:
	@echo "==> Building..."
	$(GO) build $(GOFLAGS) ./...

## build-all: Build for all platforms
build-all:
	@echo "==> Building for all platforms..."
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) ./...
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) ./...
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) ./...
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) ./...
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) ./...

## clean: Remove build artifacts
clean:
	@echo "==> Cleaning..."
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -f gosec-results.sarif
	$(GO) clean

## tools: Install development tools
tools:
	@echo "==> Installing development tools..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install golang.org/x/tools/cmd/goimports@latest

## check: Run all checks (lint, security, test)
check: lint security test
	@echo "==> All checks passed!"

## ci: Run CI pipeline locally
ci: deps fmt-check lint security test build
	@echo "==> CI pipeline completed successfully!"

## pre-commit: Run checks before committing
pre-commit: fmt lint test-short
	@echo "==> Pre-commit checks passed!"

## update-deps: Update all dependencies to latest versions
update-deps:
	@echo "==> Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy
	$(GO) mod verify

## outdated: List outdated dependencies
outdated:
	@echo "==> Checking for outdated dependencies..."
	$(GO) list -u -m all 2>/dev/null | grep '\[' || echo "All dependencies are up to date"

## verify: Verify module dependencies
verify:
	@echo "==> Verifying module..."
	$(GO) mod verify
	$(GO) mod why -m all
