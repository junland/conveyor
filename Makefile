PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
PKG_NAME := "conveyor"
GIT_COMMIT:=$(shell git rev-parse --verify HEAD --short=7)
GO_VERSION:=$(shell go version | grep -o "go1\.[0-9|\.]*")
VERSION ?= 0.0.0
BIN_NAME := conveyor
APP_NAME ?= conveyor
CGO_ENABLED := 0 
APP_NAME_UPPER := `echo $(APP_NAME) | tr '[:lower:]' '[:upper:]'`

.PHONY: binary
binary: clean fmt
	@echo "Building binary for commit $(GIT_COMMIT)"
	go build 

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf ./,*
	@rm -rf workspace*
	@rm -rf worker*
	@rm -f ./conveyor
	@rm -rf ./conveyor-*
	@rm -rf ./*.tar.gz
	@rm -rf ./conveyor_*
	@rm -rf ./*.txt
	@rm -rf ./*.pem
	@rm -rf ./jobs.d
	@echo "Done cleaning..."

.PHONY: fmt
fmt:
	@echo "Running $@"
	go fmt ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test ./...