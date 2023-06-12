# syntax=docker/dockerfile:experimental

FROM golang:1.20 AS builder

COPY ./build /build

RUN GOOS=$(go env GOOS) && \
	GOARCH=$(go env GOARCH) && \
	cp /build/${GOOS}-${GOARCH} /bin/main


FROM alpine:latest

COPY --from=builder /bin/main /main

ENTRYPOINT ["/main"]
CMD [""]
