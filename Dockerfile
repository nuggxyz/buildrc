# syntax=docker/dockerfile:1

ARG GO_VERSION=
ARG XX_VERSION=1.2.1

ARG DOCKER_VERSION=24.0.2
ARG GOTESTSUM_VERSION=v1.9.0
ARG REGISTRY_VERSION=2.8.0
ARG BUILDKIT_VERSION=v0.11.6

# xx is a helper for cross-compilation
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS golatest

FROM golatest AS gobase
COPY --from=xx / /
RUN apk add --no-cache file git bash
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0
WORKDIR /src

FROM registry:$REGISTRY_VERSION AS registry

FROM moby/buildkit:$BUILDKIT_VERSION AS buildkit

FROM docker/buildx-bin:latest AS buildx-bin

FROM gobase AS docker
ARG TARGETPLATFORM
ARG DOCKER_VERSION
ARG VERSION
ARG BIN_NAME
WORKDIR /opt/docker
RUN <<EOT
CASE=${TARGETPLATFORM:-linux/amd64}
DOCKER_ARCH=$(
	case ${CASE} in
	"linux/amd64") echo "x86_64" ;;
	"linux/arm/v6") echo "armel" ;;
	"linux/arm/v7") echo "armhf" ;;
	"linux/arm64/v8") echo "aarch64" ;;
	"linux/arm64") echo "aarch64" ;;
	"linux/ppc64le") echo "ppc64le" ;;
	"linux/s390x") echo "s390x" ;;
	*) echo "" ;; esac
)
echo "DOCKER_ARCH=$DOCKER_ARCH" &&
wget -qO- "https://download.docker.com/linux/static/stable/${DOCKER_ARCH}/docker-${DOCKER_VERSION}.tgz" | tar xvz --strip 1
EOT
RUN ./dockerd --version && ./containerd --version && ./ctr --version && ./runc --version

FROM gobase AS gotestsum
ARG GOTESTSUM_VERSION
ENV GOFLAGS=
RUN --mount=target=/root/.cache,type=cache <<EOT
	GOBIN=/out/ go install "gotest.tools/gotestsum@${GOTESTSUM_VERSION}" &&
	/out/gotestsum --version
EOT

FROM gobase AS meta
RUN --mount=type=bind,target=. <<EOT
  set -e
  mkdir /meta
  echo -n "$(./hack/git-meta version)" | tee /meta/version
  echo -n "$(./hack/git-meta revision)" | tee /meta/revision
EOT

FROM gobase AS builder
ARG TARGETPLATFORM
ARG GO_PKG
ARG BIN_NAME
RUN --mount=type=bind,target=. \
	--mount=type=cache,target=/root/.cache \
	--mount=type=cache,target=/go/pkg/mod \
	--mount=type=bind,from=meta,source=/meta,target=/meta <<EOT
  set -e
  echo "Building for ${TARGETPLATFORM}"
  xx-go --wrap
  DESTDIR=/usr/bin VERSION=$(cat /meta/version) REVISION=$(cat /meta/revision) GO_EXTRA_LDFLAGS="-s -w" ./hack/build
  xx-verify --static /usr/bin/${BIN_NAME}
EOT

FROM gobase AS test
ENV SKIP_INTEGRATION_TESTS=1
RUN --mount=type=bind,target=. \
	--mount=type=cache,target=/root/.cache \
	--mount=type=cache,target=/go/pkg/mod <<EOT
	go test -v -coverprofile=/tmp/coverage.txt -covermode=atomic ./... &&
	go tool cover -func=/tmp/coverage.txt
EOT

FROM scratch AS test-coverage
COPY --from=test /tmp/coverage.txt /coverage.txt

FROM scratch AS binaries-unix
ARG BIN_NAME
COPY --link --from=builder /usr/bin/${BIN_NAME} /${BIN_NAME}

FROM binaries-unix AS binaries-darwin
FROM binaries-unix AS binaries-linux

FROM scratch AS binaries-windows
ARG BIN_NAME
COPY --link --from=builder /usr/bin/${BIN_NAME} /${BIN_NAME}.exe

FROM binaries-$TARGETOS AS binaries
# enable scanning for this stage
ARG BUILDKIT_SBOM_SCAN_STAGE=true
ARG BIN_NAME

FROM gobase AS integration-test-base
# https://github.com/docker/docker/blob/master/project/PACKAGERS.md#runtime-dependencies
RUN apk add --no-cache \
	btrfs-progs \
	e2fsprogs \
	e2fsprogs-extra \
	ip6tables \
	iptables \
	openssl \
	shadow-uidmap \
	xfsprogs \
	xz
COPY --link --from=gotestsum /out/gotestsum /usr/bin/
COPY --link --from=registry /bin/registry /usr/bin/
COPY --link --from=docker /opt/docker/* /usr/bin/
COPY --link --from=buildkit /usr/bin/buildkitd /usr/bin/
COPY --link --from=buildkit /usr/bin/buildctl /usr/bin/
COPY --link --from=binaries /${BIN_NAME} /usr/bin/
COPY --link --from=buildx-bin /buildx /usr/libexec/docker/cli-plugins/docker-buildx

FROM integration-test-base AS integration-test
COPY . .

# Release
FROM --platform=$BUILDPLATFORM alpine AS releaser
WORKDIR /work
ARG TARGETPLATFORM
ARG BIN_NAME
RUN --mount=from=binaries \
	--mount=type=bind,from=meta,source=/meta,target=/meta <<EOT
  set -e
  mkdir -p /out
  cp ${BIN_NAME}* "/out/${BIN_NAME}-$(cat /meta/version).$(echo $TARGETPLATFORM | sed 's/\//-/g')$(ls ${BIN_NAME}* | sed -e 's/^${BIN_NAME}//')"
EOT

FROM scratch AS release
COPY --from=releaser /out/ /

# Shell
FROM docker:$DOCKER_VERSION AS dockerd-release
FROM alpine AS shell
ARG BIN_NAME
RUN apk add --no-cache iptables tmux git vim less openssh
RUN mkdir -p /usr/local/lib/docker/cli-plugins && ln -s /usr/local/bin/${BIN_NAME} /usr/local/lib/docker/cli-plugins/docker-${BIN_NAME}
COPY ./hack/demo-env/entrypoint.sh /usr/local/bin
COPY ./hack/demo-env/tmux.conf /root/.tmux.conf
COPY --from=dockerd-release /usr/local/bin /usr/local/bin
WORKDIR /work
COPY ./hack/demo-env/examples .
COPY --from=binaries / /usr/local/bin/
VOLUME /var/lib/docker
ENTRYPOINT ["entrypoint.sh"]

FROM binaries
