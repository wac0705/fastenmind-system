#!/bin/bash

# Alternative startup script using Docker (no sudo required)

echo "🐳 Starting FastenMind with Docker"
echo "=================================="

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not available in this environment"
    echo "Please use the PostgreSQL installation method instead"
    exit 1
fi

# Start PostgreSQL using Docker
echo "Starting PostgreSQL with Docker..."
docker run -d \
  --name fastenmind-postgres \
  -e POSTGRES_USER=fastenmind \
  -e POSTGRES_PASSWORD=fastenmind123 \
  -e POSTGRES_DB=fastenmind_db \
  -p 5432:5432 \
  postgres:15-alpine

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to start..."
sleep 10

# Initialize database
echo "Initializing database..."
for script in /workspaces/fastenmind-system/database/init/*.sql; do
    echo "Executing: $(basename "$script")"
    docker exec -i fastenmind-postgres psql -U fastenmind -d fastenmind_db < "$script"
done

# Start backend
echo "Starting backend..."
cd /workspaces/fastenmind-system/backend
cp .env.example .env
chmod +x fastenmind-api
./fastenmind-api &

# Start frontend
echo "Starting frontend..."
cd /workspaces/fastenmind-system/frontend
npm install
npm run dev &

echo "✅ All services started!"
echo "Frontend: http://localhost:3000"
echo "Backend: http://localhost:8080"