SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

# Require the Go compiler/toolchain to be installed
ifeq (, $(shell which go 2>/dev/null))
$(error No 'go' found in $(PATH), please install the Go compiler for your system)
endif

.DEFAULT_GOAL: generate

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test -race ./...

.PHONY: testv
testv:
	go test -v -race ./...

.PHONY: modprobe
modprobe:
ifeq ($(shell id -u),0)
	modprobe -a nf_conntrack_ipv4 nf_conntrack_ipv6
else
	sudo modprobe -a nf_conntrack_ipv4 nf_conntrack_ipv6
endif

.PHONY: integration
integration: modprobe
ifeq ($(shell id -u),0)
	go test -v -race -coverprofile=cover-int.out -covermode=atomic -tags=integration ./...
else
	$(info Running integration tests under sudo..)
	go test -v -race -coverprofile=cover-int.out -covermode=atomic -tags=integration -exec sudo ./...
endif

.PHONY: coverhtml-integration
coverhtml-integration: integration
	go tool cover -html=cover-int.out

.PHONY: bench
bench:
	go test -bench=. ./...

.PHONY: bench-integration
bench-integration: modprobe
	go test -bench=. -tags=integration -exec sudo ./...

cover: cover.out
cover.out: $(SOURCES)
	go test -coverprofile=cover.out -covermode=atomic ./...
	go tool cover -func=cover.out

.PHONY: coverhtml
coverhtml: cover
	go tool cover -html=cover.out

.PHONY: check
check: test cover
	go vet ./...
	megacheck ./...
	golint -set_exit_status ./...
