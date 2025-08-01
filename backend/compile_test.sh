#!/bin/bash

echo "FastenMind Backend Compilation Test"
echo "==================================="
echo

cd "$(dirname "$0")"

echo "Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

go version
echo

echo "Building the application..."
if go build -v ./... 2>&1; then
    echo
    echo "Build successful!"
else
    echo
    echo "Build failed! Please check the errors above."
    exit 1
fi