# Build stage
FROM golang:alpine AS build

RUN apk update && apk add --no-cache git

ADD . /src

WORKDIR /src

ENV CGO_ENABLED 0

RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w" \
    -o /node_exporter .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /node_exporter /node_exporter

EXPOSE 9100

ENTRYPOINT ["/node_exporter"]
