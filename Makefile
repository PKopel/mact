PROGRAM_NAME ?= mact
VERSION ?= v0.0.3
BIN_DIR ?= bin

SUPPORTED_ARCH ?= amd64 arm64
SUPPORTED_OS ?= linux windows darwin

.PHONY: build
build:
	@echo "building $(BIN_DIR)/$(PROGRAM_NAME) executable"
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(PROGRAM_NAME) main.go

rename-win = $(shell mv $(BIN_DIR)/$(1) $(BIN_DIR)/$(1).exe)
build-cmd = $(shell GOOS=$(1) GOARCH=$(2) go build -o $(BIN_DIR)/$(PROGRAM_NAME)-$(VERSION)-$(1)-$(2) main.go)
build-os = $(foreach arch, $(SUPPORTED_ARCH), $(call build-cmd,$(1),$(arch)))

.PHONY: build-all
build-all:
	@echo "building executables for all architectures and OSs"
	$(foreach os, $(SUPPORTED_OS), $(call build-os,$(os)))
	$(foreach file, $(shell ls $(BIN_DIR) | grep windows), $(call rename-win,$(file)))

.PHONY: run
run:
	@go run main.go

.PHONY: install
install:
	@echo "installing $(PROGRAM_NAME) executable"
	@go install

.PHONY: test
test:
	@go test -v ./...

.PHONY: vendor
vendor:
	@go mod vendor

.PHONY: tidy
tidy:
	@go mod tidy