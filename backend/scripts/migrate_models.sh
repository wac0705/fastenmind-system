#!/bin/bash

# Script to migrate all imports from old model paths to new unified model path
# This script will update all Go files to use the new domain/models package

echo "Starting model migration..."

# Define the base directory
BASE_DIR="."

# Find all Go files and update imports
find $BASE_DIR -name "*.go" -type f | while read file; do
    # Skip the new models directory itself
    if [[ $file == *"internal/domain/models"* ]]; then
        continue
    fi
    
    # Replace old import paths with new one
    sed -i 's|"github.com/fastenmind/fastener-api/internal/model"|"github.com/fastenmind/fastener-api/internal/domain/models"|g' "$file"
    sed -i 's|"github.com/fastenmind/fastener-api/internal/models"|"github.com/fastenmind/fastener-api/internal/domain/models"|g' "$file"
    
    # Update type references
    sed -i 's|model\.|models\.|g' "$file"
    sed -i 's|models\.|models\.|g' "$file"
done

echo "Model migration completed!"
echo "Please review the changes and run tests to ensure everything works correctly."