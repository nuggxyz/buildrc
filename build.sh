#!/usr/bin/env bash

OS=$1
ARCH=$2
OUTPUT_FILE=$3

export CGO_ENABLED=0
export GO111MODULE=on
# export GOFLAGS="-mod=vendor"
export GOOS=$OS
export GOARCH=$ARCH

go build -pgo=auto -v -installsuffix cgo -o "$OUTPUT_FILE" "./cmd"
