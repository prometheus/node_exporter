JSONNET_FMT := jsonnetfmt -n 2 --max-blank-lines 2 --string-style s --comment-style s

all: fmt node_alerts.yaml node_rules.yaml dashboards_out lint

fmt:
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		xargs -n 1 -- $(JSONNET_FMT) -i

node_alerts.yaml: mixin.libsonnet config.libsonnet $(wildcard alerts/*)
	jsonnet -S alerts.jsonnet > $@

node_rules.yaml: mixin.libsonnet config.libsonnet $(wildcard rules/*)
	jsonnet -S rules.jsonnet > $@

dashboards_out: mixin.libsonnet config.libsonnet $(wildcard dashboards/*)
	@mkdir -p dashboards_out
	jsonnet -J vendor -m dashboards_out dashboards.jsonnet

lint: node_alerts.yaml node_rules.yaml
	find . -name 'vendor' -prune -o -name '*.libsonnet' -print -o -name '*.jsonnet' -print | \
		while read f; do \
			$(JSONNET_FMT) "$$f" | diff -u "$$f" -; \
		done

	promtool check rules node_alerts.yaml node_rules.yaml

clean:
	rm -rf dashboards_out node_alerts.yaml node_rules.yaml
