# Default values used by tests
GITHUB_EMAIL ?= foo@form3.tech
GITHUB_USERNAME ?= foo
COMMIT_MESSAGE_PREFIX ?= '[foo]'

default: vet test build

.PHONY: build
build:
	go build -o bin/terraform-provider-githubdeploymentenvironment

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	GITHUB_TOKEN=$(GITHUB_TOKEN) \
	GITHUB_OWNER=$(GITHUB_OWNER) \
	GITHUB_ORGANIZATION=$(GITHUB_ORGANIZATION) \
	go test -count 1 -v ./...
