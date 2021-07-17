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

mockgen_version            := v1.6.0

# ----

all: lint test
.PHONY: all

# ----

install_mockgen:
	@bin/install_mockgen $(mockgen_version)
.PHONY: install_mockgen

install: install_mockgen

.PHONY: install

# ----

fmt:
	@printf "\nRunning go fmt...\n"
	@go fmt ./...
.PHONY: fmt


tidy:
	@printf "\nRunning go mod tidy...\n"
	@go mod tidy
.PHONY: tidy

generate:
	@printf "\nRunning go generate...\n"
	@go generate -x ./...
.PHONY: generate

lint: fmt tidy generate
	@printf "\nChecking for changed files...\n"
	@git status --porcelain
	@printf "\n"
	@if [[ "$${CI}" == "true" ]]; then $$(exit $$(git status --porcelain | wc -l)); fi
.PHONY: lint

# ----

test:
	@printf "\nRunning go test...\n"
	@go test -cover $$(go list ./... | grep -v /mock)
.PHONY: test
