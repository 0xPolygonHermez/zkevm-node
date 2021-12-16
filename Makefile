STARTDB = docker run --rm --name hermez-polygon-db -p 5432:5432 -e POSTGRES_DB="polygon-hermez" -e POSTGRES_USER="hermez" -e POSTGRES_PASSWORD="polygon" -d postgres
STOPDB = docker stop hermez-polygon-db

VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist
GOENVVARS := GOBIN=$(GOBIN)
GOBINARY := hezcore
GOCMD := $(GOBASE)/cmd

LINT := $$(go env GOPATH)/bin/golangci-lint run --timeout=5m -E whitespace -E gosec -E gci -E misspell -E gomnd -E gofmt -E goimports -E golint --exclude-use-default=false --max-same-issues 0
BUILD := $(GOENVVARS) go build $(LDFLAGS) -o $(GOBIN)/$(GOBINARY) $(GOCMD)
RUN := $(GOBIN)/$(GOBINARY) run --network local --cfg ./cmd/config.toml

.PHONY: build
build: ## Build the binary
	$(BUILD)

.PHONY: test
test: ## runs only short tests without checking race conditions
	$(STOPDB) || true
	$(STARTDB)
	go test -short -p 1 ./...
	$(STOPDB)

.PHONY: test-full
test-full: ## runs all tests checking race conditions
	$(STOPDB) || true
	$(STARTDB)
	go test -race -p 1 -timeout 180s ./...
	$(STOPDB)

.PHONY: install-linter
install-linter: ## install linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.30.0

.PHONY: lint
lint: ## runs linter
	$(LINT)

.PHONY: deploy
deploy: lint test-full build ## Validate and create the binary to be deployed

.PHONY: run
run: ## Runs the application
	$(RUN)

.PHONY: start-db
start-db: ## starts a docker container to run the db instance
	$(STARTDB)

.PHONY: stop-db
stop-db: ## stops the docker container running the db instance
	$(STOPDB)

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
		