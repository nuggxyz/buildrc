# syntax=docker/dockerfile:experimental

FROM golang:latest AS builder

ARG name
ARG dir

COPY ${dir} /build

RUN GOOS=$(go env GOOS) && GOARCH=$(go env GOARCH) && cp /build/${name}-${GOOS}-${GOARCH} /bin/main

FROM alpine:latest

COPY --from=builder /bin/main /main

ENTRYPOINT ["/main"]
CMD [""]
