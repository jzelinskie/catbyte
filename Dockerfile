FROM golang:1.13-alpine AS build-env

WORKDIR /go/src/catbyte

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
RUN go install -v .

FROM alpine:3.11
RUN apk add --no-cache ca-certificates

WORKDIR /opt/catbye
COPY . .
COPY --from=build-env /go/bin/* /usr/local/bin/
EXPOSE 80
ENTRYPOINT catbyte
