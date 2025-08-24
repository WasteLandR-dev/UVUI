.PHONY: build run clean install deps test

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

# Run tests
test:
	go test ./... -v

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
check: fmt vet lint test