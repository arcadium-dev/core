# Copyright 2021 arcadium.dev <info@arcadium.dev>
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

export SHELL := /bin/bash

mockgen_version := v1.6.0

# ----

.PHONY: all
all: lint test

# ----

.PHONY: install_mockgen
install_mockgen:
	@bin/install_mockgen $(mockgen_version)

.PHONY: install
install: install_mockgen

# ----

.PHONY: fmt
fmt:
	@printf "\nRunning go fmt...\n"
	@go fmt ./...

.PHONY: tidy
tidy:
	@printf "\nRunning go mod tidy...\n"
	@go mod tidy

#.PHONY: generate
#generate:
#	@printf "\nRunning go generate...\n"
#	@go generate -x ./...

.PHONY: lint
lint: fmt tidy
	@printf "\nChecking for changed files...\n"
	@git status --porcelain
	@printf "\n"
	@if [[ "$${CI}" == "true" ]]; then $$(exit $$(git status --porcelain | wc -l)); fi

# ----

.PHONY: test
unit_test:
	@printf "\nRunning go test...\n"
	@go test -cover -race $$(go list ./... | grep -v /mock)

.PHONY: test
test: unit_test
