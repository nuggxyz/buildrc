#!/usr/bin/env bash

arches=$(echo "$BUILDRC_PACKAGE_JSON" | jq -r '.arch[]')
oses=$(echo "$BUILDRC_PACKAGE_JSON" | jq -r '.os[]')

function ROLL() {

	local name=$1
	local os=$2
	local arch=$3
	local override=$4

	local output_file

	if [ -z "$override" ]; then
		output_file="./$name-$os-$arch"
	else
		output_file="$override"
	fi

	GO111MODULE=on CGO_ENABLED=0 GOOS=$os GOARCH=$arch \
		go build -pgo=auto -v -installsuffix cgo -ldflags "-s -w" -o "$output_file" "./cmd"

	if [ ! -f "$output_file" ]; then
		echo "❌ build failed: $output_file not found"
		exit 1
	else
		"$output_file" version || echo "not a valid binary (this is expected)"
		echo "✅ build succeeded: $output_file"
	fi

	echo "🚀 moving $output_file to $BUILDRC_TARGZ"
	cp -r "$output_file" "$BUILDRC_TARGZ"
	echo "🚀 moving $output_file to $BUILDRC_TARGZ"
	cp -r "$output_file" "$BUILDRC_SHA256"
}

if [ -z "$BUILDRC_PACKAGE_NAME" ]; then
	echo "⚠️ BUILDRC_PACKAGE_NAME is not set, building for exec override"
	ROLL "early" "$(go env GOOS)" "$(go env GOARCH)" "$BUILDRC_EXEC_OVERRIDE"
	exit 0
else
	for os in $oses; do
		for arch in $arches; do
			echo "📦 building $BUILDRC_PACKAGE_NAME for arch='$arch' os='$os'"
			ROLL "$BUILDRC_PACKAGE_NAME" "$os" "$arch"
		done
	done
fi
