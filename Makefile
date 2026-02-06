.PHONY: help build run docker-build docker-up docker-down docker-logs

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application locally"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start Docker container"
	@echo "  make docker-down    - Stop Docker container"
	@echo "  make docker-logs    - View Docker logs"
	@echo "  make docker-rebuild - Rebuild and restart Docker"

build:
	go build -o pintu-backend .

run:
	go run main.go

docker-build:
	docker build -t pintu-backend:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

docker-rebuild:
	docker-compose down
	docker build -t pintu-backend:latest .
	docker-compose up -d

dev:
	go mod tidy
	go run main.go
