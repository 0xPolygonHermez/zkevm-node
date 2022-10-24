DOCKERCOMPOSE := docker-compose -f docker-compose.yml
DOCKERCOMPOSEAPPSEQ := zkevm-sequencer
DOCKERCOMPOSEAPPAGG := zkevm-aggregator
DOCKERCOMPOSEAPPRPC := zkevm-json-rpc
DOCKERCOMPOSEAPPSYNC := zkevm-sync
DOCKERCOMPOSEAPPBROADCAST := zkevm-broadcast
DOCKERCOMPOSESTATEDB := zkevm-state-db
DOCKERCOMPOSEPOOLDB := zkevm-pool-db
DOCKERCOMPOSERPCDB := zkevm-rpc-db
DOCKERCOMPOSENETWORK := zkevm-mock-l1-network
DOCKERCOMPOSEEXPLORERL1 := zkevm-explorer-l1
DOCKERCOMPOSEEXPLORERL1DB := zkevm-explorer-l1-db
DOCKERCOMPOSEEXPLORERL2 := zkevm-explorer-l2
DOCKERCOMPOSEEXPLORERL2DB := zkevm-explorer-l2-db
DOCKERCOMPOSEEXPLORERRPC := zkevm-explorer-json-rpc
DOCKERCOMPOSEZKPROVER := zkevm-prover
DOCKERCOMPOSEZKPROVERMOCK := zkprover-mock
DOCKERCOMPOSEPERMISSIONLESSDB := zkevm-permissionless-db
DOCKERCOMPOSEPERMISSIONLESSNODE := zkevm-permissionless-node
DOCKERCOMPOSENODEAPPROVE := zkevm-approve

RUNSTATEDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSESTATEDB)
RUNPOOLDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPOOLDB)
RUNRPCDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSERPCDB)
RUNSEQUENCER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQ)
RUNAGGREGATOR := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPAGG)
RUNJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPRPC)
RUNSYNC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSYNC)
RUNBROADCAST := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPBROADCAST)

RUNL1NETWORK := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENETWORK)
RUNEXPLORERL1 := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL1)
RUNEXPLORERL1DB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL1DB)
RUNEXPLORERL2 := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL2)
RUNEXPLORERL2DB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL2DB)
RUNEXPLORERJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERRPC)
RUNZKPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEZKPROVER)
RUNZKPROVERMOCK := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEZKPROVERMOCK)

RUNPERMISSIONLESSDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPERMISSIONLESSDB)
RUNPERMISSIONLESSNODE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPERMISSIONLESSNODE)

RUNAPPROVE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENODEAPPROVE)

RUN := $(DOCKERCOMPOSE) up -d

STOPSTATEDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSESTATEDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSESTATEDB)
STOPPOOLDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPOOLDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPOOLDB)
STOPRPCDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSERPCDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSERPCDB)
STOPSEQUENCER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQ) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQ)
STOPAGGREGATOR := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPAGG) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPAGG)
STOPJSONRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPRPC)
STOPSYNC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSYNC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSYNC)
STOPBROADCAST := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPBROADCAST) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPBROADCAST)

STOPNETWORK := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENETWORK) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENETWORK)
STOPEXPLORERL1 := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL1) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL1)
STOPEXPLORERL1DB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL1DB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL1DB)
STOPEXPLORERL2 := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL2)
STOPEXPLORERL2DB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL2DB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL2DB)
STOPEXPLORERRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERRPC)
STOPZKPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEZKPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEZKPROVER)
STOPZKPROVERMOCK := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEZKPROVERMOCK) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEZKPROVERMOCK)

STOPPERMISSIONLESSDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPERMISSIONLESSDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPERMISSIONLESSDB)
STOPPERMISSIONLESSNODE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPERMISSIONLESSNODE) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPERMISSIONLESSNODE)

STOPAPPROVE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENODEAPPROVE) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENODEAPPROVE)

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
	export CONFIG_MODE="test"	
	$(STOPSTATEDB)
	$(STOPPOOLDB)
	$(STOPRPCDB)
	$(STOPZKPROVER)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNRPCDB)
	$(RUNZKPROVER); sleep 5
	trap '$(STOPSTATEDB) && $(STOPPOOLDB) && $(STOPRPCDB) && $(STOPZKPROVER)' EXIT; go test -short -race -p 1 ./...

.PHONY: test-full
test-full: build-docker compile-scs ## Runs all tests checking race conditions
	export CONFIG_MODE="test"
	$(STOPSTATEDB)
	$(STOPPOOLDB)
	$(STOPRPCDB)
	$(STOPZKPROVER)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNRPCDB)
	$(RUNZKPROVER); sleep 5
	$(RUNZKPROVERMOCK)
	trap '$(STOPSTATEDB) && $(STOPPOOLDB) && $(STOPRPCDB) && $(STOPZKPROVER) && $(STOPZKPROVERMOCK)' EXIT; MallocNanoZone=0 go test -race -v -p 1 -timeout 1200s `go list ./... | grep -v \/ci\/e2e-group`

.PHONY: test-full-non-e2e
test-full-non-e2e: build-docker compile-scs ## Runs non-e2e tests checking race conditions
	export CONFIG_MODE="test"	
	$(STOPSTATEDB)
	$(STOPPOOLDB)
	$(STOPRPCDB)
	$(STOPZKPROVER)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNRPCDB)
	$(RUNZKPROVER); sleep 5
	$(RUNZKPROVERMOCK)
	sleep 2
	$(RUNL1NETWORK)
	sleep 15
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOPSTATEDB) && $(STOPPOOLDB) && $(STOPRPCDB) && $(STOPZKPROVER) && $(STOPZKPROVERMOCK) && $(STOPNETWORK)' EXIT; MallocNanoZone=0 go test -short -race -p 1 -timeout 60s ./...

.PHONY: test-e2e-group-1
test-e2e-group-1: build-docker compile-scs ## Runs group 1 e2e tests checking race conditions
	export CONFIG_MODE="test"	
	$(STOPSTATEDB)
	$(STOPPOOLDB)
	$(STOPRPCDB)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNRPCDB); sleep 5
	$(RUNZKPROVER)
	trap '$(STOPSTATEDB) && $(STOPPOOLDB) && $(STOPRPCDB) && $(STOPZKPROVER)' EXIT; MallocNanoZone=0 go test -race -v -p 1 -timeout 600s ./ci/e2e-group1/...

.PHONY: test-e2e-group-2
test-e2e-group-2: build-docker compile-scs ## Runs group 2 e2e tests checking race conditions
	export CONFIG_MODE="test"	
	$(STOPSTATEDB)
	$(STOPPOOLDB)
	$(STOPRPCDB)
	$(STOPZKPROVER)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNRPCDB); sleep 5
	${RUNL1NETWORK}
	CONFIG_MODE="test" $(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOPSTATEDB) && $(STOPPOOLDB) && $(STOPRPCDB) && $(STOPZKPROVER)' EXIT; MallocNanoZone=0 go test -race -v -p 1 -timeout 600s ./ci/e2e-group2/...

.PHONY: test-e2e-group-3
test-e2e-group-3: build-docker compile-scs ## Runs group 3 e2e tests checking race conditions
	export CONFIG_MODE="test"	
	$(STOPSTATEDB)
	$(STOPPOOLDB)
	$(STOPRPCDB)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNRPCDB); sleep 5
	$(RUNZKPROVER); sleep 2
	trap '$(STOPSTATEDB) && $(STOPPOOLDB) && $(STOPRPCDB)' EXIT; MallocNanoZone=0 go test -race -v -p 1 -timeout 600s ./ci/e2e-group3/...

.PHONY: install-linter
install-linter: ## Installs the linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.46.2

.PHONY: lint
lint: ## Runs the linter
	$(LINT)

.PHONY: check
check: stop lint build build-docker test-full-non-e2e test-e2e-group-2 ## lint, build and essential tests

.PHONY: validate
validate: lint build test-full ## lint, build, unit and e2e tests

.PHONY: run-db
run-db: ## Runs the node database
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNRPCDB)

.PHONY: stop-db
stop-db: ## Stops the node database
	$(STOPRPCDB)
	$(STOPPOOLDB)
	$(STOPSTATEDB)

.PHONY: run-node
run-node: ## Runs the node
	$(RUNSYNC)
	$(RUNSEQUENCER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)

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

.PHONY: run-zkprover
run-zkprover: ## Runs zkprover
	$(RUNZKPROVER)

.PHONY: stop-zkprover
stop-zkprover: ## Stops zkprover
	$(STOPZKPROVER)

.PHONY: run-zkprover-mock
run-zkprover-mock: ## Runs zkprover-mock
	$(RUNZKPROVERMOCK)

.PHONY: stop-zkprover-mock
stop-zkprover-mock: ## Stops zkprover-mock
	$(STOPZKPROVERMOCK)

.PHONY: run-explorer
run-explorer: ## Runs the explorer
	$(RUNEXPLORERL1DB)
	$(RUNEXPLORERL2DB)
	$(RUNEXPLORERJSONRPC)
	$(RUNEXPLORERL1)
	$(RUNEXPLORERL2)

.PHONY: stop-explorer
stop-explorer: ## Stops the explorer
	$(STOPEXPLORERL2)
	$(STOPEXPLORERL1)
	$(STOPEXPLORERRPC)
	$(STOPEXPLORERL2DB)
	$(STOPEXPLORERL1DB)

.PHONY: run-explorer-db
run-explorer-db: ## Runs the explorer database
	$(RUNEXPLORERL1DB)
	$(RUNEXPLORERL2DB)

.PHONY: stop-explorer-db
stop-explorer-db: ## Stops the explorer database
	$(STOPEXPLORERL2DB)
	$(STOPEXPLORERL1DB)

.PHONY: run
run: ## Runs all the services
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNRPCDB)
	$(RUNL1NETWORK)
	sleep 2
	$(RUNZKPROVER)
	sleep 5
	$(RUNSEQUENCER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)
	$(RUNSYNC)

.PHONY: run-broadcast
run-broadcast: ## Runs the broadcast service
	$(RUNBROADCAST)

run-seq:
	$(RUNSEQUENCER)

.PHONY: stop-broadcast
stop-broadcast: ## Stops the broadcast service
	$(STOPBROADCAST)

.PHONY: run-permissionless
run-permissionless: ## Runs the permissionless node
	$(RUNPERMISSIONLESSDB)
	$(RUNPERMISSIONLESSNODE)

.PHONY: stop-permissionless
stop-permissionless: ## Stops the permissionless node
	$(STOPPERMISSIONLESSNODE)
	$(STOPPERMISSIONLESSDB)

.PHONY: run-approve-matic
run-approve-matic: ## Runs approve in node container
	$(RUNAPPROVE)

.PHONY: stop-approve-matic
stop-approve-matic: ## Stops approve in node container
	$(STOPAPPROVE)

#.PHONY: init-network
#init-network: ## Initializes the network
#	go run ./scripts/init_network/main.go .

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

	mockery --name=txManager --dir=sequencer --output=sequencer/mocks --outpkg=mocks --structname=TxmanagerMock --filename=mock_txmanager.go
	mockery --name=etherman --dir=sequencer --output=sequencer/mocks --outpkg=mocks --structname=EthermanMock --filename=mock_etherman.go
	mockery --name=etherman --dir=sequencer/profitabilitychecker --output=sequencer/profitabilitychecker/mocks --outpkg=mocks --structname=EthermanMock --filename=mock_etherman.go
	mockery --name=stateInterface --dir=sequencer/broadcast --output=sequencer/broadcast/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go

	mockery --name=ethermanInterface --dir=synchronizer --output=synchronizer --outpkg=synchronizer --structname=ethermanMock --filename=mock_etherman.go
	mockery --name=stateInterface --dir=synchronizer --output=synchronizer --outpkg=synchronizer --structname=stateMock --filename=mock_state.go
	mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=synchronizer --outpkg=synchronizer --structname=dbTxMock --filename=mock_dbtx.go

	## mocks for the aggregator tests
	mockery --name=stateInterface --dir=aggregator --output=aggregator/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go
	mockery --name=proverClientInterface --dir=aggregator --output=aggregator/mocks --outpkg=mocks --structname=ProverClientMock --filename=mock_proverclient.go
	mockery --name=etherman --dir=aggregator --output=aggregator/mocks --outpkg=mocks --structname=Etherman --filename=mock_etherman.go
	mockery --name=ethTxManager --dir=aggregator --output=aggregator/mocks --outpkg=mocks --structname=EthTxManager --filename=mock_ethtxmanager.go

.PHONY: generate-code-from-proto
generate-code-from-proto: ## Generates code from proto files
	cd proto/src/proto/statedb/v1 && protoc --proto_path=. --proto_path=../../../../include --go_out=../../../../../merkletree/pb --go-grpc_out=../../../../../merkletree/pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative statedb.proto
	cd proto/src/proto/zkprover/v1 && protoc --proto_path=. --go_out=../../../../../proverclient/pb --go-grpc_out=../../../../../proverclient/pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative zk_prover.proto
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
