PKG := github.com/d3mondev/puredns/v2
PKG_LIST := $(shell go list ./... | grep -v /vendor/)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')
REVISION := $(shell git rev-parse --short HEAD)

.SILENT: ;
.PHONY: all

all: build

lint: ## Lint the files
	golint -set_exit_status $(PKG_LIST)
	staticcheck ./...

test: ## Run unit tests
	go fmt $(PKG_LIST)
	go vet $(PKG_LIST)
	go test -race -timeout 30s -cover -count 1 $(PKG_LIST)

msan: ## Run memory sanitizer
	go test -msan $(PKG_LIST)

build: ## Build the binary file
	go build -trimpath -ldflags="-s -w"

cover: ## Code coverage
	go test -coverprofile=cover.out $(PKG_LIST)

clean: ## Remove previous build
	rm -f cover.out
	go clean

help: ## Display this help screen
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
