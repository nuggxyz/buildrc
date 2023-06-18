#!/usr/bin/env bash

OUTPUT_FILE=$1

export CGO_ENABLED=0
export GO111MODULE=on

echo "🚀 building $OUTPUT_FILE"

go build -pgo=auto -v -installsuffix cgo -ldflags "${GO_LDFLAGS}" -o "$OUTPUT_FILE" "./cmd"

if [ ! -f "$OUTPUT_FILE" ]; then
	echo "❌ build failed: $OUTPUT_FILE not found"
	exit 1
else
	$OUTPUT_FILE version || echo "not a valid binary (this is expected)"
	echo "✅ build succeeded: $OUTPUT_FILE"
fi
