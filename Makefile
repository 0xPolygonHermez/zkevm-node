DOCKERCOMPOSE := docker-compose -f docker-compose.yml
DOCKERCOMPOSEAPP := hez-core
DOCKERCOMPOSEDB := hez-postgres
DOCKERCOMPOSENETWORK := hez-network
DOCKERCOMPOSEPROVER := hez-prover
DOCKERCOMPOSEEXPLORER := hez-explorer
DOCKERCOMPOSEEXPLORERDB := hez-explorer-postgres

RUNDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEDB)
RUNCORE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPP)
RUNNETWORK := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENETWORK)
RUNPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPROVER)
RUNEXPLORER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORER)
RUNEXPLORERDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERDB)
RUN := $(DOCKERCOMPOSE) up -d

STOPDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEDB)
STOPCORE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPP) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPP)
STOPNETWORK := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENETWORK) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENETWORK)
STOPPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPROVER)
STOPEXPLORER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORER)
STOPEXPLORERDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERDB)
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
build: ## Builds the binary locally into ./dist
	$(BUILD)

.PHONY: build-docker
build-docker: ## Builds a docker image with the core binary
	docker build -t hezcore -f ./Dockerfile .

.PHONY: test
test: ## Runs only short tests without checking race conditions
	$(STOPDB) || true
	$(RUNDB); sleep 5
	trap '$(STOPDB)' EXIT; go test -short -p 1 ./...

.PHONY: test-full
test-full: build-docker ## Runs all tests checking race conditions
	$(STOPDB) || true
	$(RUNDB); sleep 5
	trap '$(STOPDB)' EXIT; MallocNanoZone=0 go test -race -p 1 -timeout 600s ./...

.PHONY: install-linter
install-linter: ## Installs the linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.39.0

.PHONY: lint
lint: ## Runs the linter
	$(LINT)

.PHONY: validate
validate: lint build test-full ## Validates the whole integrity of the code base

.PHONY: run-db
run-db: ## Runs the node database
	$(RUNDB)

.PHONY: stop-db
stop-db: ## Stops the node database
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
stop-prover: ## Stops the zk prover
	$(STOPPROVER)

.PHONY: run-explorer
run-explorer: ## Runs the explorer
	$(RUNEXPLORER)

.PHONY: stop-explorer
stop-explorer: ## Stops the explorer
	$(STOPEXPLORER)

.PHONY: run-explorer-db
run-explorer-db: ## Runs the explorer database
	$(RUNEXPLORERDB)

.PHONY: stop-explorer-db
stop-explorer-db: ## Stops the explorer database
	$(STOPEXPLORERDB)

.PHONY: run
run: ## Runs all the services
	$(RUNDB)
	$(RUNEXPLORERDB)
	$(RUNNETWORK)
	sleep 5
	$(RUNPROVER)
	sleep 2
	$(RUNCORE)
	sleep 3
	$(RUNEXPLORER)
	sleep 3
	go run ./test/init_network.go .

.PHONY: stop
stop: ## Stops all services
	$(STOP)

.PHONY: restart
restart: stop run ## Executes `make stop` and `make run` commands

.PHONY: run-db-scripts
run-db-scripts: ## Executes scripts on the db after it has been initialized, potentially using info from the environment
	./scripts/postgres/run.sh

.PHONY: install-git-hooks
install-git-hooks: ## Moves hook files to the .git/hooks directory
	cp .github/hooks/* .git/hooks

.PHONY: generate-mocks
generate-mocks: ## Generates mocks for the tests, using mockery tool
	mockery --name=etherman --dir=sequencer/strategy/txprofitabilitychecker --output=sequencer/strategy/txprofitabilitychecker --outpkg=txprofitabilitychecker_test --filename=etherman-mock_test.go
	mockery --name=batchProcessor --dir=sequencer/strategy/txselector --output=sequencer/strategy/txselector --outpkg=txselector_test --filename=batchprocessor-mock_test.go

.PHONY: generate-code-from-proto
generate-code-from-proto: ## Generates code from proto files
	cd state/tree/pb && protoc --proto_path=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative mt.proto

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
