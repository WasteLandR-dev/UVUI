# UV CLI Manager - Installation and Usage Guide

## Prerequisites

- Go 1.21 or higher
- Git (for cloning the repository)
- Terminal with color support (recommended)

## Installation

### Option 1: Build from Source

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd uvui
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Build the application:**
   ```bash
   go build -o bin/uvui main.go
   ```

4. **Run the application:**
   ```bash
   ./bin/uvui
   ```

### Option 2: Using Make (if Makefile is present)

1. **Build:**
   ```bash
   make build
   ```

2. **Run:**
   ```bash
   make run
   ```

3. **Install system-wide (optional):**
   ```bash
   make install
   ```

### Option 3: Direct Go Run

```bash
go run main.go
```

## Development Setup

### Live Reload with Air (Recommended for Development)

1. **Install Air:**
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. **Run with live reload:**
   ```bash
   make dev
   # or
   air
   ```

## Usage

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `Tab` | Navigate to next panel |
| `Shift+Tab` | Navigate to previous panel |
| `Arrow Keys` | Navigate within panels (future phases) |
| `i` | Install UV (when not installed) |
| `r` | Refresh UV status |
| `?` | Show help message |
| `q` or `Ctrl+C` | Quit application |