.PHONY: all fmt-check lint vet staticcheck test container-shell container-test

GO ?= go

# Variables for container targets
GO_VERSION ?= latest
CONTAINER = golang:$(GO_VERSION)
PKG = github.com/lorenzosaino/go-sysctl
DOCKER_RUN_FLAGS = --rm -v $$(pwd):/go/src/$(PKG) -w /go/src/$(PKG)

PKGS = $(shell go list ./... | grep -v '/vendor/')

all: fmt-check lint vet staticcheck test

# Ensure that all source files pass "go fmt"
fmt-check:
	exit $(shell $(GO) fmt ./... | wc -l)

vet:
	$(GO) vet ./...

staticcheck:
	[ -x "$(shell which staticcheck)" ] || $(GO) install ./vendor/honnef.co/go/tools/cmd/staticcheck 2>/dev/null || $(GO) get -u honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...

lint:
	[ -x "$(shell which golint)" ] || $(GO) install ./vendor/golang.org/x/lint/golint 2>/dev/null || $(GO) get -u golang.org/x/lint/golint
	# We need to explicitly exclude ./vendor becuase of https://github.com/golang/lint/issues/320
	golint -set_exit_status $(shell $(GO)  list ./... | grep -v '/vendor/')

test:
	$(GO) test -v ./...

container-shell:
	docker run $(DOCKER_RUN_FLAGS) $(CONTAINER) /bin/bash

container-test:
	docker run $(DOCKER_RUN_FLAGS) $(CONTAINER) make all
