.PHONY: all fmt-check lint vet staticcheck test container-debug container-all

GO_VERSION ?= latest
CONTAINER = golang:$(GO_VERSION)
LOCAL ?= $(PWD)
MOUNT ?= /go/src/github.com/lorenzosaino/go-sysctl

# Ensure that all source files pass "go fmt"
fmt-check:
	exit `go fmt ./... | wc -l`

lint: $(GOPATH)/bin/golint
	golint -set_exit_status ./...

$(GOPATH)/bin/golint:
	go get github.com/golang/lint/golint

vet:
	go vet ./...

staticcheck: $(GOPATH)/bin/staticcheck
	staticcheck ./...

$(GOPATH)/bin/staticcheck:
	go get honnef.co/go/tools/cmd/staticcheck

test:
	go test ./...

container-debug:
	docker run -it -v $(LOCAL):$(MOUNT) -w $(MOUNT) $(CONTAINER) /bin/bash

container-all:
	docker run -it -v $(LOCAL):$(MOUNT) -w $(MOUNT) $(CONTAINER) make all

all: fmt-check lint vet staticcheck test
