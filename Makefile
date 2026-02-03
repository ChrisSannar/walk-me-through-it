.PHONY: build test lint fmt clean install

# Variables
BINARY_NAME=wmti
BUILD_DIR=./build
MAIN_PACKAGE=./cmd/wmti

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofumpt -w .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod tidy

# Run the application locally
run:
	go run $(MAIN_PACKAGE)

# Run the init command
run-init:
	go run $(MAIN_PACKAGE) init

# Development mode with auto-reload (requires air)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air is not installed. Install with: go install github.com/cosmtrek/air@latest"; \
	fi
