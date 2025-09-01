# Variables
APP_NAME = url-shortener
DOCKER_IMAGE = $(APP_NAME):latest
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_PROD = docker-compose -f docker-compose.prod.yml

# Default target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development commands
.PHONY: dev
dev: ## Start development environment
	$(DOCKER_COMPOSE) up --build

.PHONY: dev-detach
dev-detach: ## Start development environment in background
	$(DOCKER_COMPOSE) up --build -d

.PHONY: dev-logs
dev-logs: ## Show development logs
	$(DOCKER_COMPOSE) logs -f

.PHONY: dev-stop
dev-stop: ## Stop development environment
	$(DOCKER_COMPOSE) down

# Production commands
.PHONY: prod
prod: ## Start production environment
	$(DOCKER_COMPOSE_PROD) up --build -d

.PHONY: prod-logs
prod-logs: ## Show production logs
	$(DOCKER_COMPOSE_PROD) logs -f

.PHONY: prod-stop
prod-stop: ## Stop production environment
	$(DOCKER_COMPOSE_PROD) down

.PHONY: prod-restart
prod-restart: ## Restart production environment
	$(DOCKER_COMPOSE_PROD) restart

# Docker commands
.PHONY: build
build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

.PHONY: run
run: ## Run Docker container
	docker run -p 9808:9808 --name $(APP_NAME) $(DOCKER_IMAGE)

.PHONY: stop
stop: ## Stop Docker container
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# Utility commands
.PHONY: clean
clean: ## Clean up Docker resources
	docker system prune -f
	docker volume prune -f

.PHONY: logs
logs: ## Show container logs
	docker logs -f $(APP_NAME)

.PHONY: shell
shell: ## Access container shell
	docker exec -it $(APP_NAME) /bin/sh

.PHONY: redis-cli
redis-cli: ## Access Redis CLI
	docker exec -it url-shortener-redis redis-cli

# Testing
.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: test-docker
test-docker: ## Run tests in Docker
	docker run --rm -v $(PWD):/app -w /app golang:1.24-alpine sh -c "apk add --no-cache git && go test ./..."

# Database commands
.PHONY: db-backup
db-backup: ## Backup Redis data
	docker exec url-shortener-redis redis-cli BGSAVE
	docker cp url-shortener-redis:/data/dump.rdb ./backup/

.PHONY: db-restore
db-restore: ## Restore Redis data
	docker cp ./backup/dump.rdb url-shortener-redis:/data/

# Monitoring
.PHONY: status
status: ## Show service status
	$(DOCKER_COMPOSE) ps

.PHONY: health
health: ## Check service health
	curl -f http://localhost:9808/ || echo "Service not healthy"

# Development setup
.PHONY: setup
setup: ## Initial setup
	mkdir -p backup
	chmod +x scripts/*.sh || true

# Clean everything
.PHONY: clean-all
clean-all: ## Clean everything including volumes
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_COMPOSE_PROD) down -v
	docker system prune -af
	docker volume prune -f
