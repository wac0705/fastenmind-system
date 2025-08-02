#!/bin/bash

# Build script for FastenMind backend

echo "Building the application..."

# Build the application
go build -o server cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
else
    echo "Build failed!"
    exit 1
fi