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

STATICCHECK_IGNORE =

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
