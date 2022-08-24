# disable implicit rules
.SUFFIXES:

SHELL = /bin/bash

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
GO_BIN_RESIZER = bin/image-resizer

.PHONY: build
build:
	go build $(GO_BUILD_FLAGS) -o "$(GO_BIN_RESIZER)" "."

.PHONY: run
run: build
	$(GO_BIN_RESIZER)

##### TOOLS #####

export TOOLS_PATH = $(CURDIR)/$(TOOLS_BIN_DIR)

TOOLS_MODFILE = tools/go.mod
define install-go-tool =
	go build \
		-o $(TOOLS_BIN_DIR) \
		-ldflags "-s -w" \
		-modfile $(TOOLS_MODFILE)
endef

## COMMON TOOLS

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

IMAGE_RESIZER_IMAGE_TAG = 0.1.0
IMAGE_RESIZER_IMAGE_NAME = abramlab/image-resizer
export IMAGE_RESIZER_IMAGE ?= $(IMAGE_RESIZER_IMAGE_NAME):$(IMAGE_RESIZER_IMAGE_TAG)

.PHONY: build-image
build-image:
	docker build -t $(IMAGE_RESIZER_IMAGE) .

.PHONY: push-image
push-image:
	docker push $(IMAGE_RESIZER_IMAGE)

CONTAINER_DIR = /app

.PHONY: run-image
run-image:
	docker run \
    	-it \
    	--rm \
    	-w $(CONTAINER_DIR) \
    	--mount type=bind,source=$(CURDIR)/images,target=$(CONTAINER_DIR)/images \
    	--mount type=bind,source=$(CURDIR)/resized-images,target=$(CONTAINER_DIR)/resized-images \
		$(IMAGE_RESIZER_IMAGE)
