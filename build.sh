#!/usr/bin/env bash

PACKAGE_NAME=$1
OUTPUT_FILE=$2
CUSTOM_DATA=$3

export CGO_ENABLED=0
export GO111MODULE=on

echo "🚀 building $PACKAGE_NAME to $OUTPUT_FILE"

echo "custom data: $CUSTOM_DATA"

go build -pgo=auto -v -installsuffix cgo -ldflags "${GO_LDFLAGS}" -o "$OUTPUT_FILE" "./cmd"

if [ ! -f "$OUTPUT_FILE" ]; then
	echo "❌ build failed: $OUTPUT_FILE not found"
	exit 1
else
	$OUTPUT_FILE version || echo "not a valid binary (this is expected)"
	echo "✅ build succeeded: $OUTPUT_FILE"
fi
