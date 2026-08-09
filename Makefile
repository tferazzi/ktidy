VERSION ?= dev

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/ktidy .

test:
	go test ./...

lint:
	go vet ./...
	golangci-lint run

snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean

.PHONY: build test lint snapshot release
