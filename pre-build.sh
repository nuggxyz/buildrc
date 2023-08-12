#!/usr/bin/env bash

OUTPUT_FILE=$1

if [ -z "$OUTPUT_FILE" ]; then
	OUTPUT_FILE="${BUILDRC_WORKING_DIR}"
fi

echo "🚀 building $BUILDRC_PACKAGE_NAME to $OUTPUT_FILE"

export CGO_ENABLED=0
export GO111MODULE=on

go build -pgo=auto -v -installsuffix cgo -ldflags "-s -w" -o "$OUTPUT_FILE" "./cmd"

if [ ! -f "$OUTPUT_FILE" ]; then
	echo "❌ build failed: $OUTPUT_FILE not found"
	exit 1
else
	$OUTPUT_FILE version || echo "not a valid binary (this is expected)"
	echo "✅ build succeeded: $OUTPUT_FILE"
fi

cp -r "$OUTPUT_FILE" "$BUILDRC_TARGZ"
cp -r "$OUTPUT_FILE" "$BUILDRC_SHA256"
