ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"
LABEL org.opencontainers.image.authors="The Prometheus Authors"
LABEL org.opencontainers.image.vendor="Prometheus"
LABEL org.opencontainers.image.title="node_exporter"
LABEL org.opencontainers.image.description="Prometheus exporter for hardware and OS metrics exposed by *NIX kernels"
LABEL org.opencontainers.image.source="https://github.com/prometheus/node_exporter"
LABEL org.opencontainers.image.url="https://github.com/prometheus/node_exporter"
LABEL org.opencontainers.image.documentation="https://github.com/prometheus/node_exporter"
LABEL org.opencontainers.image.licenses="Apache License 2.0"
LABEL io.prometheus.image.variant="busybox"

COPY .build/${OS}-${ARCH}/node_exporter /bin/node_exporter

EXPOSE      9100
USER        nobody
ENTRYPOINT  [ "/bin/node_exporter" ]
