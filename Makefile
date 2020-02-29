SHELL = /bin/bash -euo pipefail

GO ?= go

# Variables for container targets
GO_VERSION ?= latest
CONTAINER = golang:$(GO_VERSION)
PKG = github.com/lorenzosaino/go-sysctl
DOCKER_RUN_FLAGS = --rm -it -v $$(pwd):/go/src/$(PKG) -w /go/src/$(PKG)

export GO111MODULE=on

.PHONY: all mod-upgrade mod-update fmt-check lint vet staticcheck test container-shell container-test

all: fmt-check lint vet staticcheck test

mod-upgrade:
	$(GO) get -u -t ./...
	$(GO) mod tidy
	$(GO) mod vendor

mod-update:
	$(GO) mod tidy
	$(GO) mod vendor

# Ensure that all source files pass "go fmt"
fmt-check:
	exit $(shell $(GO) fmt ./... | wc -l)

lint:
	[ -x "$(shell which golint)" ] || $(GO) install ./vendor/golang.org/x/lint/golint 2>/dev/null || $(GO) get -u golang.org/x/lint/golint
	# We need to explicitly exclude ./vendor because of https://github.com/golang/lint/issues/320
	golint -set_exit_status $(shell $(GO)  list ./... | grep -v '/vendor/')

vet:
	$(GO) vet ./...

staticcheck:
	[ -x "$(shell which staticcheck)" ] || $(GO) install ./vendor/honnef.co/go/tools/cmd/staticcheck 2>/dev/null || $(GO) get -u honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...

test:
	$(GO) test -v ./...

container-shell:
	docker run $(DOCKER_RUN_FLAGS) $(CONTAINER) /bin/bash

container-test:
	docker run $(DOCKER_RUN_FLAGS) $(CONTAINER) make all
