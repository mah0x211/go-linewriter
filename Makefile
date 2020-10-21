#
# Makefile
# Created by Masatoshi Fukunaga on 19/07/17
#
BUILD_DIR := $(PWD)/build
DEPS_DIR := $(BUILD_DIR)/deps
GOCMD:=GOPATH=$(DEPS_DIR) go
COVER_PATH := $(BUILD_DIR)/cover.out
GOTEST := $(GOCMD) test -timeout 1m
GOTOOL := $(GOCMD) tool
PKGS=$(addprefix ./,$(filter-out _%/ build/,$(sort $(dir $(wildcard */*)))))
LINT_OPT=--issues-exit-code=0 \
		--enable-all \
		--tests=false \
		--disable=wsl \
		--disable=nlreturn

.EXPORT_ALL_VARIABLES:

.PHONY: all test lint coverage

all: lint coverage

lint:
	golangci-lint run $(LINT_OPT) . $(PKGS)

test:
	$(GOTEST) -coverprofile=$(COVER_PATH) -covermode=atomic . $(PKGS)
	$(GOTOOL) cover -html $(COVER_PATH) -o $(COVER_PATH).html

coverage: test
	$(GOTOOL) cover -func=$(COVER_PATH)
