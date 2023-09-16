FROM golang:1.21.0-alpine

ARG PORT
ARG HEALTH_CHECK_PORT

RUN apk add --update make musl-dev gcc libc-dev binutils-gold

ENV GO111MODULE=on
ENV GOPATH=/go

WORKDIR /go/src/worker

COPY go.mod go.mod
COPY go.sum go.sum
COPY . .
RUN go mod download

# Enable this if some tools needed in docker
# RUN GOBIN=$GOBIN make install-tools

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN TARGET_DIR=/app make build

EXPOSE $PORT
EXPOSE $HEALTH_CHECK_PORT

ENTRYPOINT make migrate up && TARGET_DIR=/app make run

