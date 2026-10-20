.PHONY: java-fmt java-fmt-check java-spotbugs java-sec java-test java-check
java-fmt: ## Format Java code with Spotless
	cd $(JAVA_MODULE) && mvn -q spotless:apply

java-fmt-check: ## Check Java formatting
	cd $(JAVA_MODULE) && mvn -q spotless:check

java-spotbugs: ## Run SpotBugs static analysis
	cd $(JAVA_MODULE) && mvn -q compile spotbugs:check

java-sec: ## Scan Java dependencies for CVEs (OWASP)
	cd $(JAVA_MODULE) && mvn -q org.owasp:dependency-check-maven:check -DfailBuildOnCVSS=7

java-test: ## Run Java tests
	cd $(JAVA_MODULE) && mvn -q test

java-check: java-fmt-check java-spotbugs java-test ## Run all Java checks