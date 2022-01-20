PACKAGE    = pxlmtc
DATE      ?= $(shell date +%FT%T%z)
VERSION   ?= $(shell echo $(shell cat $(PWD)/.version)-$(shell git describe --tags --always))
DIR        = $(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GO         = go
GOROOT     ?= $(shell go env GOROOT)
GODOC      = godoc
GOFMT      = gofmt

# For CI
ifneq ($(wildcard ./bin/golangci-lint),)
	GOLINT = ./bin/golangci-lint
else
	GOLINT = golangci-lint
endif

V          = 0
Q          = $(if $(filter 1,$V),,@)
M          = $(shell printf "\033[0;35m▶\033[0m")

GO_PACKAGE        = github.com/elojah/pxlmtc
SOLVER            = solver
GENERATOR         = generator

.PHONY: all
all: solver generator

.PHONY: solver
solver:  ## Build solver binary
	$(info $(M) building executable solver…) @
	$Q cd cmd/$(SOLVER) && $(GO) build \
		-mod=readonly \
		-tags release \
		-ldflags '-X main.version=$(VERSION) -X main.prog=$(SOLVER)'\
		-o ../../bin/$(PACKAGE)_$(SOLVER)_$(VERSION)
	$Q cp bin/$(PACKAGE)_$(SOLVER)_$(VERSION) bin/$(PACKAGE)_$(SOLVER)

.PHONY: generator
generator:  ## Build generator binary
	$(info $(M) building executable generator…) @
	$Q cd cmd/$(GENERATOR) && $(GO) build \
		-mod=readonly \
		-tags release \
		-ldflags '-X main.version=$(VERSION) -X main.prog=$(GENERATOR)'\
		-o ../../bin/$(PACKAGE)_$(GENERATOR)_$(VERSION)
	$Q cp bin/$(PACKAGE)_$(GENERATOR)_$(VERSION) bin/$(PACKAGE)_$(GENERATOR)

# Vendor
.PHONY: vendor
vendor:
	$(info $(M) running go mod vendor…) @
	$Q $(GO) mod vendor

# Tidy
.PHONY: tidy
tidy:
	$(info $(M) running go mod tidy…) @
	$Q $(GO) mod tidy

# Check
.PHONY: check
check: vendor test lint

# Lint
.PHONY: lint
lint:
	$(info $(M) running $(GOLINT)…)
	$Q $(GOLINT) run

# Test
.PHONY: test
test:
	$(info $(M) running go test…) @
	$Q $(GO) test -cover -race -v ./...

# Clean
.PHONY: clean
clean:
	$(info $(M) cleaning bin…) @
	$Q rm -rf bin/$(PACKAGE)_$(SOLVER)_*

## Helpers

.PHONY: go-version
go-version: ## Print go version used in this makefile
	$Q echo $(GO)

.PHONY: fmt
fmt: ## Format code
	$(info $(M) running $(GOFMT)…) @
	$Q $(GOFMT) ./...

.PHONY: doc
doc: ## Generate project documentation
	$(info $(M) running $(GODOC)…) @
	$Q $(GODOC) ./...

.PHONY: version
version: ## Print current project version
	@echo $(VERSION)

.PHONY: help
help: ## Print this
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
