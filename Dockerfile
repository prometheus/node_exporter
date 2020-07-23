FROM        quay.io/prometheus/busybox:glibc
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

COPY node_exporter /bin/node_exporter

EXPOSE      9100
USER        nobody
ENTRYPOINT  [ "/bin/node_exporter" ]
