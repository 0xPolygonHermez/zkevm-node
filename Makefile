STARTDB = docker run --rm --name hermez-polygon-db -p 5432:5432 -e POSTGRES_DB="polygon-hermez" -e POSTGRES_USER="hermez" -e POSTGRES_PASSWORD="polygon" -d postgres
STOPDB = docker stop hermez-polygon-db

.PHONY: build
build: lint test ## Build the binary
	go build -o ./dist/hezcore ./cmd/main.go

.PHONY: test
test: ## runs tests
	$(STOPDB) || true
	$(STARTDB)
	go test -race -p 1 ./...
	$(STOPDB)

.PHONY: install-linter
install-linter: ## install linter
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.30.0

.PHONY: lint
lint: ## runs linter
	$$(go env GOPATH)/bin/golangci-lint run --timeout=5m -E whitespace -E gosec -E gci -E misspell -E gomnd -E gofmt -E goimports -E golint --exclude-use-default=false --max-same-issues 0

.PHONY: start-db
start-db: ## starts a docker container to run the db instance
	$(STARTDB)

.PHONY: stop-db
stop-db: ## stops the docker container running the db instance
	$(STOPDB)

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
		