#!/bin/bash

# FastenMind Codespaces Quick Start Script
# This script quickly sets up and starts all services for testing

set -e

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# ASCII Art Banner
echo -e "${BLUE}"
cat << "EOF"
 ___         _            __  __ _         _ 
| __|_ _ ___| |_ ___ _ _ |  \/  (_)_ _  __| |
| _/ _` (_-<  _/ -_) ' \| |\/| | | ' \/ _` |
|_|\__,_/__/\__\___|_||_|_|  |_|_|_||_\__,_|
                                             
EOF
echo -e "${NC}"
echo "🚀 Codespaces Quick Start Script"
echo "================================"

# Function to print status
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to wait for a service to be ready
wait_for_service() {
    local host=$1
    local port=$2
    local service=$3
    local max_attempts=30
    local attempt=1
    
    echo -n "  Waiting for $service to be ready"
    while ! nc -z $host $port >/dev/null 2>&1; do
        if [ $attempt -eq $max_attempts ]; then
            echo ""
            print_error "$service failed to start"
            return 1
        fi
        echo -n "."
        sleep 1
        ((attempt++))
    done
    echo ""
    print_status "$service is ready"
    return 0
}

# Step 1: Check environment
echo ""
echo "Step 1: Checking environment..."

if [ -n "$CODESPACES" ]; then
    print_status "Running in GitHub Codespaces"
else
    print_warning "Not running in Codespaces - some features may not work as expected"
fi

# Check required tools
for cmd in psql go npm; do
    if command_exists $cmd; then
        print_status "$cmd is installed"
    else
        print_error "$cmd is not installed"
        exit 1
    fi
done

# Step 2: Initialize database
echo ""
echo "Step 2: Initializing database..."

# Check if database is already initialized
if sudo -u postgres psql -t -c "SELECT 1 FROM pg_database WHERE datname='fastenmind_db';" | grep -q 1; then
    print_info "Database already exists"
    read -p "Do you want to reinitialize the database? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        /workspaces/fastenmind-system/scripts/init-database.sh
    else
        print_status "Using existing database"
        sudo service postgresql start
    fi
else
    /workspaces/fastenmind-system/scripts/init-database.sh
fi

# Step 3: Setup backend
echo ""
echo "Step 3: Setting up backend..."

cd /workspaces/fastenmind-system/backend

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    print_info "Creating .env file..."
    cp .env.example .env
    print_status ".env file created"
else
    print_status ".env file already exists"
fi

# Make sure the binary is executable
chmod +x fastenmind-api

# Step 4: Setup frontend
echo ""
echo "Step 4: Setting up frontend..."

cd /workspaces/fastenmind-system/frontend

# Install dependencies if needed
if [ ! -d node_modules ]; then
    print_info "Installing frontend dependencies..."
    npm install
    print_status "Dependencies installed"
else
    print_status "Dependencies already installed"
fi

# Create .env.local if it doesn't exist
if [ ! -f .env.local ]; then
    echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local
    print_status ".env.local file created"
else
    print_status ".env.local file already exists"
fi

# Step 5: Start services
echo ""
echo "Step 5: Starting services..."

# Kill any existing processes
print_info "Cleaning up existing processes..."
pkill -f fastenmind-api || true
pkill -f "next dev" || true
sleep 2

# Start backend
print_info "Starting backend API..."
cd /workspaces/fastenmind-system/backend
nohup ./fastenmind-api > /tmp/backend.log 2>&1 &
BACKEND_PID=$!
print_status "Backend started (PID: $BACKEND_PID)"

# Wait for backend to be ready
wait_for_service localhost 8080 "Backend API"

# Start frontend
print_info "Starting frontend..."
cd /workspaces/fastenmind-system/frontend
nohup npm run dev > /tmp/frontend.log 2>&1 &
FRONTEND_PID=$!
print_status "Frontend started (PID: $FRONTEND_PID)"

# Wait for frontend to be ready
wait_for_service localhost 3000 "Frontend"

# Step 6: Display information
echo ""
echo "============================================"
echo -e "${GREEN}✅ All services started successfully!${NC}"
echo ""
echo "🌐 Access Points:"
echo "   Frontend:    http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo "   Health Check: http://localhost:8080/health"
echo ""
echo "📝 Default Login:"
echo "   Email:    admin@fastenmind.com"
echo "   Password: password123"
echo ""
echo "📊 Service Status:"
echo "   PostgreSQL: $(sudo service postgresql status | grep -o 'online\|offline')"
echo "   Backend PID: $BACKEND_PID"
echo "   Frontend PID: $FRONTEND_PID"
echo ""
echo "📁 Log Files:"
echo "   Backend:  /tmp/backend.log"
echo "   Frontend: /tmp/frontend.log"
echo ""
echo "🛑 To stop all services, run:"
echo "   pkill -f fastenmind-api && pkill -f 'next dev'"
echo ""

# Create stop script
cat > /workspaces/fastenmind-system/stop-services.sh << 'EOF'
#!/bin/bash
echo "Stopping FastenMind services..."
pkill -f fastenmind-api && echo "✓ Backend stopped"
pkill -f "next dev" && echo "✓ Frontend stopped"
sudo service postgresql stop && echo "✓ PostgreSQL stopped"
echo "All services stopped."
EOF

chmod +x /workspaces/fastenmind-system/stop-services.sh

# Create status check script
cat > /workspaces/fastenmind-system/check-status.sh << 'EOF'
#!/bin/bash
echo "FastenMind Service Status"
echo "========================"
echo -n "PostgreSQL: "
sudo service postgresql status | grep -o 'online\|offline' || echo "offline"
echo -n "Backend API: "
curl -s http://localhost:8080/health > /dev/null 2>&1 && echo "running" || echo "stopped"
echo -n "Frontend: "
curl -s http://localhost:3000 > /dev/null 2>&1 && echo "running" || echo "stopped"
EOF

chmod +x /workspaces/fastenmind-system/check-status.sh

print_info "Helper scripts created:"
print_info "  ./stop-services.sh  - Stop all services"
print_info "  ./check-status.sh   - Check service status"
echo ""

# Show recent logs
echo "📜 Recent Backend Logs:"
echo "----------------------"
tail -n 10 /tmp/backend.log 2>/dev/null || echo "No logs yet"
echo ""

# Final message
echo -e "${BLUE}Happy testing! 🎉${NC}"