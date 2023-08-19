# syntax=docker/dockerfile:1

ARG GO_VERSION=

FROM alpine:latest AS base
RUN apk add --no-cache git rsync
WORKDIR /src

FROM walteh/buildrc:pr-25 AS buildrc
RUN --mount=target=/over/here,type=tmpfs <<EOT
set -e
buildrc --git-dir=/over/here --hit-dir /go/pkg/mod
EOT
