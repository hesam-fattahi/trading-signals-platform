.PHONY: proto-generate proto-lint proto-format proto-format-check proto-breaking proto-check proto-clean
proto-generate: ## Generate Go protobuf stubs from .proto files
	buf generate

proto-lint: ## Lint .proto files
	buf lint

proto-format: ## Format .proto files in place
	buf format -w

proto-format-check: ## Check .proto formatting (fails if not formatted)
	@buf format -d --exit-code

proto-breaking: ## Check .proto files for breaking changes against main
	buf breaking --against '.git#branch=main'

proto-check: proto-lint proto-format-check ## Run all proto checks

proto-clean: ## Remove generated proto output
	rm -rf gen/go/trading