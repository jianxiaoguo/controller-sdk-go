# the filepath to this repository, relative to $GOPATH/src
REPO_PATH = github.com/drycc/controller-sdk-go

DEV_ENV_IMAGE := golang:1.14
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}

# Enable vendor/ directory support.
export GO15VENDOREXPERIMENT=1

PKG_DIRS := ./...

DEV_ENV_CMD := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}

bootstrap:
	${DEV_ENV_CMD} go mod vendor

build:
	${DEV_ENV_CMD} go build ${PKG_DIRS}

test: build
	${DEV_ENV_CMD} go test ${PKG_DIRS}
