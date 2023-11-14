JSONNET_FMT := jsonnetfmt -n 2 --max-blank-lines 1 --string-style s --comment-style s

.PHONY: all
all: build

vendor: jsonnetfile.json
	jb install

.PHONY: build
build: vendor

.PHONY: fmt
fmt:
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		xargs -n 1 -- $(JSONNET_FMT) -i

.PHONY: lint
lint: build
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		while read f; do \
			$(JSONNET_FMT) "$$f" | diff -u "$$f" -; \
		done

