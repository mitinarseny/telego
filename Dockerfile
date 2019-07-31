ARG _build_path=/bin/bot

FROM golang:1.12.6-alpine3.10 AS build-env

RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
        git \
    && go get -u \
        github.com/derekparker/delve/cmd/dlv

ARG _project_path=github.com/mitinarseny/telego
WORKDIR ${GOPATH}/src/${_project_path}

ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download


FROM build-env AS builder

COPY . .

ARG _path=.
WORKDIR ${_path}

ARG _build_path
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -o ${_build_path} .


FROM alpine:latest AS base_runner

RUN apk add --no-cache \
    ca-certificates \
    libc6-compat

ARG _build_path
COPY --from=builder ${_build_path} /bin/


FROM base_runner AS server

ENTRYPOINT ["/bin/bot"]


FROM base_runner AS debugger

COPY --from=build-env  /go/bin/dlv /bin/

EXPOSE 40000

ENTRYPOINT ["/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/bin/bot"]
