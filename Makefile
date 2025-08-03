.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/ecommerce-api .

# Run the application
run:
	go run .

# Run tests
test:
	go test -v ./...

# Run tests with Ginkgo
test-ginkgo:
	ginkgo -v

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Docker build
docker-build:
	docker build -t ecommerce-api .

# Docker run
docker-run:
	docker-compose up -d

# Docker stop
docker-stop:
	docker-compose down

# Database migration (if using external DB)
migrate:
	@echo "Running database migrations..."
	go run . -migrate

# Seed database
seed:
	@echo "Seeding database..."
	go run . -seed
