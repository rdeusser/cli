# Setup name variables for the package/tool.
NAME := cli
PKG := github.com/rdeusser/$(NAME)
BUILD_PATH := $(PKG)/cmd/$(NAME)
VERSION := $(shell grep -oE "[0-9]+[.][0-9]+[.][0-9]+" version/version.go)

SEMVER := patch

OLDPWD := $(PWD)
export OLDPWD

FILES_TO_FMT ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

GOBIN		   ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO111MODULE	   ?= on
export GO111MODULE

GOIMPORTS_VERSION       ?= master
GOIMPORTS               ?= $(GOBIN)/goimports

REVIVE_VERSION          ?= v1.1.4
REVIVE                  ?= $(GOBIN)/revive

GEN_ENUM_VERSION        ?= main
GEN_ENUM                ?= $(GOBIN)/gen-enum

UPDATE_TESTDATA_VERSION ?= main
UPDATE_TESTDATA         ?= $(GOBIN)/update-testdata

.DEFAULT_GOAL := help

define fetch_go_bin_version
	@cd /tmp
	@go install $(1)@$(2)
	@cd -
endef

.PHONY: help
help: ## Display this help text.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nAvailable targets:\n"} /^[\/0-9a-zA-Z_-]+:.*?##/ { printf "  \x1b[32;01m%-20s\x1b[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: tidy
tidy: $(GOIMPORTS) ## Formats Go code including imports and cleans up noise.
	@echo ">> formatting code"
	@$(GOIMPORTS) -local github.com/rdeusser/$(NAME) -w $(FILES_TO_FMT)
	@echo ">> cleaning up noise"
	@find . -type f \( -name "*.md" -o -name "*.go" \) | SED_BIN="$(SED)" xargs scripts/cleanup-noise.sh
	@echo ">> running 'go mod tidy'"
	@go mod tidy

.PHONY: generate
generate: $(GEN_ENUM) ## Generate code
	@echo ">> generating code"
	@go generate ./...

.PHONY: lint
lint: $(REVIVE) ## Run static analysis tools.
	@echo ">> running linting tools"
	@revive -config revive.toml ./...

.PHONY: test
test: ## Runs all cli's unit tests. This excludes tests in ./test/e2e.
	@echo ">> running unit tests (without /test/e2e)"
	@go test -v -coverprofile=coverage.out $(shell go list ./... | grep -v /test/e2e);

.PHONY: test/e2e
test/e2e: ## Runs all cli's e2e tests from test/e2e.
	@echo ">> running e2e tests"
	@go test -v -tags=e2e -coverprofile=coverage.out ./test/e2e/...

.PHONY: update-testdata
update-testdata: $(UPDATE_TESTDATA) ## Updates all files in testdata directories.
	@echo ">> updating files in testdata directories"
	@update-testdata

.PHONY: bump-version
bump-version: ## Bump the version in the version file. Set SEMVER to [ patch (default) | major | minor ].
	@./scripts/bump-version.sh $(SEMVER)

.PHONY: tag
tag: ## Create and push a new git tag (creates tag using version/version.go file).
	@./scripts/tag.sh

$(GOIMPORTS):
	$(call fetch_go_bin_version,golang.org/x/tools/cmd/goimports,$(GOIMPORTS_VERSION))

$(REVIVE):
	$(call fetch_go_bin_version,github.com/mgechev/revive,$(REVIVE_VERSION))

$(GEN_ENUM):
	$(call fetch_go_bin_version,github.com/rdeusser/x/tools/gen-enum,$(GEN_ENUM_VERSION))

$(UPDATE_TESTDATA):
	$(call fetch_go_bin_version,github.com/rdeusser/x/tools/update-testdata,$(UPDATE_TESTDATA_VERSION))
