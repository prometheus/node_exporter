JSONNET_FMT := jsonnetfmt -n 2 --max-blank-lines 1 --string-style s --comment-style s

.PHONY: all
all: build dashboards_out

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
	mixtool lint mixin.libsonnet

dashboards_out: mixin.libsonnet $(wildcard dashboards/*)
	@mkdir -p dashboards_out
	mixtool generate dashboards -J lib -J vendor mixin.libsonnet -d dashboards_out

.PHONY: clean
clean:
	rm -rf dashboards_out prometheus_alerts.yaml
