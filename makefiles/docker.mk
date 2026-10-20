.PHONY: build-generator build-ingestion build-indicator build-signal build-bff build-all \
        load-generator load-ingestion load-indicator load-signal load-bff load-all

build-generator: ## Build the generator image
	docker buildx build --file $(GENERATOR_DIR)/Dockerfile --tag $(GENERATOR_IMAGE) --load $(GENERATOR_DIR)

build-ingestion: ## Build the ingestion image
	docker buildx build --file $(INGESTION_DIR)/Dockerfile --tag $(INGESTION_IMAGE) --load .

build-indicator: ## Build the indicator-calculator image
	docker buildx build --file $(INDICATOR_DIR)/Dockerfile --tag $(INDICATOR_IMAGE) --load .

build-signal: ## Build the signal-decision image
	docker buildx build --file $(SIGNAL_DIR)/Dockerfile --tag $(SIGNAL_IMAGE) --load .

build-bff: ## Build the bff image (embeds frontend)
	docker buildx build --file $(BFF_DIR)/Dockerfile --tag $(BFF_IMAGE) --load .

build-all: build-generator build-ingestion build-indicator build-signal build-bff ## Build all Docker images

load-generator: build-generator ## Build and push the generator image
	docker push $(GENERATOR_IMAGE)

load-ingestion: build-ingestion ## Build and push the ingestion image
	docker push $(INGESTION_IMAGE)

load-indicator: build-indicator ## Build and push the indicator-calculator image
	docker push $(INDICATOR_IMAGE)

load-signal: build-signal ## Build and push the signal-decision image
	docker push $(SIGNAL_IMAGE)

load-bff: build-bff ## Build, push, and roll out the bff image
	docker push $(BFF_IMAGE)
	-kubectl -n $(APP_NAMESPACE) rollout restart deploy/bff

load-all: load-generator load-ingestion load-indicator load-signal load-bff ## Build and push all Docker images