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

# Ensure that 'all' is the default target otherwise it will be the first target from Makefile.common.
all::

include Makefile.common

PROMTOOL_VERSION ?= 2.5.0
PROMTOOL_URL     ?= https://github.com/prometheus/prometheus/releases/download/v$(PROMTOOL_VERSION)/prometheus-$(PROMTOOL_VERSION).$(GO_BUILD_PLATFORM).tar.gz
PROMTOOL         ?= $(FIRST_GOPATH)/bin/promtool

DOCKER_IMAGE_NAME       ?= node-exporter
MACH                    ?= $(shell uname -m)
DOCKERFILE              ?= Dockerfile

STATICCHECK_IGNORE =

ifeq ($(GOHOSTOS), linux)
	test-e2e := test-e2e
else
	test-e2e := skip-test-e2e
endif

# Use CGO for non-Linux builds.
ifeq ($(GOOS), linux)
	PROMU_CONF ?= .promu.yml
else
	ifndef GOOS
		ifeq ($(GOHOSTOS), linux)
			PROMU_CONF ?= .promu.yml
		else
			PROMU_CONF ?= .promu-cgo.yml
		endif
	else
		PROMU_CONF ?= .promu-cgo.yml
	endif
endif

PROMU := $(FIRST_GOPATH)/bin/promu --config $(PROMU_CONF)

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
	ifeq ($$(GOHOSTOS),linux)
		ifeq ($$(GOHOSTARCH),$1)
			GOARCH_CROSS = $2
			cross-test = test-32bit
		endif
	endif
endef

# By default, "cross" test with ourselves to cover unknown pairings.
$(eval $(call goarch_pair,amd64,386))
$(eval $(call goarch_pair,mips64,mips))
$(eval $(call goarch_pair,mips64el,mipsel))

all:: vet checkmetrics checkrules common-all $(cross-test) $(test-e2e)

.PHONY: test
test: collector/fixtures/sys/.unpacked
	@echo ">> running tests"
	$(GO) test -short $(test-flags) $(pkgs)

.PHONY: test-32bit
test-32bit: collector/fixtures/sys/.unpacked
	@echo ">> running tests in 32-bit mode"
	@env GOARCH=$(GOARCH_CROSS) $(GO) test $(pkgs)

.PHONY: skip-test-32bit
skip-test-32bit:
	@echo ">> SKIP running tests in 32-bit mode: not supported on $(GOHOSTOS)/$(GOHOSTARCH)"

collector/fixtures/sys/.unpacked: collector/fixtures/sys.ttar
	@echo ">> extracting sysfs fixtures"
	if [ -d collector/fixtures/sys ] ; then rm -r collector/fixtures/sys ; fi
	./ttar -C collector/fixtures -x -f collector/fixtures/sys.ttar
	touch $@

.PHONY: test-e2e
test-e2e: build collector/fixtures/sys/.unpacked
	@echo ">> running end-to-end tests"
	./end-to-end-test.sh

.PHONY: skip-test-e2e
skip-test-e2e:
	@echo ">> SKIP running end-to-end tests on $(GOHOSTOS)"

.PHONY: checkmetrics
checkmetrics: $(PROMTOOL)
	@echo ">> checking metrics for correctness"
	./checkmetrics.sh $(PROMTOOL) $(e2e-out)

.PHONY: checkrules
checkrules: $(PROMTOOL)
	@echo ">> checking rules for correctness"
	find . -name "*rules*.yml" | xargs -I {} $(PROMTOOL) check rules {}

.PHONY: docker
docker:
ifeq ($(MACH), ppc64le)
	$(eval DOCKERFILE=Dockerfile.ppc64le)
endif
	@echo ">> building docker image from $(DOCKERFILE)"
	@docker build --file $(DOCKERFILE) -t "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

.PHONY: test-docker
test-docker:
	@echo ">> testing docker image"
	./test_image.sh "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" 9100

.PHONY: promtool
promtool: $(PROMTOOL)

$(PROMTOOL):
	$(eval PROMTOOL_TMP := $(shell mktemp -d))
	curl -s -L $(PROMTOOL_URL) | tar -xvzf - -C $(PROMTOOL_TMP)
	mkdir -p $(FIRST_GOPATH)/bin
	cp $(PROMTOOL_TMP)/prometheus-$(PROMTOOL_VERSION).$(GO_BUILD_PLATFORM)/promtool $(FIRST_GOPATH)/bin/promtool
	rm -r $(PROMTOOL_TMP)

