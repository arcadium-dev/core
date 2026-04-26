export SHELL := /bin/bash

# ____ all __________________________________________________________________

.PHONY: all

all: test lint

# ____ lint __________________________________________________________________

.PHONY: fmt tidy vet staticcheck vuln lint

fmt:
	@printf "\nRunning go fmt...\n"
	go fmt ./...

tidy:
	@printf "\nRunning go mod tidy...\n"
	go mod tidy

vet:
	@printf "\nRunning go vet...\n"
	go vet ./...

staticcheck:
	@printf "\nRunning staticcheck...\n"
	@go tool staticcheck ./...

vuln:
	@printf "\nRunning govulncheck...\n"
	@go tool govulncheck ./...

lint: fmt tidy vet staticcheck vuln
	@printf "\nChecking for changed files...\n"
	git status --porcelain
	@printf "\n"
	@if [[ "$${CI}" == "true" ]]; then $$(exit $$(git status --porcelain | wc -l)); fi

# ____ test __________________________________________________________________

.PHONY: unit_test test

unit_test:
	@printf "\nRunning go test...\n"
	go test -cover -race ./...

test: unit_test

# ____ clean artifacts _______________________________________________________

.PHONY: clean

clean:
	@printf "\nClean...\n"
	-go clean -testcache
	-rm -f $$(go env GOPATH)/bin/staticcheck
	-rm -f $$(go env GOPATH)/bin/govulncheck
	-rm -f $$(go env GOPATH)/bin/swagger
