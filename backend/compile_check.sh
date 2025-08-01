#!/bin/bash

echo "FastenMind Backend - Compilation Check"
echo "======================================"
echo ""

# Navigate to the backend directory
cd "$(dirname "$0")"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed in this environment."
    echo "Please run this script in GitHub Codespace where Go is available."
    exit 1
fi

# Run go build
echo "Building the application..."
if go build ./...; then
    echo ""
    echo "✅ Build completed successfully!"
else
    echo ""
    echo "❌ Build failed with errors."
    exit 1
fi