.PHONY: build install test run

build:
	go build -o gx .

install:
	go install .

test:
	go test ./...

run:
	go run .
