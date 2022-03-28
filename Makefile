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
test: compile-scs ## Runs only short tests without checking race conditions
	$(STOPDB) || true
	$(RUNDB); sleep 5
	trap '$(STOPDB)' EXIT; go test -short -p 1 ./...

.PHONY: test-full
test-full: build-docker compile-scs ## Runs all tests checking race conditions
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
run: compile-scs ## Runs all the services
	$(RUNDB)
	$(RUNEXPLORERDB)
	$(RUNNETWORK)
	sleep 5
	$(RUNPROVER)
	sleep 2
	$(RUNCORE)
	sleep 3
	$(RUNEXPLORER)

.PHONY: init-network
init-network: ## Inits network and deploys test smart contract
	go run ./scripts/init_network/main.go .
	sleep 5
	go run ./scripts/deploy_sc/main.go .

.PHONY: deploy-uniswap
deploy-uniswap: ## Deploy the uniswap environment to the network
	go run ./scripts/uniswap/main.go .

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
	cd proto/src/proto/mt/v1 && protoc --proto_path=. --go_out=../../../../../state/tree/pb --go-grpc_out=../../../../../state/tree/pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative mt.proto
	cd proto/src/proto/zkprover/v1 && protoc --proto_path=. --go_out=../../../../../proverclient --go-grpc_out=../../../../../proverclient --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative zk-prover.proto
	cd proto/src/proto/zkprover/v1 && protoc --proto_path=. --go_out=../../../../../proverservice/pb --go-grpc_out=../../../../../proverservice/pb --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative zk-prover.proto

.PHONY: update-external-dependencies
update-external-dependencies: ## Updates external dependencies like images, test vectors or proto files
	go run ./scripts/cmd/... updatedeps

.PHONY: run-benchmarks
run-benchmarks: run-db ## Runs benchmars
	go test -bench=. ./state/tree

DOCKER_CMD := docker run --rm
INPUT_DIR := /contracts
OUTPUT_DIR := $(INPUT_DIR)/bin
OUTPUT_TYPE := --abi --bin
CONTRACTS_DIR := $$(pwd)/test/contracts
CONTRACTS_VOLUME := -v $(CONTRACTS_DIR):$(INPUT_DIR)
SOLC_IMAGE_PREFIX := ethereum/solc:
SOLC_IMAGE_SUFFIX := -alpine
FLAGS := --overwrite --optimize
ABIGEN_DOCKER_IMAGE := ethereum/client-go:alltools-latest
COMPILE_CMD := eval $(DOCKER_CMD) $(CONTRACTS_VOLUME) -e SC_NAME='$$SC_NAME' $(SOLC_IMAGE_PREFIX)'$$SOLC_VERSION'$(SOLC_IMAGE_SUFFIX) -o $(OUTPUT_DIR)/'$$SC_OUTPUT_PATH''$$SC_NAME' $(OUTPUT_TYPE) $(INPUT_DIR)/'$$SC_INPUT_PATH''$$SC_NAME'.sol $(FLAGS)
GENERATE_CMD := eval $(DOCKER_CMD) $(CONTRACTS_VOLUME) $(ABIGEN_DOCKER_IMAGE) abigen --bin=$(OUTPUT_DIR)/'$$SC_OUTPUT_PATH''$$SC_NAME'/'$$SC_NAME'.bin --abi=$(OUTPUT_DIR)/'$$SC_OUTPUT_PATH''$$SC_NAME'/'$$SC_NAME'.abi --pkg='$$SC_NAME' --out=$(OUTPUT_DIR)/'$$SC_OUTPUT_PATH''$$SC_NAME'/'$$SC_NAME'.go

.PHONY: compile-scs
compile-scs: ## Compiles smart contracts used in tests and local deployments
	SC_NAME=counter SOLC_VERSION=0.8.13 $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=destruct SOLC_VERSION=0.8.13 $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=double SOLC_VERSION=0.8.13 $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=emitLog SOLC_VERSION=0.8.13 $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=erc20 SOLC_VERSION=0.8.13 $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=interaction SOLC_VERSION=0.8.13 $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=storage SOLC_VERSION=0.8.13 $(COMPILE_CMD) && $(GENERATE_CMD)

	SC_NAME=UniswapV2ERC20 SOLC_VERSION=0.5.16 SC_INPUT_PATH=uniswap/v2/ SC_OUTPUT_PATH=uniswap/v2/core/ $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=UniswapV2Factory SOLC_VERSION=0.5.16 SC_INPUT_PATH=uniswap/v2/ SC_OUTPUT_PATH=uniswap/v2/core/ $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=UniswapV2Pair SOLC_VERSION=0.5.16 SC_INPUT_PATH=uniswap/v2/ SC_OUTPUT_PATH=uniswap/v2/core/ $(COMPILE_CMD) && $(GENERATE_CMD)

	SC_NAME=UniswapV2Migrator SOLC_VERSION=0.6.6 SC_INPUT_PATH=uniswap/v2/ SC_OUTPUT_PATH=uniswap/v2/periphery/ $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=UniswapV2Router01 SOLC_VERSION=0.6.6 SC_INPUT_PATH=uniswap/v2/ SC_OUTPUT_PATH=uniswap/v2/periphery/ $(COMPILE_CMD) && $(GENERATE_CMD)
	SC_NAME=UniswapV2Router02 SOLC_VERSION=0.6.6 SC_INPUT_PATH=uniswap/v2/ SC_OUTPUT_PATH=uniswap/v2/periphery/ $(COMPILE_CMD) && $(GENERATE_CMD)

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
