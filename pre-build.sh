#!/usr/bin/env bash

OUTPUT_FILE_OVERRIDE=$1

export CGO_ENABLED=0
export GO111MODULE=on

arches=$(jq -r '.arch' "$BUILDRC_PACKAGE_JSON")
oses=$(jq -r '.os' "$BUILDRC_PACKAGE_JSON")

function ROLL() {

	local os=$1
	local arch=$2
	local output_file=$3-"$BUILDRC_PACKAGE_NAME-$os-$arch"

	echo "🚀 building $BUILDRC_PACKAGE_NAME for $os/$arch"

	GOOS=$os GOARCH=$arch go build -pgo=auto -v -installsuffix cgo -ldflags "-s -w" -o "$output_file" "./cmd"

	if [ ! -f "$output_file" ]; then
		echo "❌ build failed: $output_file not found"
		exit 1
	else
		$output_file version || echo "not a valid binary (this is expected)"
		echo "✅ build succeeded: $output_file"
	fi

	cp -r "$output_file" "$BUILDRC_TARGZ"
	cp -r "$output_file" "$BUILDRC_SHA256"
}

if [ -n "$OUTPUT_FILE_OVERRIDE" ]; then
	ROLL "${GOOS} ${GOARCH}" "$OUTPUT_FILE_OVERRIDE"
	exit 0
else
	for os in $oses; do
		for arch in $arches; do
			ROLL "$os" "$arch"
		done
	done
fi
