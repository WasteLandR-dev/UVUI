.PHONY: build run clean install deps

# Build the application
build:
	go build -o bin/uvui main.go

# Run the application
run:
	go run main.go

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

# Development mode with live reload (requires air)
dev:
	air
