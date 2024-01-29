#!/bin/bash -eo pipefail

docker run --entrypoint "rcodesign" --rm -v $(pwd)/.build:/build quay.io/prometheus/golang-builder:1.21-main sign /build/darwin-arm64/node_exporter
