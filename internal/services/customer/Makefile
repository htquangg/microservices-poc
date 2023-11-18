GOPATH:=$(shell go env GOPATH)

.PHONY: dev
dev: ## run dev service
	@go run main.go

.PHONY: build
build: ## build service
	@go build main.go

.PHONY: air
air: ## run dev service and reloading
	@air -c .air.toml

.PHONY: lint
lint: ## lint
	staticcheck ./...
	# golangci-lint run ./...

.PHONY: format
format: ## format
	golines -m 120 -w --ignore-generated .
	gofumpt -l -w .

.PHONY: tidy
tidy: ## add missing and remove unused modules
	@go mod tidy

.PHONY: proto
proto: ## compile proto file
	@echo Generating proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/*.proto


# ==============================================================================
# Go migrate mysql https://github.com/pressly/goose

DB_NAME = dev-local-customer-001
DB_HOST = localhost
DB_PORT = 3306
DB_USER = root
DB_PASS = 1
DB_PARSE_TIME = true

.PHONY: migrate-create
migrate-create: ## create new migration file
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c create ${name}

.PHONY: migrate-up
migrate-up: ## migrate the DB to the most recent version available
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c up -o "${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=${DB_PARSE_TIME}"

.PHONY: migrate-down
migrate-down: ## roll back the version by 1
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c down -o "${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=${DB_PARSE_TIME}"

.PHONY: migrate-status
migrate-status: ## check the migration status for the current DB
	@./scripts/goose-migrate.sh -p ./migrations/schemas -c status -o "${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=${DB_PARSE_TIME}"

.PHONY: help
help: ## print help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_0-9-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help
