#!/bin/bash

# Install all dependencies for FastenMind in Codespaces

set -e

echo "📦 Installing Dependencies for FastenMind"
echo "========================================"

# Update package list
echo "Updating package list..."
sudo apt-get update

# Install PostgreSQL
echo ""
echo "Installing PostgreSQL..."
sudo apt-get install -y postgresql postgresql-contrib
echo "✓ PostgreSQL installed"

# Install additional tools
echo ""
echo "Installing additional tools..."
sudo apt-get install -y netcat-openbsd htop
echo "✓ Additional tools installed"

# Check Go installation
echo ""
if command -v go &> /dev/null; then
    echo "✓ Go is already installed: $(go version)"
else
    echo "Installing Go..."
    wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    rm go1.21.5.linux-amd64.tar.gz
    echo "✓ Go installed"
fi

# Check Node.js installation
echo ""
if command -v node &> /dev/null; then
    echo "✓ Node.js is already installed: $(node --version)"
else
    echo "Installing Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt-get install -y nodejs
    echo "✓ Node.js installed"
fi

echo ""
echo "✅ All dependencies installed successfully!"
echo ""
echo "Now run: ./scripts/start-codespaces.sh"