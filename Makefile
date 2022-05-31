
PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | tr '[:upper:]' '[:lower:]')
PWD := $(shell pwd)

ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

COMMIT_HASH :=$(shell git rev-parse --short HEAD)

USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

all: update build

update:
	go get ./...

build:
	go build -ldflags "-X github.com/sethlivingston/repowatch.Version=${VERSION}" -o bin/repowatch github.com/sethlivingston/repowatch/cmd/repowatch

.PHONY: check
check: update
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=70.0 ./lint-project.sh

run: update build
	./bin/repowatch

test: update
	go test -cover github.com/sethlivingston/repowatch/...

.PHONY: clean
clean:
	@rm -rf cover.out coverage.txt misspell* staticcheck*
	@rm -rf ./bin/