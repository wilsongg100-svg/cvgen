# ==================================================================================
# cvgen - Makefile
# ==================================================================================

APP_NAME    = cvgen
BINARY      = ./bin/$(APP_NAME)
CMD         = ./cmd/main.go
COMPOSE_DEV  := docker compose -f docker-compose.yml
COMPOSE_PROD := docker compose -f docker-compose.prod.yml


CYAN   := \033[36m
GREEN  := \033[32m
YELLOW := \033[33m
RED    := \033[31m
RESET  := \033[0m

.PHONY: help build run dev fmt lint test tidy clean \
        docker-up docker-down docker-rebuild docker-logs docker-logs-api docker-clean \
        prod-up prod-down prod-rebuild prod-logs \
        health commit commit-push first-commit-push

# ==================================================================================
# HELP
# ==================================================================================

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ==================================================================================
# LOCAL DEV (API native, only postgres in docker)
# ==================================================================================

build: ## Build the binary into ./bin/
	@echo "$(CYAN)Building...$(RESET)"
	@mkdir -p bin
	@go build -o $(BINARY) $(CMD)
	@echo "$(GREEN)Built: $(BINARY)$(RESET)"

run: build ## Build and run the binary locally
	@echo "$(GREEN)Running...$(RESET)"
	@$(BINARY)

dev: ## Run with live reload API via air
	@echo "$(GREEN)Starting API with air...$(RESET)"
	@air

fmt:  ## Format all Go code
	@go fmt ./...

lint: ## Lint (requires golangci-lint)
	@golangci-lint run ./...

test: ## Run all tests
	@go test ./... -v

test-cover: ## Run tests with coverage report
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out

tidy: ## Tidy go.mod and go.sum
	@go mod tidy

clean: ## Remove built binary
	@rm -rf bin/
	@echo "$(GREEN)Cleaned$(RESET)"

# ==================================================================================
# DOCKER DEV
# ==================================================================================

docker-up: ## Start all services in Docker (dev)
	@$(COMPOSE_DEV) up --build -d
	@echo "$(GREEN)Dev services running$(RESET)"

docker-down: ## Stop dev Docker services
	@$(COMPOSE_DEV) down

docker-rebuild: ## Rebuild and restart dev services
	@$(COMPOSE_DEV) down
	@$(COMPOSE_DEV) up --build -d

docker-logs: ## Tail all dev container logs
	@$(COMPOSE_DEV) logs -f

docker-logs-api: ## Tail dev API logs only
	@$(COMPOSE_DEV) logs -f api

docker-clean: ## Stop dev containers and wipe volumes
	@$(COMPOSE_DEV) down -v
	@echo "$(RED)Dev volumes removed$(RESET)"

# ==================================================================================
# DOCKER PROD
# ==================================================================================

prod-up: ## Start all services in Docker (prod)
	@$(COMPOSE_PROD) up --build -d
	@echo "$(GREEN)Prod services running$(RESET)"

prod-down: ## Stop prod Docker services
	@$(COMPOSE_PROD) down

prod-rebuild: ## Rebuild and restart prod services
	@$(COMPOSE_PROD) down
	@$(COMPOSE_PROD) up --build -d

prod-logs: ## Tail all prod container logs
	@$(COMPOSE_PROD) logs -f

prod-logs-api: ## Tail prod API logs only
	@$(COMPOSE_PROD) logs -f api

prod-clean: ## Stop prod containers and wipe volumes (careful!)
	@$(COMPOSE_PROD) down -v
	@echo "$(RED)Prod volumes removed$(RESET)"

# ==================================================================================
# HEALTH
# ==================================================================================

health: ## Check if the API is responding
	@curl -s http://localhost:$(PORT)/health | echo " <- API response"

# ==================================================================================
# GIT
# ==================================================================================

commit: ## Commit all changes. Usage: make commit msg='your message'
	@if [ -z "$(msg)" ]; then \
		echo "$(RED)Error: provide msg='your message'$(RESET)"; \
		exit 1; \
	fi
	git add .
	git commit -m "$(msg)"
	@echo "$(GREEN)Committed$(RESET)"

commit-push: commit ## Commit and push
	git push
	@echo "$(GREEN)Pushed$(RESET)"

first-commit-push: commit ## Commit and push (first time, sets upstream)
	git push -u origin HEAD
	@echo "$(GREEN)Pushed$(RESET)"