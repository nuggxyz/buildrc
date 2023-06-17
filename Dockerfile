# syntax=docker/dockerfile:experimental

FROM golang:latest AS builder

ARG NAME
ARG DIR

# COPY ${DIR} /build

RUN GOOS=$(go env GOOS) && GOARCH=$(go env GOARCH) && cp ./${NAME}-${GOOS}-${GOARCH} /bin/main

FROM alpine:latest

COPY --from=builder /bin/main /main

ENTRYPOINT ["/main"]
CMD [""]
