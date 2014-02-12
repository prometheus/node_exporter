FROM       ubuntu:13.10
MAINTAINER Prometheus Team <prometheus-developers@googlegroups.com>

ENV        COLLECTORS NativeCollector
RUN        apt-get update && apt-get install -yq curl git mercurial make
RUN        curl -s https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV        PATH    /usr/local/go/bin:$PATH
ENV        GOPATH  /go

ADD        . /usr/src/node_exporter
RUN        cd /usr/src/node_exporter && make && cp node_exporter /
RUN        printf '{ "scrapeInterval": 10, "attributes": {} }' > \
                  node_exporter.conf

ENTRYPOINT [ "/node_exporter" ]
EXPOSE     8080
