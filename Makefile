# Ensure Make is run with bash shell as some syntax below is bash-specific
SHELL:=/usr/bin/env bash

.DEFAULT_GOAL := help

VERSION := $(shell git rev-parse --short HEAD)

# Define Docker related variables. Releases should modify and double check these vars.
REGISTRY := uhub.service.ucloud.cn/infra
IMAGE := statefulset-pingcap-controller-manager
CONTROLLER_IMG := $(REGISTRY)/$(IMAGE)

# Use GOPROXY environment variable if set
GOPROXY := $(shell go env GOPROXY)
ifeq ($(GOPROXY),)
GOPROXY := https://goproxy.cn
endif
export GOPROXY
# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

.PHONY: server

server:
	@echo "version: $(VERSION)"
	docker build --no-cache -t $(CONTROLLER_IMG):$(VERSION) -f Dockerfile .
	docker push $(CONTROLLER_IMG):$(VERSION)
