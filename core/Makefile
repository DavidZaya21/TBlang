# TBLang Core Makefile

.PHONY: all build clean test fmt vet install dev-setup parser

# Build directories
BUILD_DIR := build

# Go build flags
GO_BUILD_FLAGS := -ldflags="-s -w"

all: build

# Build TBLang core
build:
	@echo "Building TBLang core..."
	@mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/tblang ./cmd/tblang

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

# Install core binary globally
install: build
	@echo "Installing TBLang core globally..."
	sudo cp $(BUILD_DIR)/tblang /usr/local/bin/

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go mod tidy
	go mod download
	@mkdir -p .tblang/plugins
	@echo "Development environment ready"

# Generate parser from grammar
parser:
	@echo "Generating parser from grammar..."
	@if ! command -v antlr4 >/dev/null 2>&1; then \
		echo "Error: antlr4 not found. Please install ANTLR4."; \
		echo "Visit: https://github.com/antlr/antlr4/blob/master/doc/getting-started.md"; \
		exit 1; \
	fi
	cd grammar && antlr4 -Dlanguage=Go -o ../parser tblang.g4

# Run example
example: build
	@echo "Running example..."
	@mkdir -p .tblang/plugins
	./$(BUILD_DIR)/tblang plan ../main.tbl

# Show help
help:
	@echo "TBLang Core Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all              - Build TBLang core"
	@echo "  build            - Build TBLang core"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run tests"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  install          - Install core binary globally"
	@echo "  dev-setup        - Set up development environment"
	@echo "  parser           - Generate parser from grammar"
	@echo "  example          - Run example with main.tbl"
	@echo "  help             - Show this help"

# Default target
.DEFAULT_GOAL := help