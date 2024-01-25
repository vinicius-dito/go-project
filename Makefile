-include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

api:
	@ echo
	@ echo "Running API..."
	@ go run main.go