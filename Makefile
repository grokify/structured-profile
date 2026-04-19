.PHONY: all build test lint clean tidy fmt vet cli migrate install

GO ?= go
GOLANGCI_LINT ?= golangci-lint

# Build targets
all: fmt lint test build

build:
	$(GO) build -v ./...

# Test targets
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

test-short:
	$(GO) test -v -short ./...

coverage: test
	$(GO) tool cover -html=coverage.out -o coverage.html

# Lint targets
lint:
	$(GOLANGCI_LINT) run ./...

lint-fix:
	$(GOLANGCI_LINT) run --fix ./...

vet:
	$(GO) vet ./...

# Format targets
fmt:
	$(GO) fmt ./...
	gofmt -s -w .

# Dependency targets
tidy:
	$(GO) mod tidy

deps:
	$(GO) mod download

# Clean targets
clean:
	$(GO) clean
	rm -f coverage.out coverage.html

# CLI build
cli:
	$(GO) build -o bin/sprofile ./cmd/sprofile

migrate:
	$(GO) build -o bin/migrate ./cmd/migrate

# Install
install:
	$(GO) install ./cmd/sprofile
