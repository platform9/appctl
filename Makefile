# Copyright 2021 The pf9ctl authors.
#
# Usage:
# make                 # builds the artifacts
# make clean           # removes the artifact and the vendored packages

SHELL := /usr/bin/env bash
GITHASH := $(shell git rev-parse --short HEAD)
BIN_DIR := $(shell pwd)/bin
BIN := appctl
REPO := appctl

SEGMENT_KEY := -X github.com/platform9/appctl/pkg/segment.APPCTL_SEGMENT_WRITE_KEY=$(APPCTL_SEGMENT_WRITE_KEY)
APPURL := -X github.com/platform9/appctl/pkg/constants.APPURL=$(APPURL)
DOMAIN := -X github.com/platform9/appctl/pkg/constants.DOMAIN=$(DOMAIN)
CLIENTID := -X github.com/platform9/appctl/pkg/constants.CLIENTID=$(CLIENTID)
GRANT_TYPE := -X github.com/platform9/appctl/pkg/constants.GrantType=$(GRANT_TYPE)

LD_FLAGS := $(SEGMENT_KEY) $(APPURL) $(DOMAIN) $(CLIENTID) $(GRANT_TYPE)

.PHONY: clean format test build-all build-linux64 build-win64 build-mac

build-all: build-linux64 build-win64 build-mac

format:
	gofmt -w -s *.go
	gofmt -w -s */*.go

clean:
	rm -rf $(BIN_DIR)

build-mac: $(BIN_DIR)/$(BIN)-mac
$(BIN_DIR)/$(BIN)-mac: test
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o $(BIN_DIR)/$(BIN)-mac -ldflags '$(LD_FLAGS)' main.go

build-win64: $(BIN_DIR)/$(BIN)-win64
$(BIN_DIR)/$(BIN)-win64: test
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o $(BIN_DIR)/$(BIN)-win64 -ldflags '$(LD_FLAGS)' main.go

build-linux64: $(BIN_DIR)/$(BIN)-linux64
$(BIN_DIR)/$(BIN)-linux64: test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(BIN_DIR)/$(BIN)-linux64 -ldflags '$(LD_FLAGS)' main.go

test:
	go test -v ./...
