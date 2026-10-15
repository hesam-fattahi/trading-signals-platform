# =============================================================================
# Configuration
# =============================================================================
REGISTRY  := 192.168.56.241:5000
TAG       := dev

GENERATOR_DIR := services/generator
INGESTION_DIR := services/ingestion

GENERATOR_IMAGE := $(REGISTRY)/trading-generator:$(TAG)
INGESTION_IMAGE := $(REGISTRY)/trading-ingestion:$(TAG)

PLATFORM_NAMESPACE := platform
PLATFORM_RELEASE   := platform
PLATFORM_CHART     := charts/platform

APP_NAMESPACE := default
APP_RELEASE   := app
APP_CHART     := charts/app

# =============================================================================
# Help
# =============================================================================
.PHONY: help
help: ## Show this help
	@echo "Trading signals platform — available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# =============================================================================
# Protobuf (buf)
# =============================================================================
.PHONY: proto-generate
proto-generate: ## Generate protobuf stubs (go) from .proto files
	buf generate

.PHONY: proto-lint
proto-lint: ## Lint .proto files
	buf lint

.PHONY: proto-format
proto-format: ## Format .proto files in place
	buf format -w

.PHONY: proto-format-check
proto-format-check: ## Check .proto formatting (fails if not formatted)
	@buf format -d --exit-code

.PHONY: proto-breaking
proto-breaking: ## Check .proto files for breaking changes against main
	buf breaking --against '.git#branch=main'

.PHONY: proto-check
proto-check: proto-lint proto-format-check ## Run all proto checks

.PHONY: proto-clean
proto-clean: ## Remove generated proto output
	rm -rf gen/go/trading

# =============================================================================
# Go tooling (ingestion module)
# =============================================================================
GO_MODULES := services/ingestion gen/go

.PHONY: fmt
fmt: ## Format Go code
	cd services/ingestion && gofmt -s -w .

.PHONY: fmt-check
fmt-check: ## Check Go formatting (fails if not formatted)
	@cd services/ingestion && \
	  out=$$(gofmt -s -l .) && \
	  if [ -n "$$out" ]; then \
	    echo "Files not formatted:"; \
	    echo "$$out"; \
	    exit 1; \
	  fi

.PHONY: vet
vet: ## Run go vet on Go code
	cd services/ingestion && go vet ./...

.PHONY: tidy
tidy: ## Run go mod tidy on all Go modules
	@for m in $(GO_MODULES); do \
	  echo "==> tidy $$m"; \
	  (cd $$m && go mod tidy); \
	done

.PHONY: download
download: ## Download Go module dependencies
	cd services/ingestion && go mod download

.PHONY: test
test: ## Run Go tests
	cd services/ingestion && go test ./...

.PHONY: build
build: ## Build the ingestion Go binary locally (no Docker)
	cd services/ingestion && go build -o ../../bin/ingestion ./cmd

.PHONY: lint
lint: fmt-check vet ## Run Go linters

.PHONY: check
check: fmt-check vet test ## Run all Go checks

# =============================================================================
# Python (generator)
# =============================================================================
GENERATOR_VENV    := services/generator/venv
GENERATOR_PYTHON  := $(GENERATOR_VENV)/bin/python
GENERATOR_PIP     := $(GENERATOR_VENV)/bin/pip
GENERATOR_PYFILES := $(shell find services/generator -name '*.py' -not -path '*/venv/*')

.PHONY: python-venv
python-venv: ## Create the generator virtualenv if it doesn't exist
	@test -d $(GENERATOR_VENV) || python3 -m venv $(GENERATOR_VENV)
	@$(GENERATOR_PIP) install --quiet --upgrade pip

.PHONY: python-install
python-install: python-venv ## Install generator dependencies into the venv
	$(GENERATOR_PIP) install -r services/generator/requirements.txt

.PHONY: python-lint
python-lint: python-install ## Run ruff on the generator source
	$(GENERATOR_VENV)/bin/ruff check services/generator

.PHONY: python-format
python-format: python-install ## Format generator source with ruff
	$(GENERATOR_VENV)/bin/ruff format services/generator

.PHONY: python-format-check
python-format-check: python-install ## Check generator formatting (fails if not formatted)
	$(GENERATOR_VENV)/bin/ruff format --check services/generator

.PHONY: python-check
python-check: python-lint python-format-check ## Run all Python checks

# =============================================================================
# Aggregating all the checks
# =============================================================================
.PHONY: check-all
check-all: proto-check python-check check ## Run all proto, Python, and Go checks

# =============================================================================
# Build images (buildx → local Docker daemon)
# =============================================================================
.PHONY: build-generator
build-generator: ## Build the generator image locally
	docker buildx build \
		--file $(GENERATOR_DIR)/Dockerfile \
		--tag $(GENERATOR_IMAGE) \
		--load \
		$(GENERATOR_DIR)

.PHONY: build-ingestion
build-ingestion: ## Build the ingestion image locally
	docker buildx build \
		--file $(INGESTION_DIR)/Dockerfile \
		--tag $(INGESTION_IMAGE) \
		--load \
		.

.PHONY: build-all
build-all: build-generator build-ingestion ## Build all images locally

# =============================================================================
# Load images (push to the local registry)
# =============================================================================
.PHONY: load-generator
load-generator: build-generator ## Build and push the generator image
	docker push $(GENERATOR_IMAGE)

.PHONY: load-ingestion
load-ingestion: build-ingestion ## Build and push the ingestion image
	docker push $(INGESTION_IMAGE)

.PHONY: load-all
load-all: load-generator load-ingestion ## Build and push all images

# =============================================================================
# Platform chart (Kafka + Strimzi + topics)
# =============================================================================
.PHONY: platform-install
platform-install: ## Install or upgrade the platform chart
	helm upgrade --install $(PLATFORM_RELEASE) $(PLATFORM_CHART) \
		-n $(PLATFORM_NAMESPACE) --create-namespace

.PHONY: platform-uninstall
platform-uninstall: ## Uninstall the platform chart
	helm uninstall $(PLATFORM_RELEASE) -n $(PLATFORM_NAMESPACE)

# =============================================================================
# App chart (generator, ingestion, stream-processor, signal-service)
# =============================================================================
.PHONY: app-install
app-install: ## Install or upgrade the app chart
	helm upgrade --install $(APP_RELEASE) $(APP_CHART) \
		-n $(APP_NAMESPACE) --create-namespace

.PHONY: app-uninstall
app-uninstall: ## Uninstall the app chart
	helm uninstall $(APP_RELEASE) -n $(APP_NAMESPACE)
