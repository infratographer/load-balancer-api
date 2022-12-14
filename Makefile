all: lint test
PHONY: test coverage lint golint clean vendor docker-up docker-down unit-test
GOOS=linux
DB_STRING=host=crdb port=26257 user=root sslmode=disable
DB=load_balancer_api
DEV_DB=${DB}_dev
TEST_DB=${DB}_test
DEV_URI="postgresql://root@crdb:26257/${DEV_DB}?sslmode=disable"
TEST_URI="postgresql://root@crdb:26257/${TEST_DB}?sslmode=disable"

APP_NAME=loadbalancer-api

ci: | dev-database test

test: | unit-test

unit-test:
	@echo Running unit tests...
	@go test -cover -failfast ./...

coverage:
	@echo Generating coverage report...
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic -tags testtools -p 1
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

lint: golint

golint:
	@echo Running golint...
	@golangci-lint run

clean:
	@echo Cleaning...
	@rm -rf ./dist/
	@rm -rf coverage.out
	@go clean -testcache

models: dev-database
	@echo Generating models...
	@sqlboiler crdb --add-soft-deletes --config sqlboiler.toml --always-wrap-errors --wipe --output internal/models --no-tests --tag query,param
	@go mod tidy

binary: | models
	@echo Building binary...
	@go build -o bin/${APP_NAME} main.go

vendor:
	@echo Downloading dependencies...
	@go mod tidy
	@go mod download

dev-database: | vendor
	@echo Creating dev database...
	@cockroach sql -e "drop database if exists ${DEV_DB}"
	@cockroach sql -e "create database ${DEV_DB}"
	@LOADBALANCERAPI_DB_URI="${DEV_URI}" go run main.go migrate up
