.PHONY: all

PKG := "github.com/kochabonline/kcloud"

install: ## Install dependencies
	@go install github.com/google/wire/cmd/wire@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

wire: ## Generate wire code
	@wire ./...

swag: ## Generate swagger docs
	@swag init

build: ## Build the project
	@GOFLAGS=-buildvcs=false go build -trimpath .

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help