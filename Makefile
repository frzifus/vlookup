APP = vlookup
DATE = $(shell date +%FT%T%Z)
BUILD_DIR = build/bin
GIT_VER=$(shell git rev-parse HEAD)

LDFLAGS=-ldflags "-X github.com/frzifus/vlookup/pkg/version.hash=${GIT_VER} \
									-X github.com/frzifus/vlookup/pkg/version.buildtimestamp=${DATE}"

.PHONY: test clean arm amd64 disclean mrproper

# Build the project
all: amd64 arm

build_deps:
	go install golang.org/x/lint/golint@latest
amd64:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP}-linux-amd64 -v cmd/${APP}/*.go
arm:
	GOOS=linux GOARCH=arm go build ${LDFLAGS} -o ${BUILD_DIR}/${APP}-linux-arm -v cmd/${APP}/*.go

lint:
	golint -set_exit_status ./pkg/... ./cmd/...

test:
	go test -v ./...

clean:
	-rm -f ${BUILD_DIR}/${BINARY}-*

distclean:
	rm -rf ./build

mrproper: distclean
	git ls-files --others | xargs rm -rf
