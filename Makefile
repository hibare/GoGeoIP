SHELL=/bin/bash

UI := $(shell id -u)
GID := $(shell id -g)
MAKEFLAGS += -s
GODOTENV_CMD_PATH = $$GOPATH/bin/godotenv
DOCKER_COMPOSE_PREFIX = HOST_UID=${UID} HOST_GID=${GID} docker-compose -f docker-compose.dev.yml

# Bold
BCYAN=\033[1;36m
BBLUE=\033[1;34m

# No color (Reset)
NC=\033[0m

.DEFAULT_GOAL := help

.PHONY: api-up
api-up: ## Start API (dev)
	go mod tidy
	${DOCKER_COMPOSE_PREFIX} up api

.PHONY: api-down
api-down: ## Stop API (dev)
	${DOCKER_COMPOSE_PREFIX} rm -fsv api

.PHONY: test
test: ## Run tests
	go test ./... -cover

.PHONY: clean
clean: ## Cleanup
	${DOCKER_COMPOSE_PREFIX} down
	go mod tidy

.PHONY: prod-up
prod-up: ## start prod API
	$(MAKE) docker-build
	docker compose up

.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t hibare/go-geo-ip . 

.PHONY: help
help: ## Disply this help
		echo -e "\n$(BBLUE)GoGeoIP: IP Geolocation Service$(NC)\n"
		@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BCYAN)%-18s$(NC)%s\n", $$1, $$2}'