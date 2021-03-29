#
# Makefile
# Created by Masatoshi Fukunaga on 19/07/17
#
LINT_OPT=--issues-exit-code=0 \
		--enable-all \
		--tests=false \
		--disable=wsl \
		--disable=nlreturn

.EXPORT_ALL_VARIABLES:

.PHONY: all test lint coverage clean

all: lint coverage

lint:
	golangci-lint run $(LINT_OPT) ./...

test:
	go test -timeout 1m -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html coverage.out -o coverage.out.html

coverage: test
	go tool cover -func=coverage.out

clean:
	go clean
