-include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

api:
	@ echo
	@ echo "Running API..."
	@ go run cmd/server/main.go

tests:
	@ echo
	@ echo "Running tests..."
	@ echo
	@ ginkgo -r --randomize-all --coverprofile=.coverage-report.out