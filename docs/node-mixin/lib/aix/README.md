# AIX exporter observability lib

This jsonnet observability lib can be used to generate observability package for node exporter(AIX).

## Import

```sh
jb init
jb install https://github.com/prometheus/node_exporter/docs/node-mixin/lib/aix
```

## Examples

### Example 1: Basic example

You can use observ-lib to fill in monitoring-mixin structure:

```jsonnet
// mixin.libsonnet file
local aixlib = import 'aix/main.libsonnet';

local aix =
  aixlib.new()
  + aixlib.withConfigMixin({
    filteringSelector: 'job=~".*aix.*"',
    groupLabels: ['job'],
    instanceLabels: ['instance'],
    dashboardNamePrefix: 'AIX / ',
    dashboardTags: ['aix-mixin'],
    uid: 'aix',
    // enable loki logs
    enableLokiLogs: true,
  });

{
  grafanaDashboards+:: aix.grafana.dashboards,
  prometheusAlerts+:: aix.prometheus.alerts,
  prometheusRules+:: aix.prometheus.recordingRules,
}

```
For more examples see [node-mixin/lib/linux](../linux).
