all: lint test

PHONY: help all test coverage lint golint clean vendor docker-up docker-down unit-test
GOOS=linux
DB=load_balancer_api
DEV_DB=${DB}_dev
TEST_DB=${DB}_test
DEV_URI="postgresql://root@crdb:26257/${DEV_DB}?sslmode=disable"
TEST_URI="postgresql://root@crdb:26257/${TEST_DB}?sslmode=disable"

APP_NAME=loadbalancer-api

help: Makefile ## Print help
	@grep -h "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/:.*##/#/' | column -c 2 -t -s#

tests: | unit-test

unit-tests: ## Runs unit tests
	@echo --- Running unit tests...
	@date --rfc-3339=seconds
	@go test -race -cover -failfast -tags testtools -p 1 -v ./...

coverage: ## Generates coverage report
	@echo --- Generating coverage report...
	@date --rfc-3339=seconds
	@go test -race -coverprofile=coverage.out -covermode=atomic -tags testtools -p 1 ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

lint: golint ## Runs linting

golint:
	@echo --- Running golint...
	@date --rfc-3339=seconds
	@golangci-lint run

clean: ## Clean up all the things
	@echo --- Cleaning...
	@date --rfc-3339=seconds
	@rm -rf ./dist/
	@rm -rf coverage.out
	@go clean -testcache

binary: | generate ## Builds the binary
	@echo --- Building binary...
	@date --rfc-3339=seconds
	@go build -o bin/${APP_NAME} main.go

vendor: ## Vendors dependencies
	@echo --- Downloading dependencies...
	@date --rfc-3339=seconds
	@go mod tidy
	@go mod download

dev-nats: ## Initializes nats
	@echo --- Initializing nats
	@date --rfc-3339=seconds
	@.devcontainer/scripts/nats_account.sh

generate: vendor ## Generates code
	@echo --- Generating code...
	@date --rfc-3339=seconds
	@go generate ./...