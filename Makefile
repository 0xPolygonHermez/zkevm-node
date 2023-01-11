include version.mk

ARCH := $(shell arch)

ifeq ($(ARCH),x86_64)
	ARCH = amd64
else 
	ifeq ($(ARCH),aarch64)
		ARCH = arm64
	endif
endif
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist
GOENVVARS := GOBIN=$(GOBIN) CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH)
GOBINARY := zkevm-node
GOCMD := $(GOBASE)/cmd

LDFLAGS += -X 'github.com/0xPolygonHermez/zkevm-node.Version=$(VERSION)'
LDFLAGS += -X 'github.com/0xPolygonHermez/zkevm-node.GitRev=$(GITREV)'
LDFLAGS += -X 'github.com/0xPolygonHermez/zkevm-node.GitBranch=$(GITBRANCH)'
LDFLAGS += -X 'github.com/0xPolygonHermez/zkevm-node.BuildDate=$(DATE)'

.PHONY: build
build: ## Builds the binary locally into ./dist
	$(GOENVVARS) go build -ldflags "all=$(LDFLAGS)" -o $(GOBIN)/$(GOBINARY) $(GOCMD)

.PHONY: build-docker
build-docker: ## Builds a docker image with the node binary
	docker build -t zkevm-node -f ./Dockerfile .

.PHONY: build-docker-nc
build-docker-nc: ## Builds a docker image with the node binary - but without build cache
	docker build --no-cache=true -t zkevm-node -f ./Dockerfile .

.PHONY: run-rpc
run-rpc: ## Runs all the services need to run a local zkEMV RPC node
	docker-compose up -d zkevm-state-db zkevm-pool-db
	sleep 2
	docker-compose up -d zkevm-prover
	sleep 5
	docker-compose up -d zkevm-sync
	sleep 2
	docker-compose up -d zkevm-rpc

.PHONY: stop
stop: ## Stops all services
	docker-compose down

.PHONY: install-linter
install-linter: ## Installs the linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.46.2

.PHONY: lint
lint: ## Runs the linter
	$$(go env GOPATH)/bin/golangci-lint run

.PHONY: update-external-dependencies
update-external-dependencies: ## Updates external dependencies like images, test vectors or proto files
	go run ./scripts/cmd/... updatedeps

.PHONY: install-git-hooks
install-git-hooks: ## Moves hook files to the .git/hooks directory
	cp .github/hooks/* .git/hooks

.PHONY: generate-code-from-proto
generate-code-from-proto: ## Generates code from proto files
	cd proto/src/proto/statedb/v1 && protoc --proto_path=. --proto_path=../../../../include --go_out=../../../../../merkletree/pb --go-grpc_out=../../../../../merkletree/pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative statedb.proto
	cd proto/src/proto/executor/v1 && protoc --proto_path=. --go_out=../../../../../state/runtime/executor/pb --go-grpc_out=../../../../../state/runtime/executor/pb --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative executor.proto
	cd proto/src/proto/broadcast/v1 && protoc --proto_path=. --proto_path=../../../../include --go_out=../../../../../sequencer/broadcast/pb --go-grpc_out=../../../../../sequencer/broadcast/pb --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative broadcast.proto
	cd proto/src/proto/aggregator/v1 && protoc --proto_path=. --proto_path=../../../../include --go_out=../../../../../aggregator/pb --go-grpc_out=../../../../../aggregator/pb --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative aggregator.proto

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
