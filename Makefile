.DEFAULT_GOAL : help

GOLINT_CI_COMMAND = docker run -t --rm -v $(PWD):/app -v ~/.cache/golangci-lint/v1.60.1:/root/.cache -w /app golangci/golangci-lint:v1.60.1 golangci-lint

TAG = $$(git describe --tags --abbrev=0)

.PHONY: help
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[32m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: run
run: ## make run main.go
	go run main.go

.PHONY: lint
lint: ## golangci-lint run
	golangci-lint run

.PHONY: test
test: ## go test ./...
	go test ./...

.PHONY: ci-go-lint-run
ci-go-lint-run: ## golangci-lint run
	$(GOLINT_CI_COMMAND) run

.PHONY: pkg-update
pkg-update: ## Update package tag version on pkg.go.dev
	curl https://sum.golang.org/lookup/github.com/thumbrise/validrator@$(TAG)