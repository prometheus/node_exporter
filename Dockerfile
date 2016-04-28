FROM        quay.io/prometheus/busybox:glibc
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY node_exporter /bin/node_exporter

EXPOSE      9100
ENTRYPOINT  [ "/bin/node_exporter" ]
