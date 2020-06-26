ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:glibc
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/node_exporter /bin/node_exporter

# Add license files to the image
COPY LICENSE README.md /license/

EXPOSE      9100
USER        nobody
ENTRYPOINT  [ "/bin/node_exporter" ]
