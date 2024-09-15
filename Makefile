# Default build target
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

.PHONY: help test lint build

help: ## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST) | column -tl 2

test: ## Test go files and report coverage.
	go test -v -race -cover ./...

lint: ## List all the linting issues.
	golangci-lint run
