# Build Instructions

The base tag this release is branched from is `v0.18.1`

---
## Build

`docker image build -f ./Dockerfile_verrazzano -t <docker-image-repo>/node-exporter:0.18.1-1 .`

## Push to OCIR

`docker image push <docker-image-repo>/node-exporter:0.18.1-1`
