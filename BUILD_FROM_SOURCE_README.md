# Build Instructions
---
## Build

`docker image build -f ./Dockerfile_sauron --build-arg GOLANG_IMG="<docker-image-repo>/oraclelinux-golang:1.12.10_1" --build-arg ORACLE_LINUX_IMG="<docker-image-repo>/oraclelinux-mirror:7-slim" -t <docker-image-repo>/node-exporter:0.18.1-1 .`

## Push to OCIR

`docker image push <docker-image-repo>/node-exporter:0.18.1-1`