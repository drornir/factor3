SHELL:=/bin/zsh

BINARY=factor3
VERSION=$(strip $(shell cat version.txt))

default: help


##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Show help for each of the Makefile recipes.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<command>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


.PHONY: test
test: build ## Run tests
	go test ./...
	# go run ./example/app/

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...

.PHONY: setup
setup: setup-tools ## setup dev env
	mkdir -p bin
	go mod download

.PHONY: setup-tools
setup-tools: ## install dev deps 
	go install github.com/pquerna/ffjson@latest
	go install golang.org/x/tools/cmd/stringer@latest
	go install github.com/campoy/jsonenums@latest

.PHONY: gen 
gen: ## run go generate
	go generate ./...

.PHONY: build 
build: ## full build including generate, go get
	go get ./...
	go mod tidy
	go build -o bin/${BINARY} .

.PHONY: publish
publish: ## Publish a new version to github
	git tag $(VERSION)
	git push origin $(VERSION)
	gh release create $(VERSION)  --prerelease --title $(VERSION) --notes "$(VERSION)"
	GOPROXY=proxy.golang.org go list -m github.com/drornir/factor3@${VERSION}

.PHONY: version
version:
	@echo $(VERSION)