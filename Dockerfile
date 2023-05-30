FROM golang:1.20.3 AS build
WORKDIR /go/src/github.com/synchthia/altair

ENV GOOS linux
ENV CGO_ENABLED 0

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a -v -o /altair cmd/altair/main.go

FROM alpine

RUN apk --no-cache add tzdata
COPY --from=build /altair /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/altair"]
