.PHONY: build
build: lint test ## Build the binary
	go build -o ./dist/hezcore ./cmd/main.go

.PHONY: test
test: ## runs tests
	go test ./... -p 1

.PHONY: install-linter
install-linter: ## install linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.30.0

.PHONY: lint
lint: ## runs linter
	$$(go env GOPATH)/bin/golangci-lint run --timeout=5m -E whitespace -E gosec -E gci -E misspell -E gomnd -E gofmt -E goimports -E golint --exclude-use-default=false --max-same-issues 0

.PHONY: run-db
run-db: ## runs the db instance
	docker run --rm -p 5432:5432 -e POSTGRES_DB="polygon-hermez" -e POSTGRES_USER="hermez" -e POSTGRES_PASSWORD="polygon" -d postgres

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
		