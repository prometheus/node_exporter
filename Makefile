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

include Makefile.common

GO     ?= GO15VENDOREXPERIMENT=1 go
GOARCH := $(shell $(GO) env GOARCH)
GOHOSTARCH := $(shell $(GO) env GOHOSTARCH)

PROMTOOL    ?= $(FIRST_GOPATH)/bin/promtool

DOCKER_IMAGE_NAME       ?= node-exporter
MACH                    ?= $(shell uname -m)
DOCKERFILE              ?= Dockerfile

STATICCHECK_IGNORE =

ifeq ($(OS),Windows_NT)
    OS_detected := Windows
else
    OS_detected := $(shell uname -s)
endif

ifeq ($(GOHOSTARCH),amd64)
	ifeq ($(OS_detected),$(filter $(OS_detected),Linux FreeBSD Darwin Windows))
                # Only supported on amd64
                test-flags := -race
        endif
endif

ifeq ($(OS_detected), Linux)
    test-e2e := test-e2e
else
    test-e2e := skip-test-e2e
endif

e2e-out = collector/fixtures/e2e-output.txt
ifeq ($(MACH), ppc64le)
	e2e-out = collector/fixtures/e2e-64k-page-output.txt
endif
ifeq ($(MACH), aarch64)
	e2e-out = collector/fixtures/e2e-64k-page-output.txt
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
$(eval $(call goarch_pair,mips64,mips))
$(eval $(call goarch_pair,mips64el,mipsel))

all: style vet staticcheck checkmetrics build test $(cross-test) $(test-e2e)

test: collector/fixtures/sys/.unpacked
	@echo ">> running tests"
	$(GO) test -short $(test-flags) $(pkgs)

test-32bit: collector/fixtures/sys/.unpacked
	@echo ">> running tests in 32-bit mode"
	@env GOARCH=$(GOARCH_CROSS) $(GO) test $(pkgs)

skip-test-32bit:
	@echo ">> SKIP running tests in 32-bit mode: not supported on $(OS_detected)/$(GOARCH)"

collector/fixtures/sys/.unpacked: collector/fixtures/sys.ttar
	@echo ">> extracting sysfs fixtures"
	if [ -d collector/fixtures/sys ] ; then rm -r collector/fixtures/sys ; fi
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

docker:
ifeq ($(MACH), ppc64le)
	$(eval DOCKERFILE=Dockerfile.ppc64le)
endif
	@echo ">> building docker image from $(DOCKERFILE)"
	@docker build --file $(DOCKERFILE) -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

test-docker:
	@echo ">> testing docker image"
	./test_image.sh "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" 9100

$(FIRST_GOPATH)/bin/promtool promtool:
	@GOOS= GOARCH= $(GO) get -u github.com/prometheus/prometheus/cmd/promtool

.PHONY: test-e2e promtool checkmetrics

# Declaring the binaries at their default locations as PHONY targets is a hack
# to ensure the latest version is downloaded on every make execution.
# If this is not desired, copy/symlink these binaries to a different path and
# set the respective environment variables.
.PHONY: $(FIRST_GOPATH)/bin/promtool
