.PHONY: platform-install platform-uninstall app-install app-uninstall deploy-all undeploy-all helm-lint
platform-install: ## Install or upgrade the platform chart (Kafka + Strimzi)
	helm upgrade --install $(PLATFORM_RELEASE) $(PLATFORM_CHART) \
		-n $(PLATFORM_NAMESPACE) --create-namespace

platform-uninstall: ## Uninstall the platform chart
	helm uninstall $(PLATFORM_RELEASE) -n $(PLATFORM_NAMESPACE)

app-install: ## Install or upgrade the app chart
	helm upgrade --install $(APP_RELEASE) $(APP_CHART) \
		-n $(APP_NAMESPACE) --create-namespace

app-uninstall: ## Uninstall the app chart
	helm uninstall $(APP_RELEASE) -n $(APP_NAMESPACE)

deploy-all: platform-install app-install ## Install platform then app charts

undeploy-all: ## Uninstall all Helm releases (keeps VMs)
	-helm uninstall $(APP_RELEASE) -n $(APP_NAMESPACE) 2>/dev/null
	-helm uninstall $(PLATFORM_RELEASE) -n $(PLATFORM_NAMESPACE) 2>/dev/null

helm-lint: ## Lint all Helm charts
	helm lint $(PLATFORM_CHART)
	helm lint $(APP_CHART)