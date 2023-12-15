FROM golang:1.21 AS builder

ENV GOOS=linux 
ENV GOARCH=arm64

WORKDIR /build
ADD . /build

RUN go version && go mod download && \
    go build -a -o dht22-exporter .

FROM  alpine:latest
LABEL maintainer="Max Oreshnikov <m.oreshnikov@gmail.com>"

RUN apk add libc6-compat

COPY --from=builder /build/dht22-exporter /usr/local/bin

EXPOSE 9543

ENTRYPOINT ["/usr/local/bin/dht22-exporter"]
