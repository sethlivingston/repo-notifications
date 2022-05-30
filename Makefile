
PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | tr '[:upper:]' '[:lower:]')
PWD := $(shell pwd)

ifndef VERSION
	VERSION := $(shell git describe --tags --abbrev=0)
endif

COMMIT_HASH :=$(shell git rev-parse --short HEAD)
DEV_VERSION := dev-${COMMIT_HASH}

USERID := $(shell id -u $$USER)
GROUPID:= $(shell id -g $$USER)

all: update build

update:
	go get ./...

build:
	go build -ldflags "-X github.com/sethlivingston/repo-notifications.Version=${VERSION}" -o bin/repo-notifications github.com/sethlivingston/repo-notifications/cmd/repo-notifications

.PHONY: check
check: update
@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
@chmod +x ./lint-project.sh
COVER_THRESHOLD=70.0 ./lint-project.sh

dockerx: update
ifeq ($(shell docker manifest inspect moovfinancial/busint:${VERSION} > /dev/null ; echo $$?), 0)
	$(error docker tag already exists)
else
	docker buildx build \
		--push \
		--platform linux/arm64/v8,linux/amd64 \
		--pull --build-arg VERSION=${VERSION} \
		-t moovfinancial/busint:${VERSION} \
		-t moovfinancial/busint:latest \
		-f Dockerfile .
endif

.PHONY: dev-docker
dev-docker: update
# Build a docker image for our local platform
	docker build --pull \
		--build-arg VERSION=${DEV_VERSION} \
		-t moovfinancial/busint:${DEV_VERSION} \
		-f Dockerfile .
ifeq ($(GITHUB_ACTIONS),true)
	docker push moovfinancial/busint:${DEV_VERSION}
endif

docker-run:
	docker run -v ${PWD}/data:/data -v ${PWD}/configs:/configs --env APP_CONFIG="/configs/config.yml" -it --rm moovfinancial/busint:${VERSION}

test: update
	go test -cover github.com/sethlivingston/repo-notifications/...

.PHONY: clean
clean:
@rm -rf cover.out coverage.txt misspell* staticcheck*
@rm -rf ./bin/