FROM openshift/origin-release:golang-1.10 AS builder

WORKDIR /go/src/github.com/prometheus/node_exporter
COPY . .
RUN make build

FROM  openshift/origin-base

MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY --from=builder /go/src/github.com/prometheus/node_exporter/node_exporter /bin/node_exporter

EXPOSE      9100
USER        nobody
ENTRYPOINT  [ "/bin/node_exporter" ]
