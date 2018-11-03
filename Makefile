.PHONY: all fmt-check lint vet staticcheck test container-debug container-all

GO_VERSION ?= latest
CONTAINER = golang:$(GO_VERSION)
LOCAL ?= $(PWD)
MOUNT ?= /go/src/github.com/lorenzosaino/go-sysctl
PKGS_NOVENDOR = $(shell go list ./... | grep -v '/vendor/')

# Ensure that all source files pass "go fmt"
fmt-check:
	exit `go fmt $(PKGS_NOVENDOR) | wc -l`

lint: $(GOPATH)/bin/golint
	golint -set_exit_status ./...

$(GOPATH)/bin/golint:
	go get golang.org/x/lint/golint

vet:
	go vet $(PKGS_NOVENDOR)

staticcheck: $(GOPATH)/bin/staticcheck
	staticcheck $(PKGS_NOVENDOR)

$(GOPATH)/bin/staticcheck:
	go get honnef.co/go/tools/cmd/staticcheck

test:
	go test $(PKGS_NOVENDOR)

container-shell:
	docker run -it -v $(LOCAL):$(MOUNT) -w $(MOUNT) $(CONTAINER) /bin/bash

container-all:
	docker run -it -v $(LOCAL):$(MOUNT) -w $(MOUNT) $(CONTAINER) make all

all: fmt-check lint vet staticcheck test
