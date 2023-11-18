GOPATH:=$(shell go env GOPATH)

.PHONY: lint
lint: ## lint
	staticcheck ./...
	golangci-lint run ./...

.PHONY: format
format: ## format
	golines -m 120 -w --ignore-generated .
	gofumpt -l -w .

.PHONY: tidy
tidy: ## add missing and remove unused modules
	@go mod tidy

.PHONY: deps
deps: ## download modules to local cache
	@go mod download

.PHONY: dev-customer
dev-customer: ## run dev customer service
	@./scripts/run.sh customer

.PHONY: docker-up
docker-up: ## docker-compose up and detach
	@docker-compose -f ./deployments/docker-compose/docker-compose.dev.yml up -d

.PHONY: docker-down
docker-down: ## docker-compose down
	@docker-compose -f ./deployments/docker-compose/docker-compose.dev.yml down

.PHONY: docker-down-volumes
docker-down-volumes: ## docker-compose down and delele volumes
	@docker-compose -f ./deployments/docker-compose/docker-compose.dev.yml down --volumes

.PHONY: help
help: ## print help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_0-9-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
