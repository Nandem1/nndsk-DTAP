# DTAP Makefile - Normal and WoE Versions

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=dtap
BINARY_NAME_WOE=dtap-woe
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_UNIX_WOE=$(BINARY_NAME_WOE)_unix

# Build directories
CMD_DIR=cmd
NORMAL_CMD=$(CMD_DIR)/dtap
WOE_CMD=$(CMD_DIR)/dtap-woe

# Output directories
OUTPUT_DIR=.
CONFIGS_DIR=configs

.PHONY: all build build-normal build-woe clean test deps help

# Default target
all: build

# Build both versions
build: build-normal build-woe

# Build normal version
build-normal:
	@echo "Building normal version..."
	$(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME).exe -v $(NORMAL_CMD)

# Build WoE version
build-woe:
	@echo "Building WoE version..."
	$(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_NAME_WOE).exe -v $(WOE_CMD)

# Build normal version for Linux
build-linux-normal:
	@echo "Building normal version for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_UNIX) -v $(NORMAL_CMD)

# Build WoE version for Linux
build-linux-woe:
	@echo "Building WoE version for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUTPUT_DIR)/$(BINARY_UNIX_WOE) -v $(WOE_CMD)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(OUTPUT_DIR)/$(BINARY_NAME).exe
	rm -f $(OUTPUT_DIR)/$(BINARY_NAME_WOE).exe
	rm -f $(OUTPUT_DIR)/$(BINARY_UNIX)
	rm -f $(OUTPUT_DIR)/$(BINARY_UNIX_WOE)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Get dependencies
deps:
	@echo "Getting dependencies..."
	$(GOGET) -d -v ./...

# Install dependencies
install:
	@echo "Installing dependencies..."
	$(GOGET) -v ./...

# Run normal version
run-normal: build-normal
	@echo "Running normal version..."
	./$(BINARY_NAME).exe

# Run WoE version
run-woe: build-woe
	@echo "Running WoE version..."
	./$(BINARY_NAME_WOE).exe

# Setup project structure
setup:
	@echo "Setting up project structure..."
	@mkdir -p $(CONFIGS_DIR)
	@echo "Project structure ready!"

# Show help
help:
	@echo "Available targets:"
	@echo "  all              - Build both versions (default)"
	@echo "  build            - Build both versions"
	@echo "  build-normal     - Build normal version only"
	@echo "  build-woe        - Build WoE version only"
	@echo "  build-linux-normal - Build normal version for Linux"
	@echo "  build-linux-woe  - Build WoE version for Linux"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run tests"
	@echo "  deps             - Get dependencies"
	@echo "  install          - Install dependencies"
	@echo "  run-normal       - Build and run normal version"
	@echo "  run-woe          - Build and run WoE version"
	@echo "  setup            - Setup project structure"
	@echo "  help             - Show this help message"