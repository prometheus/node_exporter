# MacOS exporter observability lib

This jsonnet observability lib can be used to generate observability package for node exporter(MacOS).

## Import

```sh
jb init
jb install https://github.com/prometheus/node_exporter/docs/node-mixin/lib/macos
```

## Examples

### Example 1: Basic example

You can use observ-lib to fill in monitoring-mixin structure:

```jsonnet
// mixin.libsonnet file
local macoslib = import 'macos/main.libsonnet';

local mac =
  macoslib.new()
  + macoslib.withConfigMixin({
    filteringSelector: 'job=~".*mac.*"',
    groupLabels: ['job'],
    instanceLabels: ['instance'],
    dashboardNamePrefix: 'MacOS / ',
    dashboardTags: ['macos-mixin'],
    uid: 'darwin',
    // enable loki logs
    enableLokiLogs: true,
  });

{
  grafanaDashboards+:: mac.grafana.dashboards,
  prometheusAlerts+:: mac.prometheus.alerts,
  prometheusRules+:: mac.prometheus.recordingRules,
}

```
For more examples see [node-mixin/lib/linux](../linux).

## Collectors used:

Grafana Agent or combination of node_exporter/promtail can be used in order to collect data required.

### Logs collection

Loki logs are used to populate logs dashboard and also for annotations.

To use logs, you need to opt-in, with setting `enableLokiLogs: true` in config.

See example above.

The following scrape snippet can be used in grafana-agent/promtail:

```yaml
    - job_name: integrations/node_exporter_direct_scrape
      static_configs:
      - targets:
        - localhost
        labels:
          __path__: /var/log/*.log
          instance: '<your-instance-name>'
          job: integrations/macos-node
      pipeline_stages:
      - multiline:
          firstline: '^([\w]{3} )?[\w]{3} +[\d]+ [\d]+:[\d]+:[\d]+|[\w]{4}-[\w]{2}-[\w]{2} [\w]{2}:[\w]{2}:[\w]{2}(?:[+-][\w]{2})?'
      - regex:
          expression: '(?P<timestamp>([\w]{3} )?[\w]{3} +[\d]+ [\d]+:[\d]+:[\d]+|[\w]{4}-[\w]{2}-[\w]{2} [\w]{2}:[\w]{2}:[\w]{2}(?:[+-][\w]{2})?) (?P<hostname>\S+) (?P<sender>.+?)\[(?P<pid>\d+)\]:? (?P<message>(?s:.*))$'
      - labels:
          sender:
          hostname:
          pid:
      - match:
          selector: '{sender!="", pid!=""}'
          stages:
            - template:
                source: message
                template: '{{ .sender }}[{{ .pid }}]: {{ .message }}'
            - labeldrop:
                - pid
            - output:
                source: message
```
