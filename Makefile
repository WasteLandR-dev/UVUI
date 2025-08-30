.PHONY: build run clean install deps test test-pretty test-color test-coverage test-oneline ch


# Build the application
build:
	go build -o bin/uvui ./cmd

# Run the application
run:
	go run ./cmd

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install the binary system-wide (optional)
install: build
	sudo cp bin/uvui /usr/local/bin/

test:
	go test ./... -v

# Beautified test output with colors and formatting
test-pretty:
	@echo "🧪 Running tests with pretty output..."
	@go test ./... -v 2>&1 | sed \
		-e 's/PASS/✅ PASS/g' \
		-e 's/FAIL/❌ FAIL/g' \
		-e 's/RUN/🏃 RUN/g' \
		-e 's/=== /\n=== /g' \
		-e 's/--- /    --- /g'
    
# Simple one-line test summary
test-oneline:
	@echo "🧪 Running tests..."
	@go test ./... -v > /tmp/go_test_output 2>&1 && \
		echo "✅ All tests PASS!" || \
		(echo "❌ Some tests FAILED!" && \
		 echo "Failed tests:" && \
		 grep "FAIL:" /tmp/go_test_output | head -10 | sed 's/^/  - /')
	@rm -f /tmp/go_test_output

# Run tests with coverage
test-coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests for specific package
test-package:
	@echo "Usage: make test-package PKG=./internal/app"
	@if [ -z "$(PKG)" ]; then echo "Please specify PKG=package_path"; exit 1; fi
	go test $(PKG) -v

# Run tests with race detection
test-race:
	go test -race ./... -v

# Development mode with live reload (requires air)
dev:
	air -c .air.toml

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Run all quality checks
check: fmt vet lint test-oneline

ch:
	golangci-lint run --fix
	staticcheck -checks all ./...
	revive -formatter friendly -exclude ./vendor/... ./...
	./bin/gosec -fmt=golint -quiet ./...