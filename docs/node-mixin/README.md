# Node Mixin

_This is a work in progress. We aim for it to become a good role model for alerts
and dashboards eventually, but it is not quite there yet._

The Node Mixin is a set of configurable, reusable, and extensible alerts and
dashboards based on the metrics exported by the Node Exporter and logs by Loki(optional). The mixin creates
recording and alerting rules for Prometheus and suitable dashboard descriptions
for Grafana.

To use them, you need to have `jsonnet` (v0.16+) and `jb` installed. If you
have a working Go development environment, it's easiest to run the following:
```bash
$ go install github.com/google/go-jsonnet/cmd/jsonnet@latest
$ go install github.com/google/go-jsonnet/cmd/jsonnetfmt@latest
$ go install github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb@latest
```

Next, install the dependencies by running the following command in this
directory:
```bash
$ jb install
```

You can then build the Prometheus rules files `node_alerts.yaml` and
`node_rules.yaml`:
```bash
$ make node_alerts.yaml node_rules.yaml
```

You can also build a directory `dashboard_out` with the JSON dashboard files
for Grafana:
```bash
$ make dashboards_out
```

Note that some of the generated dashboards require recording rules specified in
the previously generated `node_rules.yaml`.

For more advanced uses of mixins, see
https://github.com/monitoring-mixins/docs.

## Loki Logs configuration

To enable logs support in Node Mixin, enable them in config.libsonnet first:

```
{
  _config+:: {
    enableLokiLogs: true,
  },
}

```

then run
```bash
$ make build
```

This would generate extra Logs row on dashboards.

For proper logs correlation, you need to make sure that `job` and `instance` labels values match for both node_exporter metrics and logs, collected by [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) or [Grafana Agent](https://grafana.com/docs/grafana-cloud/agent/).

To scrape system logs, the following promtail config snippet can be used for `job=integrations/node` and `instance=host-01`:

```yaml
configs:
  - name: integrations
    scrape_configs:
    - job_name: integrations/node_exporter_journal_scrape
      journal:
        max_age: 24h
        labels:
          instance: host-01
          job: integrations/node
      relabel_configs:
      - source_labels: ['__journal__systemd_unit']
        target_label: 'unit'
      - source_labels: ['__journal__boot_id']
        target_label: 'boot_id'
      - source_labels: ['__journal__transport']
        target_label: 'transport'
      - source_labels: ['__journal_priority_keyword']
        target_label: 'level'
```
