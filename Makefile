BINARY = secrets
INSTALL_DIR ?= $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)
GOARCH = amd64
RELEASE_VERSION = v$(shell sed -n 's/^const VERSION = "\([^" ]*\)"/\1/p' secrets.go)
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_DIRTY = $(shell git diff --quiet 2>/dev/null && git diff --cached --quiet 2>/dev/null || echo -dirty)
GIT_TAG = $(shell git describe --tags --exact-match HEAD 2>/dev/null)
BUILD_VERSION ?= $(if $(filter $(RELEASE_VERSION),$(GIT_TAG)),$(RELEASE_VERSION)$(GIT_DIRTY),$(RELEASE_VERSION)$(if $(GIT_COMMIT),-g$(GIT_COMMIT),-development)$(GIT_DIRTY))
GO_BUILD = GOARCH=${GOARCH} go build -mod=vendor -ldflags="-X main.version=$(BUILD_VERSION)"

all: test clean linux freebsd openbsd darwin windows

clean:
	rm -rf bin/

vendor:
	go mod tidy
	go mod vendor

linux: vendor
	GOOS=linux ${GO_BUILD} -o bin/linux_${GOARCH}/${BINARY}

freebsd: vendor
	GOOS=freebsd ${GO_BUILD} -o bin/freebsd_${GOARCH}/${BINARY}

openbsd: vendor
	GOOS=openbsd ${GO_BUILD} -o bin/openbsd_${GOARCH}/${BINARY}

darwin: vendor
	GOOS=darwin ${GO_BUILD} -o bin/darwin_${GOARCH}/${BINARY}

windows: vendor
	GOOS=windows ${GO_BUILD} -o bin/windows_${GOARCH}/${BINARY}.exe

test: vendor
	script/run_tests.sh

install: vendor
	install -d "$(INSTALL_DIR)"
	go build -mod=vendor -ldflags="-X main.version=$(BUILD_VERSION)" -o "$(INSTALL_DIR)/$(BINARY)" .

release:
	script/release.sh --check
	$(MAKE) all BUILD_VERSION=$(RELEASE_VERSION)
	script/release.sh

.PHONY: all clean vendor linux freebsd openbsd darwin windows test install release
