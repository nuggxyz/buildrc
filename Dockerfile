# syntax=docker/dockerfile:experimental

FROM golang:latest AS builder

COPY . /exec

RUN GOOS=$(go env GOOS) && GOARCH=$(go env GOARCH) && /exec/buildrc-${GOOS}-${GOARCH} /bin/main

FROM alpine:latest

COPY --from=builder /bin/main /main

ENTRYPOINT ["/main"]

CMD [""]
