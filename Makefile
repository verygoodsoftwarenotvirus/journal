.PHONY: build test clean install help test-entry

# Binary name
BINARY_NAME=journal

# Test artifacts directory
ARTIFACTS_DIR=./artifacts

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "Build complete!"

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install .
	@echo "Installation complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts and test artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(ARTIFACTS_DIR)
	@echo "Clean complete!"

# Test the CLI by creating a journal entry in the artifacts folder
try: build
	@echo "Creating test journal entry in $(ARTIFACTS_DIR)..."
	@mkdir -p $(ARTIFACTS_DIR)
	@JOURNAL_PATH=$(ARTIFACTS_DIR) ./$(BINARY_NAME) new -t "test,makefile"
	@echo "Test entry created! Check $(ARTIFACTS_DIR) for the result."

# Test interactive mode (requires manual input)
test-interactive: build
	@echo "Testing interactive mode..."
	@mkdir -p $(ARTIFACTS_DIR)
	@JOURNAL_PATH=$(ARTIFACTS_DIR) ./$(BINARY_NAME) new -i

# Show help
help:
	@echo "Available targets:"
	@echo "  make build          - Build the journal binary"
	@echo "  make install        - Build and install the binary to GOPATH/bin"
	@echo "  make test           - Run Go tests"
	@echo "  make test-entry     - Create a test journal entry in ./artifacts"
	@echo "  make test-interactive - Test interactive mode (requires input)"
	@echo "  make clean          - Remove build artifacts and test artifacts"
	@echo "  make help           - Show this help message"

# Default target
.DEFAULT_GOAL := help

