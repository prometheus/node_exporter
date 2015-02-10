FROM       golang:onbuild
MAINTAINER Prometheus Team <prometheus-developers@googlegroups.com>

ENTRYPOINT [ "go-wrapper", "run" ]
CMD        [ "-logtostderr" ]
EXPOSE     9100
