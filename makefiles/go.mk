.PHONY: go-fmt go-fmt-check go-vet go-tidy go-download go-test go-check

go-fmt: ## Format Go code in all modules
	@for m in $(GO_MODULES); do \
	  echo "==> gofmt $$m"; \
	  (cd $$m && gofmt -s -w .); \
	done

go-fmt-check: ## Check Go formatting (fails if any file is unformatted)
	@failed=0; \
	for m in $(GO_MODULES); do \
	  out=$$(cd $$m && gofmt -s -l .); \
	  if [ -n "$$out" ]; then \
	    echo "Unformatted in $$m:"; \
	    echo "$$out"; \
	    failed=1; \
	  fi; \
	done; \
	if [ $$failed -eq 1 ]; then exit 1; fi

go-vet: ## Run `go vet` on all modules
	@for m in $(GO_MODULES); do \
	  echo "==> vet $$m"; \
	  (cd $$m && go vet ./...); \
	done

go-tidy: ## Run `go mod tidy` on all modules
	@for m in $(GO_MODULES); do \
	  echo "==> tidy $$m"; \
	  (cd $$m && go mod tidy); \
	done

go-download: ## Download Go module dependencies
	@for m in $(GO_MODULES); do \
	  echo "==> download $$m"; \
	  (cd $$m && go mod download); \
	done

go-test: ## Run Go tests in all modules
	@for m in $(GO_MODULES); do \
	  echo "==> test $$m"; \
	  (cd $$m && CGO_LDFLAGS="-lm" go test ./...); \
	done

go-check: go-fmt-check go-vet go-test ## Run all Go checks