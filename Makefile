VERSION  := 0.7.1

SRC      := $(wildcard *.go)
TARGET   := node_exporter

OS         := $(subst Darwin,darwin,$(subst Linux,linux,$(subst FreeBSD,freebsd,$(shell uname))))
ARCH       := $(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m)))

# The release engineers apparently need to key their binary artifacts to the
# Mac OS X release family.
MAC_OS_X_VERSION ?= 10.8

GOOS       ?= $(OS)
GOARCH     ?= $(ARCH)

ifeq ($(GOOS),darwin)
RELEASE_SUFFIX ?= -osx$(MAC_OS_X_VERSION)
else
RELEASE_SUFFIX ?=
endif

GO_VERSION ?= 1.4.1
GOURL      ?= https://golang.org/dl
GOPKG      ?= go$(GO_VERSION).$(GOOS)-$(GOARCH)$(RELEASE_SUFFIX).tar.gz
GOROOT     := $(CURDIR)/.deps/go
GOPATH     := $(CURDIR)/.deps/gopath
GOCC       := $(GOROOT)/bin/go
GOLIB      := $(GOROOT)/pkg/$(GOOS)_$(GOARCH)
GO         := GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GOCC)

SUFFIX  := $(GOOS)-$(GOARCH)
BINARY  := $(TARGET)
ARCHIVE := $(TARGET)-$(VERSION).$(SUFFIX).tar.gz
SELFLINK := $(GOPATH)/src/github.com/prometheus/node_exporter

default: $(BINARY)

.deps/$(GOPKG):
	mkdir -p .deps
	curl -o .deps/$(GOPKG) -L $(GOURL)/$(GOPKG)

$(GOCC): .deps/$(GOPKG)
	tar -C .deps -xzf .deps/$(GOPKG)
	touch $@

$(SELFLINK):
	mkdir -p $(GOPATH)/src/github.com/prometheus
	ln -s $(CURDIR) $(SELFLINK)

dependencies: $(SRC) $(SELFLINK)
	$(GO) get -d

$(BINARY): $(GOCC) $(SRC) dependencies
	$(GO) build $(GOFLAGS) -o $@

$(ARCHIVE): $(BINARY)
	tar -czf $@ $<

release: REMOTE     ?= $(error "can't release, REMOTE not set")
release: REMOTE_DIR ?= $(error "can't release, REMOTE_DIR not set")
release: $(ARCHIVE)
	scp $< $(REMOTE):$(REMOTE_DIR)/$(ARCHIVE)

test: $(GOCC) dependencies
	$(GO) test ./...

clean:
	rm -rf node_exporter .deps

.PHONY: clean default dependencies release test
