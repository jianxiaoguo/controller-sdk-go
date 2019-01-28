# the filepath to this repository, relative to $GOPATH/src
REPO_PATH = github.com/drycc/controller-sdk-go

DEV_ENV_IMAGE := quay.io/drycc/go-dev:v0.22.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}

# Enable vendor/ directory support.
export GO15VENDOREXPERIMENT=1

PKG_DIRS := ./...

DEV_ENV_CMD := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}

bootstrap:
	${DEV_ENV_CMD} dep ensure

build:
	${DEV_ENV_CMD} go build ${PKG_DIRS}

test-cover:
	${DEV_ENV_CMD} test-cover.sh
test-style:
	${DEV_ENV_CMD} lint
test: build test-style test-cover
	${DEV_ENV_CMD} go test ${PKG_DIRS}
