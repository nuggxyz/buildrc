# syntax=docker/dockerfile:labs

ARG GO_VERSION

ARG BUILDRC_VERSION=
FROM walteh/buildrc:${BUILDRC_VERSION} AS buildrc

FROM --platform=$BUILDPLATFORM debian:bookworm-slim AS bookworm

# Base tools image with required packages
FROM bookworm AS tools
COPY --from=buildrc /usr/bin/buildrc /usr/bin/buildrc
RUN apt-get update && apt-get --no-install-recommends install -y git unzip

# Set common working directory
WORKDIR /wrk

# Buf stage
FROM tools AS bufgen
COPY --from=bufbuild/buf:latest /usr/local/bin/buf /usr/bin/
RUN --mount=type=bind,target=.,rw \
	--mount=type=cache,target=/root/.cache \
	--mount=type=cache,target=/go/pkg/mod <<EOT
	set -ex
	buf generate --exclude-path ./vendor --output . --include-imports --include-wkt || echo "buf generate failed - ignoring"
	mkdir /out
	git ls-files -m --others -- ':!vendor' '**/*.pb.go' | tar -cf - --files-from - | tar -C /out -xf -
EOT

# Final update stage
FROM scratch AS generate
COPY --from=bufgen /out /

FROM tools AS validate
ARG DESTDIR TARGETARCH
COPY --from=generate . /out
RUN --mount=target=/base <<EOT
	set -e
	cd /base
	./tmp/buildrc-${TARGETARCH}-tmp-diff diff --current="${DESTDIR}" --correct="/out" --glob="**/*.pb.go"
EOT
