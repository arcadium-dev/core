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
	@if [[ ! -x "$$(go env GOPATH)/bin/staticcheck" ]]; then \
		printf "\nInstalling staticcheck...\n"; \
		go get "honnef.co/go/tools/cmd/staticcheck"; \
		go install "honnef.co/go/tools/cmd/staticcheck"; \
	fi
	@printf "\nRunning staticcheck...\n"
	$$(go env GOPATH)/bin/staticcheck ./...

vuln:
	@if [[ ! -x "$$(go env GOPATH)/bin/govulncheck" ]]; then \
		printf "\nInstalling govulncheck...\n"; \
		go get "golang.org/x/vuln/cmd/govulncheck"; \
		go install "golang.org/x/vuln/cmd/govulncheck"; \
	fi
	@printf "\nRunning govulncheck...\n"
	$$(go env GOPATH)/bin/govulncheck ./...

lint: fmt tidy vet staticcheck vuln
	@printf "\nChecking for changed files...\n"
	git status --porcelain
	@printf "\n"
	git diff go.mod go.sum
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
