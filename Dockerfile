FROM golang:1.22 AS builder
ARG TARGETARCH
ADD . /src
WORKDIR /src
RUN go mod vendor
RUN mkdir -p /src/bin
RUN GOOS=linux GOARCH="${TARGETARCH}" go build -o bin/konterfai cmd/konterfai/main.go

FROM debian:sid-slim
COPY entrypoint.sh /entrypoint.sh
COPY --from=builder /src/bin/konterfai /usr/local/bin/konterfai

EXPOSE 8080/tcp
EXPOSE 8081/tcp

ENTRYPOINT ["/entrypoint.sh"]