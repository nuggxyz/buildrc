#!/usr/bin/env bash

ENTRY=$1
OS=$2
ARCH=$3

export CGO_ENABLED=0
export GO111MODULE=on
# export GOFLAGS="-mod=vendor"
export GOOS=$OS
export GOARCH=$ARCH

OUTPUT_FILE=./build/$GOOS-$GOARCH

go build -pgo=auto -v -installsuffix cgo -o "$OUTPUT_FILE" "./${ENTRY}"

tar -czvf "$OUTPUT_FILE".tar.gz "$OUTPUT_FILE"

shasum -a 256 "$OUTPUT_FILE".tar.gz | awk '{ print $1 }' >"$OUTPUT_FILE".sha256
