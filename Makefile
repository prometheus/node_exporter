VERSION  := 0.7.1

SRC      := $(wildcard *.go)
TARGET   := node_exporter

OS   := $(subst Darwin,darwin,$(subst Linux,linux,$(shell uname)))
ARCH := $(subst x86_64,amd64,$(shell uname -m))

GOOS   ?= $(OS)
GOARCH ?= $(ARCH)
GOPKG  := go1.3.3.$(OS)-$(ARCH).tar.gz
GOROOT ?= $(CURDIR)/.deps/go
GOPATH ?= $(CURDIR)/.deps/gopath
GOCC   := $(GOROOT)/bin/go
GOLIB  := $(GOROOT)/pkg/$(GOOS)_$(GOARCH)
GO     := GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GOCC)

SUFFIX  := $(GOOS)-$(GOARCH)
BINARY  := $(TARGET)
ARCHIVE := $(TARGET)-$(VERSION).$(SUFFIX).tar.gz

default:
	go build $(GOFLAGS)

.deps/$(GOPKG):
	mkdir -p .deps
	curl -o .deps/$(GOPKG) https://storage.googleapis.com/golang/$(GOPKG)

$(GOCC): .deps/$(GOPKG)
	tar -C .deps -xzf .deps/$(GOPKG)
	touch $@

dependencies: $(SRC)
	$(GO) get -d

$(BINARY): $(GOCC) $(SRC) dependencies
	$(GO) build $(GOFLAGS) -o $@

$(ARCHIVE): $(BINARY)
	tar -czf $@ $<

release: REMOTE     ?= $(error "can't release, REMOTE not set")
release: REMOTE_DIR ?= $(error "can't release, REMOTE_DIR not set")
release: $(ARCHIVE)
	scp $< $(REMOTE):$(REMOTE_DIR)/$(ARCHIVE)

clean:
	rm -rf bin

.PHONY: dependencies clean release
