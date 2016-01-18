FROM wehkamp/alpine:3.2

LABEL container.name="wehkamp/prometheus-node-exporter:latest"

ENV  GOPATH /go
ENV APPPATH $GOPATH/src/

WORKDIR $APPPATH

RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc

ADD *.go $APPPATH/

RUN go get -d
RUN go build -o /bin/node-exporter

RUN apk del --purge build-deps && rm -rf $GOPATH

EXPOSE      9100
ENTRYPOINT ["/bin/node-exporter", "-collector.filesystem.ignored-mount-points", "^/(sys|proc|dev|host|etc)($|/)"]
#ENTRYPOINT ["/bin/node-exporter"]
