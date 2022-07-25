DOCKERCOMPOSE := docker-compose -f docker-compose.yml
DOCKERCOMPOSEAPPSEQ := zkevm-sequencer
DOCKERCOMPOSEAPPAGG := zkevm-aggregator
DOCKERCOMPOSEAPPRPC := zkevm-json-rpc
DOCKERCOMPOSEAPPSYNC := zkevm-sync
DOCKERCOMPOSEAPPBROADCAST := zkevm-broadcast
DOCKERCOMPOSEDB := zkevm-db
DOCKERCOMPOSENETWORK := zkevm-mock-l1-network
DOCKERCOMPOSEPROVER := zkevm-mock-prover
DOCKERCOMPOSEEXPLORER := zkevm-explorer
DOCKERCOMPOSEEXPLORERDB := zkevm-explorer-db
DOCKERCOMPOSEEXPLORERRPC := zkevm-explorer-json-rpc
DOCKERCOMPOSEZKPROVER := zkevm-prover

RUNDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEDB)
RUNSEQUENCER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQ)
RUNAGGREGATOR := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPAGG)
RUNJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPRPC)
RUNSYNC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSYNC)
RUNBROADCAST := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPBROADCAST)

RUNL1NETWORK := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENETWORK)
RUNPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPROVER)
RUNEXPLORER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORER)
RUNEXPLORERDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERDB)
RUNEXPLORERJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERRPC)
RUNZKPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEZKPROVER)

RUN := $(DOCKERCOMPOSE) up -d

STOPDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEDB)
STOPSEQUENCER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQ) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQ)
STOPAGGREGATOR := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPAGG) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPAGG)
STOPJSONRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPRPC)
STOPSYNC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSYNC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSYNC)
STOPBROADCAST := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPBROADCAST) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPBROADCAST)

STOPNETWORK := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENETWORK) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENETWORK)
STOPPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPROVER)
STOPEXPLORER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORER)
STOPEXPLORERDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERDB)
STOPEXPLORERRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERRPC)
STOPZKPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEZKPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEZKPROVER)

STOP := $(DOCKERCOMPOSE) down --remove-orphans

VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist
GOENVVARS := GOBIN=$(GOBIN)
GOBINARY := zkevm-node
GOCMD := $(GOBASE)/cmd

LINT := $$(go env GOPATH)/bin/golangci-lint run

BUILD := $(GOENVVARS) go build $(LDFLAGS) -o $(GOBIN)/$(GOBINARY) $(GOCMD)

.PHONY: build
build: ## Builds the binary locally into ./dist
	$(BUILD)

.PHONY: build-docker
build-docker: ## Builds a docker image with the node binary
	docker build -t zkevm-node -f ./Dockerfile .

.PHONY: build-docker-nc
build-docker-nc: ## Builds a docker image with the node binary - but without build cache
	docker build --no-cache=true -t zkevm-node -f ./Dockerfile .

.PHONY: test
test: compile-scs ## Runs only short tests without checking race conditions
	$(STOPDB)
	$(STOPZKPROVER)
	$(RUNDB); sleep 5
	$(RUNZKPROVER); sleep 5
	trap '$(STOPDB) && $(STOPZKPROVER)' EXIT; go test -short -race -p 1 ./...

.PHONY: test-full
test-full: build-docker compile-scs ## Runs all tests checking race conditions
	$(STOPDB)
	$(STOPZKPROVER)
	$(RUNDB); sleep 7
	$(RUNZKPROVER); sleep 5
	trap '$(STOPDB) && $(STOPZKPROVER)' EXIT; MallocNanoZone=0 go test -race -p 1 -timeout 1200s `go list ./... | grep -v \/ci\/e2e-group`

.PHONY: test-full-non-e2e
test-full-non-e2e: build-docker compile-scs ## Runs non-e2e tests checking race conditions
	$(STOPDB)
	$(STOPZKPROVER)
	$(RUNDB); sleep 7
	$(RUNZKPROVER)
	sleep 5
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOPDB) && $(STOPZKPROVER)' EXIT; MallocNanoZone=0 go test -short -race -p 1 -timeout 600s ./...

.PHONY: test-e2e-group-1
test-e2e-group-1: build-docker compile-scs ## Runs group 1 e2e tests checking race conditions
	$(STOPDB)
	$(RUNDB); sleep 7
	trap '$(STOPDB)' EXIT; MallocNanoZone=0 go test -race -p 1 -timeout 600s ./ci/e2e-group1/...

.PHONY: test-e2e-group-2
test-e2e-group-2: build-docker compile-scs ## Runs group 2 e2e tests checking race conditions
	$(STOPDB)
	$(STOPZKPROVER)
	$(RUNDB); sleep 7
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOPDB) && $(STOPZKPROVER)' EXIT; MallocNanoZone=0 go test -race -p 1 -timeout 600s ./ci/e2e-group2/...

.PHONY: test-e2e-group-3
test-e2e-group-3: build-docker compile-scs ## Runs group 3 e2e tests checking race conditions
	$(STOPDB)
	$(RUNDB); sleep 7
	trap '$(STOPDB)' EXIT; MallocNanoZone=0 go test -race -p 1 -timeout 600s ./ci/e2e-group3/...

.PHONY: install-linter
install-linter: ## Installs the linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.46.2

.PHONY: lint
lint: ## Runs the linter
	$(LINT)

.PHONY: check
check: lint build test ## lint, build and unit tests

.PHONY: validate
validate: lint build test-full ## lint, build, unit and e2e tests

.PHONY: run-db
run-db: ## Runs the node database
	$(RUNDB)

.PHONY: stop-db
stop-db: ## Stops the node database
	$(STOPDB)

.PHONY: run-node
run-node: ## Runs the node
	$(RUNSEQUENCER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)
	$(RUNSYNC)

.PHONY: stop-node
stop-node: ## Stops the node
	$(STOPSEQUENCER)
	$(STOPJSONRPC)
	$(STOPAGGREGATOR)
	$(STOPSYNC)

.PHONY: run-network
run-network: ## Runs the l1 network
	$(RUNL1NETWORK)

.PHONY: stop-network
stop-network: ## Stops the l1 network
	$(STOPNETWORK)

.PHONY: run-prover
run-prover: ## Runs the zk prover
	$(RUNPROVER)

.PHONY: stop-prover
stop-prover: ## Stops the zk prover
	$(STOPPROVER)

.PHONY: run-zkprover
run-zkprover: ## Runs zkprover
	$(RUNZKPROVER)

.PHONY: stop-zkprover
stop-zkprover: ## Stops zkprover
	$(STOPZKPROVER)

.PHONY: run-explorer
run-explorer: ## Runs the explorer
	$(RUNEXPLORERJSONRPC)
	$(RUNEXPLORER)

.PHONY: stop-explorer
stop-explorer: ## Stops the explorer
	$(STOPEXPLORER)
	$(STOPEXPLORERRPC)

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
	$(RUNL1NETWORK)
	sleep 5
	$(RUNZKPROVER)
	sleep 5
	$(RUNPROVER)
	sleep 2
	$(RUNSEQUENCER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)
	$(RUNEXPLORERJSONRPC)
	$(RUNSYNC)
	$(RUNEXPLORER)

.PHONY: run-broadcast
run-broadcast: ## Runs the broadcast service
	$(RUNBROADCAST)

run-seq:
	$(RUNSEQUENCER)

.PHONY: stop-broadcast
stop-broadcast: ## Stops the broadcast service
	$(STOPBROADCAST)

.PHONY: init-network
init-network: ## Initializes the network
	go run ./scripts/init_network/main.go .

.PHONY: deploy-sc
deploy-sc: ## deploys some examples of transactions and smart contracts
	go run ./scripts/deploy_sc/main.go .

.PHONY: deploy-uniswap
deploy-uniswap: ## deploy the uniswap environment to the network
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
	mockery --name=storageInterface --dir=jsonrpc --output=jsonrpc --outpkg=jsonrpc --inpackage --structname=storageMock --filename=mock_storage_test.go
	mockery --name=jsonRPCTxPool --dir=jsonrpc --output=jsonrpc --outpkg=jsonrpc --inpackage --structname=poolMock --filename=mock_pool_test.go
	mockery --name=gasPriceEstimator --dir=jsonrpc --output=jsonrpc --outpkg=jsonrpc --inpackage --structname=gasPriceEstimatorMock --filename=mock_gasPriceEstimator_test.go
	mockery --name=stateInterface --dir=jsonrpc --output=jsonrpc --outpkg=jsonrpc --inpackage --structname=stateMock --filename=mock_state_test.go
	mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=jsonrpc --outpkg=jsonrpc --structname=dbTxMock --filename=mock_dbtx_test.go

	mockery --name=txManager --dir=sequencer --output=sequencer --outpkg=sequencer --structname=txmanagerMock --filename=txmanager-mock_test.go
	mockery --name=etherman --dir=sequencer --output=sequencer --outpkg=sequencer --structname=ethermanMock --filename=etherman-mock_test.go
	mockery --name=etherman --dir=sequencer/profitabilitychecker --output=sequencer/profitabilitychecker --outpkg=profitabilitychecker_test --structname=ethermanMock --filename=etherman-mock_test.go
	mockery --name=stateInterface --dir=sequencer/broadcast --output=sequencer/broadcast --outpkg=broadcast_test --structname=stateMock --filename=state-mock_test.go

	mockery --name=ethermanInterface --dir=synchronizer --output=synchronizer --outpkg=synchronizer --structname=ethermanMock --filename=mock_etherman.go
	mockery --name=stateInterface --dir=synchronizer --output=synchronizer --outpkg=synchronizer --structname=stateMock --filename=mock_state.go
	mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=synchronizer --outpkg=synchronizer --structname=dbTxMock --filename=mock_dbtx.go


.PHONY: generate-code-from-proto
generate-code-from-proto: ## Generates code from proto files
	cd proto/src/proto/statedb/v1 && protoc --proto_path=. --proto_path=../../../../include --go_out=../../../../../merkletree/pb --go-grpc_out=../../../../../merkletree/pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative statedb.proto
	cd proto/src/proto/zkprover/v1 && protoc --proto_path=. --go_out=../../../../../proverclient/pb --go-grpc_out=../../../../../proverclient/pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative zk-prover.proto
	cd proto/src/proto/zkprover/v1 && protoc --proto_path=. --go_out=../../../../../proverservice/pb --go-grpc_out=../../../../../proverservice/pb --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative zk-prover.proto
	cd proto/src/proto/executor/v1 && protoc --proto_path=. --go_out=../../../../../state/runtime/executor/pb --go-grpc_out=../../../../../state/runtime/executor/pb --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative executor.proto
	cd proto/src/proto/broadcast/v1 && protoc --proto_path=. --proto_path=../../../../include --go_out=../../../../../sequencer/broadcast/pb --go-grpc_out=../../../../../sequencer/broadcast/pb --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative broadcast.proto

.PHONY: update-external-dependencies
update-external-dependencies: ## Updates external dependencies like images, test vectors or proto files
	go run ./scripts/cmd/... updatedeps

.PHONY: run-benchmarks
run-benchmarks: run-db ## Runs benchmars
	go test -bench=. ./state/tree

.PHONY: compile-scs
compile-scs: ## Compiles smart contracts, configuration in test/contracts/index.yaml
	go run ./scripts/cmd... compilesc --input ./test/contracts

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
