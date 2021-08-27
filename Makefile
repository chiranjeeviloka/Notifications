.PHONY: help build build-race test run-dev run-stg 

# Setup these variables to be used by commands later.
# --------------------------------------------------------------------------------------------------
SERVICE_NAME         := my-service

.EXPORT_ALL_VARIABLES:
GOPRIVATE              = github.com/FreedomCentral/*
CGO_ENABLED           ?= 0
CGO_CFLAGS             = -g -O2 -Wno-return-local-addr
# BR_VAULT.. values below should be replaced once Vault is setup.
BR_VAULT_ADDR         ?= http://127.0.0.1:8200
BR_VAULT_TOKEN        ?= 
BR_ENV                ?= dev 

# Sets a default for make
.DEFAULT_GOAL := help

help:; ## Output help
	@printf "%s\\n" \
		"The following targets are available:" \
		""
	@awk 'BEGIN {FS = ":.*?## "} /^[\/.%a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m - %s\n", $$1, $$2}' ${MAKEFILE_LIST}

	@printf "%s\\n" "" "" \
		"Examples:" \
		"make build" \
		"    build the service"

build: ## Build binary
	go generate ./...
	go build 

build-race: ## Build with race detector turned on.
	go generate ./...
	go build -race 

test: ## Test all packages
	go test -v ./...

run-dev: build ## Run local instance of the service pointing to dev region
	BR_ENV=dev ./notification-service

run-stg: build ## Run local instance of the service pointing to stg region
	BR_ENV=stg ./notification-service
