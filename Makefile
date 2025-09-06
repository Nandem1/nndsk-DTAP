# DTAP-Go Makefile - WoE Optimized Edition

# Variables
BINARY_NAME=dtap
WOE_BINARY_NAME=dtap-woe
ADDRESS_FINDER=find_addresses
VERSION=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M:%S)
GO_VERSION=$(shell go version | cut -d' ' -f3)

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# Default target
.PHONY: all
all: build-woe

# Build standard version
.PHONY: build
build:
	@echo "Building standard DTAP-Go..."
	go build $(LDFLAGS) -o $(BINARY_NAME).exe .
	@echo "✅ Standard build complete: $(BINARY_NAME).exe"

# Build WoE optimized version
.PHONY: build-woe
build-woe:
	@echo "Building WoE Optimized DTAP-Go..."
	go build $(LDFLAGS) -o $(WOE_BINARY_NAME).exe main_woe.go
	@echo "✅ WoE build complete: $(WOE_BINARY_NAME).exe"

# Build address finder
.PHONY: build-finder
build-finder:
	@echo "Building address finder..."
	go build $(LDFLAGS) -o $(ADDRESS_FINDER).exe ./cmd/find_addresses
	@echo "✅ Address finder build complete: $(ADDRESS_FINDER).exe"

# Build all versions
.PHONY: build-all
build-all: build build-woe build-finder
	@echo "✅ All builds complete"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@if exist $(BINARY_NAME).exe del /q $(BINARY_NAME).exe
	@if exist $(WOE_BINARY_NAME).exe del /q $(WOE_BINARY_NAME).exe
	@if exist $(ADDRESS_FINDER).exe del /q $(ADDRESS_FINDER).exe
	@echo "✅ Clean complete"

# Run standard version
.PHONY: run
run: build
	@echo "Running standard DTAP-Go..."
	./$(BINARY_NAME).exe

# Run WoE version
.PHONY: run-woe
run-woe: build-woe
	@echo "Running WoE Optimized DTAP-Go..."
	./$(WOE_BINARY_NAME).exe

# Run address finder
.PHONY: run-finder
run-finder: build-finder
	@echo "Running address finder..."
	./$(ADDRESS_FINDER).exe

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "✅ Tests complete"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@if exist go.exe (go vet ./...) else (echo "Go not found, skipping lint")
	@echo "✅ Lint complete"

# Show help
.PHONY: help
help:
	@echo "DTAP-Go Makefile - WoE Optimized Edition"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build standard version"
	@echo "  build-woe    - Build WoE optimized version (recommended)"
	@echo "  build-finder - Build address finder tool"
	@echo "  build-all    - Build all versions"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Run standard version"
	@echo "  run-woe      - Run WoE optimized version"
	@echo "  run-finder   - Run address finder"
	@echo "  deps         - Install dependencies"
	@echo "  test         - Run tests"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Recommended for WoE: make build-woe && make run-woe"

# Development targets
.PHONY: dev
dev: fmt lint test build-woe
	@echo "✅ Development build complete"

# Production build with optimizations
.PHONY: prod
prod:
	@echo "Building production WoE version with optimizations..."
	go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o $(WOE_BINARY_NAME).exe main_woe.go
	@echo "✅ Production build complete: $(WOE_BINARY_NAME).exe"

# Quick build for testing
.PHONY: quick
quick:
	@echo "Quick build WoE version..."
	go build -o $(WOE_BINARY_NAME).exe main_woe.go
	@echo "✅ Quick build complete"