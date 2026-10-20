GENERATOR_VENV    := services/generator/venv
GENERATOR_PYTHON  := $(GENERATOR_VENV)/bin/python
GENERATOR_PIP     := $(GENERATOR_VENV)/bin/pip

.PHONY: python-venv python-install python-lint python-format python-format-check python-check
python-venv: ## Create the generator virtualenv if missing
	@test -d $(GENERATOR_VENV) || python3 -m venv $(GENERATOR_VENV)
	@$(GENERATOR_PIP) install --quiet --upgrade pip

python-install: python-venv ## Install generator dependencies
	$(GENERATOR_PIP) install -r services/generator/requirements.txt

python-lint: python-install ## Lint generator source with ruff
	$(GENERATOR_VENV)/bin/ruff check services/generator

python-format: python-install ## Format generator source with ruff
	$(GENERATOR_VENV)/bin/ruff format services/generator

python-format-check: python-install ## Check generator formatting (fails if unformatted)
	$(GENERATOR_VENV)/bin/ruff format --check services/generator

python-check: python-lint python-format-check ## Run all Python checks