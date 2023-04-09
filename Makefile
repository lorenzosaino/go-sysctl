SHELL = /bin/bash -euo pipefail

# Go binary to use in non-container targets
GO ?= go

# Variables for container targets
GO_VERSION ?= latest
CONTAINER ?= golang:$(GO_VERSION)
PKG = github.com/lorenzosaino/go-sysctl
DOCKER_RUN_FLAGS = --rm -it -v $$(pwd):/go/src/$(PKG) -w /go/src/$(PKG)

export GO111MODULE=on

all: fmt-check vet nilness staticcheck test ## Run all checks and tests

.PHONY: mod-upgrade
mod-upgrade: ## Upgrade all vendored dependencies
	$(GO) get -d -u -t ./...
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: mod-update
mod-update: ## Ensure all used dependencies are tracked in go.{mod|sum} and vendored
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: fmt-check
fmt-check: ## Validate that all source files pass "go fmt"
	exit $(shell $(GO) fmt ./... | wc -l)

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: staticcheck
staticcheck: ## Run staticcheck
	$(GO) install ./vendor/honnef.co/go/tools/cmd/staticcheck 2>/dev/null || @[ -x "$(shell which staticcheck)" ] || $(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

.PHONY: nilness
nilness: ## Run nilness
	$(GO) install ./vendor/golang.org/x/tools/go/analysis/passes/nilness/cmd/nilness 2>/dev/null || [ -x "$(shell which nilness)" ] || $(GO) install golang.org/x/tools/go/analysis/passes/nilness/cmd/nilness@latest
	nilness ./...

.PHONY: govulncheck
govulncheck: ## Run govulncheck
	$(GO) install ./vendor/golang.org/x/vuln/cmd/govulncheck 2>/dev/null || [ -x "$(shell which govulncheck)" ] || $(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: test
test: ## Run all tests
	$(GO) test -race ./...

.PHONY: container-shell
container-shell: ## Open a shell on a Docker container
	docker run $(DOCKER_RUN_FLAGS) $(CONTAINER) /bin/bash

.PHONY: container-%
container-%: ## Run any target of this Makefile in a Docker container
	docker run $(DOCKER_RUN_FLAGS) $(CONTAINER) make $*

.PHONY: help
help: ## Print help
	@(grep -E '^[a-zA-Z0-9_%-]+:.*?## .*$$' Makefile || true )| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
