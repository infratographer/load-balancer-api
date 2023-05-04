all: lint test

PHONY: help all test coverage lint golint clean vendor docker-up docker-down unit-test
GOOS=linux
DB_STRING=host=crdb port=26257 user=root sslmode=disable
DB=load_balancer_api
DEV_DB=${DB}_dev
TEST_DB=${DB}_test
DEV_URI="postgresql://root@crdb:26257/${DEV_DB}?sslmode=disable"
TEST_URI="postgresql://root@crdb:26257/${TEST_DB}?sslmode=disable"

APP_NAME=loadbalancer-api

help: Makefile ## Print help
	@grep -h "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/:.*##/#/' | column -c 2 -t -s#

ci: | lint

test: | unit-test

unit-test: ## Runs unit tests
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


binary: | models ## Builds the binary
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

dev-database: | vendor ## Initializes the dev database
	@echo --- Creating dev database...
	@date --rfc-3339=seconds
	@cockroach sql -e "drop database if exists ${DEV_DB}"
	@cockroach sql -e "create database ${DEV_DB}"
	@LOADBALANCERAPI_DB_URI="${DEV_URI}" go run main.go migrate up
