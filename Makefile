# Copyright 2015 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO     ?= GO15VENDOREXPERIMENT=1 go
GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
GOARCH := $(shell $(GO) env GOARCH)
GOHOSTARCH := $(shell $(GO) env GOHOSTARCH)

PROMTOOL    ?= $(GOPATH)/bin/promtool
PROMU       ?= $(GOPATH)/bin/promu
STATICCHECK ?= $(GOPATH)/bin/staticcheck
pkgs         = $(shell $(GO) list ./... | grep -v /vendor/)

PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)
DOCKER_IMAGE_NAME       ?= node-exporter
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
MACH                    ?= $(shell uname -m)
DOCKERFILE              ?= Dockerfile

STATICCHECK_IGNORE =

ifeq ($(GOHOSTARCH),amd64)
	# Only supported on amd64
	test-flags := -race
endif

ifeq ($(OS),Windows_NT)
    OS_detected := Windows
else
    OS_detected := $(shell uname -s)
endif

ifeq ($(OS_detected), Linux)
    test-e2e := test-e2e
else
    test-e2e := skip-test-e2e
endif

ifeq ($(MACH), ppc64le)
	e2e-out = collector/fixtures/e2e-ppc64le-output.txt
else
	e2e-out = collector/fixtures/e2e-output.txt
endif

# 64bit -> 32bit mapping for cross-checking. At least for amd64/386, the 64bit CPU can execute 32bit code but not the other way around, so we don't support cross-testing upwards.
cross-test = skip-test-32bit
define goarch_pair
	ifeq ($$(OS_detected),Linux)
		ifeq ($$(GOARCH),$1)
			GOARCH_CROSS = $2
			cross-test = test-32bit
		endif
	endif
endef

# By default, "cross" test with ourselves to cover unknown pairings.
$(eval $(call goarch_pair,amd64,386))
$(eval $(call goarch_pair,arm64,arm))
$(eval $(call goarch_pair,mips64,mips))
$(eval $(call goarch_pair,mips64el,mipsel))

all: format vet staticcheck checkmetrics build test $(cross-test) $(test-e2e)

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

test: collector/fixtures/sys/.unpacked
	@echo ">> running tests"
	$(GO) test -short $(test-flags) $(pkgs)

test-32bit: collector/fixtures/sys/.unpacked
	@echo ">> running tests in 32-bit mode"
	@env GOARCH=$(GOARCH_CROSS) $(GO) test $(pkgs)

skip-test-32bit:
	@echo ">> SKIP running tests in 32-bit mode: not supported on $(OS_detected)/$(GOARCH)"

collector/fixtures/sys/.unpacked: collector/fixtures/sys.ttar
	./ttar -C collector/fixtures -x -f collector/fixtures/sys.ttar
	touch $@

test-e2e: build collector/fixtures/sys/.unpacked
	@echo ">> running end-to-end tests"
	./end-to-end-test.sh

skip-test-e2e:
	@echo ">> SKIP running end-to-end tests on $(OS_detected)"

checkmetrics: $(PROMTOOL)
	@echo ">> checking metrics for correctness"
	./checkmetrics.sh $(PROMTOOL) $(e2e-out)

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

staticcheck: $(STATICCHECK)
	@echo ">> running staticcheck"
	@$(STATICCHECK) -ignore "$(STATICCHECK_IGNORE)" $(pkgs)

build: $(PROMU)
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(PREFIX)

tarball: $(PROMU)
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

docker:
ifeq ($(MACH), ppc64le)
	$(eval DOCKERFILE=Dockerfile.ppc64le)
endif
	@echo ">> building docker image from $(DOCKERFILE)"
	@docker build --file $(DOCKERFILE) -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

test-docker:
	@echo ">> testing docker image"
	./test_image.sh "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" 9100

$(GOPATH)/bin/promtool promtool:
	@GOOS= GOARCH= $(GO) get -u github.com/prometheus/prometheus/cmd/promtool

$(GOPATH)/bin/promu promu:
	@GOOS= GOARCH= $(GO) get -u github.com/prometheus/promu

$(GOPATH)/bin/staticcheck:
	@GOOS= GOARCH= $(GO) get -u honnef.co/go/tools/cmd/staticcheck


.PHONY: all style format build test test-e2e vet tarball docker promtool promu staticcheck checkmetrics

# Declaring the binaries at their default locations as PHONY targets is a hack
# to ensure the latest version is downloaded on every make execution.
# If this is not desired, copy/symlink these binaries to a different path and
# set the respective environment variables.
.PHONY: $(GOPATH)/bin/promtool $(GOPATH)/bin/promu $(GOPATH)/bin/staticcheck
