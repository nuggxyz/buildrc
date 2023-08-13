# syntax=docker/dockerfile:experimental

FROM golang:latest AS builder

# COPY . /exec

COPY . /exec

# RUN cd /exec && GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -pgo=auto -v -installsuffix cgo -ldflags "-s -w" -o /bin/main "./cmd"
RUN GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) /exec/buildrc-${GOOS}-${GOARCH} /bin/main

FROM alpine:latest

COPY --from=builder /bin/main /main
# COPY --from=tester /exec/test-status /test-status

ENTRYPOINT ["/main"]

CMD [""]
