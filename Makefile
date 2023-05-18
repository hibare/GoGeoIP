SHELL=/bin/bash

UI := $(shell id -u)
GID := $(shell id -g)
MAKEFLAGS += -s
GODOTENV_CMD_PATH = $$GOPATH/bin/godotenv
DOCKER_COMPOSE_PREFIX = HOST_UID=${UID} HOST_GID=${GID} docker-compose -f docker-compose.dev.yml

all: api-up

api-up:
	go mod tidy
	${DOCKER_COMPOSE_PREFIX} up api

api-down:
	${DOCKER_COMPOSE_PREFIX} rm -fsv api

test:
	if [ ! -f ${GODOTENV_CMD_PATH} ]; then \
		echo "Missing godotenv cmd"; \
		echo "Installing..."; \
		go install github.com/joho/godotenv/cmd/godotenv@latest; \
		echo "Installed"; \
	fi
	
	${GODOTENV_CMD_PATH} -f .env go test ./... -cover

clean: 
	${DOCKER_COMPOSE_PREFIX} down
	go mod tidy

prod-up:
	docker build -t hibare/go-geo-ip . 
	docker compose up

docker-build:
	docker build -t hibare/go-geo-ip . 

.PHONY = all clean api-up api-down test prod-up docker-build
