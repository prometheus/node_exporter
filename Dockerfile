ARG ARCH="amd64"
ARG OS="linux"
FROM ubuntu:20.04
RUN \
    apt-get update && apt-get install -y --no-install-recommends libvpx  1.13.1
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/node_exporter /bin/node_exporter

EXPOSE      9100
USER        nobody
ENTRYPOINT  [ "/bin/node_exporter" ]
