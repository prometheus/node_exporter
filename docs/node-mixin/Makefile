JSONNET_FMT := jsonnetfmt -n 2 --max-blank-lines 2 --string-style s --comment-style s

all: fmt node_alerts.yaml node_rules.yaml dashboards_out lint

fmt:
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		xargs -n 1 -- $(JSONNET_FMT) -i

node_alerts.yaml: mixin.libsonnet lib/linux/config.libsonnet $(wildcard lib/linux/alerts/*)
	jsonnet -J vendor -S -e 'std.manifestYamlDoc((import "mixin.libsonnet").prometheusAlerts)' > $@

node_rules.yaml: mixin.libsonnet lib/linux/config.libsonnet $(wildcard lib/linux/rules/*)
	jsonnet -J vendor -S -e 'std.manifestYamlDoc((import "mixin.libsonnet").prometheusRules)' > $@

dashboards_out: mixin.libsonnet lib/linux/config.libsonnet lib/linux/dashboards.libsonnet $(wildcard lib/linux/*)
	@mkdir -p dashboards_out/linux
	jsonnet -J vendor -m dashboards_out/linux -e '(import "mixin.libsonnet").grafanaDashboards'

lint: node_alerts.yaml node_rules.yaml
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		while read f; do \
			$(JSONNET_FMT) "$$f" | diff -u "$$f" -; \
		done

	promtool check rules node_alerts.yaml node_rules.yaml

.PHONY: jb_install
jb_install:
	jb install

clean:
	rm -rf dashboards_out node_alerts.yaml node_rules.yaml
