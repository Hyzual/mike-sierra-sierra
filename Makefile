# FROM https://gitlab.com/pantomath-io/demo-tools/-/blob/master/Makefile
PROJECT_NAME := "mike-sierra-sierra"
PKG := "github.com/hyzual/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

get_ip_addr = `docker inspect -f '{{.NetworkSettings.Networks.tuleap_default.IPAddress}}' mike-dev`

.PHONY: all help dep clean coverage start watch race

.SILENT: clean-dev-container

all: help

lint-go: ## Lint the files
	@golangci-lint run

test-go: ## Run Go unit tests
	@go test -short ${PKG_LIST}

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

coverage-go: ## Generate global code coverage report
	./tools/coverage.sh;

coverage-go-html: ## Generate global code coverage report in HTML
	./tools/coverage.sh html;

build-ci: ## Check that the Go code can be compiled
	go build -v ${PKG_LIST}

dep: ## Get the Go dependencies
	go mod download

npm-dep: ## Get the NPM dependencies
	npm install

build-assets: ## Build the frontend assets for production
	npm run build

coverage-jest: ## Run Jest unit tests with coverage
	npm run test --silent -- --coverage

test-jest: ## Run Jest unit tests
	npm run test --silent

eslint-ci: ## Checks TypeScript and Javascript files for errors and formatting
	npm run eslint --silent -- .

stylelint-ci: ## Checks CSS for errors and formatting
	npm run stylelint --silent -- ./styles

prettier-ci: ## Checks whether HTML templates and JS configurations are well-formatted
	npm run prettier --silent -- --list-different ./templates .eslintrc.js .stylelintrc.js

build-docker-image: ## Builds the production Docker image
	docker build . --file Dockerfile --tag hyzual/mike-sierra-sierra:latest

dgoss-ci: ## Run goss tests on the production Docker image
	dgoss run -e MIKE_DISABLE_HTTPS=1 hyzual/mike-sierra-sierra:latest

# Use || true because it's not a problem if the container does not exist
clean-dev-container:
	docker container stop mike || true \
	&& docker container rm mike || true

# Bind-mount the $HOME/go/pkg folder so that the container does not re-download packages all the time
start: clean-dev-container ## Build and run the dev docker container.
	@docker build --file ./tools/docker-dev/Dockerfile.dev --tag hyzual/mike-sierra-sierra:dev ./tools/docker-dev \
	&& docker volume create mike_music || true \
	&& docker run --detach --publish 8443:8443 --name mike \
		--mount type=bind,source=`pwd`,destination=/app,readonly \
		--mount type=bind,source=`pwd`/database/file,destination=/app/database/file \
		--mount type=bind,source=`pwd`/secrets,destination=/app/secrets \
		--mount type=bind,source=$$HOME/go/pkg,destination=/go/pkg,readonly \
		--mount source=mike_music,destination=/music,readonly \
		hyzual/mike-sierra-sierra:dev
	@echo "Go to https://localhost:8443"

watch: ## Build and watch frontend assets
	npm run watch

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
