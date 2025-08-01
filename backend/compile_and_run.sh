#!/bin/bash

echo "FastenMind Backend Compilation and Run Script"
echo "============================================="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo -e "${RED}Go is not installed. Please install Go first.${NC}"
    echo "You can install Go in GitHub Codespace by running:"
    echo "  sudo apt update && sudo apt install -y golang-go"
    exit 1
fi

echo -e "${GREEN}Go is installed: $(go version)${NC}"

# Navigate to backend directory
cd "$(dirname "$0")" || exit

# Download dependencies
echo -e "${YELLOW}Downloading dependencies...${NC}"
go mod download

# Verify dependencies
echo -e "${YELLOW}Verifying dependencies...${NC}"
go mod verify

# Run go mod tidy to clean up
echo -e "${YELLOW}Tidying up go.mod...${NC}"
go mod tidy

# Build the application
echo -e "${YELLOW}Building the application...${NC}"
if go build -o server cmd/server/main.go; then
    echo -e "${GREEN}Build successful!${NC}"
    
    echo ""
    echo -e "${YELLOW}To run the server, execute:${NC}"
    echo "  ./server"
    echo ""
    echo -e "${YELLOW}Make sure you have the following environment variables set:${NC}"
    echo "  - Database connection settings"
    echo "  - JWT secret"
    echo "  - Other configuration"
    echo ""
    echo -e "${YELLOW}You can copy .env.example to .env and modify it:${NC}"
    echo "  cp .env.example .env"
else
    echo -e "${RED}Build failed! Please check the errors above.${NC}"
    exit 1
fi