.PHONY: install-tools docker-login

install-tools: ## Report installed CLI tools and how to install missing ones
	@echo "==> Checking tool versions:"
	@for t in go python3 node npm mvn docker buf kubectl helm; do \
	  if ! command -v $$t >/dev/null 2>&1; then \
	    printf "  %-8s \033[33mMISSING\033[0m\n" "$$t"; \
	    continue; \
	  fi; \
	  v=$$($$t --version 2>&1 | head -1); \
	  printf "  %-8s %s\n" "$$t" "$$v"; \
	done
	@echo ""
	@echo "Install links:"
	@echo "  buf:     https://buf.build/docs/installation"
	@echo "  helm:    https://helm.sh/docs/intro/install/"
	@echo "  kubectl: https://kubernetes.io/docs/tasks/tools/"

docker-login: ## Log Docker in to the local registry (skips if no auth)
	@if [ -z "$$REGISTRY_USER" ] || [ -z "$$REGISTRY_PASS" ]; then \
	  echo "Registry does not require auth; skipping docker login."; \
	  exit 0; \
	fi
	@echo "$$REGISTRY_PASS" | docker login $(REGISTRY) -u "$$REGISTRY_USER" --password-stdin