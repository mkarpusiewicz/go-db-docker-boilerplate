ARG GOLANG_VERSION=1.12.9-alpine

FROM golang:$GOLANG_VERSION

RUN apk update && \
    apk add --no-cache \
    git make tree

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN pwd && tree -ahC