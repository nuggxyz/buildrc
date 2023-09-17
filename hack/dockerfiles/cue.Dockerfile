
# syntax=docker/dockerfile:labs

ARG GO_VERSION
ARG CUE_VERSION
ARG BUILDRC_VERSION

FROM walteh/buildrc:${BUILDRC_VERSION} AS buildrc
FROM cuelang/cue:${CUE_VERSION} AS cue

FROM --platform=$BUILDPLATFORM alpine:latest AS alpinelatest

# Base tools image with required packages
FROM alpinelatest AS tools
COPY --from=buildrc /usr/bin/buildrc /usr/bin/buildrc
COPY --from=cue /usr/bin/cue /usr/bin/cue
RUN apk add --no-cache git curl

# Set common working directory
WORKDIR /wrk

# Buf stage
FROM tools AS pkggen
RUN --mount=type=bind,target=.,rw \
	--mount=type=cache,target=/root/.cache \
	--mount=type=cache,target=/go/pkg/mod <<SHELL
	#!/bin/sh
	set -e -o pipefail

	scheme="https://json.schemastore.org/github-workflow.json"
	curl -sSL "${scheme}" -o github-workflow.json
	cue import -f -p github -l '#GithubWorkflowSpec:' github-workflow.json > github-workflow.cue
	mkdir -p /out/pkg/json.schemastore.org/github
	cp github-workflow.cue /out/pkg/json.schemastore.org/github/workflow.cue
	# exit 1
SHELL

# Final update stage
FROM scratch AS generate
COPY --from=pkggen /out /

FROM tools AS validate
ARG DESTDIR TARGETARCH
COPY --from=generate . /out
RUN --mount=target=/base <<EOT
	set -e
	cd /base
	buildrc diff --current="${DESTDIR}" --correct="/out" --glob="**/*.pb.go"
EOT
