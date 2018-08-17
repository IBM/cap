# Copyright 2018 The cap Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO := go
EPOCH_TEST_COMMIT := f98810348b4b26d8c4dd64e36bde613049b2b0b9
PROJECT := github.com/IBM/cap
BUILD_DIR := _output

default: help

help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'binaries'      - Build binaries"
	@echo " * 'clean'         - Clean build artifacts"
	@echo " * 'test'          - Run Unit and Integration Tests"
	@echo " * 'unit'          - Run Unit Tests"
	@echo " * 'coverage'      - Determine Unit and Integration Test Coverage, generates profile in coverage.out"
	@echo " * 'verify'        - Execute the source code verification tools"
	@echo " * 'install.tools' - Install tools used by verify"

.PHONY: binaries clean test unit integration coverage $(BUILD_DIR)/captn

binaries: $(BUILD_DIR)/captn

$(BUILD_DIR)/captn:
	$(GO) build -o $@ \
		-gcflags '$(GO_GCFLAGS)' \
		$(PROJECT)/go/captn

clean:
	rm -rf $(BUILD_DIR)/*

test: ## test the go packages unit and integration
	$(GO) test ./go/atom ./go/cap ./go/shared -v -tags=integration

unit: ## test the go packages
		$(GO) test ./go/atom ./go/cap ./go/shared -v

coverage: ## test and determine coverage of the go packages
	$(GO) test ./go/atom ./go/cap ./go/shared -tags=integration -covermode=count -coverprofile=$(BUILD_DIR)/coverage.out

.PHONY: verify gofmt golint

verify: .gitvalidation gofmt golint

.PHONY: .gitvalidation
# When this is running in travis, it will only check the travis commit range.
# When running outside travis, it will check from $(EPOCH_TEST_COMMIT)..HEAD.
.gitvalidation:
ifeq ($(TRAVIS),true)
	@echo "checking for DCO in TRAVIS"
	git-validation -q -run DCO,short-subject
else
	@echo "checking for DCO"
	git-validation -v -run DCO,short-subject -range $(EPOCH_TEST_COMMIT)..HEAD
endif

gofmt:
	@echo "checking gofmt"
	@./hack/verify-gofmt.sh

golint:
	@echo "checking golint"
	@./hack/verify-golint.sh

.PHONY: install.tools .install.gitvalidation .install.gometalinter

install.tools: .install.gitvalidation .install.gometalinter

.install.gitvalidation:
	$(GO) get -u github.com/vbatts/git-validation

.install.gometalinter:
	$(GO) get -u github.com/alecthomas/gometalinter
	gometalinter --install
