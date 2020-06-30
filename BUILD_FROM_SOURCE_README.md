# Build Instructions
---
## Build

<<<<<<< HEAD
`docker image build -f ./Dockerfile_sauron -t <docker-image-repo>/node-exporter:0.18.1-1 .`

## Push to OCIR

`docker image push <docker-image-repo>/node-exporter:0.18.1-1`
=======
`docker image build -f ./Dockerfile_verrazzano --build-arg GOLANG_IMG="<docker-image-repo>/oraclelinux-golang:1.12.10_1" -t <docker-image-repo>//node-exporter:0.18.1-1 .`

## Push to OCIR

`docker image push <docker-image-repo>/node-exporter:0.18.1-1`
>>>>>>> 4488d588e41e5a5eb0e544a1f1fceba2fa5070e9
