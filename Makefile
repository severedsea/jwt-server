SHELL := /bin/bash
export TERM=xterm-256color

.PHONY: test build run setup

# APP_NAME is used as a naming convention for resources to the local environment
APP_NAME := jwt-server
# APP_PATH is the project/app directory in the container
APP_PATH := /${APP_NAME}

NO_COLOR=\x1b[0m
OK_COLOR=\x1b[32;1m
ERROR_COLOR=\x1b[31;1m

# ----------------------------
# ----------------------------
# Commands
# ----------------------------
# ----------------------------
COMPOSE := docker-compose -f docker-compose.yml
GO_COMPOSE = $(COMPOSE) run -T --rm -w $(APP_PATH) -v $(shell pwd):$(APP_PATH) $(COMPOSE_ARGS) go
RUN_COMPOSE = $(COMPOSE) run -T --rm --service-ports -w $(APP_PATH) -v $(shell pwd):$(APP_PATH) go

# ----------------------------
# ----------------------------
# Targets
# ----------------------------
# ----------------------------
all: setup test

# ----------------------------
# JWT keys
# ----------------------------
jwt-keys: jwt.rsa jwt.rsa.pub
	@printf "\n$(OK_COLOR)Generated JWT key-pair$(NO_COLOR)\n"

jwt.rsa:
	openssl genrsa -out jwt.rsa 2048

jwt.rsa.pub:
	openssl rsa -in jwt.rsa -pubout > jwt.rsa.pub

# ----------------------------
# Setup and Teardown
# ----------------------------
# setup creates/initializes development environment dependencies for run task/s
setup: jwt-keys go redis

# go builds the go service defined in the compose file
go:
	@printf "\n$(OK_COLOR)Building Go image$(NO_COLOR)\n"
	@$(COMPOSE) up -d go

# redis runs the redis service defined in the compose file
redis:
	@printf "\n$(OK_COLOR)Spinning up Redis$(NO_COLOR)\n"
	$(COMPOSE) up -d redis

# teardown stops and removes all containers and resources associated to docker-compose.yml
teardown:
	@printf "\n$(OK_COLOR)Destroying everything$(NO_COLOR)\n"
	$(COMPOSE) down --remove-orphans

# ----------------------------
# Run
# ----------------------------
run: setup
	@printf "\n$(OK_COLOR)Running serverd$(NO_COLOR)\n"
	$(RUN_COMPOSE) go run cmd/serverd/main.go

# ----------------------------
# Test
# ----------------------------
# test executes project tests in a golang container
test:
	@if $(GO_COMPOSE) env $(shell cat .env.test | egrep -v '^#|^REDIS_' | xargs -0) \
	make go-test; \
	then printf "\n$(OK_COLOR)\n$(OK_COLOR)[Test okay -- `date`]$(NO_COLOR)\n"; \
	else printf "\n$(OK_COLOR)\n$(ERROR_COLOR)[Test FAILED -- `date`]$(NO_COLOR\n)\n"; exit 1; fi

# go-test executes test for all packages
go-test:
	@printf "\n$(OK_COLOR)Running tests$(NO_COLOR)\n"
	@go test -coverprofile=c.out -failfast -shuffle=on -timeout 5m $(shell go list ./... | egrep -v 'generated')
	@printf "\n$(OK_COLOR)Coverage report$(NO_COLOR)\n"
	@go tool cover -func=c.out
	@go tool cover -html=c.out -o=coverage.html

# ----------------------------
# Build
# ----------------------------
# build executes build in a golang container
build:
	$(GO_COMPOSE) make go-build

# go-build executes the go build process
go-build:
	go build -o cmd/serverd/bin/serverd -v ./cmd/serverd