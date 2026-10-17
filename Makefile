GOLANGCI_LINT_VERSION ?= v2.15.0

# Library packages only; the demo programs under examples/ have no tests and
# would skew the coverage total.
PKGS := $(shell go list ./... | grep -v /examples/)

# Use an installed golangci-lint if there is one, otherwise run the pinned
# version through the Go toolchain.
GOLANGCI_LINT := $(shell command -v golangci-lint 2>/dev/null)
ifeq ($(GOLANGCI_LINT),)
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
endif

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: fmt vet lint test ## Format, vet, lint and test

.PHONY: test
test: ## Run tests with the race detector
	go test -race ./...

.PHONY: cover
cover: ## Run tests and report coverage (writes coverage.out and coverage.html)
	go test -race -covermode=atomic -coverprofile=coverage.out $(PKGS)
	go tool cover -func=coverage.out | tail -n 1
	go tool cover -html=coverage.out -o coverage.html

.PHONY: bench
bench: ## Run benchmarks
	go test -run='^$$' -bench=. -benchmem $(PKGS)

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	$(GOLANGCI_LINT) run ./...

.PHONY: fmt
fmt: ## Format code with gofmt and goimports rules
	gofmt -s -w .
	$(GOLANGCI_LINT) fmt ./...

.PHONY: examples
examples: ## Run every example program
	@for d in examples/*/; do echo "== $$d"; go run ./$$d || exit 1; echo; done

.PHONY: clean
clean: ## Remove generated files
	rm -f coverage.out coverage.html
