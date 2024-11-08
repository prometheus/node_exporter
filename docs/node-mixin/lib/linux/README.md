# Node exporter observability lib

This jsonnet observability lib can be used to generate observability package for node exporter.

## Import

```sh
jb init
jb install https://github.com/prometheus/node_exporter/docs/node-mixin/lib/linux
```

## Examples

### Example 1: Basic example

You can use observ-lib to fill in monitoring-mixin structure:

```jsonnet
// mixin.libsonnet file
local nodelib = import 'linux/main.libsonnet';

local linux =
  nodelib.new()
  + nodelib.withConfigMixin({
    filteringSelector: 'job=~".*node.*"',
    groupLabels: ['job'],
    instanceLabels: ['instance'],
    dashboardNamePrefix: 'Node exporter / ',
    dashboardTags: ['node-exporter-mixin'],
    uid: 'node',
    // enable loki logs
    enableLokiLogs: true,
  });

{
  grafanaDashboards+:: linux.grafana.dashboards,
  prometheusAlerts+:: linux.prometheus.alerts,
  prometheusRules+:: linux.prometheus.recordingRules,
}

```

### Example 2: Fill in monitoring-mixin with default config values and enable loki logs:


```jsonnet
// mixin.libsonnet file
local nodelib = import 'linux/main.libsonnet';

local linux =
  nodelib.new()
  + nodelib.withConfigMixin({
    enableLokiLogs: true,
  });

{
  grafanaDashboards+:: linux.grafana.dashboards,
  prometheusAlerts+:: linux.prometheus.alerts,
  prometheusRules+:: linux.prometheus.recordingRules,
}

```

### Example 3: Override some of default config values from file:


```jsonnet
// overrides.libsonnet
{
  // Memory utilzation (%) level on which to trigger the
  // 'NodeMemoryHighUtilization' alert.
  memoryHighUtilizationThreshold: 80,

  // Threshold for the rate of memory major page faults to trigger
  // 'NodeMemoryMajorPagesFaults' alert.
  memoryMajorPagesFaultsThreshold: 1000,

  // Disk IO queue level above which to trigger
  // 'NodeDiskIOSaturation' alert.
  diskIOSaturationThreshold: 20,
}

// mixin.libsonnet file
local configOverride = import './overrides.libsonnet';
local nodelib = import 'linux/main.libsonnet';

local linux =
  nodelib.new()
  + nodelib.withConfigMixin(configOverride);

{
  grafanaDashboards+:: linux.grafana.dashboards,
  prometheusAlerts+:: linux.prometheus.alerts,
  prometheusRules+:: linux.prometheus.recordingRules,
}

```

### Example 4: Modify specific panel before rendering dashboards

```jsonnet
local g = import './g.libsonnet';
// mixin.libsonnet file
local nodelib = import 'linux/main.libsonnet';

local linux =
  nodelib.new()
  + nodelib.withConfigMixin({
    filteringSelector: 'job=~".*node.*"',
    groupLabels: ['job'],
    instanceLabels: ['instance'],
    dashboardNamePrefix: 'Node exporter / ',
    dashboardTags: ['node-exporter-mixin'],
    uid: 'node',
  })
  + {
      grafana+: {
        panels+: {
          networkSockstatAll+:
            + g.panel.timeSeries.fieldConfig.defaults.custom.withDrawStyle('bars')
        }
      }
    };

{
  grafanaDashboards+:: linux.grafana.dashboards,
  prometheusAlerts+:: linux.prometheus.alerts,
  prometheusRules+:: linux.prometheus.recordingRules,
}

```

## Collectors used:

Grafana Agent or combination of node_exporter/promtail can be used in order to collect data required.

### Logs collection

Loki logs are used to populate logs dashboard and also for annotations.

To use logs, you need to opt-in, with setting `enableLokiLogs: true` in config.

See example above.

The following scrape snippet can be used in grafana-agent/promtail:

```yaml
    - job_name: integrations/node_exporter_journal_scrape
      journal:
        max_age: 24h
        labels:
          instance: '<your-instance-name>'
          job: integrations/node_exporter
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
