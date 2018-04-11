# Run tests
FROM        quay.io/prometheus/golang-builder:1.9-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY  . /go/src/github.com/prometheus/node_exporter
WORKDIR /go/src/github.com/prometheus/node_exporter
RUN make promu
RUN make

# Run Build
FROM        quay.io/prometheus/golang-builder:1.9-base
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY  . /go/src/github.com/prometheus/node_exporter
WORKDIR /go/src/github.com/prometheus/node_exporter
RUN make promu
RUN make build

# Make docker image
FROM        quay.io/prometheus/busybox:glibc
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY --from=1 /go/src/github.com/prometheus/node_exporter/node_exporter /bin/node_exporter

EXPOSE      9100
USER        nobody
ENTRYPOINT  [ "/bin/node_exporter" ]
