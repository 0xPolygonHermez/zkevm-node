# Check dependencies
# Check for Go
.PHONY: check-go
check-go:
	@which go > /dev/null || (echo "Error: Go is not installed" && exit 1)

# Targets that require the checks
generate-file: check-go

.PHONY: generate-file
generate-file: ## Runs the tool to populate the binary file
	go run main.go generate -cfg config/tool.config.toml

rebuild: ## Runs the tool to rebuild state roots from a given block
	go run main.go rebuild -cfg config/tool.config.toml

decode-entry: ## Runs the tool to decode a given entry
	go run main.go decode -cfg config/tool.config.toml -entry 1234

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
