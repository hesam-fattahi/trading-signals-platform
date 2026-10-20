.PHONY: cluster-up cluster-halt cluster-down cluster-status provision
cluster-up: ## Boot VMs and provision via Ansible
	$(VAGRANT) up

cluster-halt: ## Shut down VMs (keeps disk state)
	$(VAGRANT) halt

cluster-down: ## Destroy VMs (wipes disks, irreversible)
	$(VAGRANT) destroy -f

cluster-status: ## Show Vagrant VM status
	$(VAGRANT) status

provision: ## Re-run Ansible provisioning (no VM recreation)
	$(VAGRANT) provision