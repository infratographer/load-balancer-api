all: lint test
PHONY: test coverage lint golint clean vendor docker-up docker-down unit-test
GOOS=linux
DB_STRING=host=crdb port=26257 user=root sslmode=disable
DB=load_balancer_api
DEV_DB=${DB}_dev
TEST_DB=${DB}_test
DEV_URI=dbname=${DEV_DB} ${DB_STRING}
TEST_URI=dbname=${TEST_DB} ${DB_STRING}
# use the working dir as the app name, this should be the repo name
APP_NAME=$(shell basename $(CURDIR))

test: | unit-test

unit-test: | lint
	@echo Running unit tests...
	@go test -cover -short -tags testtools ./...

coverage:
	@echo Generating coverage report...
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic -tags testtools -p 1
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

lint: golint

clean:
	@echo Cleaning...
	@rm -rf ./dist/
	@rm -rf coverage.out
	@go clean -testcache

models: dev-database
	@sqlboiler crdb --add-soft-deletes --config sqlboilder.toml
	@go mod tidy

vendor:
	@go mod tidy
	@go mod download

dev-database: | vendor
	@cockroach sql -e "drop database if exists ${DEV_DB}"
	@cockroach sql -e "create database ${DEV_DB}"
	@LOADBALANCERAPI_DB_URI="${DEV_URI}" go run main.go migrate up

test-database: | vendor
	@cockroach sql -e "drop database if exists ${TEST_DB}"
	@cockroach sql -e "create database ${TEST_DB}"
	@LOADBALANCERAPI_DB_URI="${TEST_URI}" go run main.go migrate up
	@cockroach sql -e "use ${TEST_DB};"
