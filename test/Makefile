DOCKERCOMPOSE := docker-compose -f docker-compose.yml
DOCKERCOMPOSEAPPSEQ := zkevm-sequencer
DOCKERCOMPOSEAPPAGG := zkevm-aggregator
DOCKERCOMPOSEAPPRPC := zkevm-json-rpc
DOCKERCOMPOSEAPPSYNC := zkevm-sync
DOCKERCOMPOSEAPPBROADCAST := zkevm-broadcast
DOCKERCOMPOSEAPPETHTXMANAGER := zkevm-eth-tx-manager
DOCKERCOMPOSESTATEDB := zkevm-state-db
DOCKERCOMPOSEPOOLDB := zkevm-pool-db
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
DOCKERCOMPOSEMETRICS := zkevm-metrics
DOCKERCOMPOSEGRAFANA := grafana

RUNSTATEDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSESTATEDB)
RUNPOOLDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPOOLDB)
RUNSEQUENCER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQ)
RUNAGGREGATOR := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPAGG)
RUNJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPRPC)
RUNSYNC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSYNC)
RUNBROADCAST := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPBROADCAST)
RUNETHTXMANAGER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPETHTXMANAGER)
RUNGRAFANA := DOCKERGID=`stat -c '%g' /var/run/docker.sock` $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEGRAFANA)

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

RUNMETRICS := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEMETRICS)

RUN := $(DOCKERCOMPOSE) up -d

STOPSTATEDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSESTATEDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSESTATEDB)
STOPPOOLDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPOOLDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPOOLDB)
STOPSEQUENCER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQ) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQ)
STOPAGGREGATOR := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPAGG) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPAGG)
STOPJSONRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPRPC)
STOPSYNC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSYNC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSYNC)
STOPBROADCAST := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPBROADCAST) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPBROADCAST)
STOPETHTXMANAGER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPETHTXMANAGER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPETHTXMANAGER)
STOPGRAFANA := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEGRAFANA) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEGRAFANA)

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

STOPMETRICS := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEMETRICS) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEMETRICS)

STOP := $(DOCKERCOMPOSE) down --remove-orphans

.PHONY: test-full-non-e2e
test-full-non-e2e: stop ## Runs non-e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNZKPROVER)
	sleep 5
	$(RUNZKPROVERMOCK)
	sleep 2
	$(RUNL1NETWORK)
	sleep 15
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -short -race -p 1 -timeout 60s ../...

.PHONY: test-e2e-group-1
test-e2e-group-1: stop ## Runs group 1 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNETHTXMANAGER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -race -v -p 1 -timeout 600s ../ci/e2e-group1/...

.PHONY: test-e2e-group-2
test-e2e-group-2: stop ## Runs group 2 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNETHTXMANAGER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -race -v -p 1 -timeout 600s ../ci/e2e-group2/...

.PHONY: test-e2e-group-3
test-e2e-group-3: stop ## Runs group 3 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNETHTXMANAGER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -race -v -p 1 -timeout 600s ../ci/e2e-group3/...

.PHONY: test-e2e-group-4
test-e2e-group-4: stop ## Runs group 4 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNETHTXMANAGER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -race -v -p 1 -timeout 600s ../ci/e2e-group4/...

.PHONY: benchmark-sequencer
benchmark-sequencer: stop
	$(RUNL1NETWORK)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNJSONRPC)
	$(RUNSYNC)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	@ cd benchmarks/sequencer ; \
 	mkdir -p results ; \
 	touch ./results/out.dat ; \
	go test -bench=. -timeout=600m | tee ./results/out.dat ;

.PHONY: run-db
run-db: ## Runs the node database
	$(RUNSTATEDB)
	$(RUNPOOLDB)

.PHONY: stop-db
stop-db: ## Stops the node database
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

.PHONY: run-seq
run-seq: ## runs the sequencer
	$(RUNSEQUENCER)

.PHONY: stop-seq
stop-seq: ## stops the sequencer
	$(STOPSEQUENCER)
	
.PHONY: run-agg
run-agg: ## runs the aggregator
	$(RUNAGGREGATOR)

.PHONY: stop-agg
stop-agg: ## stops the aggregator
	$(STOPAGGREGATOR)

.PHONY: run-sync
run-sync: ## runs the synchronizer
	$(RUNSYNC)

.PHONY: stop-sync
stop-sync: ## stops the synchronizer
	$(STOPSYNC)

.PHONY: run-json-rpc
run-json-rpc: ## runs the JSON-RPC
	$(RUNJSONRPC)

.PHONY: stop-json-rpc
stop-json-rpc: ## stops the JSON-RPC
	$(STOPJSONRPC)

.PHONY: run-broadcast
run-broadcast: ## Runs the broadcast service
	$(RUNBROADCAST)

.PHONY: stop-broadcast
stop-broadcast: ## Stops the broadcast service
	$(STOPBROADCAST)

.PHONY: run-eth-tx-manager
run-eth-tx-manager: ## Runs the eth tx manager service
	$(RUNETHTXMANAGER)

.PHONY: stop-eth-tx-manager
stop-eth-tx-manager: ## Stops the eth tx manager service
	$(STOPETHTXMANAGER)

.PHONY: run-grafana
run-grafana: ## Runs the grafana service
	$(RUNGRAFANA)

.PHONY: stop-grafana
stop-grafana: ## Stops the grafana service
	$(STOPGRAFANA)

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

.PHONY: run
run: ## Runs a full node
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNL1NETWORK)
	sleep 1
	$(RUNZKPROVER)
	sleep 3
	$(RUNSEQUENCER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)
	$(RUNSYNC)
	$(RUNETHTXMANAGER)

.PHONY: stop
stop: ## Stops all services
	$(STOP)

.PHONY: sail
sail: ## Builds docker images and run them
	cd .. && make build-docker && cd ./test && make run

.PHONY: resail
resail: stop sail ## Executes `make stop` and `make sail` commands

.PHONY: restart
restart: stop run ## Executes `make stop` and `make run` commands

.PHONY: run-metrics
run-metrics: ## Runs the metrics container
	$(RUNMETRICS)

.PHONY: stop-metrics
stop-metrics: ## Stops the metrics container
	$(STOPMETRICS)

.PHONY: init-network
init-network: ## Initializes the network
	go run ./scripts/init_network/main.go .

.PHONY: deploy-sc
deploy-sc: ## deploys some examples of transactions and smart contracts
	go run ./scripts/deploy_sc/main.go .

.PHONY: send-transfers
send-transfers: ## sends some ETH transfers txs to test the network
	go run ./scripts/send_transfers/main.go .

.PHONY: deploy-uniswap
deploy-uniswap: ## deploy the uniswap environment to the network
	go run ./scripts/uniswap/main.go .

.PHONY: run-db-scripts
run-db-scripts: ## Executes scripts on the db after it has been initialized, potentially using info from the environment
	./scripts/postgres/run.sh

.PHONY: generate-mocks
generate-mocks: ## Generates mocks for the tests, using mockery tool
	mockery --name=storageInterface --dir=../jsonrpc --output=../jsonrpc --outpkg=jsonrpc --inpackage --structname=storageMock --filename=mock_storage_test.go
	mockery --name=jsonRPCTxPool --dir=../jsonrpc --output=../jsonrpc --outpkg=jsonrpc --inpackage --structname=poolMock --filename=mock_pool_test.go
	mockery --name=gasPriceEstimator --dir=../jsonrpc --output=../jsonrpc --outpkg=jsonrpc --inpackage --structname=gasPriceEstimatorMock --filename=mock_gasPriceEstimator_test.go
	mockery --name=stateInterface --dir=../jsonrpc --output=../jsonrpc --outpkg=jsonrpc --inpackage --structname=stateMock --filename=mock_state_test.go
	mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../jsonrpc --outpkg=jsonrpc --structname=dbTxMock --filename=mock_dbtx_test.go

	mockery --name=ethTxManager --dir=../sequencer --output=../sequencer/mocks --outpkg=mocks --structname=EthTxManager --filename=mock_ethtxmanager.go
	mockery --name=etherman --dir=../sequencer --output=../sequencer/mocks --outpkg=mocks --structname=EthermanMock --filename=mock_etherman.go
	mockery --name=stateInterface --dir=../sequencer --output=../sequencer/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go
	mockery --name=txPool --dir=../sequencer --output=../sequencer/mocks --outpkg=mocks --structname=PoolMock --filename=mock_pool.go
	mockery --name=gasPriceEstimator --dir=../sequencer --output=../sequencer/mocks --outpkg=mocks --structname=GasPriceEstimatorMock --filename=mock_gaspriceestimator.go
	mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../sequencer/mocks --outpkg=mocks --structname=DbTxMock --filename=mock_dbtx.go
	mockery --name=etherman --dir=../sequencer/profitabilitychecker --output=../sequencer/profitabilitychecker/mocks --outpkg=mocks --structname=EthermanMock --filename=mock_etherman.go
	mockery --name=stateInterface --dir=../sequencer/broadcast --output=../sequencer/broadcast/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go

	mockery --name=ethermanInterface --dir=../synchronizer --output=../synchronizer --outpkg=synchronizer --structname=ethermanMock --filename=mock_etherman.go
	mockery --name=stateInterface --dir=../synchronizer --output=../synchronizer --outpkg=synchronizer --structname=stateMock --filename=mock_state.go
	mockery --name=ethTxManager --dir=../synchronizer --output=../synchronizer --outpkg=synchronizer --structname=ethTxManagerMock --filename=mock_ethtxmanager.go
	mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../synchronizer --outpkg=synchronizer --structname=dbTxMock --filename=mock_dbtx.go

	mockery --name=GasPricer --srcpkg=github.com/ethereum/go-ethereum --output=../etherman --outpkg=etherman --structname=etherscanMock --filename=mock_etherscan.go
	mockery --name=GasPricer --srcpkg=github.com/ethereum/go-ethereum --output=../etherman --outpkg=etherman --structname=ethGasStationMock --filename=mock_ethgasstation.go

	mockery --name=ethermanInterface --dir=../ethtxmanager --output=../ethtxmanager --outpkg=ethtxmanager --structname=ethermanMock --filename=mock_etherman_test.go
	mockery --name=stateInterface --dir=../ethtxmanager --output=../ethtxmanager --outpkg=ethtxmanager --structname=stateMock --filename=mock_state_test.go

	## mocks for the aggregator tests
	mockery --name=stateInterface --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go
	mockery --name=proverInterface --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=ProverMock --filename=mock_prover.go
	mockery --name=etherman --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=Etherman --filename=mock_etherman.go
	mockery --name=ethTxManager --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=EthTxManager --filename=mock_ethtxmanager.go
	mockery --name=aggregatorTxProfitabilityChecker --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=ProfitabilityCheckerMock --filename=mock_profitabilitychecker.go

.PHONY: run-benchmarks
run-benchmarks: run-db ## Runs benchmars
	go test -bench=. ./state/tree

.PHONY: compile-scs
compile-scs: ## Compiles smart contracts, configuration in test/contracts/index.yaml
	go run ./scripts/cmd... compilesc --input ./contracts

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
