#!/bin/bash

# Database initialization script for Codespaces
# This script sets up PostgreSQL and initializes the FastenMind database

set -e  # Exit on error

echo "🗄️  FastenMind Database Initialization Script"
echo "============================================"

# Color codes for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check if running in Codespaces
if [ -z "$CODESPACES" ]; then
    print_warning "This script is optimized for GitHub Codespaces environment"
fi

# Step 1: Install PostgreSQL if not already installed
echo ""
echo "Step 1: Checking PostgreSQL installation..."
if ! command -v psql &> /dev/null; then
    print_warning "PostgreSQL not found. Installing..."
    sudo apt-get update
    sudo apt-get install -y postgresql postgresql-contrib
    print_status "PostgreSQL installed successfully"
else
    print_status "PostgreSQL is already installed"
fi

# Step 2: Start PostgreSQL service
echo ""
echo "Step 2: Starting PostgreSQL service..."
sudo service postgresql start
if sudo service postgresql status | grep -q "online"; then
    print_status "PostgreSQL service is running"
else
    print_error "Failed to start PostgreSQL service"
    exit 1
fi

# Step 3: Create database user and database
echo ""
echo "Step 3: Creating database and user..."
sudo -u postgres psql << EOF
-- Drop existing connections to the database
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = 'fastenmind_db'
AND pid <> pg_backend_pid();

-- Drop database if exists
DROP DATABASE IF EXISTS fastenmind_db;

-- Drop user if exists
DROP USER IF EXISTS fastenmind;

-- Create user
CREATE USER fastenmind WITH PASSWORD 'fastenmind123';

-- Create database
CREATE DATABASE fastenmind_db OWNER fastenmind;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE fastenmind_db TO fastenmind;

-- Show results
\echo 'Database created:'
\l fastenmind_db
\echo 'User created:'
\du fastenmind
EOF

if [ $? -eq 0 ]; then
    print_status "Database and user created successfully"
else
    print_error "Failed to create database or user"
    exit 1
fi

# Step 4: Initialize database schema
echo ""
echo "Step 4: Initializing database schema..."

# Find the database directory
DB_DIR="/workspaces/fastenmind-system/database/init"
if [ ! -d "$DB_DIR" ]; then
    DB_DIR="./database/init"
fi

if [ ! -d "$DB_DIR" ]; then
    print_error "Cannot find database initialization scripts"
    print_warning "Please run this script from the project root directory"
    exit 1
fi

# Execute initialization scripts in order
for script in "$DB_DIR"/*.sql; do
    if [ -f "$script" ]; then
        echo "  Executing: $(basename "$script")"
        PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db -f "$script" > /dev/null 2>&1
        if [ $? -eq 0 ]; then
            print_status "$(basename "$script") executed successfully"
        else
            print_error "Failed to execute $(basename "$script")"
            PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db -f "$script"
        fi
    fi
done

# Step 5: Verify database setup
echo ""
echo "Step 5: Verifying database setup..."
TABLE_COUNT=$(PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
TABLE_COUNT=$(echo $TABLE_COUNT | xargs)  # Trim whitespace

if [ "$TABLE_COUNT" -gt 0 ]; then
    print_status "Database initialized with $TABLE_COUNT tables"
else
    print_error "No tables found in database"
    exit 1
fi

# Step 6: Test database connection
echo ""
echo "Step 6: Testing database connection..."
PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db -c "SELECT version();" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    print_status "Database connection test successful"
else
    print_error "Database connection test failed"
    exit 1
fi

# Print summary
echo ""
echo "============================================"
echo -e "${GREEN}✅ Database initialization completed!${NC}"
echo ""
echo "Database Details:"
echo "  Host: localhost"
echo "  Port: 5432"
echo "  Database: fastenmind_db"
echo "  Username: fastenmind"
echo "  Password: fastenmind123"
echo ""
echo "Connection string:"
echo "  postgres://fastenmind:fastenmind123@localhost:5432/fastenmind_db"
echo ""

# Create a test query script
cat > /tmp/test-db.sql << 'EOF'
-- Test query to verify database structure
SELECT 
    'Companies' as table_name, COUNT(*) as row_count FROM companies
UNION ALL
SELECT 
    'Accounts' as table_name, COUNT(*) as row_count FROM accounts
UNION ALL
SELECT 
    'Customers' as table_name, COUNT(*) as row_count FROM customers
UNION ALL
SELECT 
    'Inquiries' as table_name, COUNT(*) as row_count FROM inquiries
ORDER BY table_name;
EOF

echo "Run this command to check table row counts:"
echo "  PGPASSWORD=fastenmind123 psql -h localhost -U fastenmind -d fastenmind_db -f /tmp/test-db.sql"
echo ""