ifndef GOPATH
	GOPATH := $(shell go env GOPATH)
endif
ifndef GOBIN # derive value from gopath (default to first entry, similar to 'go get')
	GOBIN := $(shell go env GOPATH | sed 's/:.*//')/bin
endif

tools = $(addprefix $(GOBIN)/, golangci-lint gofumpt golines govulncheck protoc-gen-go protoc-gen-go-grpc air)
deps = $(addprefix $(GOBIN)/, goose)


###############################################################################
#
# Initialization
#
###############################################################################
.PHONY: deps
deps: $(deps) ## download modules to local cache
	@echo "Installing dependencies"
	@go mod download

tools: $(tools) ## install tools required for the build
	@echo "Installed tools"

.PHONY: tidy
tidy: ## add missing and remove unused modules
	@go mod tidy

###############################################################################
#
# Build and testing rules
#
###############################################################################
.PHONY: dev-customer
dev-customer: ## run dev customer service
	@./scripts/run.sh customer

.PHONY: test
test: ## run the go tests
	@echo "Running tests"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out


###############################################################################
#
# Code formatting and linting
#
###############################################################################
.PHONY: lint
lint: tools ## lint
	@echo "Linting"
	golangci-lint run ./...

.PHONY: format
format: tools ## format go code
	@echo "Formating ..."
	golines -m 120 -w --ignore-generated .
	gofumpt -l -w .
	@echo "Formatting complete"

.PHONY: sec
sec: ## detect vulnerability go packages
	@echo "Vulnerability detection $(1)"
	@govulncheck ./...

###############################################################################
# Code Generation
#
# Some code generation can be slow, so we only run it if
# the source file has changed.
###############################################################################
.PHONY: proto
proto: ## compile proto file
	@echo Generating intergrate message  proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/am/proto/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/msq/proto/*.proto

###############################################################################
#
# Infrastructure
#
###############################################################################
.PHONY: docker-up
docker-up: ## docker-compose up and detach
	@docker-compose -f ./deployments/docker-compose/docker-compose.dev.yml up -d

.PHONY: docker-down
docker-down: ## docker-compose down
	@docker-compose -f ./deployments/docker-compose/docker-compose.dev.yml down

.PHONY: docker-down-volumes
docker-down-volumes: ## docker-compose down and delele volumes
	@docker-compose -f ./deployments/docker-compose/docker-compose.dev.yml down --volumes

###############################################################################
# Install Tools and deps
#
# These targets specify the full path to where the tool is installed
# If the tool already exists it wont be re-installed.
###############################################################################
update-tools: delete-tools $(tools) ## update the tools by deleting and re-installing

delete-tools: ## delete the tools
	@rm $(tools) || true

# Install golangci-lint
$(GOBIN)/golangci-lint:
	@echo "🔘 Installing golangci-lint... (`date '+%H:%M:%S'`)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN)

# Install gofumpt to format code
$(GOBIN)/gofumpt:
	@echo "🔘 Installing gofumpt ... (`date '+%H:%M:%S'`)"
	@go install mvdan.cc/gofumpt@latest

$(GOBIN)/golines:
	go install github.com/segmentio/golines@latest

$(GOBIN)/govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest

# Install goose to perform db migrations
$(GOBIN)/goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

$(GOBIN)/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31

$(GOBIN)/protoc-gen-go-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# Install air to live reloading
$(GOBIN)/air:
	go install github.com/cosmtrek/air@latest

.PHONY: help
help: ## print help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
