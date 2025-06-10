# Project configuration
APP_NAME := myapp
BUILD_DIR := cmd

# Default target
.PHONY: all
all: test build run

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v
	@echo "Tests passed."

# Build the application
.PHONY: build
build:
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "Build successful. Binary: $(BUILD_DIR)/$(APP_NAME)"

# Run the application
.PHONY: run
run:
	@echo "Running application..."
	@./$(BUILD_DIR)/$(APP_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
