# FROM https://gitlab.com/pantomath-io/demo-tools/-/blob/master/Makefile
PROJECT_NAME := "mike-sierra-sierra"
PKG := "github.com/hyzual/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

get_ip_addr = `docker inspect -f '{{.NetworkSettings.Networks.tuleap_default.IPAddress}}' mike-dev`

.PHONY: all dep build clean test coverage lint start

.SILENT: clean-dev-container

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unit tests
	@go test -short ${PKG_LIST}

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

# test -z returns code 1 when its argument is not empty. gofmt -l returns the files that are not well-formatted
goformat-ci: ## Test whether the Go code is well-formatted
	@test -z $$(gofmt -l .)

coverage: ## Generate global code coverage report
	./tools/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./tools/coverage.sh html;

build-ci: ## Check that the Go code can be compiled
	go build -v ${PKG_LIST}

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
	dgoss run -e MIKE_DISABLE_HTTPS=1 hyzual/mike-sierra-sierra:latest

# Use || true because it's not a problem if the container does not exist
clean-dev-container:
	docker container stop mike || true \
	&& docker container rm mike || true

# Bind-mount the $GOPATH/pkg folder so that the container does not re-download packages all the time
start: clean-dev-container ## Build and run the dev docker container.
	@docker build --file ./tools/docker-dev/Dockerfile.dev --tag hyzual/mike-sierra-sierra:dev ./tools/docker-dev \
	&& docker volume create mike_music || true \
	&& docker run --detach --publish 8443:8443 --name mike \
		--mount type=bind,source=`pwd`,destination=/app,readonly \
		--mount type=bind,source=`pwd`/database/file,destination=/app/database/file \
		--mount type=bind,source=`pwd`/secrets,destination=/app/secrets \
		--mount type=bind,source=$$GOPATH/pkg,destination=/go/pkg,readonly \
		--mount source=mike_music,destination=/music,readonly \
		hyzual/mike-sierra-sierra:dev
	@echo "Go to https://localhost:8443"

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
