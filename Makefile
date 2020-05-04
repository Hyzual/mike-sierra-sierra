# FROM https://gitlab.com/pantomath-io/demo-tools/-/blob/master/Makefile
PROJECT_NAME := "mike-sierra-sierra"
PKG := "github.com/hyzual/$(PROJECT_NAME)"
PROJECT_EXEC_NAME := "mike"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean test coverage coverhtml lint start npm-dep \
stylelint-ci prettier-ci

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unit tests
	@go test -short ${PKG_LIST}

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

# test -z returns code 1 when its argument is not empty. gofmt -l returns the files that are not well-formatted
goformat: ## Test whether the Go code is well-formatted
	@test -z $$(gofmt -l .)

coverage: ## Generate global code coverage report
	./tools/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./tools/coverage.sh html;

dep: ## Get the Go dependencies
	go mod download

npm-dep: ## Get the NPM dependencies
	npm install

stylelint-ci: ## Checks whether CSS stylesheets are well-formatted
	npm run stylelint --silent -- ./assets

prettier-ci: ## Checks whether HTML templates are well-formatted
	npm run prettier --silent -- --list-different ./templates

build-docker-image: ## Builds the production Docker image
	docker build . --file Dockerfile --tag hyzual/mike-sierra-sierra:latest

dgoss-ci: ## Run goss tests on the production Docker image
	dgoss run hyzual/mike-sierra-sierra:latest

build: ## Build the binary file
	@go build -i -v -o $(PROJECT_EXEC_NAME) $(PKG)

clean: ## Remove previous build
	@rm -f $(PROJECT_EXEC_NAME)

start: build ## Run the built executable
	./$(PROJECT_EXEC_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
