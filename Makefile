-include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

api:
	@ echo
	@ echo "Running API..."
	@ go run main.go

tests:
	@ echo
	@ echo "Running tests..."
	@ echo
	@ FIRESTORE_EMULATOR_HOST=localhost:8080 go test -v ./... -coverprofile=coverage.out