# OS
OSNAME 			  :=
BINARY_NAME_FILE  :=
ifeq ($(OS),Windows_NT)
	OSNAME=windows
else
	UNAME_S :=$(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OSNAME=linux
	endif
	ifeq ($(UNAME_S),Darwin)
		OSNAME=darwin
	endif
endif

# env
CGO_ENABLED=1
GOCMD=go
GOARCH=amd64
SERVICE_PORT=5000
BINARY_NAME=server
BINARY_NAME_FILE =./dist/$(OSNAME)/
BINARY_NAME_LINUX=./dist/linux/
BINARY_NAME_MACOS=./dist/darwin/
BINARY_NAME_WIN=./dist/windows/
GIT_COMMIT=$(shell git rev-list -1 HEAD)
VERSION=$(shell date "+%Y.%m.%d.%H:%M:%S")
GIT_TAG=$(shell git describe --tags --abbrev=0)
BUILD_FLAGS=-i -v -ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(GIT_TAG) -X main.BuildDate=$(VERSION) -X main.mode=release"

all: help

help:
	@echo "Utilização: make [comando]"
	@echo ""
	@echo "Comandos:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

watch: ## reload app while you develop
	docker run -it --rm -w /go/src/crk -v $(shell pwd):/go/src/crk -p $(SERVICE_PORT):$(SERVICE_PORT) cosmtrek/air -c .air.toml

clean: ## clean dirs and reports
	@rm -fr ./tmp
	@rm -fr ./dist
	@rm -fr ./vendor

prebuild:
	@mkdir -p ./dist/$(OSNAME)/

prebuild-all:
	mkdir -p $(BINARY_NAME_LINUX)
	mkdir -p $(BINARY_NAME_MACOS)
	mkdir -p $(BINARY_NAME_WIN)

_dist_os:
	GIN_MODE=release $(GOCMD) build $(BUILD_FLAGS) -o $(BINARY_NAME_FILE) ./command/...

build: prebuild _dist_os

build-linux:
	CC="x86_64-linux-musl-gcc" CXX="x86_64-linux-musl-g++" GOOS=linux $(GOCMD) build $(BUILD_FLAGS) -o $(BINARY_NAME_LINUX) ./command/...

build-mac:
	GOOS=darwin $(GOCMD) build $(BUILD_FLAGS) -o $(BINARY_NAME_MACOS) ./command/...

build-win:
	CC="x86_64-w64-mingw32-gcc" GOOS=windows $(GOCMD) build $(BUILD_FLAGS) -o $(BINARY_NAME_WIN) ./command/...

build-all: prebuild-all build-mac build-win build-linux
