# disable implicit rules
.SUFFIXES:

SHELL = /bin/bash

OS ?= linux
GO   ?= go
RESIZER_APP_DIR ?= resizer
TOOLS_BIN_DIR ?= tools/bin
CACHE_DIR ?= cache

DIRS_TO_CREATE += $(TOOLS_BIN_DIR) $(CACHE_DIR)
$(DIRS_TO_CREATE):
	-mkdir -p "$@"

define format-go-code =
	$(GO_FORMAT_TOOL) -w
endef

.PHONY: clean-cache
clean-cache:
	-rm -rf $(CACHE_DIR)/*

##### BUILD AND RUN TARGETS #####

GO_BUILD_FLAGS ?= -trimpath
GO_BIN_RESIZER = bin/$(RESIZER_APP_DIR)

.PHONY: build
build:
	$(GO) build $(GO_BUILD_FLAGS) -o "$(GO_BIN_RESIZER)" "./$(RESIZER_APP_DIR)/$*"

.PHONY: run
run: build
	$(GO_BIN_RESIZER)

##### TOOLS #####

export TOOLS_PATH = $(CURDIR)/$(TOOLS_BIN_DIR)

TOOLS_MODFILE = tools/go.mod
define install-go-tool =
	$(GO) build \
		-o $(TOOLS_BIN_DIR) \
		-ldflags "-s -w" \
		-modfile $(TOOLS_MODFILE)
endef

## COMMON TOOLS

ENUMER_TOOL = $(TOOLS_BIN_DIR)/enumer
COMMON_TOOLS += $(ENUMER_TOOL)
$(ENUMER_TOOL): | $(TOOLS_BIN_DIR)
	$(install-go-tool) github.com/alvaroloes/enumer

GO_LINT_TOOL = $(TOOLS_BIN_DIR)/golangci-lint
COMMON_TOOLS += $(GO_LINT_TOOL)
$(GO_LINT_TOOL): | $(TOOLS_BIN_DIR)
	$(install-go-tool) github.com/golangci/golangci-lint/cmd/golangci-lint

GO_FORMAT_TOOL = $(TOOLS_BIN_DIR)/gofumpt
COMMON_TOOLS += $(GO_FORMAT_TOOL)
$(GO_FORMAT_TOOL): | $(TOOLS_BIN_DIR)
	$(install-go-tool) mvdan.cc/gofumpt

$(COMMON_TOOLS): $(TOOLS_MODFILE) | $(TOOLS_BIN_DIR)

##### TOOLS TARGETS #####

.PHONY: install-tools
install-tools: $(COMMON_TOOLS)

.PHONY: lint
lint: $(GO_LINT_TOOL)
	$(GO_LINT_TOOL) run --sort-results

.PHONY: format-go
format-go: $(GO_FORMAT_TOOL)
	$(format-go-code) .

.PHONY: generate-go
generate-go: $(GO_FORMAT_TOOL) $(COMMON_TOOLS)
	$(GO) generate ./...
	$(format-go-code) .