FROM       golang:onbuild
MAINTAINER Prometheus Team <prometheus-developers@googlegroups.com>

ENTRYPOINT [ "go-wrapper", "run" ]
EXPOSE     9100
