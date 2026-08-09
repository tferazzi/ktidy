VERSION ?= dev

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/ktidy .

test:
	go test ./...

lint:
	go vet ./...

.PHONY: build test lint
