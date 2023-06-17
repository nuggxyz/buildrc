# syntax=docker/dockerfile:experimental

FROM golang:latest AS builder

ARG NAME
ARG DIR

RUN ls -la

COPY . /exec

RUN GOOS=$(go env GOOS) && GOARCH=$(go env GOARCH) && cp ./exec/${NAME}-${GOOS}-${GOARCH} /bin/main

FROM alpine:latest

COPY --from=builder /bin/main /main

ENTRYPOINT ["/main"]
CMD [""]
