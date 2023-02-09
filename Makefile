.PHONY: build test

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-weather.version=${VERSION}'" -o conduit-connector-weather cmd/connector/main.go

test:
	go test $(GOTEST_FLAGS) -v -race ./...
