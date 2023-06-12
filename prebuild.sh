#!/usr/bin/env bash

ENTRY=$1
OS=$2
ARCH=$3

export CGO_ENABLED=0
export GO111MODULE=on
# export GOFLAGS="-mod=vendor"
export GOOS=$OS
export GOARCH=$ARCH

go build -pgo=auto -v -installsuffix cgo -o "./build/$GOOS-$GOARCH" "./${ENTRY}"
