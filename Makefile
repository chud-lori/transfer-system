APP_NAME := transfer-system
BUILD_DIR := bin

all: test build run

test:
	@echo "Running tests..."
	@go test ./... -v
	@echo "Tests passed."

deps:
	@echo "Downloading dependencies..."
	@go mod download

build:
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd
	@echo "Build successful. Binary: $(BUILD_DIR)/$(APP_NAME)"

swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/main.go

run: build
	@echo "Running application..."
	@./$(BUILD_DIR)/$(APP_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

help:
	@echo "Usage:"
	@echo "  make           Run tests, build, and start the app"
	@echo "  make test      Run all tests"
	@echo "  make deps      Download Go dependencies"
	@echo "  make build     Build the Go binary into $(BUILD_DIR)/$(APP_NAME)"
	@echo "  make swagger   Generate Swagger documentation"
	@echo "  make run       Run the application"
	@echo "  make clean     Remove built binaries"

.PHONY:
	all test deps build swagger run clean help
