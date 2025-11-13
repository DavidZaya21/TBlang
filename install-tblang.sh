#!/bin/bash

# TBLang Installation Script
# This script installs TBLang CLI and AWS provider plugin

set -e

echo "🚀 Installing TBLang Infrastructure as Code CLI..."

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   echo "❌ This script should not be run as root"
   exit 1
fi

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is required but not installed. Please install Go 1.22+ first."
    exit 1
fi

# Check AWS CLI
if ! command -v aws &> /dev/null; then
    echo "⚠️  AWS CLI not found. Install it for AWS provider functionality."
fi

# Create temporary build directory
BUILD_DIR=$(mktemp -d)
echo "📁 Using build directory: $BUILD_DIR"

# Clone or copy TBLang source (for this demo, we'll use the current directory)
echo "📦 Building TBLang from source..."

# Build TBLang core
echo "🔨 Building TBLang core..."
cd core
go build -ldflags="-s -w" -o "$BUILD_DIR/tblang" ./cmd/tblang

# Build AWS provider plugin
echo "🔨 Building AWS provider plugin..."
cd ../plugin/aws
go build -ldflags="-s -w" -o "$BUILD_DIR/tblang-provider-aws" .

# Install binaries
echo "📥 Installing TBLang..."
sudo cp "$BUILD_DIR/tblang" /usr/local/bin/
sudo chmod +x /usr/local/bin/tblang

# Install plugins
echo "📥 Installing plugins..."
sudo mkdir -p /usr/local/lib/tblang/plugins
sudo cp "$BUILD_DIR/tblang-provider-aws" /usr/local/lib/tblang/plugins/
sudo chmod +x /usr/local/lib/tblang/plugins/tblang-provider-aws

# Cleanup
rm -rf "$BUILD_DIR"

# Verify installation
echo "✅ Verifying installation..."
if command -v tblang &> /dev/null; then
    echo "🎉 TBLang installed successfully!"
    echo ""
    echo "📖 Usage:"
    echo "  tblang version          - Show version"
    echo "  tblang plugins list     - List available plugins"
    echo "  tblang plan <file.tbl>  - Plan infrastructure changes"
    echo "  tblang apply <file.tbl> - Apply infrastructure changes"
    echo "  tblang show             - Show current state"
    echo "  tblang destroy <file.tbl> - Destroy infrastructure"
    echo ""
    echo "🚀 Get started:"
    echo "  mkdir my-infrastructure"
    echo "  cd my-infrastructure"
    echo "  # Create your infrastructure.tbl file"
    echo "  tblang plan infrastructure.tbl"
    echo ""
    tblang version
else
    echo "❌ Installation failed. TBLang not found in PATH."
    exit 1
fi