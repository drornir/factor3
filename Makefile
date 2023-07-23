SHELL:=/bin/zsh

BINARY=factor3

default: help

.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done


.PHONY: test
test: build # Run tests
	go test ./...
	go run ./example/app/

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...

.PHONY: setup
setup: setup-tools # Setup tools required for local development.
	mkdir -p bin
	go mod download

.PHONY: setup-tools
setup-tools:
	go install github.com/pquerna/ffjson@latest
	go install golang.org/x/tools/cmd/stringer@latest
	go install github.com/campoy/jsonenums@latest

.PHONY: build 
build: # Build the binary.
	go generate ./..
	go get ./...
	go mod tidy
	go build -o bin/${BINARY}
