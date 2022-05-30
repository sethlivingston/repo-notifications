
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
	go build -ldflags "-X github.com/sethlivingston/reponotifications.Version=${VERSION}" -o bin/reponotifications github.com/sethlivingston/reponotifications/cmd/reponotifications

.PHONY: check
check: update
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=70.0 ./lint-project.sh

dockerx: update
ifeq ($(shell docker manifest inspect sethlivingston/reponotifications:${VERSION} > /dev/null ; echo $$?), 0)
	$(error docker tag already exists)
else
	docker buildx build \
		--push \
		--platform linux/arm64/v8,linux/amd64 \
		--pull --build-arg VERSION=${VERSION} \
		-t sethlivingston/reponotifications:${VERSION} \
		-t sethlivingston/reponotifications:latest \
		-f Dockerfile .
endif

.PHONY: dev-docker
dev-docker: update
# Build a docker image for our local platform
	docker build --pull \
		--build-arg VERSION=${DEV_VERSION} \
		-t sethlivingston/reponotifications:${DEV_VERSION} \
		-f Dockerfile .
ifeq ($(GITHUB_ACTIONS),true)
	docker push sethlivingston/reponotifications:${DEV_VERSION}
endif

run: update build
	./bin/reponotifications

docker-run:
	docker run -v ${PWD}/data:/data -v ${PWD}/configs:/configs --env APP_CONFIG="/configs/config.yml" -it --rm sethlivingston/reponotifications:${VERSION}

test: update
	go test -cover github.com/sethlivingston/reponotifications/...

.PHONY: clean
clean:
	@rm -rf cover.out coverage.txt misspell* staticcheck*
	@rm -rf ./bin/