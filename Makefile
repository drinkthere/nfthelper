# strategy version
strategy ?= 0.1.0

# Repo info
GIT_COMMIT ?= git-$(shell git rev-parse --short HEAD)

TARGETS := darwin/amd64 linux/amd64 windows/amd64

$(info $(GIT_COMMIT) )

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

$(info $(GOBIN))

# Run tests
test: fmt vet lint

build: fmt vet lint
	go mod vendor
	go build -o bin/nft ./main/...
fmt:
	go fmt $$(go list ./...| grep -v /vendor/)

vet:
	go vet ./...

lint: golangci
	gofmt -s -w main config controller logger model router
	$(GOLANGCILINT) -config config/config.toml -exclude vendor/... -formatter default ./...

golangci:
ifeq (, $(shell which revive))
	@{ \
	set -e ;\
	echo 'installing revive' ;\
	echo 'Install succeed' ;\
	go install github.com/mgechev/revive@latest;\
	}
GOLANGCILINT=$(GOBIN)/revive
else
GOLANGCILINT=$(shell which revive)
endif
