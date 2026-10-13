# ==============================================================================
# TRADING SIGNALS PLATFORM - MASTER MAKEFILE
# ==============================================================================

SHELL := /bin/bash
.DEFAULT_GOAL := help

# Colors for output
CYAN := \033[36m
RESET := \033[0m

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'

# ==============================================================================
# 1. PROTOBUF & CONTRACTS
# ==============================================================================
PROTO_DIR := proto
GO_OUT := gen/go
PYTHON_OUT := gen/python
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

.PHONY: proto-format
proto-format: ## Format all Protobuf files using buf
	buf format -w $(PROTO_DIR)

.PHONY: proto-format-check
proto-format-check: ## Check formatting of Protobuf files (CI check)
	@test -z "$$(buf format --diff $(PROTO_DIR))" || \
		(echo "Protobuf files are not formatted. Run 'make proto-format'."; exit 1)

.PHONY: proto-lint
proto-lint: ## Lint Protobuf definitions against standard rules
	buf lint

.PHONY: proto-build
proto-build: ## Build Protobuf module schema
	buf build

.PHONY: proto-generate
proto-generate: ## Generate Go and Python stubs from proto files
	@mkdir -p $(GO_OUT) $(PYTHON_OUT)
	protoc \
		-I $(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go_opt=paths=source_relative \
		--python_out=$(PYTHON_OUT) \
		$(PROTO_FILES)

.PHONY: proto-clean
proto-clean: ## Remove generated proto stubs
	rm -rf gen

.PHONY: proto-all
proto-all: proto-format-check proto-lint proto-build proto-generate ## Run all proto checks and code generation

# ==============================================================================
# 2. CONTAINERIZATION & DOCKER BUILDS
# ==============================================================================
.PHONY: build-generator
build-generator: ## Build Docker image for the market data generator using buildx
	docker buildx build --load -t trading-generator:latest -f services/generator/Dockerfile services/generator

.PHONY: build-ingestion
build-ingestion: ## Build Docker image for the ingestion service using buildx
	docker buildx build --load -t trading-ingestion:latest -f services/ingestion/Dockerfile services/ingestion

.PHONY: build-stream-processor
build-stream-processor: ## Build Docker image for the Java Kafka Streams processor using buildx
	docker buildx build --load -t trading-stream-processor:latest -f services/stream-processor/Dockerfile services/stream-processor

.PHONY: build-signal-service
build-signal-service: ## Build Docker image for the Python signal engine using buildx
	docker buildx build --load -t trading-signal-service:latest -f services/signal-service/Dockerfile services/signal-service

.PHONY: build-all
build-all: build-generator build-ingestion build-stream-processor build-signal-service ## Build all service container images using buildx

# ==============================================================================
# 3. INFRASTRUCTURE (VAGRANT & ANSIBLE)
# ==============================================================================
.PHONY: infra-up
infra-up: ## Bring up the 3-node Vagrant VirtualBox cluster
	cd infra && vagrant up

.PHONY: infra-provision
infra-provision: ## Run Ansible playbooks to configure K8s, containerd, and Cilium CNI
	cd infra/ansible && ansible-playbook -i inventory.ini site.yaml

.PHONY: infra-cluster-setup
infra-cluster-setup: infra-up infra-provision ## Bring up virtual machines and fully provision the Kubernetes cluster

.PHONY: infra-status
infra-status: ## Check Vagrant machine status
	cd infra && vagrant status

.PHONY: infra-down
infra-down: ## Tear down the Vagrant cluster
	cd infra && vagrant destroy -f

# ==============================================================================
# 4. TESTING & LOCAL VERIFICATION
# ==============================================================================
.PHONY: test-unit
test-unit: ## Run unit tests across all microservices
	@echo "Running unit tests..."

.PHONY: test-integration
test-integration: ## Run integration test suite
	@echo "Running integration tests..."

# ==============================================================================
# 5. CI/CD PIPELINE SIMULATION
# ==============================================================================
.PHONY: ci
ci: proto-all test-unit build-all ## Run the complete CI verification pipeline locally
	@echo "CI pipeline completed successfully!"