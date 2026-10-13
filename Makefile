SHELL := /bin/bash

# ==============================================================================
# CONFIGURATION & BINARIES
# ==============================================================================
GO ?= go
GOFMT ?= gofmt
BUF ?= buf
PROTOC ?= protoc
RUFF ?= ruff
DOCKER ?= docker
GOVULNCHECK ?= govulncheck

# Paths & Services
PROTO_DIR := proto
GO_OUT := gen/go
PYTHON_OUT := gen/python
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')
SERVICES := generator ingestion stream-processor signal-service

.PHONY: all
all: check

.PHONY: help
help: ## Show available Make targets
	@echo "Trading Signals Platform development commands:"
	@echo
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make <target>\n\nTargets:\n"} \
		/^[a-zA-Z0-9_-]+:.*##/ {printf "  %-28s %s\n", $$1, $$2}' \
		$(MAKEFILE_LIST)

# --------------------------------------------------------------------
# Protobuf & Contracts
# --------------------------------------------------------------------
.PHONY: proto-format proto-format-check proto-lint proto-build proto-generate proto-check
proto-format: ## Format all Protobuf files using buf
	$(BUF) format -w $(PROTO_DIR)

proto-format-check: ## Check formatting of Protobuf files
	@test -z "$$($(BUF) format --diff $(PROTO_DIR))" || \
		(echo "Protobuf files are not formatted. Run 'make proto-format'." && exit 1)

proto-lint: ## Lint Protobuf definitions against standard rules
	$(BUF) lint $(PROTO_DIR)

proto-build: ## Build Protobuf module schema
	$(BUF) build $(PROTO_DIR)

proto-generate: ## Generate Go and Python stubs from proto files
	@mkdir -p $(GO_OUT) $(PYTHON_OUT)
	$(PROTOC) \
		-I $(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go_opt=paths=source_relative \
		--python_out=$(PYTHON_OUT) \
		$(PROTO_FILES)

proto-check: proto-format-check proto-lint proto-build proto-generate ## Run all proto validation and generation

# --------------------------------------------------------------------
# Code Quality & Validation
# --------------------------------------------------------------------
.PHONY: go-fmt-check go-vet go-check python-fmt-check python-lint python-check
go-fmt-check: ## Check Go source file formatting
	@test -z "$$($(GOFMT) -l services/ingestion)" || (echo "Unformatted Go files:" && $(GOFMT) -l services/ingestion && exit 1)

go-vet: ## Run go vet on ingestion service
	cd services/ingestion && $(GO) vet ./...

python-fmt-check: ## Check Python formatting using ruff
	$(RUFF) format --check services/generator services/signal-service

python-lint: ## Run ruff linter on Python services
	$(RUFF) check services/generator services/signal-service

go-check: proto-check go-fmt-check go-vet ## Run Go & Proto quality checks
python-check: python-fmt-check python-lint ## Run Python quality checks

# --------------------------------------------------------------------
# Unit Tests
# --------------------------------------------------------------------
.PHONY: unit-test go-unit-test python-unit-test
go-unit-test: ## Run Go unit tests with race detector
	cd services/ingestion && $(GO) test -race -count=1 -coverprofile=coverage.out ./...

python-unit-test: ## Run Python unit tests with pytest
	@for d in services/generator services/signal-service; do \
		if [ -d "$$d" ]; then (cd $$d && pytest); fi; \
	done

unit-test: proto-generate go-unit-test python-unit-test ## Run all service unit tests

# --------------------------------------------------------------------
# Security & Vulnerability Scans
# --------------------------------------------------------------------
.PHONY: vuln-check
vuln-check: ## Scan Go dependencies for known vulnerabilities
	cd services/ingestion && $(GOVULNCHECK) ./...

# --------------------------------------------------------------------
# Universal Container Builds (Docker Buildx)
# --------------------------------------------------------------------
.PHONY: build-%
build-%: proto-generate ## Universally build any service container image using buildx
	$(DOCKER) buildx build --load -t trading-$*:latest -f services/$*/Dockerfile services/$*

.PHONY: build-all
build-all: $(addprefix build-,$(SERVICES)) ## Build all service container images universally

# --------------------------------------------------------------------
# Pipeline Aliases
# --------------------------------------------------------------------
.PHONY: check ci
check: proto-check go-check python-check ## Run static checks
ci: check unit-test vuln-check build-all ## Run complete verification pipeline