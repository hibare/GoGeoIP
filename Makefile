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

.PHONY: init
init: ## Initialize the project
	$(MAKE) install-golangci-lint
	$(MAKE) install-pre-commit
	go mod download

.PHONY: install-golangci-lint
install-golangci-lint: ## Install golangci-lint
ifeq (, $(shell which golangci-lint))
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
endif

.PHONY: install-pre-commit
install-pre-commit: ## Install pre-commit
	pre-commit install

.PHONY: dev
dev: ## Start API (dev)
	go mod tidy
	${DOCKER_COMPOSE_PREFIX} up api

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

.PHONY: build
build: ## Build docker image
	docker build -t hibare/go-geo-ip .

.PHONY: help
help: ## Display this help
		echo -e "\n$(BBLUE)GoGeoIP: IP Geolocation Service$(NC)\n"
		@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BCYAN)%-30s$(NC)%s\n", $$1, $$2}'
