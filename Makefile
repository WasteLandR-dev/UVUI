BINARY_NAME=uvui
BUILD_DIR=bin
CMD_DIR=./cmd
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

.DEFAULT_GOAL := help

.PHONY: help build run clean deps install test test-pretty test-oneline test-coverage test-package test-race dev fmt vet lint staticcheck check tools

## help: Show this help message
help:
	@echo "Available targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/##//g' | column -t -s ':'

## build: Build the application binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

## run: Run the application directly
run:
	@echo "🚀 Running application..."
	go run $(CMD_DIR)

## clean: Remove build artifacts and temporary files
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -f /tmp/go_test_output

## deps: Download and tidy Go dependencies
deps:
	@echo "📦 Managing dependencies..."
	go mod download
	go mod tidy

## install: Install the binary system-wide (requires sudo)
install: build
	@echo "📋 Installing $(BINARY_NAME) to /usr/local/bin/..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

## test: Run all tests with verbose output
test:
	@echo "🧪 Running tests..."
	go test ./... -v

## test-pretty: Run tests with beautified output
test-pretty:
	@echo "🧪 Running tests with pretty output..."
	@go test ./... -v 2>&1 | sed \
		-e 's/PASS/✅ PASS/g' \
		-e 's/FAIL/❌ FAIL/g' \
		-e 's/RUN/🏃 RUN/g' \
		-e 's/=== /\n=== /g' \
		-e 's/--- /    --- /g'

## test-oneline: Run tests with summary output
test-oneline:
	@echo "🧪 Running tests..."
	@go test ./... -v > /tmp/go_test_output 2>&1 && \
		echo "✅ All tests PASS!" || \
		(echo "❌ Some tests FAILED!" && \
		 echo "Failed tests:" && \
		 grep "FAIL:" /tmp/go_test_output | head -10 | sed 's/^/  - /')
	@rm -f /tmp/go_test_output

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test ./... -v -coverprofile=$(COVERAGE_FILE)
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "📊 Coverage report generated: $(COVERAGE_HTML)"

## test-package: Run tests for specific package (usage: make test-package PKG=./internal/app)
test-package:
	@if [ -z "$(PKG)" ]; then \
		echo "❌ Please specify PKG=package_path"; \
		echo "Usage: make test-package PKG=./internal/app"; \
		exit 1; \
	fi
	@echo "🧪 Running tests for $(PKG)..."
	go test $(PKG) -v

## test-race: Run tests with race condition detection
test-race:
	@echo "🧪 Running tests with race detection..."
	go test -race ./... -v

## dev: Start development mode with live reload (requires air)
dev:
	@echo "🔄 Starting development mode with live reload..."
	@if ! command -v air > /dev/null; then \
		echo "❌ 'air' not found. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi
	air -c .air.toml

## fmt: Format Go code
fmt:
	@echo "📐 Formatting code..."
	go fmt ./...

## vet: Run go vet
vet:
	@echo "🔍 Running go vet..."
	go vet ./...

## lint: Run golangci-lint
lint:
	@echo "🔍 Running golangci-lint..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "❌ 'golangci-lint' not found. Install from: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi
	golangci-lint run

## staticcheck: Run staticcheck
staticcheck:
	@echo "🔍 Running staticcheck..."
	@if ! command -v staticcheck > /dev/null; then \
		echo "❌ 'staticcheck' not found. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
		exit 1; \
	fi
	staticcheck -checks all ./...

## check: Run all quality checks (fmt, vet, lint, staticcheck, test)
check: fmt vet lint staticcheck test-oneline
	@echo "✅ All quality checks completed!"

## tools: Install development tools
tools:
	@echo "🛠️  Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/mgechev/revive@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "📋 Install golangci-lint manually from: https://golangci-lint.run/usage/install/"

## fix: Automatically fix code issues where possible
fix:
	@echo "🔧 Auto-fixing code issues..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --fix; \
	fi
	go fmt ./...
	go mod tidy