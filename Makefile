.PHONY: build test

build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-weather.version=v0.1.0'" -o conduit-connector-weather cmd/connector/main.go

test:
	go test $(GOTEST_FLAGS) -v -race ./...
