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
LDFLAGS := "" 

.PHONY: clean format test build-all build-linux64 build-win64 build-mac

build-all: build-linux64 build-win64 build-mac

format:
	gofmt -w -s *.go
	gofmt -w -s */*.go

clean:
	rm -rf $(BIN_DIR)

build-mac: $(BIN_DIR)/$(BIN)-mac
$(BIN_DIR)/$(BIN)-mac: test
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o $(BIN_DIR)/$(BIN)-mac -ldflags $(LDFLAGS) main.go

build-win64: $(BIN_DIR)/$(BIN)-win64
$(BIN_DIR)/$(BIN)-win64: test
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o $(BIN_DIR)/$(BIN)-win64 -ldflags $(LDFLAGS) main.go

build-linux64: $(BIN_DIR)/$(BIN)-linux64
$(BIN_DIR)/$(BIN)-linux64: test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(BIN_DIR)/$(BIN)-linux64 -ldflags $(LDFLAGS) main.go

test:
	go test -v ./...
