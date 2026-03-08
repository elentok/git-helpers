.PHONY: build install test run

build:
	go build -ldflags "-X gx/cmd.version=$(shell git describe --tags --always --dirty)" -o gx .

install:
	go install .

test:
	go test ./...

run:
	go run .
