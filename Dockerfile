# syntax=docker/dockerfile:labs

##################################################################
# SETUP
##################################################################

ARG GO_VERSION=
ARG XX_VERSION=
ARG GOTESTSUM_VERSION=
ARG BUILDRC_VERSION=
ARG BIN_NAME=

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS golatest

FROM --platform=$BUILDPLATFORM walteh/buildrc:${BUILDRC_VERSION} as buildrc


FROM --platform=$BUILDPLATFORM alpine:latest AS alpinelatest
FROM --platform=$BUILDPLATFORM busybox:musl AS musl

FROM golatest AS gobase
COPY --from=xx / /
COPY --from=buildrc /usr/bin/ /usr/bin/
RUN apk add --no-cache file git bash
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0
WORKDIR /src

##################################################################
# BUILD
##################################################################

FROM gobase AS metarc
ARG TARGETPLATFORM
RUN --mount=type=bind,target=/src,readonly \
	buildrc full --git-dir=/src --files-dir=/meta

FROM scratch AS meta
COPY --link --from=metarc /meta /

FROM gobase AS builder
ARG TARGETPLATFORM
COPY --link --from=meta . /meta
RUN --mount=type=bind,target=. \
	--mount=type=cache,target=/root/.cache \
	--mount=type=cache,target=/go/pkg/mod  <<EOT
  	set -e
	export CGO_ENABLED=0
 	xx-go --wrap;
	GO_PKG=$(cat /meta/go-pkg);
	LDFLAGS="-s -w -X ${GO_PKG}/version.Version=$(cat /meta/version) -X ${GO_PKG}/version.Revision=$(cat /meta/revision) -X ${GO_PKG}/version.Package=${GO_PKG}";
	go build -mod vendor -trimpath -ldflags "$LDFLAGS" -o /out/$(cat /meta/executable) ./cmd;
  	xx-verify --static /out/$(cat /meta/executable);
EOT

FROM musl AS symlink
COPY --link --from=meta . /meta
RUN <<EOT
	set -e -x -o pipefail
	mkdir -p /out/symlink
	ln -s ../$(cat /meta/executable) /out/symlink/executable
EOT

FROM scratch AS build-unix
COPY --from=builder /out /
COPY --from=symlink /out /

FROM build-unix AS build-darwin
FROM build-unix AS build-linux
FROM build-unix AS build-freebsd
FROM build-unix AS build-openbsd
FROM build-unix AS build-netbsd
FROM build-unix AS build-ios

FROM scratch AS build-windows
COPY --from=builder /out/ /
COPY --from=symlink /out /

FROM build-$TARGETOS AS build
# enable scanning for this stage
ARG BUILDKIT_SBOM_SCAN_STAGE=true
COPY --from=meta /buildrc.json /


##################################################################
# TESTING
##################################################################

FROM gobase AS gotestsum
ARG GOTESTSUM_VERSION
ENV GOFLAGS=
RUN --mount=target=/root/.cache,type=cache <<EOT
	GOBIN=/out/ go install "gotest.tools/gotestsum@${GOTESTSUM_VERSION}" &&
	/out/gotestsum --version
EOT

FROM gobase AS test2json
ARG GOTESTSUM_VERSION
ENV GOFLAGS=
RUN --mount=target=/root/.cache,type=cache <<EOT
	CGO_ENABLED=0 go build -o /out/test2json -ldflags="-s -w" cmd/test2json
EOT

FROM gobase AS test-builder
ARG BIN_NAME
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev libc6-compat clang llvm llvm-dev llvm-static
RUN mkdir -p /out
RUN --mount=type=bind,target=. \
	--mount=type=cache,target=/root/.cache \
	--mount=type=cache,target=/go/pkg/mod \
	for dir in $(go list -test -f '{{if or .ForTest}}{{.Dir}}{{end}}' ./...); do \
	pkg=$(echo $dir | sed -e 's/.*\///') && \
	echo "========== [pkg:${pkg}] ==========" && \
	go test -c -v -cover -fuzz -race -vet='' -covermode=atomic -mod=vendor "$dir" -o /out; \
	done


FROM scratch AS test-build
COPY --from=test-builder /out /tests
COPY --from=gotestsum /out /bins
COPY --from=test2json /out /bins

FROM alpinelatest AS case
ARG NAME= ARGS= E2E= FUZZ=
COPY --from=test-build /bins /bins
COPY --from=test-build /tests /bins
COPY --from=build . /bins

RUN <<EOT
	set -e -x -o pipefail
	mkdir -p /dat

	echo "${NAME}" > /dat/name
	echo "${ARGS}" > /dat/args
	echo "${E2E}" > /dat/e2e

	for file in /bins/*; do	chmod +x $file;	done
EOT

FROM alpinelatest AS test
ARG GO_VERSION
ENV GOVERSION=${GO_VERSION}
RUN apk add --no-cache jq
COPY --from=case /bins /usr/bin
COPY --from=case /dat /dat
ENV PKGS=
ENTRYPOINT for PKG in $(echo "${PKGS}" | jq -r '.[]' || echo "$PKGS"); do \
	export E2E=$(cat /dat/e2e) && \
	funcs="-" && \
	name=$(cat /dat/name) && \
	if [ "${name}" = "fuzz" ]; then funcs=$("/usr/bin/${PKG}.test" -test.list=Fuzz 2> /dev/null); fi && \
	for FUNC in $funcs; do \
	echo ""  && \
	if [ "${name}" = "fuzz" ] && [ "${FUNC}" = "-" ]; then continue; fi && \
	echo "========== [pkg:${PKG}] [type:$(cat /dat/name)] $(if [ "${name}" = "fuzz" ]; then echo "[func:${FUNC}] "; fi) =========="; \
	filename=$(if [ "${name}" = "fuzz" ]; then echo "${PKG}-fuzz-${FUNC}"; else echo "${PKG}-${name}"; fi) && \
	fuzzfunc=$(if [ "${name}" = "fuzz" ]; then echo "-test.fuzz=${FUNC} -test.run=${FUNC}"; fi) && \
	/usr/bin/gotestsum --format=standard-verbose \
	--jsonfile=/out/go-test-report-${filename}.json \
	--junitfile=/out/junit-report-${filename}.xml \
	--raw-command --  /usr/bin/test2json -t -p ${PKG}  /usr/bin/${PKG}.test $(cat /dat/args) -test.bench=. -test.timeout=10m  ${fuzzfunc} \
	-test.v -test.coverprofile=/out/coverage-report-${filename}.txt \
	-test.outputdir=/out; done; done && echo ""

##################################################################
# RELEASE
##################################################################

FROM alpinelatest AS packager
RUN apk add --no-cache file tar jq
COPY --link --from=build . /src/
RUN <<EOT
	set -e -x -o pipefail
	if [ -f /src/buildrc.json ]; then
		searchdir="/src/"
	else
		searchdir="/src/*/"
	fi
	mkdir -p /out
	for pdir in ${searchdir}; do
		(
			cd "${pdir}"
			artifact="$(jq -r '.artifact' ./buildrc.json)"
			tar -czf "/out/${artifact}.tar.gz" .
		)
	done

	(
		cd /out
		find . -type f \( -name '*.tar.gz' \) -exec sha256sum -b {} \; >./checksums.txt
		sha256sum -c checksums.txt
	)
EOT

FROM scratch AS package
COPY --link --from=packager /out/ /

##################################################################
# IMAGE
##################################################################

FROM scratch AS entry
ARG BUILDKIT_SBOM_SCAN_STAGE=true
ARG TARGETOS TARGETARCH TARGETVARIANT BIN_NAME
ARG TGT=${TARGETOS}_${TARGETARCH}*${TARGETVARIANT}
COPY --from=build /symlink* /usr/bin/symlink/
COPY --from=build /${BIN_NAME}* /usr/bin/
COPY --from=build /*.json /usr/bin/
COPY --from=build /${TGT} /usr/bin/
ENTRYPOINT ["/usr/bin/symlink/executable"]

