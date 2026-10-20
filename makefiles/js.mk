.PHONY: js-install js-dev js-build js-typecheck js-check
js-install: ## Install frontend dependencies (npm ci)
	cd $(JS_DIR) && $(JS_NPM) ci

js-dev: ## Start the Vite dev server
	cd $(JS_DIR) && $(JS_NPM) run dev

js-build: js-install ## Build the production bundle
	cd $(JS_DIR) && $(JS_NPM) run build

js-typecheck: js-install ## Run TypeScript type checking (tsc --noEmit)
	cd $(JS_DIR) && npx tsc --noEmit

js-check: js-typecheck ## Run all JS checks