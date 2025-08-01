.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build docker image
docker-build:
	docker build -t fastenmind-backend:latest .

# Run with docker-compose
docker-run:
	docker-compose up -d

# Stop docker containers
docker-stop:
	docker-compose down

# View logs
logs:
	docker-compose logs -f backend

# Database migrations up
migrate-up:
	migrate -path ./migrations -database "postgresql://fastenmind:fastenmind123@localhost:5432/fastenmind_db?sslmode=disable" up

# Database migrations down
migrate-down:
	migrate -path ./migrations -database "postgresql://fastenmind:fastenmind123@localhost:5432/fastenmind_db?sslmode=disable" down

# Create a new migration
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy