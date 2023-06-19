DOCKERCOMPOSE := docker-compose -f docker-compose.yml
DOCKERCOMPOSEAPPSEQ := zkevm-sequencer
DOCKERCOMPOSEAPPSEQSENDER := zkevm-sequence-sender
DOCKERCOMPOSEAPPL2GASP := zkevm-l2gaspricer
DOCKERCOMPOSEAPPAGG := zkevm-aggregator
DOCKERCOMPOSEAPPRPC := zkevm-json-rpc
DOCKERCOMPOSEAPPSYNC := zkevm-sync
DOCKERCOMPOSEAPPETHTXMANAGER := zkevm-eth-tx-manager
DOCKERCOMPOSESTATEDB := zkevm-state-db
DOCKERCOMPOSEPOOLDB := zkevm-pool-db
DOCKERCOMPOSEEVENTDB := zkevm-event-db
DOCKERCOMPOSENETWORK := zkevm-mock-l1-network
DOCKERCOMPOSEEXPLORERL1 := zkevm-explorer-l1
DOCKERCOMPOSEEXPLORERL1DB := zkevm-explorer-l1-db
DOCKERCOMPOSEEXPLORERL2 := zkevm-explorer-l2
DOCKERCOMPOSEEXPLORERL2DB := zkevm-explorer-l2-db
DOCKERCOMPOSEEXPLORERRPC := zkevm-explorer-json-rpc
DOCKERCOMPOSEZKPROVER := zkevm-prover
DOCKERCOMPOSEPERMISSIONLESSDB := zkevm-permissionless-db
DOCKERCOMPOSEPERMISSIONLESSNODE := zkevm-permissionless-node
DOCKERCOMPOSEPERMISSIONLESSZKPROVER := zkevm-permissionless-prover
DOCKERCOMPOSENODEAPPROVE := zkevm-approve
DOCKERCOMPOSEMETRICS := zkevm-metrics
DOCKERCOMPOSEGRAFANA := grafana

RUNSTATEDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSESTATEDB)
RUNPOOLDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPOOLDB)
RUNEVENTDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEVENTDB)
RUNSEQUENCER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQ)
RUNSEQUENCESENDER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSEQSENDER)
RUNL2GASPRICER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPL2GASP)
RUNAGGREGATOR := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPAGG)
RUNJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPRPC)
RUNSYNC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPSYNC)
RUNETHTXMANAGER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEAPPETHTXMANAGER)
RUNGRAFANA := DOCKERGID=`stat -c '%g' /var/run/docker.sock` $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEGRAFANA)

RUNL1NETWORK := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENETWORK)
RUNEXPLORERL1 := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL1)
RUNEXPLORERL1DB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL1DB)
RUNEXPLORERL2 := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL2)
RUNEXPLORERL2DB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERL2DB)
RUNEXPLORERJSONRPC := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEEXPLORERRPC)
RUNZKPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEZKPROVER)

RUNPERMISSIONLESSDB := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPERMISSIONLESSDB)
RUNPERMISSIONLESSNODE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPERMISSIONLESSNODE)
RUNPERMISSIONLESSZKPROVER := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEPERMISSIONLESSZKPROVER)

RUNAPPROVE := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSENODEAPPROVE)

RUNMETRICS := $(DOCKERCOMPOSE) up -d $(DOCKERCOMPOSEMETRICS)

RUN := $(DOCKERCOMPOSE) up -d

STOPSTATEDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSESTATEDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSESTATEDB)
STOPPOOLDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPOOLDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPOOLDB)
STOPEVENTDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEVENTDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEVENTDB)
STOPSEQUENCER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQ) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQ)
STOPSEQUENCESENDER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSEQSENDER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSEQSENDER)
STOPL2GASPRICER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPL2GASP) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPL2GASP)
STOPAGGREGATOR := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPAGG) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPAGG)
STOPJSONRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPRPC)
STOPSYNC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPSYNC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPSYNC)
STOPETHTXMANAGER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEAPPETHTXMANAGER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEAPPETHTXMANAGER)
STOPGRAFANA := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEGRAFANA) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEGRAFANA)

STOPNETWORK := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENETWORK) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENETWORK)
STOPEXPLORERL1 := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL1) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL1)
STOPEXPLORERL1DB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL1DB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL1DB)
STOPEXPLORERL2 := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL2) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL2)
STOPEXPLORERL2DB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERL2DB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERL2DB)
STOPEXPLORERJSONRPC := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEEXPLORERRPC) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEEXPLORERRPC)
STOPZKPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEZKPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEZKPROVER)

STOPPERMISSIONLESSDB := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPERMISSIONLESSDB) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPERMISSIONLESSDB)
STOPPERMISSIONLESSNODE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPERMISSIONLESSNODE) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPERMISSIONLESSNODE)
STOPPERMISSIONLESSZKPROVER := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEPERMISSIONLESSZKPROVER) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEPERMISSIONLESSZKPROVER)

STOPAPPROVE := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSENODEAPPROVE) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSENODEAPPROVE)

STOPMETRICS := $(DOCKERCOMPOSE) stop $(DOCKERCOMPOSEMETRICS) && $(DOCKERCOMPOSE) rm -f $(DOCKERCOMPOSEMETRICS)

STOP := $(DOCKERCOMPOSE) down --remove-orphans

.PHONY: test-full-non-e2e
test-full-non-e2e: stop ## Runs non-e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 2
	$(RUNZKPROVER)
	sleep 7
	$(RUNL1NETWORK)
	sleep 15
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -short -race -p 1 -timeout 60s ../...

.PHONY: test-e2e-group-1
test-e2e-group-1: stop ## Runs group 1 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 600s ../ci/e2e-group1/...

.PHONY: test-e2e-group-2
test-e2e-group-2: stop ## Runs group 2 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 600s ../ci/e2e-group2/...

.PHONY: test-e2e-group-3
test-e2e-group-3: stop ## Runs group 3 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 600s ../ci/e2e-group3/...

.PHONY: test-e2e-group-4
test-e2e-group-4: stop ## Runs group 4 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 600s ../ci/e2e-group4/...

.PHONY: test-e2e-group-5
test-e2e-group-5: stop ## Runs group 5 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 1200s ../ci/e2e-group5/...

.PHONY: test-e2e-group-6
test-e2e-group-6: stop ## Runs group 6 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 600s ../ci/e2e-group6/...

.PHONY: test-e2e-group-7
test-e2e-group-7: stop ## Runs group 7 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 600s ../ci/e2e-group7/...

.PHONY: test-e2e-group-8
test-e2e-group-8: stop ## Runs group 8 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 1200s ../ci/e2e-group8/...


.PHONY: test-e2e-group-9
test-e2e-group-9: stop ## Runs group 9 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 2000s ../ci/e2e-group9/...


.PHONY: test-e2e-group-10
test-e2e-group-10: stop ## Runs group 10 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 2000s ../ci/e2e-group10/...


.PHONY: test-e2e-group-11
test-e2e-group-11: stop ## Runs group 11 e2e tests checking race conditions
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	trap '$(STOP)' EXIT; MallocNanoZone=0 go test -count=1 -race -v -p 1 -timeout 2000s ../ci/e2e-group11/...

.PHONY: benchmark-sequencer-eth-transfers
benchmark-sequencer-eth-transfers: stop
	$(RUNL1NETWORK)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNSYNC)
	sleep 2
	$(RUNL2GASPRICER)
	$(RUNJSONRPC)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	@ cd benchmarks/sequencer/eth-transfers ; \
 	mkdir -p results ; \
 	touch ./results/out.dat ; \
	go test -bench=. -timeout=600m | tee ./results/out.dat ;

.PHONY: benchmark-sequencer-erc20-transfers
benchmark-sequencer-erc20-transfers: stop
	$(RUNL1NETWORK)
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	sleep 5
	$(RUNZKPROVER)
	$(RUNSYNC)
	sleep 2
	$(RUNL2GASPRICER)
	$(RUNJSONRPC)
	docker ps -a
	docker logs $(DOCKERCOMPOSEZKPROVER)
	@ cd benchmarks/sequencer/erc20-transfers ; \
 	mkdir -p results ; \
 	touch ./results/out.dat ; \
	go test -bench=. -timeout=600m | tee ./results/out.dat ;

.PHONY: run-db
run-db: ## Runs the node database
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)

.PHONY: stop-db
stop-db: ## Stops the node database
	$(STOPEVENTDB)
	$(STOPPOOLDB)
	$(STOPSTATEDB)

.PHONY: run-node
run-node: ## Runs the node
	$(RUNETHTXMANAGER)
	$(RUNSYNC)
	sleep 2
	$(RUNSEQUENCER)
	$(RUNSEQUENCESENDER)
	$(RUNL2GASPRICER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)

.PHONY: stop-node
stop-node: ## Stops the node
	$(STOPSEQUENCER)
	$(STOPSEQUENCESENDER)
	$(STOPJSONRPC)
	$(STOPL2GASPRICER)
	$(STOPAGGREGATOR)
	$(STOPSYNC)
	$(STOPETHTXMANAGER)

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

.PHONY: run-l1-explorer
run-l1-explorer: ## Runs L1 blockscan explorer 
	$(RUNEXPLORERL1DB)
	$(RUNEXPLORERL1)

.PHONY: run-l2-explorer
run-l2-explorer: ## Runs L2 blockscan explorer
	$(RUNEXPLORERL2DB)
	$(RUNEXPLORERJSONRPC)
	$(RUNEXPLORERL2)

.PHONY: run-l2-explorer-json-rpc
run-l2-explorer-json-rpc: ## Runs L2 explorer json rpc
	$(RUNEXPLORERJSONRPC)

.PHONY: stop-l2-explorer-json-rpc
stop-l2-explorer-json-rpc: ## Stops L2 explorer json rpc
	$(STOPEXPLORERJSONRPC)

.PHONY: run-explorer
run-explorer: run-l1-explorer run-l2-explorer ## Runs both L1 and L2 explorers

.PHONY: stop-explorer
stop-explorer: ## Stops the explorer
	$(STOPEXPLORERL2)
	$(STOPEXPLORERL1)
	$(STOPEXPLORERJSONRPC)
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

.PHONY: run-seqsender
run-seqsender: ## runs the sequencer sender
	$(RUNSEQUENCESENDER)

.PHONY: stop-seqsender
stop-seqsender: ## stops the sequencer sender
	$(STOPSEQUENCESENDER)
	
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

.PHONY: run-l2gaspricer
run-l2gaspricer: ## runs the L2 Gas Price component
	$(RUNL2GASPRICER)

.PHONY: stop-l2gaspricer
stop-l2gaspricer: ## stops the L2 Gas Price component
	$(STOPL2GASPRICER)

.PHONY: run-eth-tx-manager
run-eth-tx-manager: ## Runs the eth tx manager service
	$(RUNETHTXMANAGER)

.PHONY: stop-eth-tx-manager
stop-eth-tx-manager: ## Stops the eth tx manager service
	$(STOPETHTXMANAGER)
	
.PHONY: run-agg
run-agg: ## Runs the aggregator service
	$(RUNAGGREGATOR)

.PHONY: stop-agg
stop-agg: ## Stops the aggregator service
	$(STOPAGGREGATOR)

.PHONY: run-grafana
run-grafana: ## Runs the grafana service
	$(RUNGRAFANA)

.PHONY: stop-grafana
stop-grafana: ## Stops the grafana service
	$(STOPGRAFANA)

.PHONY: run-permissionless
run-permissionless: run-node ## Runs the trusted and permissionless node
	$(RUNPERMISSIONLESSDB)
	sleep 3
	$(RUNPERMISSIONLESSZKPROVER)
	$(RUNPERMISSIONLESSNODE)

.PHONY: stop-permissionless
stop-permissionless: stop-node## Stops the permissionless node
	$(STOPPERMISSIONLESSNODE)
	$(STOPPERMISSIONLESSZKPROVER)
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
	$(RUNEVENTDB)
	$(RUNL1NETWORK)
	sleep 1
	$(RUNZKPROVER)
	$(RUNAPPROVE)
	sleep 3
	$(RUNSYNC)
	sleep 2
	$(RUNETHTXMANAGER)
	$(RUNSEQUENCER)
	$(RUNSEQUENCESENDER)
	$(RUNL2GASPRICER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)

.PHONY: stop
stop: ## Stops all services
	$(STOP)

.PHONY: ship
ship: ## Builds docker images and run them
	cd .. && make build-docker && cd ./test && make run

.PHONY: reship
reship: stop ship ## Executes `make stop` and `make sail` commands

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

.PHONY: install-mockery
install-mockery: ## Installs mockery with the correct version to generate the mocks
	go install github.com/vektra/mockery/v2@v2.22.1

.PHONY: generate-mocks
generate-mocks: ## Generates mocks for the tests, using mockery tool
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=storageInterface --dir=../jsonrpc --output=../jsonrpc --outpkg=jsonrpc --inpackage --structname=storageMock --filename=mock_storage.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=PoolInterface --dir=../jsonrpc/types --output=../jsonrpc/mocks --outpkg=mocks --structname=PoolMock --filename=mock_pool.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=StateInterface --dir=../jsonrpc/types --output=../jsonrpc/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../jsonrpc/mocks --outpkg=mocks --structname=DBTxMock --filename=mock_dbtx.go

	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=workerInterface --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage --structname=WorkerMock --filename=mock_worker.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage --structname=StateMock --filename=mock_state.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=txPool --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage  --structname=PoolMock --filename=mock_pool.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../sequencer --outpkg=sequencer --structname=DbTxMock --filename=mock_dbtx.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=dbManagerInterface --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage --structname=DbManagerMock --filename=mock_db_manager.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=etherman --dir=../sequencer --output=../sequencer --outpkg=sequencer --inpackage --structname=EthermanMock --filename=mock_etherman.go

	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethermanInterface --dir=../synchronizer --output=../synchronizer --outpkg=synchronizer --structname=ethermanMock --filename=mock_etherman.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface --dir=../synchronizer --output=../synchronizer --outpkg=synchronizer --structname=stateMock --filename=mock_state.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethTxManager --dir=../synchronizer --output=../synchronizer --outpkg=synchronizer --structname=ethTxManagerMock --filename=mock_ethtxmanager.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=poolInterface --dir=../synchronizer --output=../synchronizer --outpkg=synchronizer --structname=poolMock --filename=mock_pool.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=zkEVMClientInterface --dir=../synchronizer --output=../synchronizer --outpkg=synchronizer --structname=zkEVMClientMock --filename=mock_zkevmclient.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../synchronizer --outpkg=synchronizer --structname=dbTxMock --filename=mock_dbtx.go

	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=GasPricer --srcpkg=github.com/ethereum/go-ethereum --output=../etherman --outpkg=etherman --structname=etherscanMock --filename=mock_etherscan.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=GasPricer --srcpkg=github.com/ethereum/go-ethereum --output=../etherman --outpkg=etherman --structname=ethGasStationMock --filename=mock_ethgasstation.go

	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethermanInterface --dir=../ethtxmanager --output=../ethtxmanager --outpkg=ethtxmanager --structname=ethermanMock --filename=mock_etherman_test.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface --dir=../ethtxmanager --output=../ethtxmanager --outpkg=ethtxmanager --structname=stateMock --filename=mock_state_test.go

	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=poolInterface --dir=../gasprice --output=../gasprice --outpkg=gasprice --structname=poolMock --filename=mock_pool.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethermanInterface --dir=../gasprice --output=../gasprice --outpkg=gasprice --structname=ethermanMock --filename=mock_etherman.go

	## mocks for the aggregator tests
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=stateInterface --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=StateMock --filename=mock_state.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=proverInterface --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=ProverMock --filename=mock_prover.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=etherman --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=Etherman --filename=mock_etherman.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=ethTxManager --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=EthTxManager --filename=mock_ethtxmanager.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=aggregatorTxProfitabilityChecker --dir=../aggregator --output=../aggregator/mocks --outpkg=mocks --structname=ProfitabilityCheckerMock --filename=mock_profitabilitychecker.go
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/mockery --name=Tx --srcpkg=github.com/jackc/pgx/v4 --output=../aggregator/mocks --outpkg=mocks --structname=DbTxMock --filename=mock_dbtx.go

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
