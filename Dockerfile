FROM golang:1.12.6 AS build
WORKDIR /go/src/github.com/synchthia/altair

ENV GOOS linux
ENV CGO_ENABLED 0

RUN go get -u -v github.com/golang/dep/cmd/dep
ADD Gopkg.lock Gopkg.lock
ADD Gopkg.toml Gopkg.toml
RUN dep ensure -v --vendor-only
COPY . .
RUN go build -a -installsuffix cgo -v -o /altair cmd/altair/main.go

FROM alpine
WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=build /altair /app/

ENTRYPOINT ["/app/altair"]
