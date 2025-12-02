.PHONY: help build up down logs clean dev backend-dev frontend-dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## View logs from all services
	docker-compose logs -f

clean: ## Remove all containers, volumes, and images
	docker-compose down -v
	docker system prune -f

dev: ## Start services in development mode
	docker-compose up postgres milvus etcd minio -d
	@echo "Database and Milvus started. Run 'make backend-dev' and 'make frontend-dev' in separate terminals."

backend-dev: ## Run backend in development mode (requires dev services)
	cd backend && go run cmd/server/main.go

frontend-dev: ## Run frontend in development mode
	cd frontend && npm install && npm run dev

test-backend: ## Run backend tests
	cd backend && go test ./...

init-db: ## Initialize database
	docker-compose exec postgres psql -U postgres -d latex_ppt -f /docker-entrypoint-initdb.d/init-db.sql

backend-logs: ## View backend logs
	docker-compose logs -f backend

frontend-logs: ## View frontend logs
	docker-compose logs -f frontend

restart: ## Restart all services
	docker-compose restart

status: ## Show status of all services
	docker-compose ps
