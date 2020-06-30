# Build Instructions
---
## Build

`docker image build -f ./Dockerfile_sauron -t <docker-image-repo>/node-exporter:0.18.1-1 .`

## Push to OCIR

`docker image push <docker-image-repo>/node-exporter:0.18.1-1`