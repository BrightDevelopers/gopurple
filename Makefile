# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Build directory
BUILDDIR=bin

# Discover all example directories (each contains a main package)
EXAMPLES=$(shell ls -d examples/*/ 2>/dev/null | grep -v internal | sed 's|examples/||;s|/||')

.DEFAULT_GOAL := help

all: test build

build:
	@mkdir -p $(BUILDDIR)
	@for example in $(EXAMPLES); do \
		echo "building bin/$$example"; \
		$(GOBUILD) -o $(BUILDDIR)/$$example ./examples/$$example; \
	done

test:
	$(GOTEST) -v ./...

clean:
	rm -rf $(BUILDDIR)
	find . -name '*~' -delete

deps:
	$(GOMOD) download
	$(GOMOD) tidy

build-linux:
	@mkdir -p $(BUILDDIR)
	@for example in $(EXAMPLES); do \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDDIR)/$${example}_linux -v ./examples/$$example; \
	done

install:
	@for example in $(EXAMPLES); do \
		$(GOCMD) install ./examples/$$example; \
	done

install-hooks:
	git config core.hooksPath .githooks
	@chmod +x .githooks/* 2>/dev/null || true
	@echo "Git hooks installed (.githooks). Binaries and secrets will be blocked at commit time."

dev-deps: install-hooks
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run

fmt:
	goimports -w .
	$(GOCMD) fmt ./...

help:
	@echo "Available targets:"
	@echo "  make all         - Run tests and build all examples"
	@echo "  make build       - Build all example programs into bin/"
	@echo "  make test        - Run all tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make deps        - Download and tidy dependencies"
	@echo "  make build-linux - Cross-compile all examples for Linux amd64"
	@echo "  make install     - Install all examples to GOPATH/bin"
	@echo "  make install-hooks - Install git pre-commit guard (blocks binaries/secrets)"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make help           - Show this help message"
	@echo "  make list-examples  - List all example programs"

list-examples:
	@for example in $(EXAMPLES); do echo "bin/$$example"; done

.PHONY: all build test clean deps build-linux install install-hooks dev-deps lint fmt help list-examples
