# TBLang AWS Provider Plugin Makefile

.PHONY: all build clean test fmt vet install dev-setup

# Build directories
BUILD_DIR := build

# Go build flags
GO_BUILD_FLAGS := -ldflags="-s -w"

all: build

# Build AWS plugin
build:
	@echo "Building TBLang AWS provider plugin..."
	@mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/tblang-provider-aws .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Install plugin to core's plugin directory
install: build
	@echo "Installing AWS plugin..."
	@mkdir -p ../../core/.tblang/plugins
	cp $(BUILD_DIR)/tblang-provider-aws ../../core/.tblang/plugins/

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go mod tidy
	go mod download
	@echo "Development environment ready"

# Test plugin standalone
test-plugin: build
	@echo "Testing plugin standalone..."
	./$(BUILD_DIR)/tblang-provider-aws

# Show help
help:
	@echo "TBLang AWS Provider Plugin Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all              - Build AWS plugin"
	@echo "  build            - Build AWS plugin"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run tests"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  install          - Install plugin to core directory"
	@echo "  dev-setup        - Set up development environment"
	@echo "  test-plugin      - Test plugin standalone"
	@echo "  help             - Show this help"

# Default target
.DEFAULT_GOAL := help