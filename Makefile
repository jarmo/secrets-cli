BINARY = secrets
GOARCH = amd64
GO_BUILD = GOARCH=${GOARCH} go build -mod=vendor
PREFIX ?= ${GOPATH}

all: test clean linux freebsd openbsd darwin windows

clean:
	rm -rf bin/

vendor:
	go mod vendor
	go mod tidy

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

install:
	cp -Rf bin/ "${PREFIX}/bin"

release: all
	script/release.sh

.PHONY: all clean vendor linux freebsd openbsd darwin windows test install release
