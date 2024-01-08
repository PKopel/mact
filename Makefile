PROGRAM_NAME ?= mact
VERSION ?= v0.0.6

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object paths="./types"

SUPPORTED_ARCH ?= amd64 arm64
SUPPORTED_OS ?= linux windows darwin

.PHONY: build
build:
	@echo "building $(LOCALBIN)/$(PROGRAM_NAME) executable"
	@mkdir -p $(LOCALBIN)
	@go build -o $(LOCALBIN)/$(PROGRAM_NAME) main.go

rename-win = $(shell mv $(LOCALBIN)/$(1) $(LOCALBIN)/$(1).exe)
build-cmd = $(shell GOOS=$(1) GOARCH=$(2) go build -o $(LOCALBIN)/$(PROGRAM_NAME)-$(VERSION)-$(1)-$(2) main.go)
build-os = $(foreach arch, $(SUPPORTED_ARCH), $(call build-cmd,$(1),$(arch)))

.PHONY: build-all
build-all:
	@echo "building executables for all architectures and OSs"
	$(foreach os, $(SUPPORTED_OS), $(call build-os,$(os)))
	$(foreach file, $(shell ls $(LOCALBIN) | grep windows), $(call rename-win,$(file)))

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