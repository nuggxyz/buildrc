# syntax=docker/dockerfile:experimental

FROM golang:latest AS builder

ARG NAME
ARG DIR

RUN ls -la

COPY ./${NAME}-${GOOS}-${GOARCH} /bin/main

RUN GOOS=$(go env GOOS) && GOARCH=$(go env GOARCH)

FROM alpine:latest

COPY --from=builder /bin/main /main

ENTRYPOINT ["/main"]
CMD [""]
