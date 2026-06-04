BINARY_NAME := mentat
BUILD_DIR := build

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: build clean deps

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BUILD_DIR)

deps:
	go get -u ./...
	go mod tidy
