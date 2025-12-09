.PHONY: run build test clean swagger docker-build docker-up docker-down lint

# Application
APP_NAME=todo-api
MAIN_PATH=./cmd/server

# Go commands
GO=/opt/homebrew/bin/go
GOFLAGS=-v

# Run the application
run:
	$(GO) run $(MAIN_PATH)

# Build the application
build:
	$(GO) build $(GOFLAGS) -o bin/$(APP_NAME) $(MAIN_PATH)

# Run tests
test:
	$(GO) test ./... -v -cover

# Run tests with race detector
test-race:
	$(GO) test ./... -v -race

# Run tests with coverage report
test-coverage:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.db

# Install dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs

# Install Swagger CLI
swagger-install:
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

# Lint the code
lint:
	golangci-lint run

# Docker commands
docker-build:
	docker build -t $(APP_NAME) .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database commands
migrate:
	$(GO) run $(MAIN_PATH) migrate

# Development helpers
dev: deps run

# Production build
prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-s -w" -o bin/$(APP_NAME) $(MAIN_PATH)

# All in one setup
setup: deps swagger-install swagger
	@echo "Setup complete!"

# Help
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deps         - Download dependencies"
	@echo "  make swagger      - Generate Swagger docs"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make setup        - Full setup (deps + swagger)"
