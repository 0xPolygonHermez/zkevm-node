.PHONY: run-rpc
run-rpc: ## Runs all the services need to run a local zkEMV RPC node
	docker-compose -f docker-compose.yml up -d zkevm-state-db zkevm-pool-db zkevm-rpc-db
	sleep 2
	docker-compose -f docker-compose.yml up -d zkevm-prover
	sleep 5
	docker-compose -f docker-compose.yml up -d zkevm-sync zkevm-rpc

.PHONY: stop
stop: ## Stops all services
	docker-compose down

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
