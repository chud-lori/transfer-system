APP_NAME := transfer-system
BUILD_DIR := bin

# Default target
.PHONY: all
all: test build run

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v
	@echo "Tests passed."

# Install dependencies for running locally
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@go mod download

# Build the application
.PHONY: build
build:
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd
	@echo "Build successful. Binary: $(BUILD_DIR)/$(APP_NAME)"

# Generate Swagger documentation
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/main.go

# Run the application
.PHONY: run
run: build
	@echo "Running application..."
	@./$(BUILD_DIR)/$(APP_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Help message
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make           Run tests, build, and start the app"
	@echo "  make test      Run all tests"
	@echo "  make deps      Download Go dependencies"
	@echo "  make build     Build the Go binary into $(BUILD_DIR)/$(APP_NAME)"
	@echo "  make swagger   Generate Swagger documentation"
	@echo "  make run       Run the application"
	@echo "  make clean     Remove built binaries"
