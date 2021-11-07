.DEFAULT_GOAL := help

.git/hooks/pre-commit:
	@echo "make test" > $@
	@chmod +x $@

.PHONY: help
help:
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | sort | awk -F ':.*?## ' 'NF==2 {printf "\033[36m  %-25s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: .git/hooks/pre-commit test ## setup development environment

.PHONY: test
test: ## test application
	go test ./...
