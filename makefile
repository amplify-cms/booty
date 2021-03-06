.PHONY: all dep mock-release release
.DEFAULT_GOAL := all
VERSION_GITHASH = $(shell git rev-parse master | tr -d '\n')
GO_LDFLAGS = CGO_ENABLED=0 go build -ldflags "-s -w -X main.revision=${VERSION_GITHASH} -X main.version=0.0.0-dev" -a -tags netgo


all: build
all-ci: all dep mock-release
all-release: all dep release

build:
	$(GO_LDFLAGS) -o bin/booty main.go

dep:
	./bin/booty install-all

mock-release:
	./bin/booty release --rm-dist --skip-publish --snapshot

release:
	./bin/booty release release

test:
	go test -v -short=false -race -coverprofile=coverage.txt -covermode=atomic ./dep/orchestrator
