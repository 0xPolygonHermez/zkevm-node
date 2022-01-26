DOCKERCOMPOSE := docker-compose -f docker-compose.explorer.yml
DOCKERCOMPOSEAPP := hez-core
DOCKERCOMPOSEDB := hez-postgres
DOCKERCOMPOSENETWORK := hez-network
DOCKERCOMPOSEPROVER := hez-prover

RUNDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEDB)
RUNCORE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPP)
RUNNETWORK := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENETWORK)
RUNPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPROVER)
RUN := $(DOCKERCOMPOSE) up -d

STOPDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEDB)
STOPCORE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPP) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPP)
STOPNETWORK := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENETWORK) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENETWORK)
STOPPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPROVER)
STOP := $(DOCKERCOMPOSE) down --remove-orphans

VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist
GOENVVARS := GOBIN=$(GOBIN)
GOBINARY := hezcore
GOCMD := $(GOBASE)/cmd

LINT := $$(go env GOPATH)/bin/golangci-lint run --timeout=5m -E whitespace -E gosec -E gci -E misspell -E gomnd -E gofmt -E goimports -E golint --exclude-use-default=false --max-same-issues 0
BUILD := $(GOENVVARS) go build $(LDFLAGS) -o $(GOBIN)/$(GOBINARY) $(GOCMD)

.PHONY: build
build: ## Build the binary locally into ./dist
	$(BUILD)

.PHONY: build-docker
build-docker: ## Build a docker image with the core binary
	docker build -t hezcore -f ./Dockerfile .

.PHONY: test
test: ## runs only short tests without checking race conditions
	$(STOPDB) || true
	$(RUNDB); sleep 5
	trap '$(STOPDB)' EXIT; go test -short -p 1 ./...

.PHONY: test-full
test-full: build-docker ## runs all tests checking race conditions
	$(STOPDB) || true
	$(RUNDB); sleep 5
	trap '$(STOPDB)' EXIT; MallocNanoZone=0 go test -race -p 1 -timeout 600s ./...

.PHONY: install-linter
install-linter: ## install linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.39.0

.PHONY: lint
lint: ## runs linter
	$(LINT)

.PHONY: validate
validate: lint build test-full ## Validate the whole integrity of the code

.PHONY: run-db
run-db: ## starts a docker container to run the db instance
	$(RUNDB)

.PHONY: stop-db
stop-db: ## stops the docker container running the db instance
	$(STOPDB)

.PHONY: run-core
run-core: ## Runs the core
	$(RUNCORE)

.PHONY: stop-core
stop-core: ## Stops the core
	$(STOPCORE)

.PHONY: run-network
run-network: ## Runs the l1 network
	$(RUNNETWORK)

.PHONY: stop-network
stop-network: ## Stops the l1 network
	$(STOPNETWORK)

.PHONY: run-prover
run-prover: ## Runs the zk prover
	$(RUNPROVER)

.PHONY: stop-prover
stop-prover: ## Stop the zk prover
	$(STOPPROVER)

.PHONY: run
run: ## Runs all the services available in the docker-compose file
	$(RUN)

.PHONY: stop
stop: ## Stops all services available in the docker-compose file
	$(STOP)

.PHONY: restart
restart: stop run ## Executes `make stop` and `make run` commands

.PHONY: run-db-scripts
run-db-scripts: ## executes scripts on the db after it has been initialized, potentially using info from the environment
	./scripts/postgres/run.sh

.PHONY: install-git-hooks
install-git-hooks: ## Moves hook files to the .git/hooks directory
	cp .github/hooks/* .git/hooks

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
