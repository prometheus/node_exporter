local nodelib = import './linux/main.libsonnet';
local linux =
  nodelib.new()
  + nodelib.withConfigMixin({
    filteringSelector: 'job=~".*node.*"',
    groupLabels: ['job'],
    instanceLabels: ['instance'],
    dashboardNamePrefix: 'Node exporter / ',
    dashboardTags: ['node-exporter-mixin'],
    uid: 'node',
  });
{
  grafanaDashboards+:: linux.grafana.dashboards,
  prometheusAlerts+:: linux.prometheus.alerts,
  prometheusRules+:: linux.prometheus.recordingRules,
}
