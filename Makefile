SHELL := /bin/bash
.DEFAULT_GOAL := help

# Include all modular sub-makefiles
include makefiles/config.mk
include makefiles/tools.mk
include makefiles/cluster.mk
include makefiles/proto.mk
include makefiles/go.mk
include makefiles/python.mk
include makefiles/java.mk
include makefiles/js.mk
include makefiles/docker.mk
include makefiles/helm.mk

# =============================================================================
# Help
# =============================================================================
help: ## Show this help
	@echo "Trading signals platform — available targets:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# =============================================================================
# Composite flows
# =============================================================================
.PHONY: up
up: cluster-up build-all load-all deploy-all ## Boot cluster, build images, deploy everything
	@echo ""
	@echo "Cluster is up. Access the UI at:"
	@echo "  http://192.168.56.243"
	@echo ""

.PHONY: down
down: undeploy-all cluster-down ## Uninstall charts, destroy cluster
	@echo ""
	@echo "Cluster destroyed. Run 'make up' to start fresh."

# =============================================================================
# Aggregated checks — CI entrypoint
# =============================================================================
.PHONY: check-all
check-all: proto-check python-check java-check js-check go-check ## Run all checks (CI)

.PHONY: check-quick
check-quick: proto-check go-fmt-check java-fmt-check js-typecheck ## Fast checks only (no tests)

# =============================================================================
# Cleanup
# =============================================================================
.PHONY: clean
clean: ## Remove build artifacts across all services
	rm -rf gen/go/trading
	rm -rf services/generator/venv
	rm -rf services/indicator-calculator/target
	rm -rf services/frontend/node_modules services/frontend/dist services/frontend/.vite
	@for m in $(GO_MODULES); do \
	  echo "==> clean $$m"; (cd $$m && go clean); \
	  done
