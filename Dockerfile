# syntax=docker/dockerfile:experimental

FROM golang:latest AS builder

COPY . /ayes

RUN ls -la /ayes

RUN GOOS=$(go env GOOS) && GOARCH=$(go env GOARCH) && mv /ayes/buildrc-${GOOS}-${GOARCH} /bin/main

FROM alpine:latest

COPY --from=builder /bin/main /main

ENTRYPOINT ["/main"]

CMD [""]
