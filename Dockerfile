FROM golang:1.23.5 AS build

ARG CGO_ENABLED=1

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-w" -o /alert-forge ./cmd


FROM debian:12.9-slim AS final

RUN apt update && apt -y install ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

LABEL maintainer="soerenschneider"
USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /alert-forge /alert-forge

ENTRYPOINT ["/alert-forge"]
